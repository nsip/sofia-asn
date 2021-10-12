package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	jt "github.com/digisan/json-tool"
	"github.com/nsip/sofia-asn/tool"
	"github.com/tidwall/gjson"
)

var (
	mLaUri = map[string]string{
		"English":                        `http://vocabulary.curriculum.edu.au/framework/E`,
		"The Arts":                       `http://vocabulary.curriculum.edu.au/framework/A`,
		"Health and Physical Education":  `http://vocabulary.curriculum.edu.au/framework/P`,
		"Humanities and Social Sciences": `http://vocabulary.curriculum.edu.au/framework/U`,
		"Languages":                      `http://vocabulary.curriculum.edu.au/framework/L`,
		"Mathematics":                    `http://vocabulary.curriculum.edu.au/framework/M`,
		"Science":                        `http://vocabulary.curriculum.edu.au/framework/S`,
		"Technologies":                   `http://vocabulary.curriculum.edu.au/framework/T`,
		"Work Studies":                   `http://vocabulary.curriculum.edu.au/framework/W`,
	}

	mYrlvlUri = map[string]string{
		"Early years":     `http://vocabulary.curriculum.edu.au/schoolLevel/-`,
		"Foundation Year": `http://vocabulary.curriculum.edu.au/schoolLevel/0`,
		"Year 1":          `http://vocabulary.curriculum.edu.au/schoolLevel/1`,
		"Year 2":          `http://vocabulary.curriculum.edu.au/schoolLevel/2`,
		"Year 3":          `http://vocabulary.curriculum.edu.au/schoolLevel/3`,
		"Year 4":          `http://vocabulary.curriculum.edu.au/schoolLevel/4`,
		"Year 5":          `http://vocabulary.curriculum.edu.au/schoolLevel/5`,
		"Year 6":          `http://vocabulary.curriculum.edu.au/schoolLevel/6`,
		"Year 7":          `http://vocabulary.curriculum.edu.au/schoolLevel/7`,
		"Year 8":          `http://vocabulary.curriculum.edu.au/schoolLevel/8`,
		"Year 9":          `http://vocabulary.curriculum.edu.au/schoolLevel/9`,
		"Year 10":         `http://vocabulary.curriculum.edu.au/schoolLevel/10`,
		"Year 11":         `http://vocabulary.curriculum.edu.au/schoolLevel/11`,
		"Year 12":         `http://vocabulary.curriculum.edu.au/schoolLevel/12`,
	}
)

type asnjson struct {
	Id string `json:"Id"` // DIRECT

	///////////////

	Dcterms_modified struct {
		Literal string `json:"literal"` // DIRECT
	} `json:"dcterms_modified"`

	Dcterms_subject struct {
		Uri       string `json:"uri"`       // derived
		PrefLabel string `json:"prefLabel"` // derived
	} `json:"dcterms_subject"`

	Dcterms_educationLevel []struct {
		Uri       string `json:"uri"`       // derived
		PrefLabel string `json:"prefLabel"` // derived
	} `json:"dcterms_educationLevel"`

	Dcterms_description struct {
		Uri       string `json:"uri"`
		PrefLabel string `json:"prefLabel"`
	} `json:"dcterms_description"`

	Dcterms_title struct {
		Literal  string `json:"literal"`  // DIRECT
		Language string `json:"language"` // boilerplate
	} `json:"dcterms_title"`

	Dcterms_rights struct {
		Literal  string `json:"literal"`  // boilerplate
		Language string `json:"language"` // boilerplate
	} `json:"dcterms_rights"`

	Dcterms_rightsHolder struct {
		Literal  string `json:"literal"`  // boilerplate
		Language string `json:"language"` // boilerplate
	} `json:"dcterms_rightsHolder"`

	///////////////

	Asn_statementLabel struct {
		Literal  string `json:"literal"`  // DIRECT
		Language string `json:"language"` // boilerplate
	} `json:"asn_statementLabel"`

	Asn_statementNotation struct {
		Literal  string `json:"literal"`  // DIRECT
		Language string `json:"language"` // boilerplate
	} `json:"asn_statementNotation"`

	Asn_skillEmbodied []struct {
		Uri       string `json:"uri"`       // Predicate
		PrefLabel string `json:"prefLabel"` // Predicate
	} `json:"asn_skillEmbodied"`

	Asn_authorityStatus struct {
		Uri string `json:"uri"` // boilerplate
	} `json:"asn_authorityStatus"`

	Asn_indexingStatus struct {
		Uri string `json:"uri"` // boilerplate
	} `json:"asn_indexingStatus"`

	Asn_hasLevel []struct {
		Uri       string `json:"uri"`       // Predicate
		PrefLabel string `json:"prefLabel"` // Predicate
	} `json:"asn_hasLevel"`

	Asn_crossSubjectReference []struct {
		Uri       string `json:"uri"`       // Predicate
		PrefLabel string `json:"prefLabel"` // Predicate
	} `json:"asn_crossSubjectReference"`

	Asn_conceptTerm string `json:"asn_conceptTerm"` // tag key

	///////////////

	Dc_relation []struct {
		Uri       string `json:"uri"`       // Predicate
		PrefLabel string `json:"prefLabel"` // Predicate
	} `json:"dc_relation"`

	///////////////

	Cls string `json:"cls"` // boilerplate

	Leaf string `json:"leaf"` // boilerplate

	Text string `json:"text"` // DIRECT

	Children []string `json:"children"` //DIRECT
}

