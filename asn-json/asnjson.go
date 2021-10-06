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
	Id string `json:"Id"` // DIRECT

	///////////////

	Dcterms_modified struct {
		Literal string `json:"literal"` // DIRECT
	} `json:"dcterms_modified"`

	Dcterms_subject struct {
		Uri       string `json:"uri"`       // derived
		PrefLabel string `json:"prefLabel"` // derived
	} `json:"dcterms_subject"`

	Dcterms_educationLevel struct {
		Uri       string `json:"uri"`       // derived
		PrefLabel string `json:"prefLabel"` // derived
	} `json:"dcterms_educationLevel"`

	Dcterms_description struct {
		Uri       string `json:"uri"`
		PrefLabel string `json:"prefLabel"`
	} `json:"dcterms_description"`

	Dcterms_title struct {
		Literal  string `json:"literal"`  // DIRECT
		Language string `json:"language"` // boilerplate
	} `json:"dcterms_title"`

	Dcterms_rights struct {
		Literal  string `json:"literal"`  // boilerplate
		Language string `json:"language"` // boilerplate
	} `json:"dcterms_rights"`

	Dcterms_rightsHolder struct {
		Literal  string `json:"literal"`  // boilerplate
		Language string `json:"language"` // boilerplate
	} `json:"dcterms_rightsHolder"`

	///////////////

	Asn_statementLabel struct {
		Literal  string `json:"literal"`  // DIRECT
		Language string `json:"language"` // boilerplate
	} `json:"asn_statementLabel"`

	Asn_statementNotation struct {
		Literal  string `json:"literal"`  // DIRECT
		Language string `json:"language"` // boilerplate
	} `json:"asn_statementNotation"`

	Asn_skillEmbodied struct {
		Uri       string `json:"uri"`       // Predicate
		PrefLabel string `json:"prefLabel"` // Predicate
	} `json:"asn_skillEmbodied"`

	Asn_authorityStatus struct {
		Uri string `json:"uri"` // boilerplate
	} `json:"asn_authorityStatus"`

	Asn_indexingStatus struct {
		Uri string `json:"uri"` // boilerplate
	} `json:"asn_indexingStatus"`

	Asn_hasLevel struct {
		Uri       string `json:"uri"`       // Predicate
		PrefLabel string `json:"prefLabel"` // Predicate
	} `json:"asn_hasLevel"`

	///////////////

	Dc_relation struct {
		Uri       string `json:"uri"`       // Predicate
		PrefLabel string `json:"prefLabel"` // Predicate
	} `json:"dc_relation"`

	///////////////

	Cls string `json:"cls"` // boilerplate

	Leaf string `json:"leaf"` // boilerplate

	Text string `json:"text"` // DIRECT

	Children []string `json:"children"` //DIRECT
}

func nodeProcess(data []byte, outdir, sofiaTreeFile string) {

	e := bytes.LastIndexAny(data, "}")
	data = data[:e+1]

	outdir = strings.Trim(outdir, `./\`)
	parts := []string{}
	out := ""

	dataTree, err := os.ReadFile(sofiaTreeFile)
	if err != nil {
		panic(err)
	}
	mCodeParent := tool.GetCodeParentMap(dataTree)

	tool.ScanNode(data, func(i int, id, block string) bool {

		code := gjson.Get(block, "code").String()
		title := tool.GetAncestorTitle(mCodeParent, code, "")
		// fmt.Println(i, id, code)

		////////////////////////////////////////////////////////

		aj := asnjson{}

		// Direct
		aj.Id = gjson.Get(block, "id").String()
		aj.Dcterms_modified.Literal = gjson.Get(block, "created_at").String()
		aj.Dcterms_title.Literal = gjson.Get(block, "title").String()
		aj.Asn_statementLabel.Literal = gjson.Get(block, "doc.typeName").String()
		aj.Asn_statementNotation.Literal = gjson.Get(block, "code").String()
		aj.Text = gjson.Get(block, "text").String()
		for _, c := range gjson.Get(block, "children").Array() {
			aj.Children = append(aj.Children, c.String())
		}

		// Derived
		aj.Dcterms_subject.Uri = ""
		aj.Dcterms_subject.PrefLabel = title
		aj.Dcterms_educationLevel.Uri = ""
		aj.Dcterms_educationLevel.PrefLabel = ""

		// Boilerplate
		aj.Dcterms_title.Language = "en-au"
		aj.Asn_statementLabel.Language = "en-au"
		aj.Asn_statementNotation.Language = "en-au"
		aj.Dcterms_rights.Language = "en-au"
		aj.Dcterms_rightsHolder.Language = "en-au"
		aj.Asn_authorityStatus.Uri = `http://purl.org/ASN/scheme/ASNAuthorityStatus/Original`
		aj.Asn_indexingStatus.Uri = `http://purl.org/ASN/scheme/ASNIndexingStatus/No`
		aj.Dcterms_rights.Literal = `Copyright Australian Curriculum, Assessment and Reporting Authority`
		aj.Dcterms_rightsHolder.Literal = `Australian Curriculum, Assessment and Reporting Authority`
		aj.Cls = "folder"
		aj.Leaf = "true"

		// Predicate
		aj.Asn_skillEmbodied.Uri = ""       //
		aj.Asn_skillEmbodied.PrefLabel = "" //
		aj.Dc_relation.Uri = ""             // if
		aj.Dc_relation.PrefLabel = ""       // if
		aj.Asn_hasLevel.Uri = ""            //
		aj.Asn_hasLevel.PrefLabel = ""      //

		////////////////////////////////////////////////////////////////

		if bytes, err := json.Marshal(aj); err == nil {
			parts = append(parts, string(bytes))
		}

		if i == 10000000 {
			return false
		}
		return true
	})

	out = "[" + strings.Join(parts, ",") + "]"
	os.WriteFile(fmt.Sprintf("./%s/asn.json", outdir), []byte(out), os.ModePerm)
}
