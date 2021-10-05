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
	Id string `json:"Id"`

	///////////////

	Dcterms_modified struct {
		Literal string `json:"literal"`
	} `json:"dcterms_modified"`

	Dcterms_subject struct {
		Uri       string `json:"uri"`
		PrefLabel string `json:"prefLabel"`
	} `json:"dcterms_subject"`

	Dcterms_educationLevel struct {
		Uri       string `json:"uri"`
		PrefLabel string `json:"prefLabel"`
	} `json:"dcterms_educationLevel"`

	Dcterms_description struct {
		Uri       string `json:"uri"`
		PrefLabel string `json:"prefLabel"`
	} `json:"dcterms_description"`

	Dcterms_title struct {
		Literal string `json:"literal"`
	} `json:"dcterms_title"`

	///////////////

	Asn_statementLabel struct {
		Literal string `json:"literal"`
	} `json:"asn_statementLabel"`

	Asn_statementNotation struct {
		Literal string `json:"literal"`
	} `json:"asn_statementNotation"`

	Asn_skillEmbodied struct {
		Uri       string `json:"uri"`
		PrefLabel string `json:"prefLabel"`
	} `json:"asn_skillEmbodied"`

	///////////////

	Cls string `json:"cls"`

	Text string `json:"text"`

	Children []string `json:"children"`
}

func nodeProcess(data []byte, outdir string) {

	e := bytes.LastIndexAny(data, "}")
	data = data[:e+1]

	outdir = strings.Trim(outdir, `./\`)
	parts := []string{}
	out := ""

	tool.ScanNode(data, func(i int, id, block string) bool {

		aj := asnjson{}
		aj.Id = gjson.Get(block, "id").String()

		aj.Dcterms_modified.Literal = gjson.Get(block, "created_at").String()
		aj.Dcterms_title.Literal = gjson.Get(block, "title").String()

		aj.Asn_statementLabel.Literal = gjson.Get(block, "doc.typeName").String()
		aj.Asn_statementNotation.Literal = gjson.Get(block, "code").String()
		// aj.Asn_skillEmbodied.Uri = gjson.Get(block, "").String()
		// aj.Asn_skillEmbodied.PrefLabel = gjson.Get(block, "").String()

		// aj.Cls = gjson.Get(block, "").String()
		aj.Text = gjson.Get(block, "text").String()
		for _, c := range gjson.Get(block, "children").Array() {
			aj.Children = append(aj.Children, c.String())
		}

		////////////////////////////////////////////////////////////////

		if bytes, err := json.Marshal(aj); err == nil {
			parts = append(parts, string(bytes))
		}

		if i == 5 {
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
