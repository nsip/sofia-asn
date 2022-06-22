package main

import (
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// js must be Technologies Learning Area - Achievement Standard json
func reStructTech(js string) string {

	var (
		mLACode2ASCode = map[string]string{
			"TECTDEFY":   "ASTECDESFY",
			"TECTDEY12":  "ASTECDESY12",
			"TECTDEY34":  "ASTECDESY34",
			"TECTDEY56":  "ASTECDESY56",
			"TECTDEY78":  "ASTECDESY78",
			"TECTDEY910": "ASTECDESY910",
			"TECTDIFY":   "ASTECDIGFY",
			"TECTDIY12":  "ASTECDIGY12",
			"TECTDIY34":  "ASTECDIGY34",
			"TECTDIY56":  "ASTECDIGY56",
			"TECTDIY78":  "ASTECDIGY78",
			"TECTDIY910": "ASTECDIGY910",
		}

		mLA2Path = map[string]string{
			"TECTDEFY":   "",
			"TECTDEY12":  "",
			"TECTDEY34":  "",
			"TECTDEY56":  "",
			"TECTDEY78":  "",
			"TECTDEY910": "",
			"TECTDIFY":   "",
			"TECTDIY12":  "",
			"TECTDIY34":  "",
			"TECTDIY56":  "",
			"TECTDIY78":  "",
			"TECTDIY910": "",
		}

		mAS = make(map[string]string)
	)

	fSf := fmt.Sprintf

	for I := 0; I < 2; I++ {

		path := fSf("children.%d.code", I)
		code := gjson.Get(js, path).String()
		if code == "" {
			break
		}
		// fmt.Println(code)

		for i := 0; i < 100; i++ {

			path := fSf("children.%d.children.%d.code", I, i)
			code := gjson.Get(js, path).String()
			if code == "" {
				break
			}
			// fmt.Printf("\t%s\n", code)

			// only deal with Tech ******************************************
			if code != "TEC" && code != "ASTEC" {
				continue
			}

			if code == "ASTEC" {
				fmt.Println("DEBUG")
			}

			for j := 0; j < 100; j++ {
				path := fSf("children.%d.children.%d.children.%d.code", I, i, j)
				code := gjson.Get(js, path).String()
				if code == "" {
					break
				}
				// fmt.Printf("\t\t%s\n", code)

				for k := 0; k < 100; k++ {
					path := fSf("children.%d.children.%d.children.%d.children.%d.code", I, i, j, k)
					code := gjson.Get(js, path).String()
					if code == "" {
						break
					}
					// fmt.Printf("\t\t\t%s\n", code)

					for l := 0; l < 100; l++ {
						path := fSf("children.%d.children.%d.children.%d.children.%d.children.%d.code", I, i, j, k, l)
						code := gjson.Get(js, path).String()
						if code == "" {
							break
						}
						// fmt.Printf("\t\t\t\t%s\n", code)

						// fetch content from AS
						if I == 0 {
							mAS[code] = gjson.Get(js, fSf("children.%d.children.%d.children.%d.children.%d.children.%d", I, i, j, k, l)).String()
						}
					}

					// fetch LA dest path
					if I == 1 {
						// mLA2Path[code] = fSf("children.%d.children.%d.children.%d.children.%d.children", I, i, j, k)
						mLA2Path[code] = fSf("children.%d.children.%d.children.%d.children.%d.asn_hasLevel", I, i, j, k)
					}
				}
			}
		}
	}

	for laCode, path := range mLA2Path {
		// path += fmt.Sprintf(".%d", len(gjson.Get(js, path).Array())) // modify path, append to the last child
		content := mAS[mLACode2ASCode[laCode]]
		js, _ = sjson.SetRaw(js, path, content)
	}

	// remove AS part
	js, _ = sjson.Delete(js, "children.0")
	return js
}