func yearsSplit(yearstr string) (ret []string) {
	if strings.Contains(yearstr, "Foundation Year") {
		ret = append(ret, "Foundation Year")
	}
	r := regexp.MustCompile(`\d+( and \d+)*$`)
	ss := r.FindAllString(yearstr, 1)
	if len(ss) > 0 {
		s := ss[0]
		yn := strings.Split(s, "and")
		for _, y := range yn {
			y = strings.Trim(y, " ")
			ret = append(ret, "Year "+y)
		}
	}
	return
}

func scanIdTitle(data []byte) map[string]string {
	m := make(map[string]string)
	tool.ScanNode(data, func(i int, id, block string) bool {
		e := bytes.LastIndexAny(data, "}")
		data = data[:e+1]
		uid := gjson.Get(block, "id").String()
		title := gjson.Get(block, "title").String()
		m[uid] = title
		return true
	})
	return m
}

func nodeProc(data []byte, outdir, outname, sofiaTreeFile, pref4children string) {

	e := bytes.LastIndexAny(data, "}")
	data = data[:e+1]

	outdir = strings.Trim(outdir, `./\`)
	parts := []string{}
	out := ""

	dataTree, err := os.ReadFile(sofiaTreeFile)
	if err != nil {
		panic(err)
	}
	mCodeParent := tool.GetCodeParentMap(dataTree)

	mUidTitle := scanIdTitle(data)

	tool.ScanNode(data, func(i int, id, block string) bool {

		code := gjson.Get(block, "code").String()

		laTitle := tool.GetAncestorTitle(mCodeParent, code, "")
		if laTitle == "" {
			// fmt.Println("Learning area missing:", code)
		}
		subUri, okSubUri := mLaUri[laTitle]

		var years []string
		if tn := gjson.Get(block, "doc.typeName").String(); tn == "Level" {
			yrTitle := gjson.Get(block, "title").String()
			years = yearsSplit(yrTitle)
		}

		// fmt.Println(i, id, code)

		nodeType := tool.GetCodeAncestor(mCodeParent, code, 0)

		mConnUri := make(map[string]string)
		result := gjson.Get(block, "connections.*")
		if result.IsArray() {
			for _, rUri := range result.Array() {
				uri := rUri.String()
				mConnUri[uri] = mUidTitle[uri]
			}
		}

		rstChildren := gjson.Get(block, "children")
		rstTags := gjson.Get(block, "tags")

		////////////////////////////////////////////////////////

		aj := asnjson{}

		// Direct
		aj.Id = gjson.Get(block, "id").String()
		aj.Dcterms_modified.Literal = gjson.Get(block, "created_at").String()
		aj.Dcterms_title.Literal = gjson.Get(block, "title").String()
		aj.Asn_statementLabel.Literal = gjson.Get(block, "doc.typeName").String()
		aj.Asn_statementNotation.Literal = gjson.Get(block, "code").String()
		aj.Text = gjson.Get(block, "text").String()
		for _, c := range rstChildren.Array() {
			aj.Children = append(aj.Children, pref4children+c.String())
		}

		// Derived
		if okSubUri {
			aj.Dcterms_subject.Uri = subUri
			aj.Dcterms_subject.PrefLabel = laTitle
		}
		for _, y := range years {
			aj.Dcterms_educationLevel = append(aj.Dcterms_educationLevel, struct {
				Uri       string "json:\"uri\""
				PrefLabel string "json:\"prefLabel\""
			}{
				Uri:       mYrlvlUri[y],
				PrefLabel: y,
			})
		}

		// Boilerplate
		aj.Dcterms_title.Language = "en-au"
		aj.Asn_statementLabel.Language = "en-au"
		aj.Asn_statementNotation.Language = "en-au"
		aj.Dcterms_rights.Language = "en-au"
		aj.Dcterms_rightsHolder.Language = "en-au"
		aj.Asn_authorityStatus.Uri = `http://purl.org/ASN/scheme/ASNAuthorityStatus/Original`
		aj.Asn_indexingStatus.Uri = `http://purl.org/ASN/scheme/ASNIndexingStatus/No`
		aj.Dcterms_rights.Literal = `©Copyright Australian Curriculum, Assessment and Reporting Authority`
		aj.Dcterms_rightsHolder.Literal = `Australian Curriculum, Assessment and Reporting Authority`
		if rstChildren.IsArray() {
			aj.Cls = "folder"
		} else {
			aj.Leaf = "true"
		}
		if rstTags.IsObject() {
			aj.Asn_conceptTerm = "SCIENCE_TEACHER_BACKGROUND_INFORMATION"
		}

		// Predicate
		switch nodeType {
		case "GC":
			for uri, title := range mConnUri {
				aj.Asn_skillEmbodied = append(aj.Asn_skillEmbodied, struct {
					Uri       string "json:\"uri\""
					PrefLabel string "json:\"prefLabel\""
				}{
					Uri:       uri,
					PrefLabel: title,
				})
			}

		case "LA":
			for uri, title := range mConnUri {
				aj.Dc_relation = append(aj.Dc_relation, struct {
					Uri       string "json:\"uri\""
					PrefLabel string "json:\"prefLabel\""
				}{
					Uri:       uri,
					PrefLabel: title,
				})
			}

		case "AS":
			for uri, title := range mConnUri {
				aj.Asn_hasLevel = append(aj.Asn_hasLevel, struct {
					Uri       string "json:\"uri\""
					PrefLabel string "json:\"prefLabel\""
				}{
					Uri:       uri,
					PrefLabel: title,
				})
			}

		case "CCP":
			for uri, title := range mConnUri {
				aj.Asn_crossSubjectReference = append(aj.Asn_crossSubjectReference, struct {
					Uri       string "json:\"uri\""
					PrefLabel string "json:\"prefLabel\""
				}{
					Uri:       uri,
					PrefLabel: title,
				})
			}

		default:
			log.Printf("'%v' is not one of [GC CCP LA AS]", nodeType)
		}

		////////////////////////////////////////////////////////////////

		if bytes, err := json.Marshal(aj); err == nil {
			parts = append(parts, string(bytes))
		}

		return true
	})

	out = "[" + strings.Join(parts, ",") + "]" // combine whole
	out = jt.FmtStr(out, "  ")                 // format json
	out = trimNodeProc(out)                    // remove empty object, empty string

	if !strings.HasSuffix(outname, ".json") {
		outname += ".json"
	}
	os.WriteFile(fmt.Sprintf("./%s/%s", outdir, outname), []byte(out), os.ModePerm)
}

