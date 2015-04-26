package restclient

import (
	"bufio"
	"errors"
	"strconv"
	"time"
)

func SkipSpaces(reader *bufio.Reader) error {
	return SkipWhile(reader, func(b byte) bool { return b != ' ' && b != '\n' && b != '\t' && b != '\r' })
}

func SkipTill(reader *bufio.Reader, val byte) error {
	return SkipWhile(reader, func(b byte) bool { return b != val })
}

func NextIf(reader *bufio.Reader, value byte) bool {
	bytes, err := reader.Peek(1)
	if err != nil {
		return false
	}
	if bytes[0] != value {
		return false
	}
	reader.Read(bytes)
	return true
}

func SkipWhile(reader *bufio.Reader, filter func(byte) bool) error {
	var bytes [1]byte
	for {
		bytes, err := reader.Peek(1)
		if err != nil {
			return err
		}
		if filter(bytes[0]) {
			// byte matched filter so swallow it
			reader.Read(bytes)
		} else {
			// filter failed so return now
			break
		}
	}
	return nil
}

/**
 * Reads a string while a match succeeds.
 */
func ReadWhile(reader *bufio.Reader, matcher func(byte) bool) ([]byte, error) {
	var nextByte [1]byte
	var bytes []byte
	for {
		nextByte, err := reader.Peek(1)
		if err != nil {
			return nil, err
		}
		if !matcher(nextByte[0]) {
			return bytes, nil
		} else {
			bytes = append(bytes, nextByte[0])
		}
	}
	return nil, nil
}

func EnsureOSq(reader *bufio.Reader) error {
	if err := SkipTill(reader, '['); err != nil {
		return err
	}
	if err := SkipSpaces(reader); err != nil {
		return err
	}
	return nil
}

func EnsureCSq(reader *bufio.Reader) error {
	if err := SkipSpaces(reader); err != nil {
		return err
	}
	if err := SkipTill(reader, ']'); err != nil {
		return err
	}
	return nil
}

func EnsureOCurly(reader *bufio.Reader) error {
	if err := SkipTill(reader, '{'); err != nil {
		return err
	}
	if err := SkipSpaces(reader); err != nil {
		return err
	}
	return nil
}

func EnsureCCurly(reader *bufio.Reader) error {
	if err := SkipSpaces(reader); err != nil {
		return err
	}
	if err := SkipTill(reader, '}'); err != nil {
		return err
	}
	return nil
}

func Read_string(reader *bufio.Reader, arg *string) error {
	if !NextIf(reader, '"') {
		return errors.New("Expected \"")
	}
	var err error
	*arg, err = reader.ReadString('"')
	if err != nil {
		return err
	}
	*arg = (*arg)[:len(*arg)-1]
	return nil
}

func Read_bool(reader *bufio.Reader, arg bool) error {
	var err error
	if arg {
		_, err = reader.Read([]byte("true"))
	} else {
		_, err = reader.Read([]byte("false"))
	}
	return err
}

func Read_int(reader *bufio.Reader, arg *int) error {
	var value int64
	if err := Read_int64(reader, &value); err != nil {
		return err
	}
	*arg = int(value)
	return nil
}

func Read_int64(reader *bufio.Reader, arg *int64) error {
	index := 0
	bytes, err := ReadWhile(reader, func(b byte) bool {
		index++
		return b >= '0' && b <= '9' || (index == 1 && b == '-')
	})
	if err == nil {
		value, err := strconv.ParseInt(string(bytes), 10, 32)
		if err != nil {
			return err
		}
		*arg = value
	}
	return err
}

func Read_time_Time(reader *bufio.Reader, time time.Time) error {
	bytes, err := time.MarshalJSON()
	if err != nil {
		return err
	}
	_, err = reader.Read(bytes)
	return err
}
