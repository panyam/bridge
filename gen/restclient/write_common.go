package restclient

import (
	"io"
	"strings"
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
	if arg {
		writer.Write([]byte("true"))
	} else {
		writer.Write([]byte("false"))
	}
	return nil
}
