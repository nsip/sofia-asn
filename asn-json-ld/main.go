package main

import (
	"fmt"
	"os"
)

type AsnJsonLd struct {
}

var (
	jsonldctx = `{
		"@context": {
		  "asn": "http://purl.org/ASN/schema/core/",
		  "dc": "http://purl.org/dc/terms/",
		  "gem": "http://purl.org/gem/qualifiers/",
		  "skos": "http://www.w3.org/2004/02/skos/core#",
		  "xsd": "http://www.w3.org/2001/XMLSchema#",
		  @language: "en-au"
		}
	  }`

	ns = map[string]string{
		"asn": "http://purl.org/ASN/schema/core/",
		"deo": "http://purl.org/spar/deo",
		"esa": "http://vocabulary.curriculum.edu.au/",
		"dc":  "http://purl.org/dc/terms/",
		"gem": "http://purl.org/gem/qualifiers/",
	}

	prefrepl = map[string]string{
		"dc_":      "dc",
		"dcterms_": "dc",
		"asn_":     "asn",
	}

	repl = map[string]string{
		"text":     "dc:description",
		"children": "gem:hasChild",
		"id":       "@id",
	}

	ignr = []string{
		"cls",
		"leaf",
	}
)

func main() {

	data, err := os.ReadFile("../asn-json/out/asnroot.json")
	if err != nil {
		panic(err)
	}
	js := string(data)

	fmt.Println(len(js))
}
