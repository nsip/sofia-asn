package main

import (
	"fmt"

	"github.com/digisan/gotk/slice/ts"
	jt "github.com/digisan/json-tool"
)

func getYears(mData map[string]interface{}, path string) []string {
AGAIN:
	sp := jt.NewSibling(path, "doc.typeName")
	if mData[sp] == "Level" {
		return yearsSplit(mData[jt.NewSibling(path, "title")].(string))
	} else {
		path = jt.ParentPath(path)
		goto AGAIN
	}
}

func treeProc2(data []byte, uri4id, la string, mCodeParent, mUidTitle map[string]string) string {

	var (
		fSf     = fmt.Sprintf
		sibling = jt.NewSibling
		uncle   = jt.NewUncle
	)

	js := string(data)
	mLvlSiblings, _ := jt.FamilyTree(js)

	mData, err := jt.Flatten(data)
	if err != nil {
		panic(err)
	}

	jt.RegisterRule(`\.?uuid$`, func(path string, value interface{}) (ok bool, ps []string, vs []interface{}) {
		return true, []string{sibling(path, "Id")}, []interface{}{fSf("%s/%v", uri4id, value)}
	})

	jt.RegisterRule(`\.?type$`, func(path string, value interface{}) (ok bool, ps []string, vs []interface{}) {
		return true, []string{""}, []interface{}{nil}
	})

	jt.RegisterRule(`\.?created_at$`, func(path string, value interface{}) (ok bool, ps []string, vs []interface{}) {
		return true, []string{sibling(path, "dcterms_modified.literal")}, []interface{}{value}
	})

	jt.RegisterRule(`\.?title$`, func(path string, value interface{}) (ok bool, ps []string, vs []interface{}) {
		ok = true
		ps = []string{sibling(path, "dcterms_title.language"), sibling(path, "dcterms_title.literal")}
		vs = []interface{}{"en-au", value}
		return
	})

	jt.RegisterRule(`\.?doc\.typeName$`, func(path string, value interface{}) (ok bool, ps []string, vs []interface{}) {
		ok = true
		ps = []string{uncle(path, "asn_statementLabel.language"), uncle(path, "asn_statementLabel.literal")}
		vs = []interface{}{"en-au", value}
		if ts.NotIn(jt.GetStrVal(value), "Learning Area", "Subject") {
			for _, y := range getYears(mData, path) {
				ps = append(ps, uncle(path, "dcterms_educationLevel.uri"), uncle(path, "dcterms_educationLevel.prefLabel"))
				vs = append(vs, mYrlvlUri[y], y)
			}
		}
		return
	})

	jt.RegisterRule(`\.?code$`, func(path string, value interface{}) (ok bool, ps []string, vs []interface{}) {
		ok = true
		ps = []string{sibling(path, "asn_statementNotation.language"), sibling(path, "asn_statementNotation.literal")}
		vs = []interface{}{"en-au", value}

		// add for specific nodes
		sval := jt.GetStrVal(value)
		if sval == "ENG" {
			if subUri, okSubUri := mLaUri[la]; okSubUri {
				ps = append(ps, sibling(path, "dcterms_subject.prefLabel"), sibling(path, "dcterms_subject.uri"))
				vs = append(vs, la, subUri)
			}
		}
		if ts.In(sval, "root", "LA") {
			ps = append(ps, sibling(path, "dcterms_rights.language"), sibling(path, "dcterms_rights.literal"))
			vs = append(vs, "en-au", `©Copyright Australian Curriculum, Assessment and Reporting Authority`)
			ps = append(ps, sibling(path, "dcterms_rightsHolder.language"), sibling(path, "dcterms_rightsHolder.literal"))
			vs = append(vs, "en-au", `Australian Curriculum, Assessment and Reporting Authority`)
		}

		// add one for each node
		ps = append(ps, sibling(path, "asn_authorityStatus.uri"), sibling(path, "asn_indexingStatus.uri"))
		vs = append(vs, `http://purl.org/ASN/scheme/ASNAuthorityStatus/Original`, `http://purl.org/ASN/scheme/ASNIndexingStatus/No`)

		// add cls for which has children
		if jt.HasSiblings(path, mLvlSiblings, "children") {
			ps = append(ps, sibling(path, "cls"))
			vs = append(vs, "folder")
		} else {
			ps = append(ps, sibling(path, "leaf"))
			vs = append(vs, "true")
		}

		return
	})

	jt.RegisterRule(`\.?text$`, func(path string, value interface{}) (ok bool, ps []string, vs []interface{}) {
		return true, []string{sibling(path, "text")}, []interface{}{value}
	})

	jt.RegisterRule(`\.?tags$`, func(path string, value interface{}) (ok bool, ps []string, vs []interface{}) {
		return true, []string{sibling(path, "asn_conceptTerm")}, []interface{}{"SCIENCE_TEACHER_BACKGROUND_INFORMATION"}
	})

	// jt.RegisterRule(`\.?connections`, func(path string, value interface{}) (ok bool, ps []string, vs []interface{}) {
	// 	fmt.Println(path)
	// 	fmt.Println(value)

	// 	pp := parentpath(path)
	// 	metaCode := pp[sLastIndex(pp, ".")+1:]
	// 	fmt.Println(metaCode)

	// 	// uri := uri4id + "/" + jt.GetStrVal(value)
	// 	// mConnUri[uri] = mUidTitle[uri]

	// 	return
	// })

	js = jt.TransformUnderFirstRule(mData, data)

	// Predicate, deal with 'connections' array
	// cPaths := jt.GetFieldPaths("connections", mLvlSiblings)
	// for _, cp := range cPaths {

	// 	block := gjson.Get(js, jt.ParentPath(cp)).String()

	// 	mConnUri := make(map[string]string)
	// 	result := gjson.Get(block, "connections.*")
	// 	if result.IsArray() {
	// 		for _, rUri := range result.Array() {
	// 			uri := uri4id + "/" + rUri.String()
	// 			mConnUri[uri] = mUidTitle[uri]
	// 		}
	// 	}

	// 	code := gjson.Get(block, "asn_statementNotation.literal").String()
	// 	nodeType := tool.GetCodeAncestor(mCodeParent, code, 0)
	// 	pred := ""
	// 	switch nodeType {
	// 	case "GC":
	// 		pred = jt.NewSibling(cp, "asn_skillEmbodied")
	// 	case "LA":
	// 		pred = jt.NewSibling(cp, "dc_relation")
	// 	case "AS":
	// 		pred = jt.NewSibling(cp, "asn_hasLevel")
	// 	case "CCP":
	// 		pred = jt.NewSibling(cp, "asn_crossSubjectReference")
	// 	default:
	// 		log.Fatalf("'%v' is not one of [GC CCP LA AS], code is '%v'", nodeType, code)
	// 	}

	// 	i := 0
	// 	for uri, title := range mConnUri {
	// 		js, _ = sjson.Set(js, pred+fmt.Sprintf(".%d.prefLabel", i), title)
	// 		js, _ = sjson.Set(js, pred+fmt.Sprintf(".%d.uri", i), uri)
	// 		i++
	// 	}

	// 	js, _ = sjson.Delete(js, cp)
	// }

	return js
}

