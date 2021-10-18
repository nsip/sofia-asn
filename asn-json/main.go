package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/digisan/gotk/filedir"
	"github.com/nsip/sofia-asn/tool"
)

func main() {

	{
		outdir := "./out/"
		outfile := "asn-node.json"
		os.MkdirAll(outdir, os.ModePerm)
		outpath := filepath.Join(outdir, outfile)

		if !filedir.FileExists(outpath) {
			data, err := os.ReadFile("../partition/out/node-meta.json")
			if err != nil {
				panic(err)
			}
			nodeProc(data, outdir, outfile, "../data/tree.pretty.json", "http://rdf.curriculum.edu.au/202110/")
		}

		/////

		data, err := os.ReadFile(outpath)
		if err != nil {
			log.Fatalln(err)
		}

		mIdBlock, _ := getIdBlock(string(data))

		inpath4exp := outpath
		outexp := childrenRepl(inpath4exp, mIdBlock)
		// os.WriteFile("./out/asnexp.json", []byte(outexp), os.ModePerm)

		rootWholeBlock := getRootWholeObject(outexp)
		os.WriteFile("./out/asn-node-one.json", []byte(rootWholeBlock), os.ModePerm)
	}

	//////////////////////////////////////////////////////////////////////

	{
		mInputLa := map[string]string{
			"la-English.json":                                "English",
			"la-HASS.json":                                   "Humanities and Social Sciences",
			"la-HPE.json":                                    "Health and Physical Education",
			"la-Languages.json":                              "Languages",
			"la-Mathematics.json":                            "Mathematics",
			"la-Science.json":                                "Science",
			"la-Technologies.json":                           "Technologies",
			"la-The Arts.json":                               "The Arts",
			"ccp-Cross-curriculum Priorities.json":           "",
			"gc-Critical and Creative Thinking.json":         "",
			"gc-Digital Literacy.json":                       "",
			"gc-Ethical Understanding.json":                  "",
			"gc-Intercultural Understanding.json":            "",
			"gc-National Literacy Learning Progression.json": "",
			"gc-National Numeracy Learning Progression.json": "",
			"gc-Personal and Social Capability.json":         "",
		}

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

		os.MkdirAll("./out", os.ModePerm)

		wg := sync.WaitGroup{}
		wg.Add(len(mInputLa))

		for file, la := range mInputLa {
			go func(file, la string) {

				data, err := os.ReadFile(`../partition/out/` + file)
				if err != nil {
					log.Fatalln(err)
				}
				out := treeProc(data, "http://rdf.curriculum.edu.au/202110", la, mCodeParent, mUidTitle)
				os.WriteFile("./out/"+file, []byte(out), os.ModePerm)
				wg.Done()
				
			}(file, la)
		}

		wg.Wait()
	}
}
