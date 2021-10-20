package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	jsontool "github.com/digisan/json-tool"
)

type AsnJsonLd struct {
}

func addContext(js, ctx string) string {
	s, e := strings.Index(ctx, "{"), strings.LastIndex(ctx, "}")
	ctx = ctx[s+1 : e]
	s = strings.Index(js, "{")
	js = js[s+1:]
	return jsontool.FmtStr("{"+ctx+","+js, "  ")
}

func replace(js string) string {
	r := regexp.MustCompile(`"(dc_|dcterms_|asn_)[^"]+"`)
	js = r.ReplaceAllStringFunc(js, func(s string) string {
		s = strings.Trim(s, "\"")
		ss := strings.Split(s, "_")
		p1 := mPrefRepl[ss[0]+"_"]
		return "\"" + p1 + ss[1] + "\""
	})

	

	return js
}

func main() {

	data, err := os.ReadFile("../asn-json/out/la-English.json")
	if err != nil {
		panic(err)
	}
	js := string(data)

	fmt.Println(len(js))
}
