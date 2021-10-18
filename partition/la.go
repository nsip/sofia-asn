package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func la(js, outdir string) {

	outdir = strings.Trim(outdir, `./\`)

	mNameLaAs := map[string]string{
		"English":      "English",
		"HASS":         "Humanities and Social Sciences",
		"HPE":          "Health and Physical Education",
		"Languages":    "Languages",
		"Mathematics":  "Mathematics",
		"Science":      "Science",
		"Technologies": "Technologies",
		"The Arts":     "The Arts",
	}

	mRoot := map[string]interface{}{
		"code":       gjson.Get(js, "code").String(),
		"uuid":       gjson.Get(js, "uuid").String(),
		"type":       gjson.Get(js, "type").String(),
		"created_at": gjson.Get(js, "created_at").String(),
		"title":      gjson.Get(js, "title").String(),
		"children":   nil,
	}

	mASfield := map[string]interface{}{
		"code":       "",
		"uuid":       "",
		"type":       "",
		"created_at": "",
		"title":      "",
		"children":   nil,
	}

	mLAfield := map[string]interface{}{
		"code":       "",
		"uuid":       "",
		"type":       "",
		"created_at": "",
		"title":      "",
		"children":   nil,
	}

	mAS := make(map[string][]string)
	mLA := make(map[string][]string)

	valueC1 := gjson.Get(js, "children")
	if valueC1.IsArray() {
		for _, r1 := range valueC1.Array() {
			if r1.IsObject() { // "Achievement Standards", "Cross-curriculum Priorities", "General Capabilities", "Learning Areas"
				block1 := r1.String()
				valueTitle1 := gjson.Get(block1, "title")
				title1str := valueTitle1.String()
				fmt.Println(title1str, ":")

				valueC2 := gjson.Get(block1, "children")
				if valueC2.IsArray() {
					for _, r2 := range valueC2.Array() {
						if r2.IsObject() { // "English", "Mathematics", etc.
							block2 := r2.String()
							valueTitle2 := gjson.Get(block2, "title")
							title2str := valueTitle2.String()
							fmt.Println("	", title2str)

							switch title1str {
							case "Achievement Standards":
								mAS[title2str] = append(mAS[title2str], block2)
							case "Learning Areas":
								mLA[title2str] = append(mLA[title2str], block2)
							}
						}
					}
				}

				switch title1str {
				case "Achievement Standards":
					mASfield["code"] = gjson.Get(block1, "code").String()
					mASfield["uuid"] = gjson.Get(block1, "uuid").String()
					mASfield["type"] = gjson.Get(block1, "type").String()
					mASfield["created_at"] = gjson.Get(block1, "created_at").String()
					mASfield["title"] = gjson.Get(block1, "title").String()
					mASfield["children"] = mAS
				case "Learning Areas":
					mLAfield["code"] = gjson.Get(block1, "code").String()
					mLAfield["uuid"] = gjson.Get(block1, "uuid").String()
					mLAfield["type"] = gjson.Get(block1, "type").String()
					mLAfield["created_at"] = gjson.Get(block1, "created_at").String()
					mLAfield["title"] = gjson.Get(block1, "title").String()
					mLAfield["children"] = mLA
				}
			}
		}
	}

	if len(mAS) != len(mLA) {
		log.Println("[Achievement Standards] children count is NOT same as [Learning Areas] children count")
		if len(mLA) < len(mAS) {
			log.Fatalln("[Learning Areas] children count less than [Achievement Standards] children count")
		}
	}

NEXTLA:
	for la, blockLA := range mLA {
		out := ""
		for as, blockAS := range mAS {
			if la == as || mNameLaAs[la] == as {

				out, _ = sjson.Set(out, "code", mRoot["code"])
				out, _ = sjson.Set(out, "uuid", mRoot["uuid"])
				out, _ = sjson.Set(out, "type", mRoot["type"])
				out, _ = sjson.Set(out, "created_at", mRoot["created_at"])
				out, _ = sjson.Set(out, "title", mRoot["title"])

				out, _ = sjson.Set(out, "children.0.code", mASfield["code"])
				out, _ = sjson.Set(out, "children.0.uuid", mASfield["uuid"])
				out, _ = sjson.Set(out, "children.0.type", mASfield["type"])
				out, _ = sjson.Set(out, "children.0.created_at", mASfield["created_at"])
				out, _ = sjson.Set(out, "children.0.title", mASfield["title"])
				for i, bAS := range blockAS {
					path := fmt.Sprintf("children.0.children.%d", i)
					out, _ = sjson.SetRaw(out, path, bAS)
				}

				out, _ = sjson.Set(out, "children.1.code", mLAfield["code"])
				out, _ = sjson.Set(out, "children.1.uuid", mLAfield["uuid"])
				out, _ = sjson.Set(out, "children.1.type", mLAfield["type"])
				out, _ = sjson.Set(out, "children.1.created_at", mLAfield["created_at"])
				out, _ = sjson.Set(out, "children.1.title", mLAfield["title"])
				for i, bLA := range blockLA {
					path := fmt.Sprintf("children.1.children.%d", i)
					out, _ = sjson.SetRaw(out, path, bLA)
				}

				// out = jt.FmtStr(out, "  ")
				err := os.WriteFile(fmt.Sprintf("./%s/la-%s.json", outdir, la), []byte(out), os.ModePerm)
				if err != nil {
					fmt.Println(err)
				}

				continue NEXTLA
			}
		}
	}
}
