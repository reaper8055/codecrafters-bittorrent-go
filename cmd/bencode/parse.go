// bencode is my own implementation of decoding bencode
// inspired from github.com/jackpal/bencode-go
package bencode

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strconv"
)

func Unmarshall(rd *bufio.Reader) (interface{}, error) {
	char, err := rd.ReadByte()
	if err != nil {
		return nil, err
	}

	switch char {

	case 'i':
		buffer, err := optimisticReadBytes(rd, 'e')
		if err != nil {
			return nil, err
		}
		intValue, err := strconv.ParseInt(string(buffer[:len(buffer)-1]), 10, 64)
		if err != nil {
			return nil, err
		}
		return intValue, nil
	case 'l':
		list := []interface{}{}
		for {
			char, err := rd.ReadByte()
			if err != nil {
				return nil, err
			}
			if char == 'e' {
				return list, nil
			}
			rd.UnreadByte()
			value, err := Unmarshall(rd)
			if err != nil {
				return nil, err
			}
			list = append(list, value)
		}
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		rd.UnreadByte()
		buffer, err := optimisticReadBytes(rd, ':')
		if err != nil {
			return nil, err
		}
		strLen, err := strconv.ParseInt(string(buffer[:len(buffer)-1]), 10, 64)
		b := make([]byte, int(strLen))
		_, err = io.ReadFull(rd, b)
		if err != nil {
			return nil, err
		}
		return string(b), nil
	case 'd':
		dictionary := map[string]interface{}{}
		for {
			c, err2 := rd.ReadByte()
			if err2 == nil {
				if c == 'e' {
					return dictionary, nil
				} else {
					rd.UnreadByte()
				}
			}
			value, err := Unmarshall(rd)
			if err != nil {
				return nil, err
			}

			key, ok := value.(string)
			if !ok {
				return nil, errors.New("bencode: non-string dictionary key")
			}

			value, err = Unmarshall(rd)
			if err != nil {
				return nil, err
			}

			dictionary[key] = value
		}
	}

	return nil, nil
}

func optimisticReadBytes(rd *bufio.Reader, delim byte) ([]byte, error) {
	buffered := rd.Buffered()
	var (
		buffer []byte
		err    error
	)

	if buffer, err = rd.Peek(buffered); err != nil {
		return nil, err
	}

	if idx := bytes.IndexByte(buffer, delim); idx >= 0 {
		return rd.ReadSlice(delim)
	}

	return rd.ReadBytes(delim)
}
