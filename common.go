package main

import (
	"crypto/md5"
	"encoding/hex"

	"strings"
)

const whitespaceTrim = " \n\t\r"

func whiteTrim(s string) string {
	return strings.Trim(s, whitespaceTrim)
}

func multiWhiteTrim(s []string) []string {
	trimmed := []string{}
	for _, cur := range s {
		if curT := whiteTrim(cur); curT != "" {
			trimmed = append(trimmed, curT)
		}
	}
	return trimmed
}

func md5ToString(c []byte) string {
	hasher := md5.New()
	hasher.Write(c)
	return hex.EncodeToString(hasher.Sum(nil))
}
