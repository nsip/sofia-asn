package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/digisan/go-generics/str"
	jt "github.com/digisan/json-tool"
	"github.com/nsip/sofia-asn/tool"
)

func treeProc2(data []byte, uri4id, la string, mUidTitle, mCodeParent map[string]string, mNodeData map[string]interface{}) string {

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
		return true, []string{sibling(path, "id")}, []interface{}{fSf("%s/%v", uri4id, value)}
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
		if str.NotIn(jt.GetStrVal(value), "Learning Area", "Subject") {
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
		if str.In(sval, "ENG", "HAS", "HPE", "LAN", "MAT", "SCI", "TEC", "ART") {
			if subUri, okSubUri := mLaUri[la]; okSubUri {
				ps = append(ps, sibling(path, "dcterms_subject.prefLabel"), sibling(path, "dcterms_subject.uri"))
				vs = append(vs, la, subUri)
			}
		}
		if str.In(sval, "root", "LA") {
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

	jt.RegisterRule(`\.?tags`, func(path string, value interface{}) (ok bool, ps []string, vs []interface{}) {
		return true, []string{sibling(path, "asn_conceptTerm")}, []interface{}{"SCIENCE_TEACHER_BACKGROUND_INFORMATION"}
	})

	jt.RegisterRule(`\.?connections\.[\w\s\d]+\.\d+$`, func(path string, value interface{}) (ok bool, ps []string, vs []interface{}) {
		sval := jt.GetStrVal(value)
		id := sval[strings.LastIndex(sval, "/")+1:]
		// fmt.Println(id)

		code := jt.GetStrVal(mNodeData[id+"."+"code"])
		// fmt.Println(code)

		title := jt.GetStrVal(mNodeData[id+"."+"title"])
		// fmt.Println(title)

		nodeType := tool.GetCodeAncestor(mCodeParent, code, 0)
		// fmt.Println(nodeType)

		switch nodeType {
		case "GC":
			ps = []string{path + ".asn_skillEmbodied.uri", path + ".asn_skillEmbodied.prefLabel"}
		case "LA":
			ps = []string{path + ".dc_relation.uri", path + ".dc_relation.prefLabel"}
		case "AS":
			ps = []string{path + ".asn_hasLevel.uri", path + ".asn_hasLevel.prefLabel"}
		case "CCP":
			ps = []string{path + ".asn_crossSubjectReference.uri", path + ".asn_crossSubjectReference.prefLabel"}
		default:
			log.Fatalf("'%v' is not one of [GC CCP LA AS], code is '%v'", nodeType, code)
		}
		return true, ps, []interface{}{value, title}
	})

	return jt.TransformUnderFirstRule(mData, data)
}
