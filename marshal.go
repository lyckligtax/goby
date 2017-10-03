package main

import (
	"bytes"
	"errors"
	)

type NEString string

func (s *NEString) UnmarshalJSON(b []byte) error {
	*s = NEString(bytes.Trim(b, " \"\n\t\r"))
	return nil
}
