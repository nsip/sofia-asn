package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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

func TestYearSplit(t *testing.T) {
	for _, y := range yearsSplit("Years 1 and 2") {
		fmt.Println(y)
		fmt.Println(mYrlvlUri[y])
	}
	for _, y := range yearsSplit("Years 9") {
		fmt.Println(y)
		fmt.Println(mYrlvlUri[y])
	}
}

func TestMapSet(t *testing.T) {
	m := map[string]string{
		"a": "A",
		"b": "B",
	}
	fmt.Println("m:", m)
	mm := m
	mm["c"] = "C"
	fmt.Println("m:", m)
	fmt.Println("mm:", mm)
}

func TestTreeProc(t *testing.T) {
	data, err := os.ReadFile(`../partition/out/la-English.json`)
	if err != nil {
		log.Fatalln(err)
	}
	treeProc(data, "./out", "test", "http://rdf.curriculum.edu.au/202110", "English", "", "")
}
