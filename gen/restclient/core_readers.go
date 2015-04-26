package restclient

import (
	"bufio"
)

func SkipSpaces(reader *bufio.Reader) error {
	return SkipWhile(reader, func(b byte) bool { return b != ' ' && b != '\n' && b != '\t' && b != '\r' })
}

func SkipTill(reader *bufio.Reader, val byte) error {
	return SkipWhile(reader, func(b byte) bool { return b != val })
}

func NextIf(reader *bufio.Reader, value byte) bool {
	var bytes [1]byte
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

func EnsureOSq(reader *bufio.Reader) error {
	if err := SkipTill(reader, '['); err != nil {
		return err
	}
	if err := SkipSpaces(reader); err != nil {
		return err
	}
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