// asnjson need to be formatted
func trimNodeProc(asnjson string) string {

	removed := asnjson
	rNullRemove := regexp.MustCompile(`[,{](\n)(\s)+"[^"]+":(\s)+null[,\n]`)
	rEmptyStrRemove := regexp.MustCompile(`[,{](\n)(\s)+"[^"]+":(\s)+""[,\n]`)
	rEmptyObjRemove := regexp.MustCompile(`[,{](\n)(\s)+"[^"]+":(\s)+\{[\n\s]*\}[,\n]`)

	for _, re := range []*regexp.Regexp{rNullRemove, rEmptyStrRemove, rEmptyObjRemove} {
	AGAIN:
		rm := false
		removed = re.ReplaceAllStringFunc(removed, func(s string) string {
			rm = true
			if s[0] == '{' && s[len(s)-1] == ',' {
				return "{"
			}
			if s[0] == '{' && s[len(s)-1] == '\n' {
				return "{"
			}
			if s[0] == ',' && s[len(s)-1] == ',' {
				return ","
			}
			if s[0] == ',' && s[len(s)-1] == '\n' {
				return "\n"
			}
			return s
		})
		if rm {
			goto AGAIN
		}
	}

	return removed
}

func getIdBlock(nodefile string) (mIdBlock, mIdBlockLeaf map[string]string) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r, err := os.Open(nodefile)
	if err != nil {
		log.Fatalln(err)
	}

	result, ok := jt.ScanObject(ctx, r, true, true, jt.OUT_FMT)
	if !ok {
		log.Fatalln("node file is NOT JSON array")
	}

	mIdBlock = make(map[string]string)
	mIdBlockLeaf = make(map[string]string)

	for r := range result {
		if r.Err != nil {
			log.Fatalln(r.Err)
		}
		id := gjson.Get(r.Obj, "Id").String()
		mIdBlock[id] = r.Obj

		hasChildren := gjson.Get(r.Obj, "children").IsArray()
		if !hasChildren {
			mIdBlockLeaf[id] = r.Obj
		}
	}

	return
}

