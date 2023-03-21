package main

import (
	"os"
	"strings"
	"testing"

	jt "github.com/digisan/json-tool"
)

func TestFixOneDupKey(t *testing.T) {
	fPath := "./out/la-Science-fix"
	data, err := os.ReadFile(fPath + ".json")
	if err != nil {
		panic(err)
	}
	n := 30
	prefix := "\n" + strings.Repeat(" ", n) + "\"asn_skillEmbodied\":"
	fixed := jt.FixOneDupKey(string(data), prefix)
	os.WriteFile(fPath+"-fix.json", []byte(fixed), os.ModePerm)
}

// asn_skillEmbodied

func TestRmDupEleOnce(t *testing.T) {
	fPath := "./data/la-English-fix"
	data, err := os.ReadFile(fPath + ".json")
	if err != nil {
		panic(err)
	}
	fixed := jt.RmDupEleOnce(string(data), "root.0.Age")
	os.WriteFile(fPath+"1.json", []byte(fixed), os.ModePerm)
}
