package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"regexp"
)

var (
	isSemanticVersionRxp = regexp.MustCompile(`^\d+\.\d+\.\d+$`)
)

func readJSON(filename string, d interface{}) error {
	configData, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(configData, d)
}

func md5ToString(c []byte) string {
	hasher := md5.New()
	hasher.Write(c)
	return hex.EncodeToString(hasher.Sum(nil))
}
