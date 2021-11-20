package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/digisan/gotk/slice/ts"
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

func fnGetPathByProp(prop string, paths []string, info string) func() string {
	var idx int64 = -1
	paths = jt.GetLeafPathsOrderly(prop, paths)
	return func() string {
		idx++
		// fmt.Println(prop, idx, info)
		return paths[idx]
	}
}

func proc(

	js, s, name, value string,
	mLvlSiblings map[int][]string,
	mData map[string]interface{},
	la, uri4id string,
	mCodeParent map[string]string,
	mNodeData map[string]interface{},
	fnPathWithDocType func() string,
	fnPathWithCode func() string,
	// static for filling
	pPrevDocTypePath *string,
	pRetEL *string,

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

		path := fnPathWithDocType()

		// "asn_statementLabel"
		retSL := fSf(`"asn_statementLabel": { "language": "%s", "literal": "%s" }`, "en-au", value)

		// "dcterms_educationLevel"
		if ts.NotIn(la, "CCP", "GC-NLLP", "GC-NNLP") {
			if value == "Level" { // see doc.typeName: 'Level', update global retEL
				outArrs := []string{}
				for _, y := range getYears(mData, path) {
					outArrs = append(outArrs, fSf(`{ "uri": "%s", "prefLabel": "%s" }`, mYrlvlUri[y], y))
				}
				if len(outArrs) > 0 {
					*pRetEL = sJoin(outArrs, ",")
				}
				*pRetEL = fSf(`"dcterms_educationLevel": [%s]`, *pRetEL)
				*pPrevDocTypePath = path
			}
			// only children path can keep retEL
			if strings.Count(path, ".") < strings.Count(*pPrevDocTypePath, ".") {
				*pRetEL = ""
			}
		} else {
			*pRetEL = ""
		}

		return true, sTrim(sJoin([]string{retSL, *pRetEL}, ","), ",")

	case "code":

		path := fnPathWithCode()

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

		code := ""
		nodeType := ""
		outArrs := []string{}
		for _, item := range items {
			id := item[sLastIndex(item, "/")+1:]
			code = jt.GetStrVal(mNodeData[id+"."+"code"])
			title := jt.GetStrVal(mNodeData[id+"."+"title"])
			nodeType = tool.GetCodeAncestor(mCodeParent, code, 0)
			outArrs = append(outArrs, fSf(`{ "uri": "%s", "prefLabel": "%s" }`, item, title))
		}

		outArrStr := sJoin(outArrs, ",")
		ret := ""

		switch nodeType {
		case "GC":
			ret = fSf(`"%s": [%s]`, "asn_skillEmbodied", outArrStr)
		case "LA":
			ret = fSf(`"%s": [%s]`, "dc_relation", outArrStr)
		case "AS":
			ret = fSf(`"%s": [%s]`, "asn_hasLevel", outArrStr)
		case "CCP":
			ret = fSf(`"%s": [%s]`, "asn_crossSubjectReference", outArrStr)
		default:
			log.Fatalf("'%v' is not one of [GC CCP LA AS], code is '%v'", nodeType, code)
		}

		return true, ret

	default:
		return false, ""
	}
}

func treeProc3(
	data []byte,
	la string,
	mCodeParent map[string]string,
	mNodeData map[string]interface{},
	paths []string,
	// static for filling
	pPrevDocTypePath *string,
	pRetEL *string,

) string {

	var (
		uri4id = "http://rdf.curriculum.edu.au/202110"
	)

	js := string(data)
	mLvlSiblings, _ := jt.FamilyTree(js)

	mData, err := jt.Flatten(data)
	if err != nil {
		log.Fatalln(err)
	}

	re4json, mRE4Each := reMerged()
	// fmt.Println(re4json, len(mRE4Each))

	getPathWithTypeName := fnGetPathByProp("typeName", paths, "")
	getPathWithCode := fnGetPathByProp("code", paths, "")

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
					mData,
					la,
					uri4id,
					mCodeParent,
					mNodeData,
					getPathWithTypeName,
					getPathWithCode,
					pPrevDocTypePath,
					pRetEL,
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
