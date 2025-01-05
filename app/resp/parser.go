package resp

import (
	"strconv"
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
