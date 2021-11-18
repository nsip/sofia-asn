package main

import (
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// js must be Languages Learning Area - Achievement Standard json
func reStructLang(js string) string {

	var (
		mLACode2ASCode = map[string]string{
			"LANCHISEC7-1":     "ASLANCHI7_10SL",
			"LANCHISECF-1":     "ASLANCHISLF10",
			"LANFRESEC7-1Y78":  "ALANFRE7_10Y78L",
			"LANFRESEC7-1Y910": "ALANFRE7_10Y910L",
			"LANFREF-1FY":      "ASLANFREF10FYL",
			"LANFREF-1Y12":     "ASLANFREF10Y12L",
			"LANFREF-1Y34":     "ASLANFREF10Y34L",
			"LANFREF-1Y56":     "ASLANFREF10Y56L",
			"LANFREF-1Y78":     "ASLANFREF10Y178l",
			"LANFREF-1Y910":    "ASLANFREF10Y910L",
			"LANITA7-1Y78":     "ASLANITA7_10Y78L",
			"LANITA7-1Y910":    "ASLANITA7_10Y910L",
			"LANITAF-1FY":      "ASLANITAF10FYL",
			"LANITAF-1Y12":     "ASLANITAF10Y12L",
			"LANITAF-1Y34":     "ASLANITAF10Y34L",
			"LANITAF-1Y56":     "ASLANITAF10Y56L",
			"LANITAF-1Y78":     "ASLANITAF10Y78L",
			"LANITAF-1Y910":    "ASLANITAF10Y910L",
			"LANJAPSEC7-1Y78":  "ASLANJAP7_10Y78L",
			"LANJAPSEC7-1Y910": "ASLANJAP7_10Y910L",
			"LANJAPF-1FY":      "ASLANJAPF10FYL",
			"LANJAPF-1Y12":     "ASLANJAPF10Y12L",
			"LANJAPF-1Y34":     "ASLANJAPF10Y34L",
			"LANJAPF-1Y56":     "ASLANJAPF10Y56L",
			"LANJAPSECF-1Y78":  "ASLANJAPF10Y78L",
			"LANJAPSECF-1Y910": "ASLANJAPF10Y910L",
		}

		mLA2Path = map[string]string{
			"LANCHISEC7-1":     "",
			"LANCHISECF-1":     "",
			"LANFRESEC7-1Y78":  "",
			"LANFRESEC7-1Y910": "",
			"LANFREF-1FY":      "",
			"LANFREF-1Y12":     "",
			"LANFREF-1Y34":     "",
			"LANFREF-1Y56":     "",
			"LANFREF-1Y78":     "",
			"LANFREF-1Y910":    "",
			"LANITA7-1Y78":     "",
			"LANITA7-1Y910":    "",
			"LANITAF-1FY":      "",
			"LANITAF-1Y12":     "",
			"LANITAF-1Y34":     "",
			"LANITAF-1Y56":     "",
			"LANITAF-1Y78":     "",
			"LANITAF-1Y910":    "",
			"LANJAPSEC7-1Y78":  "",
			"LANJAPSEC7-1Y910": "",
			"LANJAPF-1FY":      "",
			"LANJAPF-1Y12":     "",
			"LANJAPF-1Y34":     "",
			"LANJAPF-1Y56":     "",
			"LANJAPSECF-1Y78":  "",
			"LANJAPSECF-1Y910": "",
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

			// only deal with Languages ******************************************
			if code != "ASLAN" && code != "LAN" {
				continue
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
					// fmt.Printf("\t\t\t%s - %s\n", code, path)

					for l := 0; l < 100; l++ {
						path := fSf("children.%d.children.%d.children.%d.children.%d.children.%d.code", I, i, j, k, l)
						code := gjson.Get(js, path).String()
						if code == "" {
							break
						}
						// fmt.Printf("\t\t\t\t%s - %s\n", code, path)

						// for m := 0; m < 100; m++ {
						// 	path := fSf("children.%d.children.%d.children.%d.children.%d.children.%d.children.%d.code", I, i, j, k, l, m)
						// 	code := gjson.Get(js, path).String()
						// 	if code == "" {
						// 		break
						// 	}
						// 	fmt.Printf("\t\t\t\t\t%s\n", code)
						// }

						// fetch content from AS
						if I == 0 {
							mAS[code] = gjson.Get(js, fSf("children.%d.children.%d.children.%d.children.%d.children.%d", I, i, j, k, l)).String()
						}

						// fetch LA dest path, 'languages' level is different from other Learning Area
						if I == 1 {
							// fmt.Println("code:", code)
							// mLA2Path[code] = fSf("children.%d.children.%d.children.%d.children.%d.children.%d.children", I, i, j, k, l)
							mLA2Path[code] = fSf("children.%d.children.%d.children.%d.children.%d.children.%d.asn_hasLevel", I, i, j, k, l)
						}
					}

					// fetch LA dest path
					// if I == 1 {
					// 	fmt.Println("code:", code)
					// 	mLA2Path[code] = fSf("children.%d.children.%d.children.%d.children.%d.children", I, i, j, k)
					// }
				}
			}
		}
	}

	for laCode, path := range mLA2Path {
		// fmt.Println(laCode, path)
		// path += fmt.Sprintf(".%d", len(gjson.Get(js, path).Array())) // modify path, append to the last child
		content := mAS[mLACode2ASCode[laCode]]
		js, _ = sjson.SetRaw(js, path, content)
	}

	// remove AS part
	js, _ = sjson.Delete(js, "children.0")
	return js
}

