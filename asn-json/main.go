package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {

	/////

	// outdir := "./out/"
	// outfile := "asn.json"
	// os.MkdirAll(outdir, os.ModePerm)
	// outpath := filepath.Join(outdir, outfile)

	// if !filedir.FileExists(outpath) {
	// 	data, err := os.ReadFile("../partition/out/node-meta.json")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	nodeProc(data, outdir, outfile, "../data/tree.pretty.json", "http://rdf.curriculum.edu.au/202110/")
	// }

	/////

	// mIdBlock, _ := getIdBlock(outpath)

	// inpath4exp := outpath
	// path4exp := "./out/asnexp.json"

	// outexp := childrenRepl(inpath4exp, mIdBlock)
	// os.WriteFile(path4exp, []byte(outexp), os.ModePerm)

	/////

	asnjson := "./out/asnexp.json"
	data, err := os.ReadFile(asnjson)
	if err != nil {
		panic(err)
	}
	js := string(data)
	fmt.Println(len(js))

	rId := regexp.MustCompile(`"Id": "http[^"]+"`)

	ids := rId.FindAllString(js, -1)
	fmt.Println(len(ids))
	fmt.Println(strings.Count(js, `"Id":`))

	mIdCnt := make(map[string]int)
	rId.ReplaceAllStringFunc(js, func(s string) string {
		mIdCnt[s]++
		return s
	})

	fmt.Println(len(mIdCnt))

	mIdCnt2 := make(map[string]int)
	for idstr, cnt := range mIdCnt {
		if cnt == 1 {
			mIdCnt2[idstr] = cnt
		}
	}

	fmt.Println(mIdCnt2)

	mIdBlock, _ := getIdBlock(asnjson)

	root := mIdBlock["http://rdf.curriculum.edu.au/202110/649c9d14-75b7-41e3-ac5f-c4c86fd8f57c"]

	os.WriteFile("./out/asnroot.json", []byte(root), os.ModePerm)

	// fmt.Println(len(mIdBlock))

	// exclId := []string{}
	// for id := range mIdBlock {
	// 	if strings.Count(js, id) > 1 {
	// 		exclId = append(exclId, id)
	// 	}
	// }

	// fmt.Println(len(exclId))

}
