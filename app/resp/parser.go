package resp

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

	return 0, RESP{}
}
