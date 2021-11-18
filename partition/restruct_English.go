package main

import (
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// js must be English Learning Area - Achievement Standard json
func reStructEng(js string) string {

	var (
		mLACode2ASCode = map[string]string{
			"ENGFY":     "ASENGFY",
			"ENGENGY1":  "ASENGY1",
			"ENGENGY2":  "ASENGY2",
			"ENGENGY3":  "ASENGY3",
			"ENGENGY4":  "ASENGY4",
			"ENGENGY5":  "ASENGY5",
			"ENGENGY6":  "ASENGY6",
			"ENGENGY7":  "ASENGY7",
			"ENGENGY8":  "ASENGY8",
			"ENGENGY9":  "ASENGY9",
			"ENGENGY10": "ASENGY10",
		}

		mLA2Path = map[string]string{
			"ENGFY":     "",
			"ENGENGY1":  "",
			"ENGENGY2":  "",
			"ENGENGY3":  "",
			"ENGENGY4":  "",
			"ENGENGY5":  "",
			"ENGENGY6":  "",
			"ENGENGY7":  "",
			"ENGENGY8":  "",
			"ENGENGY9":  "",
			"ENGENGY10": "",
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
		// fmt.Println(code) // AS

		for i := 0; i < 100; i++ {

			path := fSf("children.%d.children.%d.code", I, i)
			code := gjson.Get(js, path).String()
			if code == "" {
				break
			}
			// fmt.Printf("\t%s\n", code) // ASENG

			// only deal with English ******************************************
			if code != "ENG" && code != "ASENG" {
				continue
			}

			for j := 0; j < 100; j++ {
				path := fSf("children.%d.children.%d.children.%d.code", I, i, j)
				code := gjson.Get(js, path).String()
				if code == "" {
					break
				}
				// fmt.Printf("\t\t%s\n", code) // ASENGENG

				for k := 0; k < 100; k++ {
					path := fSf("children.%d.children.%d.children.%d.children.%d.code", I, i, j, k)
					code := gjson.Get(js, path).String()
					if code == "" {
						break
					}
					// fmt.Printf("\t\t\t%s\n", code) // ASENGFYL ASENGY1L ASENGY2L...

					for l := 0; l < 100; l++ {
						path := fSf("children.%d.children.%d.children.%d.children.%d.children.%d.code", I, i, j, k, l)
						code := gjson.Get(js, path).String()
						if code == "" {
							break
						}
						// fmt.Printf("\t\t\t\t%s\n", code) // ASENGFY ASENGY1 ASENGY2...

						// fetch content from AS
						if I == 0 {
							mAS[code] = gjson.Get(js, fSf("children.%d.children.%d.children.%d.children.%d.children.%d", I, i, j, k, l)).String()
						}
					}

					// fetch LA dest path
					if I == 1 {
						// mLA2Path[code] = fSf("children.%d.children.%d.children.%d.children.%d.children", I, i, j, k) // not to be one of children
						mLA2Path[code] = fSf("children.%d.children.%d.children.%d.children.%d.asn_hasLevel", I, i, j, k) // separate key as 'asn_hasLevel'
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

// ASENG
// 	ASENGENG
// 		ASENGFYL
// 			ASENGFY
// 		ASENGY1L
// 			ASENGY1
// 		ASENGY2L
// 			ASENGY2
// 		ASENGY3L
// 			ASENGY3
// 		ASENGY4L
// 			ASENGY4
// 		ASENGY5L
// 			ASENGY5
// 		ASENGY6L
// 			ASENGY6
// 		ASENGY7L
// 			ASENGY7
// 		ASENGY8L
// 			ASENGY8
// 		ASENGY9L
// 			ASENGY9
// 		ASENGY10L
// 			ASENGY10

// ENG
// 	ENGENG
// 		ENGFY
// 			ENGENGFYLIT
// 			ENGENGFYLITCY
// 			ENGFYLANG
// 		ENGENGY1
// 			ENGENGY1LANG
// 			ENGENGY1LIT
// 			ENGENGY1LITCY
// 		ENGENGY2
// 			ENGENGY2LANG
// 			ENGENGY2LIT
// 			ENGENGY2LITCY
// 		ENGENGY3
// 			ENGENGY3LANG
// 			ENGENGY3LIT
// 			ENGENGY3LITCY
// 		ENGENGY4
// 			ENGENGY4LANG
// 			ENGENGY4LIT
// 			ENGENGY4LITCY
// 		ENGENGY5
// 			ENGENGY5LANG
// 			ENGENGY5LIT
// 			ENGENGY5LITCY
// 		ENGENGY6
// 			ENGENGY6LANG
// 			ENGENGY6LIT
// 			ENGENGY6LITCY
// 		ENGENGY7
// 			ENGENGY7LANG
// 			ENGENGY7LIT
// 			ENGENGY7LITCY
// 		ENGENGY8
// 			ENGENGY8LANG
// 			ENGENGY8LIT
// 			ENGENGY8LITCY
// 		ENGENGY9
// 			ENGENGY9LANG
// 			ENGENGY9LIT
// 			ENGENGY9LITCY
// 		ENGENGY10
// 			ENGENGY10LANG
// 			ENGENGY10LIT
// 			ENGENGY10LITCY
