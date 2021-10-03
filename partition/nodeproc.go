package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func parseMeta(js string) map[string]string {
	mMeta := make(map[string]string)
	r := gjson.Get(js, "fields")
	if r.IsArray() {
		for _, ra := range r.Array() {
			if ra.IsObject() {
				mMeta[ra.Get("key").String()] = ra.Get("name").String()
			}
		}
	}
	return mMeta
}

func nodeProcess(bytes []byte, uri string, meta map[string]string, outdir string) {

	outdir = strings.Trim(outdir, `./\`)

	uri = strings.TrimSuffix(uri, "/")
	uri = strings.TrimSuffix(uri, "\\")
	uri += "/"

	// mData := make(map[string]interface{})
	// json.Unmarshal(bytes, &mData)

	js := string(bytes)
	// fmt.Sprint(js)

	r := regexp.MustCompile(`"[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}":`)
	// ids := r.FindAllString(js, -1)
	pGrp := r.FindAllStringIndex(js, -1)
	fmt.Println(len(pGrp), js[pGrp[0][0]:pGrp[0][1]])

	parts := []string{}

	out := ""
	for i := 0; i < len(pGrp); i++ {

		p, pn := pGrp[i], []int{}
		if i < len(pGrp)-1 {
			p, pn = pGrp[i], pGrp[i+1]
		}

		ids, ide := p[0]+1, p[1]-2
		id := js[ids:ide]
		// fmt.Println(id)

		blks, blke := 0, 0
		if i < len(pGrp)-1 {
			blks, blke = p[1], pn[0]-1
		} else {
			blks, blke = p[1], len(js)-1
		}

		block := js[blks:blke]
		block = strings.TrimSuffix(block, " ")
		block = strings.TrimSuffix(block, "\n")
		block = strings.TrimSuffix(block, ",")

		////////////////////////////////////////

		// "uuid": {id} => "id": "http://abc/def/{id}"
		newIdVal := fmt.Sprintf("%s%s", uri, gjson.Get(block, "uuid").String())
		block, _ = sjson.Set(block, "uuid", newIdVal)
		block = strings.Replace(block, `"uuid"`, `"id"`, 1)

		m := make(map[string]interface{})
		json.Unmarshal([]byte(gjson.Get(block, "connections").String()), &m)

		for k, v := range m {
			// "abcdeft" => "Levels" etc.
			block = strings.Replace(block, k, meta[k], 1)
			// "abc-def" => "http://abc/def/{id}"
			for _, a := range v.([]interface{}) {
				block = strings.Replace(block, a.(string), fmt.Sprintf("%s%s", uri, a), 1)
			}
		}

		part := fmt.Sprintf(`"%s": %s`, id, block)
		parts = append(parts, part)
	}

	out = "{" + strings.Join(parts, ",")

	/////////////////////////////////////////////////////////////////////////

	// I := 0

	// for id, mBlock := range mData {
	// 	// fmt.Sprintln(id)

	// 	bytes, _ := json.Marshal(mBlock)
	// 	block := string(bytes)
	// 	// block := gjson.Get(js, id).String()

	// 	// fmt.Println(block)

	// 	newIdVal := fmt.Sprintf("%s%s", uri, gjson.Get(block, "uuid").String())
	// 	block, _ = sjson.Set(block, "uuid", newIdVal)
	// 	block = strings.Replace(block, `"uuid"`, `"id"`, 1)

	// 	constr := gjson.Get(block, "connections").String()
	// 	m := make(map[string]interface{})
	// 	json.Unmarshal([]byte(constr), &m)

	// 	for k, v := range m {
	// 		// "abcdeft" => "Levels" etc.
	// 		block = strings.Replace(block, k, meta[k], 1)
	// 		// "abc-def" => "http://abc/def/{id}"
	// 		for _, a := range v.([]interface{}) {
	// 			block = strings.Replace(block, a.(string), fmt.Sprintf("%s%s", uri, a), 1)
	// 		}
	// 	}

	// 	out, _ = sjson.SetRaw(out, id, block)

	// 	// fmt.Println(out)

	// 	if I == 10000 {
	// 		break
	// 	}

	// 	I++
	// }

	os.WriteFile(fmt.Sprintf("./%s/node-meta.json", outdir), []byte(out), os.ModePerm)
}
