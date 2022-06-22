package main

import (
	"os"
	"strings"
	"testing"
	jt "github.com/digisan/json-tool"
)

func TestFixOneDupKey(t *testing.T) {
	fpath := "./out/la-The Arts"
	data, err := os.ReadFile(fpath + ".json")
	if err != nil {
		panic(err)
	}
	n := 26
	prefix := "\n" + strings.Repeat(" ", n) + "\"asn_skillEmbodied\":"
	fixed := jt.FixOneDupKey(string(data), prefix)
	os.WriteFile(fpath+"-fix.json", []byte(fixed), os.ModePerm)
}


// asn_skillEmbodied

func TestRmDupEleOnce(t *testing.T) {
	fpath := "./data/la-English-fix"
	data, err := os.ReadFile(fpath + ".json")
	if err != nil {
		panic(err)
	}
	fixed := jt.RmDupEleOnce(string(data), "root.0.Age")
	os.WriteFile(fpath+"1.json", []byte(fixed), os.ModePerm)
}