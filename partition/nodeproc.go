package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/nsip/sofia-asn/tool"
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

func nodeProcess(data []byte, uri string, meta map[string]string, outdir string) {

	e := bytes.LastIndexAny(data, "}")
	data = data[:e+1]

	outdir = strings.Trim(outdir, `./\`)
	parts := []string{}
	out := ""

	tool.ScanNode(data, func(i int, id, block string) bool {

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
		return true
	})

	out = "{" + strings.Join(parts, ",") + "}"
	// out = jt.FmtStr(out, "  ")
	os.WriteFile(fmt.Sprintf("./%s/node-meta.json", outdir), []byte(out), os.ModePerm)
}