func childrenId(cBlock string) (cid []string) {
	s := strings.Index(cBlock, "[")
	e := strings.LastIndex(cBlock, "]")
	cBlock = cBlock[s+1 : e]
	cBlock = strings.Trim(cBlock, " \n\t")
	for _, id := range strings.Split(cBlock, ",") {
		cid = append(cid, strings.Trim(id, " \n\t"))
	}
	return
}

func childrenRepl(inpath string, mIdBlock map[string]string) string {

	data, err := os.ReadFile(inpath)
	if err != nil {
		log.Fatalln(err)
	}

	rChildren := regexp.MustCompile(`"children":(\s)+\[(\n\s+"http[^"]+",?)+\n\s+\]`)
	js := string(data)
	repl := false

AGAIN:
	repl = false
	js = rChildren.ReplaceAllStringFunc(js, func(s string) string {
		for _, id := range childrenId(s) {
			id = id[1 : len(id)-1]
			if block, ok := mIdBlock[id]; ok {
				s = strings.ReplaceAll(s, "\""+id+"\"", block)
				repl = true
			}
		}
		return s
	})

	if repl {
		goto AGAIN
	}

	return jt.FmtStr(js, "  ")
}

func rmSingleLeaf(input string) string {

	r := strings.NewReader(input)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nodes := []string{}

	cRst, _ := jt.ScanObject(ctx, r, true, true, jt.OUT_FMT)
	for rst := range cRst {
		if rst.Err != nil {
			log.Fatalln(rst.Err)
		}
		block := rst.Obj

		rstLeaf := gjson.Get(block, "leaf")
		// fmt.Println(rstLeaf.String())

		if rstLeaf.String() != "true" {
			nodes = append(nodes, block)
		}
	}

	js := "[" + strings.Join(nodes, ",") + "]"
	return jt.FmtStr(js, "  ")
}
