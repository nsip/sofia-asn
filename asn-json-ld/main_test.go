package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	main()
}

func TestAddCtx(t *testing.T) {

	os.MkdirAll("./out", os.ModePerm)

	data, err := os.ReadFile("../asn-json/out/la-English.json")
	if err != nil {
		panic(err)
	}
	js := string(data)
	js = addContext(js, contextRoot)
	js = replace(js)

	os.WriteFile("./out/test-ld.json", []byte(js), os.ModePerm)
}
