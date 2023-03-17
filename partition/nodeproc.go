package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	// . "github.com/digisan/go-generics/v2"
	dt "github.com/digisan/gotk/data-type"
	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
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

	out = "{" + strings.Join(parts, "") + "}"
	// out = jt.FmtStr(out, "  ")

	lk.FailOnErrWhen(!dt.IsJSON([]byte(out)), "%v", errors.New("Invalid JSON from node & meta"))

	os.WriteFile(fmt.Sprintf("./%s/node-meta.json", outdir), []byte(out), os.ModePerm)
}

//////////////////////////////////////////////////////////////////////

func MarkUrl(ids, codes []string, mCodeUrl, mIdUrl map[string]string) {

	// for _, code := range codes {
	// 	if url, ok := mCodeUrl[code]; ok {
	// 		for i, code := range codes {
	// 			if NotIn(code, "root", "LA", "AS", "GC", "CCP") {
	// 				mCodeUrl[code] = url
	// 				mIdUrl[ids[i]] = url
	// 			}
	// 		}
	// 		return
	// 	}
	// }

	url := ""
	for i, code := range codes {
		if i == len(codes)-3 {
			switch code {
			case "HAS", "HASS", "ASHAS", "ASHASS":
				url = "http://vocabulary.curriculum.edu.au/MRAC/LA/HASS/"
			case "ENG", "ASENG":
				url = "http://vocabulary.curriculum.edu.au/MRAC/LA/ENG/"
			case "LAN", "ASLAN":
				url = "http://vocabulary.curriculum.edu.au/MRAC/LA/LAN/"
			case "SCI", "ASSCI":
				url = "http://vocabulary.curriculum.edu.au/MRAC/LA/SCI/"
			case "ART", "ASART":
				url = "http://vocabulary.curriculum.edu.au/MRAC/LA/ART/"
			case "HPE", "ASHPE":
				url = "http://vocabulary.curriculum.edu.au/MRAC/LA/HPE/"
			case "MAT", "ASMAT":
				url = "http://vocabulary.curriculum.edu.au/MRAC/LA/MAT/"
			case "TEC", "ASTEC":
				url = "http://vocabulary.curriculum.edu.au/MRAC/LA/TEC/"

			case "CCT":
				url = "http://vocabulary.curriculum.edu.au/MRAC/GC/CCT/"
			case "N":
				url = "http://vocabulary.curriculum.edu.au/MRAC/GC/N/"
			case "DL":
				url = "http://vocabulary.curriculum.edu.au/MRAC/GC/DL/"
			case "L":
				url = "http://vocabulary.curriculum.edu.au/MRAC/GC/L/"
			case "PSC":
				url = "http://vocabulary.curriculum.edu.au/MRAC/GC/PSC/"
			case "IU":
				url = "http://vocabulary.curriculum.edu.au/MRAC/GC/IU/"
			case "EU":
				url = "http://vocabulary.curriculum.edu.au/MRAC/GC/EU/"

			case "AA":
				url = "http://vocabulary.curriculum.edu.au/MRAC/CCP/AA/"
			case "S":
				url = "http://vocabulary.curriculum.edu.au/MRAC/CCP/S/"
			case "A_TSI":
				url = "http://vocabulary.curriculum.edu.au/MRAC/CCP/A_TSI/"
			}
			break
		}
		if i == len(codes)-2 {
			switch code {
			case "AS", "LA":
				url = "http://vocabulary.curriculum.edu.au/MRAC/LA/"
			case "GC":
				url = "http://vocabulary.curriculum.edu.au/MRAC/GC/"
			case "CCP":
				url = "http://vocabulary.curriculum.edu.au/MRAC/CCP/"
			}
			mCodeUrl[code] = url
			mIdUrl[ids[i]] = url
			break
		}
		if i == len(codes)-1 {
			switch code {
			case "root":
				url = "http://vocabulary.curriculum.edu.au/MRAC/"
			}
			mCodeUrl[code] = url
			mIdUrl[ids[i]] = url
			break
		}
	}

	if len(codes) > 2 && url == "" {
		panic("Need Code: " + strings.Join(codes, ","))
	}

	for i, code := range codes {
		if i < len(codes)-2 {
			if _, ok := mCodeUrl[code]; !ok {
				mCodeUrl[code] = url
			}
			if _, ok := mIdUrl[ids[i]]; !ok {
				mIdUrl[ids[i]] = url
			}
		}
	}
}

func TrackCode(ms map[string]string, code string) (codes, ids []string) {
	ID := ""
	for id, valstr := range ms {
		if gjson.Get(valstr, "code").String() == code {
			ID = id
			break
		}
	}
	if ID != "" {
		for _, id := range TrackId(ms, ID) {
			valstr := ms[id]
			codes = append(codes, gjson.Get(valstr, "code").String())
			ids = append(ids, id)
		}
	}
	return // Reverse(codes), Reverse(ids)
}

func TrackId(ms map[string]string, id string) (ids []string) {
	ids = append(ids, id)
	for parent := IsChild(ms, id); len(parent) > 0; parent = IsChild(ms, parent) {
		ids = append(ids, parent)
	}
	return // Reverse(ids)
}

func IsChild(ms map[string]string, childId string) string {
	for id := range ms {
		if HasChild(ms, id, childId) {
			return id
		}
	}
	return ""
}

func HasChild(ms map[string]string, id, childId string) bool {
	valstr := ms[id]
	if children := gjson.Get(valstr, "children").Array(); len(children) > 0 {
		for _, child := range children {
			// fmt.Println(child)
			if childId == child.String() {
				return true
			}
		}
	}
	return false
}

func Scan2map(data []byte) map[string]any {
	M := make(map[string]any)
	if err := json.Unmarshal(data, &M); err != nil {
		panic(err)
	}
	return M
}

func Scan2mapstrval(data []byte) map[string]string {
	M := Scan2map(data)
	ret := make(map[string]string)
	for k, v := range M {
		vdata, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		ret[k] = string(vdata)
	}
	return ret
}

func Scan2flatmap(data []byte) map[string]map[string]any {
	M := Scan2map(data)
	ret := make(map[string]map[string]any)
	for k, v := range M {
		vdata, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		mf, err := jt.Flatten(vdata)
		if err != nil {
			panic(err)
		}
		ret[k] = mf
	}
	return ret
}