// func treeProc(data []byte, uri4id, la string, mCodeParent, mUidTitle map[string]string) string {

// 	uri4id = strings.TrimSuffix(uri4id, "/")
// 	js := string(data)

// 	var (
// 		fSf        = fmt.Sprintf
// 		fetchValue = func(kvstr string) string {
// 			r := regexp.MustCompile(`:\s*"`)
// 			loc := r.FindAllStringIndex(kvstr, 1)[0]
// 			start := loc[1]
// 			end := strings.LastIndex(kvstr, `"`)
// 			return kvstr[start:end]
// 		}
// 		setWhenEquals = func(c0, c1 byte) string {
// 			if c0 == c1 {
// 				return string(c0)
// 			}
// 			return ""
// 		}
// 	)

// 	//////
// 	// Id
// 	rId := regexp.MustCompile(`"uuid":\s*"[\d\w]{8}-[\d\w]{4}-[\d\w]{4}-[\d\w]{4}-[\d\w]{12}"`)
// 	js = rId.ReplaceAllStringFunc(js, func(s string) string {
// 		return fSf(`"Id": "%s/%s"`, uri4id, fetchValue(s))
// 	})

// 	// type => removed
// 	rType := regexp.MustCompile(`"type":\s*"\w+",?\n?`)
// 	js = rType.ReplaceAllStringFunc(js, func(s string) string {
// 		return ""
// 	})

// 	// created_at
// 	rCreated := regexp.MustCompile(`"created_at":\s*"[^"]+"`)
// 	js = rCreated.ReplaceAllStringFunc(js, func(s string) string {
// 		return fSf(`"dcterms_modified": { "literal": "%s" }`, fetchValue(s))
// 	})

