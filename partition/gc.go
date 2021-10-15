package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/digisan/gotk/slice/ts"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func gc(js, outdir string) {

	outdir = strings.Trim(outdir, `./\`)

	var (
		gcTitles = []string{
			"National Literacy Learning Progression",
			"National Numeracy Learning Progression",
			"Critical and Creative Thinking",
			"Personal and Social Capability",
			"Digital Literacy",
			"Intercultural Understanding",
			"Ethical Understanding",
		}
	)

	// L0
	mRoot := map[string]interface{}{
		"code":       gjson.Get(js, "code").String(),
		"uuid":       gjson.Get(js, "uuid").String(),
		"type":       gjson.Get(js, "type").String(),
		"created_at": gjson.Get(js, "created_at").String(),
		"title":      gjson.Get(js, "title").String(),
		"children":   nil,
	}

	// L1
	mGC := map[string]interface{}{
		"code":       "",
		"uuid":       "",
		"type":       "",
		"created_at": "",
		"title":      "",
		"children":   nil,
	}

	// L2
	mNLLP := map[string]interface{}{
		"code":       "",
		"uuid":       "",
		"type":       "",
		"created_at": "",
		"title":      "",
		"children":   nil,
	}

	// L2
	mNNLP := map[string]interface{}{
		"code":       "",
		"uuid":       "",
		"type":       "",
		"created_at": "",
		"title":      "",
		"children":   nil,
	}

	// L2
	mCCT := map[string]interface{}{
		"code":       "",
		"uuid":       "",
		"type":       "",
		"created_at": "",
		"title":      "",
		"children":   nil,
	}

	// L2
	mPSC := map[string]interface{}{
		"code":       "",
		"uuid":       "",
		"type":       "",
		"created_at": "",
		"title":      "",
		"children":   nil,
	}

	// L2
	mDL := map[string]interface{}{
		"code":       "",
		"uuid":       "",
		"type":       "",
		"created_at": "",
		"title":      "",
		"children":   nil,
	}

	// L2
	mIU := map[string]interface{}{
		"code":       "",
		"uuid":       "",
		"type":       "",
		"created_at": "",
		"title":      "",
		"children":   nil,
	}

	// L2
	mEU := map[string]interface{}{
		"code":       "",
		"uuid":       "",
		"type":       "",
		"created_at": "",
		"title":      "",
		"children":   nil,
	}

	var (
		mL2s = []map[string]interface{}{mNLLP, mNNLP, mCCT, mPSC, mDL, mIU, mEU}
	)

	valueC1 := gjson.Get(js, "children")
	if valueC1.IsArray() {
		for _, r1 := range valueC1.Array() {
			if r1.IsObject() { // "Achievement Standards", "Cross-curriculum Priorities", "General Capabilities", "Learning Areas"
				block1 := r1.String()
				valueTitle1 := gjson.Get(block1, "title")
				title1str := valueTitle1.String()
				fmt.Println(title1str, ":")

				if title1str == "General Capabilities" {
					// mRoot["children"] = block1

					mGC["code"] = gjson.Get(block1, "code").String()
					mGC["uuid"] = gjson.Get(block1, "uuid").String()
					mGC["type"] = gjson.Get(block1, "type").String()
					mGC["created_at"] = gjson.Get(block1, "created_at").String()
					mGC["title"] = gjson.Get(block1, "title").String()

					valueC2 := gjson.Get(block1, "children")
					if valueC2.IsArray() {
						for _, r2 := range valueC2.Array() {
							if r2.IsObject() {
								block2 := r2.String()
								title2str := gjson.Get(block2, "title").String()

								if ts.In(title2str, gcTitles...) {
									fmt.Println("  ", title2str)

									var m map[string]interface{}
									switch title2str {
									case "National Literacy Learning Progression":
										m = mNLLP
									case "National Numeracy Learning Progression":
										m = mNNLP
									case "Critical and Creative Thinking":
										m = mCCT
									case "Personal and Social Capability":
										m = mPSC
									case "Digital Literacy":
										m = mDL
									case "Intercultural Understanding":
										m = mIU
									case "Ethical Understanding":
										m = mEU
									}

									m["code"] = gjson.Get(block2, "code").String()
									m["uuid"] = gjson.Get(block2, "uuid").String()
									m["type"] = gjson.Get(block2, "type").String()
									m["created_at"] = gjson.Get(block2, "created_at").String()
									m["title"] = gjson.Get(block2, "title").String()
									m["children"] = gjson.Get(block2, "children").String()
								}
							}
						}
					}
				}
			}
		}
	}

	fmt.Println(mRoot["title"])
	fmt.Println(mGC["title"])
	fmt.Println(mIU["title"])
	fmt.Println(mEU["title"])

	for _, L2 := range mL2s {
		out := ""

		out, _ = sjson.Set(out, "code", mRoot["code"])
		out, _ = sjson.Set(out, "uuid", mRoot["uuid"])
		out, _ = sjson.Set(out, "type", mRoot["type"])
		out, _ = sjson.Set(out, "created_at", mRoot["created_at"])
		out, _ = sjson.Set(out, "title", mRoot["title"])
		out, _ = sjson.Set(out, "children.0.code", mGC["code"])
		out, _ = sjson.Set(out, "children.0.uuid", mGC["uuid"])
		out, _ = sjson.Set(out, "children.0.type", mGC["type"])
		out, _ = sjson.Set(out, "children.0.created_at", mGC["created_at"])
		out, _ = sjson.Set(out, "children.0.title", mGC["title"])
		out, _ = sjson.Set(out, "children.0.children.0.code", L2["code"])
		out, _ = sjson.Set(out, "children.0.children.0.uuid", L2["uuid"])
		out, _ = sjson.Set(out, "children.0.children.0.type", L2["type"])
		out, _ = sjson.Set(out, "children.0.children.0.created_at", L2["created_at"])
		out, _ = sjson.Set(out, "children.0.children.0.title", L2["title"])
		out, _ = sjson.SetRaw(out, "children.0.children.0.children", L2["children"].(string))

		// out = jt.FmtStr(out, "  ")
		err := os.WriteFile(fmt.Sprintf("./%s/gc-%s.json", outdir, L2["title"]), []byte(out), os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
}
