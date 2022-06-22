package main

import (
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// js must be Science Learning Area - Achievement Standard json
func reStructSci(js string) string {

	var (
		mLACode2ASCode = map[string]string{
			"SCISCIFY":  "ASSCIFY",
			"SCISCIY1":  "ASSCIY1",
			"SCISCIY2":  "ASSCIY2",
			"SCISCIY3":  "ASSCIY3",
			"SCISCIY4":  "ASSCIY4",
			"SCISCIY5":  "ASSCIY5",
			"SCISCIY6":  "ASSCIY6",
			"SCISCIY7":  "ASSCIY7",
			"SCISCIY8":  "ASSCIY8",
			"SCISCIY9":  "ASSCIY9",
			"SCISCIY10": "ASSCIY10",
		}

		mLA2Path = map[string]string{
			"SCISCIFY":  "",
			"SCISCIY1":  "",
			"SCISCIY2":  "",
			"SCISCIY3":  "",
			"SCISCIY4":  "",
			"SCISCIY5":  "",
			"SCISCIY6":  "",
			"SCISCIY7":  "",
			"SCISCIY8":  "",
			"SCISCIY9":  "",
			"SCISCIY10": "",
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

			// only deal with Science ******************************************
			if code != "SCI" && code != "ASSCI" {
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

// ASSCI
// 	ASSCISCI
// 		ASSCIFYL
// 			ASSCIFY
// 		ASSCIY1L
// 			ASSCIY1
// 		ASSCIY2L
// 			ASSCIY2
// 		ASSCIY3L
// 			ASSCIY3
// 		ASSCIY4L
// 			ASSCIY4
// 		ASSCIY5L
// 			ASSCIY5
// 		ASSCIY6L
// 			ASSCIY6
// 		ASSCIY7L
// 			ASSCIY7
// 		ASSCIY8L
// 			ASSCIY8
// 		ASSCIY9L
// 			ASSCIY9
// 		ASSCIY10L
// 			ASSCIY10

// SCI
// 	SCISCI
// 		SCISCIFY
// 			SCISCIFYSHE
// 			SCISCIFYSI
// 			SCISCIFYSU
// 		SCISCIY1
// 			SCISCIY1SHE
// 			SCISCIY1SI
// 			SCISCIY1SU
// 		SCISCIY2
// 			SCISCIY2SHE
// 			SCISCIY2SI
// 			SCISCIY2SU
// 		SCISCIY3
// 			SCISCIY3SHE
// 			SCISCIY3SI
// 			SCISCIY3SU
// 		SCISCIY4
// 			SCISCIY4SHE
// 			SCISCIY4SI
// 			SCISCIY4SU
// 		SCISCIY5
// 			SCISCIY5SHE
// 			SCISCIY5SI
// 			SCISCIY5SU
// 		SCISCIY6
// 			SCISCIY6SHE
// 			SCISCIY6SI
// 			SCISCIY6SU
// 		SCISCIY7
// 			SCISCIY7SHE
// 			SCISCIY7SI
// 			SCISCIY7SU
// 		SCISCIY8
// 			SCISCIY8SHE
// 			SCISCIY8SI
// 			SCISCIY8SU
// 		SCISCIY9
// 			SCISCIY9SHE
// 			SCISCIY9SI
// 			SCISCIY9SU
// 		SCISCIY10
// 			SCISCIY10SHE
// 			SCISCIY10SI
// 			SCISCIY10SU