// ASLAN
// 	ASLANCHI
// 		ASLANCHISL
// 			ASLANCHI7_10SL
// 			ASLANCHISLF10
// 	ASLANFRE
// 		ASLANFRE7_10
// 			ALANFRE7_10Y78L
// 			ALANFRE7_10Y910L
// 		ASLANFREF10
// 			ASLANFREF10FYL
// 			ASLANFREF10Y12L
// 			ASLANFREF10Y34L
// 			ASLANFREF10Y56L
// 			ASLANFREF10Y178l
// 			ASLANFREF10Y910L
// 	ASLANITA
// 		ASLANITA7_10
// 			ASLANITA7_10Y78L
// 			ASLANITA7_10Y910L
// 		ASLANITAF10
// 			ASLANITAF10FYL
// 			ASLANITAF10Y12L
// 			ASLANITAF10Y34L
// 			ASLANITAF10Y56L
// 			ASLANITAF10Y78L
// 			ASLANITAF10Y910L
// 	ASLANJAP
// 		ASJAP7_10
// 			ASLANJAP7_10Y78L
// 			ASLANJAP7_10Y910L
// 		ASLANJAPF10
// 			ASLANJAPF10FYL
// 			ASLANJAPF10Y12L
// 			ASLANJAPF10Y34L
// 			ASLANJAPF10Y56L
// 			ASLANJAPF10Y78L
// 			ASLANJAPF10Y910L

// LAN
// 	LANCHI
// 		LANCHISEC
// 			LANCHISEC7-1
// 			LANCHISECF-1
// 	LANFRE
// 		LANFREF-1
// 			LANFREF-1FY
// 			LANFREF-1Y12
// 			LANFREF-1Y34
// 			LANFREF-1Y56
// 			LANFREF-1Y78
// 			LANFREF-1Y910
// 		LANFRESEC7-1
// 			LANFRESEC7-1Y78
// 			LANFRESEC7-1Y910
// 	LANITA
// 		LANITA7-1
// 			LANITA7-1Y78
// 			LANITA7-1Y910
// 		LANITAF-1
// 			LANITAF-1FY
// 			LANITAF-1Y12
// 			LANITAF-1Y34
// 			LANITAF-1Y56
// 			LANITAF-1Y78
// 			LANITAF-1Y910
// 	LANJAP
// 		LANJAPF-1
// 			LANJAPF-1FY
// 			LANJAPF-1Y12
// 			LANJAPF-1Y34
// 			LANJAPF-1Y56
// 			LANJAPSECF-1Y78
// 			LANJAPSECF-1Y910
// 		LANJAPSEC7-1
// 			LANJAPSEC7-1Y78
// 			LANJAPSEC7-1Y910
