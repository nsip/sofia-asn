package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/digisan/gotk"
	"github.com/digisan/gotk/slice/ti"
	"github.com/digisan/gotk/slice/ts"
	"github.com/digisan/gotk/slice/tso"
	jt "github.com/digisan/json-tool"
	"github.com/nsip/sofia-asn/tool"
	"github.com/tidwall/gjson"
)

var (
	fSf        = fmt.Sprintf
	sJoin      = strings.Join
	sTrim      = strings.Trim
	sLastIndex = strings.LastIndex
	sHasSuffix = strings.HasSuffix
	sSplit     = strings.Split
)

var (
	mRES = map[string]string{
		// "text":               `"text":\s*"[^"]+",?`,
		"uuid":               `"uuid":\s*"[\d\w]{8}-[\d\w]{4}-[\d\w]{4}-[\d\w]{4}-[\d\w]{12}",?`,
		"type":               `"type":\s*"\w+",?`,
		"created_at":         `"created_at":\s*"[^"]+",?`,
		"title":              `"title":\s*"[^"]+",?`,
		"doc.typeName":       `"doc":\s*\{[^{}]+\},?`,
		"code":               `"code":\s*"[^"]+",?`,
		"tag":                `"tags":\s*\{[^{}]+\},?`,
		"connections.Levels": `"Levels":\s*\[[^\[\]]+\],?`,
		"connections.OI":     `"Organising Ideas":\s*\[[^\[\]]+\],?`,
		"connections.ASC":    `"Achievement Standard Components":\s*\[[^\[\]]+\],?`,
		"connections.IG":     `"Indicator Groups":\s*\[[^\[\]]+\],?`,
		"connections.CD":     `"Content Descriptions":\s*\[[^\[\]]+\],?`,
	}

	mRE4Path = map[string]*regexp.Regexp{
		// "text":               regexp.MustCompile(`\.?text$`),
		"uuid":               regexp.MustCompile(`\.?uuid$`),
		"type":               regexp.MustCompile(`\.?type$`),
		"created_at":         regexp.MustCompile(`\.?created_at$`),
		"title":              regexp.MustCompile(`\.?title$`),
		"doc.typeName":       regexp.MustCompile(`\.?doc\.typeName$`),
		"code":               regexp.MustCompile(`\.?code$`),
		"tag":                regexp.MustCompile(`\.?tags\.`),
		"connections.Levels": regexp.MustCompile(`\.?connections\.Levels\.\d+$`),
		"connections.OI":     regexp.MustCompile(`\.?connections\.Organising Ideas\.\d+$`),
		"connections.ASC":    regexp.MustCompile(`\.?connections\.Achievement Standard Components\.\d+$`),
		"connections.IG":     regexp.MustCompile(`\.?connections\.Indicator Groups\.\d+$`),
		"connections.CD":     regexp.MustCompile(`\.?connections\.Content Descriptions\.\d+$`),
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
)

func fn4GetPathByProp(js, prop string) func() string {
	idx := -1
	paths, _ := jt.GetAllLeafPaths(js)
	paths = jt.GetLeafPathsOrderly(prop, paths)
	return func() string {
		idx++
		return paths[idx]
	}
}

var (
	prevDocTypePath = ""
	retEL           = `` // used by 'Level' & its descendants
)

func proc(

	js, s, name, value string,
	mLvlSiblings map[int][]string,
	dataPaths *[]string,
	dataValues *[]interface{},
	mData map[string]interface{},
	la, uri4id string,
	mCodeParent map[string]string,
	mNodeData map[string]interface{},
	getPathWithDocType func() string,
	getPathWithCode func() string,

) (bool, string) {

	// if name == "doc.typeName" {
	// 	fmt.Println(name, s)
	// }

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

		path := getPathWithDocType()

		// "asn_statementLabel"
		retSL := fSf(`"asn_statementLabel": { "language": "%s", "literal": "%s" }`, "en-au", value)

		// "dcterms_educationLevel"
		if value == "Level" { // see doc.typeName: 'Level', update global retEL
			if la != "" {
				outArrs := []string{}
				for _, y := range getYears(mData, path) {
					outArrs = append(outArrs, fSf(`{ "uri": "%s", "prefLabel": "%s" }`, mYrlvlUri[y], y))
				}
				if len(outArrs) > 0 {
					retEL = sJoin(outArrs, ",")
				}
			}
			retEL = fSf(`"dcterms_educationLevel": [%s]`, retEL)
			prevDocTypePath = path
		}

		// only children path can keep retEL
		if strings.Count(path, ".") < strings.Count(prevDocTypePath, ".") {
			retEL = ""
		}

		return true, sTrim(sJoin([]string{retSL, retEL}, ","), ",")

	case "code":

		path := getPathWithCode()

		retSN := fSf(`"asn_statementNotation": { "language": "%s", "literal": "%s" }`, "en-au", value)

		retAS := fSf(`"asn_authorityStatus": { "uri": "%s" }`, `http://purl.org/ASN/scheme/ASNAuthorityStatus/Original`)

		retIS := fSf(`"asn_indexingStatus": { "uri": "%s" }`, `http://purl.org/ASN/scheme/ASNIndexingStatus/No`)

		retTxt := ""
		if !gjson.Get(js, jt.NewSibling(path, "text")).Exists() {
			retTxt = `"text": null`
		}

		retSub := ``
		if ts.In(value, "ENG", "HAS", "HPE", "LAN", "MAT", "SCI", "TEC", "ART") {
			retS := []string{}
			if subUri, okSubUri := mLaUri[la]; okSubUri {
				retS = append(retS, fSf(`"dcterms_subject": { "prefLabel": "%s", "uri": "%s" }`, la, subUri))
			}
			retSub = sJoin(retS, ",")
		}

		retRT, retRTH := ``, ``
		if ts.In(value, "root", "LA") {
			retRT = fSf(`"dcterms_rights": { "language": "%s", "literal": "%s" }`, "en-au", `©Copyright Australian Curriculum, Assessment and Reporting Authority`)
			retRTH = fSf(`"dcterms_rightsHolder": { "language": "%s", "literal": "%s" }`, "en-au", `Australian Curriculum, Assessment and Reporting Authority`)
		}

		retCLS, retLEAF := ``, ``
		if jt.HasSiblings(path, mLvlSiblings, "children") {
			retCLS = fSf(`"cls": "folder"`)
		} else {
			retLEAF = fSf(`"leaf": "true"`)
		}

		rets := []string{}
		for _, r := range []string{retSN, retAS, retIS, retTxt, retSub, retRT, retRTH, retCLS, retLEAF} {
			if r != "" {
				rets = append(rets, r)
			}
		}
		return true, sJoin(rets, ",")

	// case "text":
	// 	return true, fSf(`"text": "%s"`, value)

	case "tag":
		return true, fSf(`"asn_conceptTerm": "%s"`, "SCIENCE_TEACHER_BACKGROUND_INFORMATION")

	case "connections.Levels",
		"connections.OI",
		"connections.ASC",
		"connections.IG",
		"connections.CD":

		items := sSplit(value, "|")
		// fmt.Println(items)

		fieldname := ""
		outArrs := []string{}
		for _, item := range items {
			id := item[sLastIndex(item, "/")+1:]
			code := jt.GetStrVal(mNodeData[id+"."+"code"])
			title := jt.GetStrVal(mNodeData[id+"."+"title"])
			nodeType := tool.GetCodeAncestor(mCodeParent, code, 0)
			switch nodeType {
			case "GC":
				outArrs = append(outArrs, fSf(`{ "asn_skillEmbodied": { "uri": "%s", "prefLabel": "%s"} }`, item, title))
			case "LA":
				outArrs = append(outArrs, fSf(`{ "dc_relation": { "uri": "%s", "prefLabel": "%s"} }`, item, title))
			case "AS":
				outArrs = append(outArrs, fSf(`{ "asn_hasLevel": { "uri": "%s", "prefLabel": "%s"} }`, item, title))
			case "CCP":
				outArrs = append(outArrs, fSf(`{ "asn_crossSubjectReference": { "uri": "%s", "prefLabel": "%s"} }`, item, title))
			default:
				log.Fatalf("'%v' is not one of [GC CCP LA AS], code is '%v'", nodeType, code)
			}
		}

		switch name {
		case "connections.Levels":
			fieldname = "Level"
		case "connections.OI":
			fieldname = "Organising Ideas"
		case "connections.ASC":
			fieldname = "Achievement Standard Components"
		case "connections.IG":
			fieldname = "Indicator Groups"
		case "connections.CD":
			fieldname = "Content Descriptions"
		}

		outArrStr := sJoin(outArrs, ",")
		ret := fSf(`"%s": [%s]`, fieldname, outArrStr)
		return true, ret

	default:
		return false, ""
	}
}

