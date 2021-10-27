package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	jt "github.com/digisan/json-tool"
)

var (
	mRES = map[string]string{
		"uuid":       `"uuid":\s*"[\d\w]{8}-[\d\w]{4}-[\d\w]{4}-[\d\w]{4}-[\d\w]{12}",?`,
		"type":       `"type":\s*"\w+",?`,
		"created_at": `"created_at":\s*"[^"]+",?`,
		"title":      `"title":\s*".+",?`,
	}

	mRE4Path = map[string]*regexp.Regexp{
		"uuid":       regexp.MustCompile(`\.?uuid$`),
		"type":       regexp.MustCompile(`\.?type$`),
		"created_at": regexp.MustCompile(`\.?created_at$`),
		"title":      regexp.MustCompile(`\.?title$`),
	}

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

func treeProc3(data []byte) string {

	js := string(data)

	mData, err := jt.Flatten(data)
	if err != nil {
		log.Fatalln(err)
	}

	proc := func(name, value string) (bool, string) {
		// fmt.Println(name, value)

		if name == "uuid" && value == "f0eaeb1b-9b65-4756-a309-9f79acf48ebd" {
			paths := []string{}
			for p, v := range mData {
				if value == v && mRE4Path[name].MatchString(p) {
					paths = append(paths, p)
				}
			}
			fmt.Println(name, value, paths)
			return true, `"uuid": ` + `"f0eaeb1b-9b65-4756-a309-9f79acf48ABC"`
		}

		if name == "type" && value == "r" {
			return true, ""
		}

		return false, ""
	}

	re4json, mRE4Each := reMerged()
	return re4json.ReplaceAllStringFunc(js, func(s string) string {
		hasComma := false
		if strings.HasSuffix(s, ",") {
			hasComma = true
		}
		for name, v := range mRE4Each {
			if v.MatchString(s) {
				if ok, repl := proc(name, fetchValue(s)); ok {
					if hasComma && repl != "" {
						return repl + ","
					}
					return repl
				}
				break
			}
		}
		return s
	})
}
