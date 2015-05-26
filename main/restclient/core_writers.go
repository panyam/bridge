package restclient

import (
	"fmt"
	"io"
	"strings"
	"time"
)

func Write_string(writer io.Writer, arg string) error {
	arg = strings.Replace(arg, "\"", "\\\"", -1)
	arg = strings.Replace(arg, "\\", "\\\\", -1)
	writer.Write([]byte("\""))
	writer.Write([]byte(arg))
	writer.Write([]byte("\""))
	return nil
}

func Write_bool(writer io.Writer, arg bool) error {
	var err error
	if arg {
		_, err = writer.Write([]byte("true"))
	} else {
		_, err = writer.Write([]byte("false"))
	}
	return err
}

func Write_int(writer io.Writer, arg int) error {
	_, err := writer.Write([]byte(fmt.Sprintf("%d", arg)))
	return err
}

func Write_int64(writer io.Writer, arg int64) error {
	_, err := writer.Write([]byte(fmt.Sprintf("%d", arg)))
	return err
}

func Write_time_Time(writer io.Writer, time time.Time) error {
	bytes, err := time.MarshalJSON()
	if err != nil {
		return err
	}
	_, err = writer.Write(bytes)
	return err
}
