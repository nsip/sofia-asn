package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/digisan/gotk/slice/ts"
	strs "github.com/digisan/gotk/strings"
	jt "github.com/digisan/json-tool"
	"github.com/tidwall/gjson"
)

var (
	mStat = make(map[string]int)
)

var (
	mRES = map[string]string{
		"uuid":         `"uuid":\s*"[\d\w]{8}-[\d\w]{4}-[\d\w]{4}-[\d\w]{4}-[\d\w]{12}",?`,
		"type":         `"type":\s*"\w+",?`,
		"created_at":   `"created_at":\s*"[^"]+",?`,
		"title":        `"title":\s*".+",?`,
		"doc.typeName": `"doc":\s*\{[^{}]+\},?`,
		"code":         `"code":\s*"[^"]+",?`,
		"text":         `"text":\s*"[^"]+",?`,
	}

	mRE4Path = map[string]*regexp.Regexp{
		"uuid":         regexp.MustCompile(`\.?uuid$`),
		"type":         regexp.MustCompile(`\.?type$`),
		"created_at":   regexp.MustCompile(`\.?created_at$`),
		"title":        regexp.MustCompile(`\.?title$`),
		"doc.typeName": regexp.MustCompile(`\.?doc\.typeName$`),
		"code":         regexp.MustCompile(`\.?code$`),
		"text":         regexp.MustCompile(`\.?text$`),
	}

	mRELocGrp = map[string][]int{}
	mRELocIdx = map[string]int{}

	reMerged = func() (*regexp.Regexp, map[string]*regexp.Regexp) {
		mRE4Each := map[string]*regexp.Regexp{}
		restr := ""
		for k, v := range mRES {
			restr += fmt.Sprintf("(%s)|", v)    // merged restr for whole regexp
			mRE4Each[k] = regexp.MustCompile(v) // init each regexp
		}
		// remove last '|' and compile to regexp
		return regexp.MustCompile(restr[:len(restr)-1]), mRE4Each
	}

	r4FetchValue = regexp.MustCompile(`:\s*"`)
	fetchValue   = func(kvstr string) string {
		loc := r4FetchValue.FindAllStringIndex(kvstr, 1)[0]
		start := loc[1]
		end := strings.LastIndex(kvstr, `"`)
		return kvstr[start:end]
	}
)

func getPathByFieldValuePos(js, name, value string, pos int, mData map[string]interface{}) string {
	for p, v := range mData {
		if value == v && mRE4Path[name].MatchString(p) {
			rst := gjson.Get(js, jt.ParentPath(p)).String()
			if candidate := strings.Index(js, rst); candidate != -1 {
				if candidate-(pos+len(name)+3) < 3 {
					delete(mData, p)
					return p
				}
			}
		}
	}
	log.Fatalln(name, value, pos)
	return ""
}

func proc(js, s, name, value string, mData, mData4Yr map[string]interface{}) (bool, string) {

	var (
		fSf    = fmt.Sprintf
		uri4id = "http://rdf.curriculum.edu.au/202110"
	)

	if mRELocGrp[s] == nil {
		mRELocGrp[s], _ = strs.IndexAll(js, s)
	}

	pos := mRELocGrp[s][mRELocIdx[s]]
	mRELocIdx[s]++

	//
	// get 'path' once, mData removes returned path record
	//
	// path := getPathByFieldValuePos(js, name, value, pos, mData)
	// fmt.Println(path)

	switch name {
	case "uuid":
		return true, fSf(`"Id": "%s/%s"`, uri4id, value)

	case "type":
		return true, ""

	case "created_at":
		return true, fSf(`"dcterms_modified": { "literal": "%s" }`, value)

	case "title":
		return true, fSf(`"dcterms_title": { "language": "%s", "literal": "%s" }`, "en-au", value)

	case "doc.typeName":

		// return true, fSf(`"asn_statementLabel": { "language": "%s", "literal": "%s" }`, "en-au", value)

		// ret1 := fSf(`"asn_statementLabel": { "language": "%s", "literal": "%s" }`, "en-au", value)
		// ret2 := ""
		if ts.NotIn(value, "Learning Area", "Subject") {
			path := getPathByFieldValuePos(js, name, value, pos, mData)
			fmt.Println(path)

			// audit, make sure one path only for once use only
			mStat[path]++

			for _, y := range getYears(mData4Yr, path) {
				fmt.Println(y)
			}
		}

		// TEST
		return true, fSf(`"asn_statementLabel": { "language": "%s", "literal": "%s" }`, "en-au", value)

		// return true, strings.Join([]string{ret1, ret2}, ",")

	case "code":
		return true, fSf(`"asn_statementNotation": { "language": "%s", "literal": "%s" }`, "en-au", value)

	case "text":
		return true, fSf(`"text": "%s"`, value)

	default:
		return false, ""
	}
}

func treeProc3(data []byte) string {

	js := string(data)

	mData, err := jt.Flatten(data)
	if err != nil {
		log.Fatalln(err)
	}

	mData4Yr := make(map[string]interface{})
	for k, v := range mData {
		mData4Yr[k] = v
	}

	re4json, mRE4Each := reMerged()
	js = re4json.ReplaceAllStringFunc(js, func(s string) string {
		hasComma := false
		if strings.HasSuffix(s, ",") {
			hasComma = true
		}
		for name, v := range mRE4Each {
			if v.MatchString(s) {

				// if name == "doc.typeName" {
				// 	fmt.Println(name, s)
				// }

				if ok, repl := proc(js, s, name, fetchValue(s), mData, mData4Yr); ok {
					if hasComma && repl != "" {
						return repl + ","
					}
					return repl
				}
			}
		}
		return s
	})

	// audit, make sure one path only for once use only
	for k, v := range mStat {
		if v > 1 {
			log.Fatalln(k, v)
		}
	}
	fmt.Println(len(mStat), "paths")

	return js
}