// 	// title
// 	rTitle := regexp.MustCompile(`"title":\s*".+",?\n`)
// 	js = rTitle.ReplaceAllStringFunc(js, func(s string) string {
// 		sfx0, sfx1 := setWhenEquals(s[len(s)-1], '\n'), setWhenEquals(s[len(s)-2], ',')
// 		return fSf(`"dcterms_title": { "language": "%s", "literal": "%s" }%s%s`, "en-au", fetchValue(s), sfx1, sfx0)
// 	})

// 	// doc.typeName
// 	rDocType := regexp.MustCompile(`"doc":\s*{[\n\s]*"typeName":\s*"[^"]+"[\n\s]*},?\n`)
// 	js = rDocType.ReplaceAllStringFunc(js, func(s string) string {
// 		sfx0, sfx1 := setWhenEquals(s[len(s)-1], '\n'), setWhenEquals(s[len(s)-2], ',')
// 		return fSf(`"asn_statementLabel": { "language": "%s", "literal": "%s" }%s%s`, "en-au", fetchValue(s), sfx1, sfx0)
// 	})

// 	// code
// 	rCode := regexp.MustCompile(`"code":\s*"[^"]+"`)
// 	js = rCode.ReplaceAllStringFunc(js, func(s string) string {
// 		return fSf(`"asn_statementNotation": { "language": "%s", "literal": "%s" }`, "en-au", fetchValue(s))
// 	})

// 	// text
// 	rText := regexp.MustCompile(`"text":\s*".+",?\n`)
// 	js = rText.ReplaceAllStringFunc(js, func(s string) string {
// 		sfx0, sfx1 := setWhenEquals(s[len(s)-1], '\n'), setWhenEquals(s[len(s)-2], ',')
// 		return fSf(`"text": "%s"%s%s`, fetchValue(s), sfx1, sfx0)
// 	})

// 	//////
// 	// dcterms_subject
// 	if subUri, okSubUri := mLaUri[la]; okSubUri {
// 		rId4uri := regexp.MustCompile(`"Id":\s*"http[^"]+",?\n`)
// 		js = rId4uri.ReplaceAllStringFunc(js, func(s string) string {
// 			sfx0, sfx1 := setWhenEquals(s[len(s)-1], '\n'), setWhenEquals(s[len(s)-2], ',')
// 			suffix := fSf(`"dcterms_subject": { "prefLabel": "%s", "uri": "%s" }%s%s`, la, subUri, sfx1, sfx0)
// 			return s + suffix
// 		})
// 	}

// 	//
// 	// using 'mLvlSiblings' for new fields, so put it here !!!
// 	//
// 	mLvlSiblings, _ := jt.FamilyTree(js)

// 	// [ dcterms_title, asn_statementLabel ] => dcterms_educationLevel
// 	mFieldSibling := jt.GetSiblingPath("dcterms_title", "asn_statementLabel", mLvlSiblings)
// 	for fp, sp := range mFieldSibling {
// 		if gjson.Get(js, sp+".literal").String() == "Level" {
// 			for _, y := range yearsSplit(gjson.Get(js, fp+".literal").String()) {
// 				fmt.Println(y)
// 				js, _ = sjson.Set(js, jt.NewSibling(fp, "dcterms_educationLevel.uri"), mYrlvlUri[y])
// 				js, _ = sjson.Set(js, jt.NewSibling(fp, "dcterms_educationLevel.prefLabel"), y)
// 			}
// 		}
// 	}

// 	// "children"? => add "cls": "folder"; else add "leaf": "true"
// 	allPaths := jt.GetFieldPaths("Id", mLvlSiblings)
// 	allPaths = ts.FM(allPaths, nil, func(i int, e string) string {
// 		return jt.ParentPath(e)
// 	})

// 	cPaths := jt.GetFieldPaths("children", mLvlSiblings)
// 	cPaths = ts.FM(cPaths, nil, func(i int, e string) string {
// 		return jt.ParentPath(e)
// 	})

// 	for _, cp := range cPaths {
// 		if cp != "" {
// 			js, _ = sjson.Set(js, cp+".cls", "folder")
// 		} else {
// 			js, _ = sjson.Set(js, "cls", "folder")
// 		}
// 	}

// 	for _, lp := range ts.Minus(allPaths, cPaths) {
// 		if lp != "" {
// 			js, _ = sjson.Set(js, lp+".leaf", "true")
// 		} else {
// 			js, _ = sjson.Set(js, "leaf", "true")
// 		}
// 	}

