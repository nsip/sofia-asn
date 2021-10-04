package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestAsnJson(t *testing.T) {
	aj := asnjson{
		Id: "abc",
	}
	aj.Dcterms_title.Literal = "def"
	bytes, err := json.Marshal(aj)
	if err == nil {
		fmt.Println(string(bytes))
	}
}
