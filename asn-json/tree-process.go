package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	. "github.com/digisan/go-generics/v2"
	jt "github.com/digisan/json-tool"
	"github.com/nsip/sofia-asn/tool"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func treeProc(data []byte, uri4id, la string, mCodeParent, mUidTitle map[string]string) string {

	uri4id = strings.TrimSuffix(uri4id, "/")
	js := string(data)

	var (
		fSf        = fmt.Sprintf
		fetchValue = func(kvstr string) string {
			r := regexp.MustCompile(`:\s*"`)
			loc := r.FindAllStringIndex(kvstr, 1)[0]
			start := loc[1]
			end := strings.LastIndex(kvstr, `"`)
			return kvstr[start:end]
		}
		setWhenEquals = func(c0, c1 byte) string {
			if c0 == c1 {
				return string(c0)
			}
			return ""
		}
	)

	//////
	// id
	rId := regexp.MustCompile(`"uuid":\s*"[\d\w]{8}-[\d\w]{4}-[\d\w]{4}-[\d\w]{4}-[\d\w]{12}"`)
	js = rId.ReplaceAllStringFunc(js, func(s string) string {
		return fSf(`"id": "%s/%s"`, uri4id, fetchValue(s))
	})

	// type => removed
	rType := regexp.MustCompile(`"type":\s*"\w+",?\n?`)
	js = rType.ReplaceAllStringFunc(js, func(s string) string {
		return ""
	})

	// created_at
	rCreated := regexp.MustCompile(`"created_at":\s*"[^"]+"`)
	js = rCreated.ReplaceAllStringFunc(js, func(s string) string {
		return fSf(`"dcterms_modified": { "literal": "%s" }`, fetchValue(s))
	})

	// title
	rTitle := regexp.MustCompile(`"title":\s*".+",?\n`)
	js = rTitle.ReplaceAllStringFunc(js, func(s string) string {
		sfx0, sfx1 := setWhenEquals(s[len(s)-1], '\n'), setWhenEquals(s[len(s)-2], ',')
		return fSf(`"dcterms_title": { "language": "%s", "literal": "%s" }%s%s`, "en-au", fetchValue(s), sfx1, sfx0)
	})

	// doc.typeName
	rDocType := regexp.MustCompile(`"doc":\s*{[\n\s]*"typeName":\s*"[^"]+"[\n\s]*},?\n`)
	js = rDocType.ReplaceAllStringFunc(js, func(s string) string {
		sfx0, sfx1 := setWhenEquals(s[len(s)-1], '\n'), setWhenEquals(s[len(s)-2], ',')
		return fSf(`"asn_statementLabel": { "language": "%s", "literal": "%s" }%s%s`, "en-au", fetchValue(s), sfx1, sfx0)
	})

	// code
	rCode := regexp.MustCompile(`"code":\s*"[^"]+"`)
	js = rCode.ReplaceAllStringFunc(js, func(s string) string {
		return fSf(`"asn_statementNotation": { "language": "%s", "literal": "%s" }`, "en-au", fetchValue(s))
	})

	// text
	rText := regexp.MustCompile(`"text":\s*".+",?\n`)
	js = rText.ReplaceAllStringFunc(js, func(s string) string {
		sfx0, sfx1 := setWhenEquals(s[len(s)-1], '\n'), setWhenEquals(s[len(s)-2], ',')
		return fSf(`"text": "%s"%s%s`, fetchValue(s), sfx1, sfx0)
	})

	//////
	// dcterms_subject
	if subUri, okSubUri := mLaUri[la]; okSubUri {
		rId4uri := regexp.MustCompile(`"id":\s*"http[^"]+",?\n`)
		js = rId4uri.ReplaceAllStringFunc(js, func(s string) string {
			sfx0, sfx1 := setWhenEquals(s[len(s)-1], '\n'), setWhenEquals(s[len(s)-2], ',')
			suffix := fSf(`"dcterms_subject": { "prefLabel": "%s", "uri": "%s" }%s%s`, la, subUri, sfx1, sfx0)
			return s + suffix
		})
	}

	//
	// using 'mLvlSiblings' for new fields, so put it here !!!
	//
	mLvlSiblings, _ := jt.FamilyTree(js)

	// [ dcterms_title, asn_statementLabel ] => dcterms_educationLevel
	mFieldSibling := jt.GetSiblingPath("dcterms_title", "asn_statementLabel", mLvlSiblings)
	for fp, sp := range mFieldSibling {
		if gjson.Get(js, sp+".literal").String() == "Level" {
			for _, y := range yearsSplit(gjson.Get(js, fp+".literal").String()) {
				fmt.Println(y)
				js, _ = sjson.Set(js, jt.NewSibling(fp, "dcterms_educationLevel.uri"), mYrlvlUri[y])
				js, _ = sjson.Set(js, jt.NewSibling(fp, "dcterms_educationLevel.prefLabel"), y)
			}
		}
	}

	// "children"? => add "cls": "folder"; else add "leaf": "true"
	allPaths := jt.GetFieldPaths("id", mLvlSiblings)
	allPaths = Map4SglTyp(allPaths, func(i int, e string) string {
		return jt.ParentPath(e)
	})

	cPaths := jt.GetFieldPaths("children", mLvlSiblings)
	cPaths = Map4SglTyp(cPaths, func(i int, e string) string {
		return jt.ParentPath(e)
	})

	for _, cp := range cPaths {
		if cp != "" {
			js, _ = sjson.Set(js, cp+".cls", "folder")
		} else {
			js, _ = sjson.Set(js, "cls", "folder")
		}
	}

	for _, lp := range Minus(allPaths, cPaths) {
		if lp != "" {
			js, _ = sjson.Set(js, lp+".leaf", "true")
		} else {
			js, _ = sjson.Set(js, "leaf", "true")
		}
	}

	// "tags" => "asn_conceptTerm"
	tPaths := jt.GetFieldPaths("tags", mLvlSiblings)
	for _, tp := range tPaths {
		js, _ = sjson.Set(js, jt.NewSibling(tp, "asn_conceptTerm"), "SCIENCE_TEACHER_BACKGROUND_INFORMATION")
		js, _ = sjson.Delete(js, tp)
	}

	// append some after "id"
	rNewId := regexp.MustCompile(`"id":\s*"http[^"]+",?`)
	// asn_authorityStatus
	js = rNewId.ReplaceAllStringFunc(js, func(s string) string {
		sfx := setWhenEquals(s[len(s)-1], ',')
		uri := `http://purl.org/ASN/scheme/ASNAuthorityStatus/Original`
		suffix := fSf(`"asn_authorityStatus": { "uri": "%s" }%s`, uri, sfx)
		if sfx == "" {
			return s + "," + suffix
		}
		return s + suffix
	})
	// asn_indexingStatus
	js = rNewId.ReplaceAllStringFunc(js, func(s string) string {
		sfx := setWhenEquals(s[len(s)-1], ',')
		uri := `http://purl.org/ASN/scheme/ASNIndexingStatus/No`
		suffix := fSf(`"asn_indexingStatus": { "uri": "%s" }%s`, uri, sfx)
		if sfx == "" {
			return s + "," + suffix
		}
		return s + suffix
	})
	// dcterms_rights
	js = rNewId.ReplaceAllStringFunc(js, func(s string) string {
		sfx := setWhenEquals(s[len(s)-1], ',')
		rights := `©Copyright Australian Curriculum, Assessment and Reporting Authority`
		suffix := fSf(`"dcterms_rights": { "language": "%s", "literal": "%s" }%s`, "en-au", rights, sfx)
		if sfx == "" {
			return s + "," + suffix
		}
		return s + suffix
	})
	// dcterms_rightsHolder
	js = rNewId.ReplaceAllStringFunc(js, func(s string) string {
		sfx := setWhenEquals(s[len(s)-1], ',')
		rh := `Australian Curriculum, Assessment and Reporting Authority`
		suffix := fSf(`"dcterms_rightsHolder": { "language": "%s", "literal": "%s" }%s`, "en-au", rh, sfx)
		if sfx == "" {
			return s + "," + suffix
		}
		return s + suffix
	})

	///////////////////////////////////////////////

	// Predicate, deal with 'connections' array
	cPaths = jt.GetFieldPaths("connections", mLvlSiblings)
	for _, cp := range cPaths {

		block := gjson.Get(js, jt.ParentPath(cp)).String()

		mConnUri := make(map[string]string)
		result := gjson.Get(block, "connections.*")
		if result.IsArray() {
			for _, rUri := range result.Array() {
				uri := uri4id + "/" + rUri.String()
				mConnUri[uri] = mUidTitle[uri]
			}
		}

		code := gjson.Get(block, "asn_statementNotation.literal").String()
		nodeType := tool.GetCodeAncestor(mCodeParent, code, 0)
		pred := ""
		switch nodeType {
		case "GC":
			pred = jt.NewSibling(cp, "asn_skillEmbodied")
		case "LA":
			pred = jt.NewSibling(cp, "dc_relation")
		case "AS":
			pred = jt.NewSibling(cp, "asn_hasLevel")
		case "CCP":
			pred = jt.NewSibling(cp, "asn_crossSubjectReference")
		default:
			log.Fatalf("'%v' is not one of [GC CCP LA AS], code is '%v'", nodeType, code)
		}

		i := 0
		for uri, title := range mConnUri {
			js, _ = sjson.Set(js, pred+fmt.Sprintf(".%d.prefLabel", i), title)
			js, _ = sjson.Set(js, pred+fmt.Sprintf(".%d.uri", i), uri)
			i++
		}

		js, _ = sjson.Delete(js, cp)
	}

	return js
}