// 	// "tags" => "asn_conceptTerm"
// 	tPaths := jt.GetFieldPaths("tags", mLvlSiblings)
// 	for _, tp := range tPaths {
// 		js, _ = sjson.Set(js, jt.NewSibling(tp, "asn_conceptTerm"), "SCIENCE_TEACHER_BACKGROUND_INFORMATION")
// 		js, _ = sjson.Delete(js, tp)
// 	}

// 	// append some after "Id"
// 	rNewId := regexp.MustCompile(`"Id":\s*"http[^"]+",?`)
// 	// asn_authorityStatus
// 	js = rNewId.ReplaceAllStringFunc(js, func(s string) string {
// 		sfx := setWhenEquals(s[len(s)-1], ',')
// 		uri := `http://purl.org/ASN/scheme/ASNAuthorityStatus/Original`
// 		suffix := fSf(`"asn_authorityStatus": { "uri": "%s" }%s`, uri, sfx)
// 		if sfx == "" {
// 			return s + "," + suffix
// 		}
// 		return s + suffix
// 	})
// 	// asn_indexingStatus
// 	js = rNewId.ReplaceAllStringFunc(js, func(s string) string {
// 		sfx := setWhenEquals(s[len(s)-1], ',')
// 		uri := `http://purl.org/ASN/scheme/ASNIndexingStatus/No`
// 		suffix := fSf(`"asn_indexingStatus": { "uri": "%s" }%s`, uri, sfx)
// 		if sfx == "" {
// 			return s + "," + suffix
// 		}
// 		return s + suffix
// 	})
// 	// dcterms_rights
// 	js = rNewId.ReplaceAllStringFunc(js, func(s string) string {
// 		sfx := setWhenEquals(s[len(s)-1], ',')
// 		rights := `©Copyright Australian Curriculum, Assessment and Reporting Authority`
// 		suffix := fSf(`"dcterms_rights": { "language": "%s", "literal": "%s" }%s`, "en-au", rights, sfx)
// 		if sfx == "" {
// 			return s + "," + suffix
// 		}
// 		return s + suffix
// 	})
// 	// dcterms_rightsHolder
// 	js = rNewId.ReplaceAllStringFunc(js, func(s string) string {
// 		sfx := setWhenEquals(s[len(s)-1], ',')
// 		rh := `Australian Curriculum, Assessment and Reporting Authority`
// 		suffix := fSf(`"dcterms_rightsHolder": { "language": "%s", "literal": "%s" }%s`, "en-au", rh, sfx)
// 		if sfx == "" {
// 			return s + "," + suffix
// 		}
// 		return s + suffix
// 	})

// 	///////////////////////////////////////////////

// 	// Predicate, deal with 'connections' array
// 	cPaths = jt.GetFieldPaths("connections", mLvlSiblings)
// 	for _, cp := range cPaths {

// 		block := gjson.Get(js, jt.ParentPath(cp)).String()

// 		mConnUri := make(map[string]string)
// 		result := gjson.Get(block, "connections.*")
// 		if result.IsArray() {
// 			for _, rUri := range result.Array() {
// 				uri := uri4id + "/" + rUri.String()
// 				mConnUri[uri] = mUidTitle[uri]
// 			}
// 		}

// 		code := gjson.Get(block, "asn_statementNotation.literal").String()
// 		nodeType := tool.GetCodeAncestor(mCodeParent, code, 0)
// 		pred := ""
// 		switch nodeType {
// 		case "GC":
// 			pred = jt.NewSibling(cp, "asn_skillEmbodied")
// 		case "LA":
// 			pred = jt.NewSibling(cp, "dc_relation")
// 		case "AS":
// 			pred = jt.NewSibling(cp, "asn_hasLevel")
// 		case "CCP":
// 			pred = jt.NewSibling(cp, "asn_crossSubjectReference")
// 		default:
// 			log.Fatalf("'%v' is not one of [GC CCP LA AS], code is '%v'", nodeType, code)
// 		}

// 		i := 0
// 		for uri, title := range mConnUri {
// 			js, _ = sjson.Set(js, pred+fmt.Sprintf(".%d.prefLabel", i), title)
// 			js, _ = sjson.Set(js, pred+fmt.Sprintf(".%d.uri", i), uri)
// 			i++
// 		}

// 		js, _ = sjson.Delete(js, cp)
// 	}

// 	return js
// }
