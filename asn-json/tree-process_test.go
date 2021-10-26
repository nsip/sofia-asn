package main

import (
	"log"
	"os"
	"testing"

	"github.com/nsip/sofia-asn/tool"
)

var (
	mInputLa = map[string]string{
		"la-English.json": "English",
		// "la-HASS.json":                                   "Humanities and Social Sciences",
		// "la-HPE.json":                                    "Health and Physical Education",
		// "la-Languages.json":                              "Languages",
		// "la-Mathematics.json":                            "Mathematics",
		// "la-Science.json":                                "Science",
		// "la-Technologies.json":                           "Technologies",
		// "la-The Arts.json":                               "The Arts",
		// "ccp-Cross-curriculum Priorities.json":           "",
		// "gc-Critical and Creative Thinking.json":         "",
		// "gc-Digital Literacy.json":                       "",
		// "gc-Ethical Understanding.json":                  "",
		// "gc-Intercultural Understanding.json":            "",
		// "gc-National Literacy Learning Progression.json": "",
		// "gc-National Numeracy Learning Progression.json": "",
		// "gc-Personal and Social Capability.json":         "",
	}
)

func TestTreeProc2(t *testing.T) {
	os.MkdirAll("./out/", os.ModePerm)

	for file := range mInputLa {
		data, err := os.ReadFile(`../partition/out/` + file)
		if err != nil {
			log.Fatalln(err)
		}
		out := treeProc2(data, "http://rdf.curriculum.edu.au/202110", "English")
		os.WriteFile("./out/"+file, []byte(out), os.ModePerm)
	}
}

func TestTreeProc(t *testing.T) {

	dataTree, err := os.ReadFile("../data/tree.pretty.json")
	if err != nil {
		log.Fatalln(err)
	}
	mCodeParent := tool.GetCodeParentMap(dataTree)

	dataNode, err := os.ReadFile("../partition/out/node-meta.json")
	if err != nil {
		log.Fatalln(err)
	}
	mUidTitle := scanNodeIdTitle(dataNode) // title should be node title

	for file, la := range mInputLa {
		data, err := os.ReadFile(`../partition/out/` + file)
		if err != nil {
			log.Fatalln(err)
		}
		out := treeProc(data, "http://rdf.curriculum.edu.au/202110", la, mCodeParent, mUidTitle)
		os.WriteFile("./out/"+file, []byte(out), os.ModePerm)
	}
}
