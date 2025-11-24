package utils

import "bytes"

type ErrorsMap map[string]error

func (e ErrorsMap) Error() string {
	buf := new(bytes.Buffer)
	for key, err := range e {
		buf.WriteString(key)
		buf.WriteString(": ")
		buf.WriteString(err.Error())
		buf.WriteString("\n")
	}
	return buf.String()
}

func (e ErrorsMap) HasError() error {
	if len(e) == 0 {
		return nil
	}
	return e
}
