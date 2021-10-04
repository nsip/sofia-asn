package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/nsip/sofia-asn/tool"
	"github.com/tidwall/gjson"
)

type asnjson struct {
	Id string `json:"id"`

	Dcterms_modified struct {
		Literal string `json:"literal"`
	} `json:"dcterms_modified"`

	Dcterms_title struct {
		Literal string `json:"literal"`
	} `json:"dcterms_title"`

	Asn_statementLabel struct {
		Literal string `json:"literal"`
	} `json:"asn_statementLabel"`

	Asn_statementNotation struct {
		Literal string `json:"literal"`
	} `json:"asn_statementNotation"`

	Cls string `json:"cls"`

	Text string `json:"text"`

	Children []asnjson `json:"children"`

	Asn_skillEmbodied struct {
		Uri       string `json:"uri"`
		PrefLabel string `json:"prefLabel"`
	} `json:"asn_skillEmbodied"`
}

func nodeProcess(data []byte, outdir string) {

	e := bytes.LastIndexAny(data, "}")
	data = data[:e+1]

	outdir = strings.Trim(outdir, `./\`)
	parts := []string{}
	out := ""

	tool.ScanNode(data, func(i int, id, block string) bool {
		aj := asnjson{}
		aj.Id = id
		aj.Dcterms_modified.Literal = gjson.Get(block, "created_at").String()

		if bytes, err := json.Marshal(aj); err == nil {
			parts = append(parts, string(bytes))
		}

		if i == 10 {
			return false
		}
		return true
	})

	out = "[" + strings.Join(parts, ",") + "]"
	os.WriteFile(fmt.Sprintf("./%s/asn.json", outdir), []byte(out), os.ModePerm)
}

func main() {

	outdir := "out"
	os.MkdirAll(fmt.Sprintf("./%s/", outdir), os.ModePerm)

	data, err := os.ReadFile("../partition/out/node-meta.json")
	if err != nil {
		panic(err)
	}

	nodeProcess(data, "./out/")
}
