package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/digisan/gotk/misc"
	jt "github.com/digisan/json-tool"
	"github.com/nsip/sofia-asn/tool"
)

var mEscStr = map[string]string{
	`\n`: "*LF*",
	`\"`: "*DQ*",
}

func removeEsc(js string) string {
	for esc, str := range mEscStr {
		js = strings.ReplaceAll(js, esc, str)
	}
	return js
}

func restoreEsc(js string) string {
	for esc, str := range mEscStr {
		js = strings.ReplaceAll(js, str, esc)
	}
	return js
}

func main() {
	defer misc.TrackTime(time.Now())

	// {
	// 	outdir := "./out/"
	// 	outfile := "asn-node.json"
	// 	os.MkdirAll(outdir, os.ModePerm)
	// 	outpath := filepath.Join(outdir, outfile)

	// 	if !filedir.FileExists(outpath) {
	// 		data, err := os.ReadFile("../partition/out/node-meta.json")
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		nodeProc(data, outdir, outfile, "../data/Sofia-API-Tree-Data-09062022.json", uri4id)
	// 	}

	// 	// 	// 	/////

	// 	// 	data, err := os.ReadFile(outpath)
	// 	// 	if err != nil {
	// 	// 		log.Fatalln(err)
	// 	// 	}

	// 	// 	mIdBlock, _ := getIdBlock(string(data))

	// 	// 	inpath4exp := outpath
	// 	// 	outexp := childrenRepl(inpath4exp, mIdBlock)
	// 	// 	// os.WriteFile("./out/asnexp.json", []byte(outexp), os.ModePerm)

	// 	// 	rootWholeBlock := getRootWholeObject(outexp)
	// 	// 	os.WriteFile("./out/asn-node-one.json", []byte(rootWholeBlock), os.ModePerm)

	// }

	//////////////////////////////////////////////////////////////////////

	{
		mInputLa := map[string]string{
			"la-Languages.json":                      "Languages",
			"la-English.json":                        "English",
			"la-Humanities and Social Sciences.json": "Humanities and Social Sciences",
			"la-Health and Physical Education.json":  "Health and Physical Education",
			"la-Mathematics.json":                    "Mathematics",
			"la-Science.json":                        "Science",
			"la-Technologies.json":                   "Technologies",
			"la-The Arts.json":                       "The Arts",
			"ccp-Cross-curriculum Priorities.json":   "CCP",
			"gc-Critical and Creative Thinking.json": "GC-CCT",
			"gc-Digital Literacy.json":               "GC-DL",
			"gc-Ethical Understanding.json":          "GC-EU",
			"gc-Intercultural Understanding.json":    "GC-IU",
			// "gc-National Literacy Learning Progression.json": "GC-NLLP",
			// "gc-National Numeracy Learning Progression.json": "GC-NNLP",
			// "gc-Personal and Social Capability.json":         "GC-PSC",
		}

		os.MkdirAll("./out", os.ModePerm)

		dataTree, err := os.ReadFile("../data/Sofia-API-Tree-Data-09062022.json")
		if err != nil {
			log.Fatalln(err)
		}
		mCodeParent := tool.GetCodeParentMap(dataTree)

		dataNode, err := os.ReadFile("../partition/out/node-meta.json")
		if err != nil {
			log.Fatalln(err)
		}
		// mUidTitle := scanNodeIdTitle(dataNode) // title should be node title

		mNodeData, err := jt.Flatten(dataNode)
		if err != nil {
			log.Fatalln(err)
		}

		wg := sync.WaitGroup{}
		wg.Add(len(mInputLa))

		for file, la := range mInputLa {

			go func(file, la string) {

				// if file != "la-English.json" {
				// 	continue
				// }

				fmt.Printf("----- %s ----- %s\n", file, la)

				var (
					prevDocTypePath = ""
					retEL           = `` // used by 'Level' & its descendants
					retPL           = `` // used by 'Level' & its descendants
					progLvlABC      = "" // indicate Level 1a, 1b or 1c
				)

				data, err := os.ReadFile(`../partition/out/` + file)
				if err != nil {
					log.Fatalln(err)
				}
				js := removeEsc(string(data))

				///
				switch {
				case strings.Contains(js, `"Level 1c"`):
					progLvlABC = "1c"
				case strings.Contains(js, `"Level 1b"`):
					progLvlABC = "1b"
				case strings.Contains(js, `"Level 1a"`):
					progLvlABC = "1a"
				}
				///

				paths, _ := jt.GetLeavesPathOrderly(js)

				js = treeProc3(
					[]byte(js),
					la,
					mCodeParent,
					mNodeData,
					paths,
					&prevDocTypePath,
					&retEL,
					&retPL,
					progLvlABC,
				)

				js = restoreEsc(js)

				os.WriteFile("./out/"+file, []byte(js), os.ModePerm)

				wg.Done()

			}(file, la)
		}

		wg.Wait()
	}
}