func treeProc3(data []byte, la string, mCodeParent map[string]string, mNodeData map[string]interface{}) string {

	var (
		uri4id = "http://rdf.curriculum.edu.au/202110"
	)

	js := string(data)
	mLvlSiblings, _ := jt.FamilyTree(js)

	mData, err := jt.Flatten(data)
	if err != nil {
		log.Fatalln(err)
	}

	sortRule := func(s1, s2 string) bool {
		a1, a2 := []int{}, []int{}
		for i, s := range []string{s1, s2} {
			var a *[]int
			if i == 0 {
				a = &a1
			} else {
				a = &a2
			}
			for _, seg := range sSplit(s, ".") {
				if gotk.IsNumeric(seg) {
					n, _ := strconv.Atoi(seg)
					*a = append(*a, n)
				}
			}
		}
		lmin := ti.Min(len(a1), len(a2))
		for i := 0; i < lmin; i++ {
			n1, n2 := a1[i], a2[i]
			switch {
			case n1 < n2:
				return true
			case n1 > n2:
				return false
			default:
				continue
			}
		}
		return true
	}

	dataPaths, dataValues := tso.Map2KVs(mData, sortRule, nil)
	// fmt.Println(len(dataPaths), len(dataValues))

	re4json, mRE4Each := reMerged()
	// fmt.Println(re4json, len(mRE4Each))

	getPathWithTypeName := fn4GetPathByProp(js, "typeName")
	getPathWithCode := fn4GetPathByProp(js, "code")

	js = re4json.ReplaceAllStringFunc(js, func(s string) string {

		hasComma := false
		if sHasSuffix(s, ",") {
			hasComma = true
		}

		for name, v := range mRE4Each {
			if v.MatchString(s) {

				// if name == "doc.typeName" {
				// 	fmt.Println(name, s)
				// }

				if ok, repl := proc(
					js,
					s,
					name,
					tool.FetchValue(s, "|"),
					mLvlSiblings,
					&dataPaths,
					&dataValues,
					mData,
					la,
					uri4id,
					mCodeParent,
					mNodeData,
					getPathWithTypeName,
					getPathWithCode,
				); ok {
					if hasComma && repl != "" {
						return repl + ","
					}
					return repl
				}
			}
		}
		return s
	})

	return js
}
