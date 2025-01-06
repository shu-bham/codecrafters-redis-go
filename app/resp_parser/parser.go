package parser

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	Integer = ':'
	String  = '+'
	Bulk    = '$'
	Array   = '*'
	Error   = '-'
)

type Type byte

type RESP struct {
	Type  Type
	Raw   []byte
	Data  []byte
	Count int
}

// Parse returns the next resp in b and returns the number of bytes the
// took up the result.
func Parse(b []byte) (n int, resp RESP) {
	if len(b) == 0 {
		return 0, RESP{} // no data to return
	}

	resp.Type = Type(b[0])
	switch resp.Type {
	case Integer, String, Bulk, Array, Error:
	default:
		return 0, RESP{} // invalid kind
	}

	i := 1
	for ; ; i++ {
		if i == len(b) {
			return 0, RESP{} // not enough data
		}
		if b[i] == '\n' {
			if b[i-1] != '\r' {
				return 0, RESP{} // missing CR character
			}
			i++
			break
		}
	}

	resp.Raw = b[0:i]
	resp.Data = b[1 : i-2]

	if resp.Type == Integer {
		if len(resp.Data) == 0 {
			return 0, RESP{} // invalid integer
		}
		var j int
		if resp.Data[0] == '-' {
			if len(resp.Data) < 2 {
				return 0, RESP{} // invalid integer
			}
			j++
		}

		for ; j < len(resp.Data); j++ {
			if resp.Data[j] < '0' || resp.Data[j] > '9' {
				return 0, RESP{} // invalid integer
			}
		}

		return len(resp.Raw), resp
	}

	if resp.Type == String || resp.Type == Error {
		return len(resp.Raw), resp
	}

	var err error
	resp.Count, err = strconv.Atoi(string(resp.Data))
	if resp.Type == Bulk {
		if err != nil {
			return 0, RESP{} // invalid no of bytes
		}

		if resp.Count < 0 {
			resp.Data = nil
			resp.Count = 0
			return len(resp.Raw), resp
		}

		if i+resp.Count+2 > len(b) {
			return 0, RESP{} // not enough data
		}

		if b[i+resp.Count] != '\r' || b[i+resp.Count+1] != '\n' {
			return 0, RESP{} // invalid end of line
		}

		resp.Data = b[i : i+resp.Count]
		resp.Raw = b[0 : i+resp.Count+2]
		resp.Count = 0
		return len(resp.Raw), resp
	}

	if resp.Type == Array {
		if err != nil {
			return 0, RESP{} // invalid number of elements
		}

		var tn int
		sdata := b[i:]
		for j := 0; j < (resp.Count); j++ {
			rn, rresp := Parse(sdata)
			if rresp.Type == 0 {
				return 0, RESP{}
			}
			tn += rn
			sdata = sdata[rn:]
		}
		resp.Data = b[i : i+tn]
		resp.Raw = b[0 : i+tn]
		return len(resp.Raw), resp
	}

	return 0, RESP{}
}

func (r *RESP) ToStringArr() ([]string, error) {
	if r.Type != Array {
		return []string{}, fmt.Errorf("expected Array RESP, got %c", r.Type)
	}

	data := r.Data
	result := make([]string, 0)
	for i := 0; i < r.Count; i++ {
		n, resp := Parse(data)
		if resp.Type == 0 {
			return nil, fmt.Errorf("failed to parse element %d", i)
		}

		switch resp.Type {
		case String, Bulk:
			result = append(result, string(resp.Data))
		default:
			return nil, fmt.Errorf("element %d is not a string type (got %c)", i, resp.Type)
		}

		data = data[n:]
	}
	return result, nil
}

// appendPrefix will append a "$3\r\n" style redis prefix for a message.
func appendPrefix(b []byte, c byte, n int64) []byte {
	if n >= 0 && n <= 9 {
		return append(b, c, byte('0'+n), '\r', '\n')
	}
	b = append(b, c)
	b = strconv.AppendInt(b, n, 10)
	return append(b, '\r', '\n')
}

// AppendInt appends a Redis protocol int64 to the input bytes.
func AppendInt(b []byte, n int64) []byte {
	return appendPrefix(b, ':', n)
}

// AppendArray appends a Redis protocol array to the input bytes.
func AppendArray(b []byte, n int) []byte {
	return appendPrefix(b, '*', int64(n))
}

// AppendBulk appends a Redis protocol bulk byte slice to the input bytes.
func AppendBulk(b []byte, bulk []byte) []byte {
	b = appendPrefix(b, '$', int64(len(bulk)))
	b = append(b, bulk...)
	return append(b, '\r', '\n')
}

// AppendBulkString appends a Redis protocol bulk string to the input bytes.
func AppendBulkString(b []byte, bulk string) []byte {
	b = appendPrefix(b, '$', int64(len(bulk)))
	b = append(b, bulk...)
	return append(b, '\r', '\n')
}

// AppendString appends a Redis protocol string to the input bytes.
func AppendString(b []byte, s string) []byte {
	b = append(b, '+')
	b = append(b, stripNewlines(s)...)
	return append(b, '\r', '\n')
}

// AppendError appends a Redis protocol error to the input bytes.
func AppendError(b []byte, s string) []byte {
	b = append(b, '-')
	b = append(b, stripNewlines(s)...)
	return append(b, '\r', '\n')
}

func stripNewlines(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] == '\r' || s[i] == '\n' {
			s = strings.Replace(s, "\r", " ", -1)
			s = strings.Replace(s, "\n", " ", -1)
			break
		}
	}
	return s
}

// AppendBulkFloat appends a float64, as bulk bytes.
func AppendBulkFloat(dst []byte, f float64) []byte {
	return AppendBulk(dst, strconv.AppendFloat(nil, f, 'f', -1, 64))
}

// AppendBulkInt appends an int64, as bulk bytes.
func AppendBulkInt(dst []byte, x int64) []byte {
	return AppendBulk(dst, strconv.AppendInt(nil, x, 10))
}
