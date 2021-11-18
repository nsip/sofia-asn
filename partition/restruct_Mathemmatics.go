package main

import (
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// js must be Mathematics Learning Area - Achievement Standard json
func reStructMath(js string) string {

	var (
		mLACode2ASCode = map[string]string{
			"MATMATFY":  "ASMATFY",
			"MATMATY1":  "ASMATY1",
			"MATMATY2":  "ASMATY2",
			"MATMATY3":  "ASMATY3",
			"MATMATY4":  "ASMATY4",
			"MATMATY5":  "ASMATY5",
			"MATMATY6":  "ASMATY6",
			"MATMATY7":  "ASMATY7",
			"MATMATY8":  "ASMATY8",
			"MATMATY9":  "ASMATY9",
			"MATMATY10": "ASMATY10",
		}

		mLA2Path = map[string]string{
			"MATMATFY":  "",
			"MATMATY1":  "",
			"MATMATY2":  "",
			"MATMATY3":  "",
			"MATMATY4":  "",
			"MATMATY5":  "",
			"MATMATY6":  "",
			"MATMATY7":  "",
			"MATMATY8":  "",
			"MATMATY9":  "",
			"MATMATY10": "",
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

			// only deal with English ******************************************
			if code != "MAT" && code != "ASMAT" {
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

// ASMAT
// 	ASMATMAT
// 		ASMATFYL
// 			ASMATFY
// 		ASMATY1L
// 			ASMATY1
// 		ASMATY2L
// 			ASMATY2
// 		ASMATY3L
// 			ASMATY3
// 		ASMATY4L
// 			ASMATY4
// 		ASMATY5L
// 			ASMATY5
// 		ASMATY6L
// 			ASMATY6
// 		ASMATY7L
// 			ASMATY7
// 		ASMATY8L
// 			ASMATY8
// 		ASMATY9L
// 			ASMATY9
// 		ASMATY10L
// 			ASMATY10

// MAT
// 	MATMAT
// 		MATMATFY
// 			MATMATFYALG
// 			MATMATFYMEA
// 			MATMATFYNUM
// 			MATMATFYPRO
// 			MATMATFYSPA
// 			MATMATFYSTA
// 		MATMATY1
// 			MATMATY1ALG
// 			MATMATY1MEA
// 			MATMATY1NUM
// 			MATMATY1PRO
// 			MATMATY1SPA
// 			MATMATY1STA
// 		MATMATY2
// 			MATMATY2ALG
// 			MATMATY2MEA
// 			MATMATY2NUM
// 			MATMATY2PRO
// 			MATMATY2SPA
// 			MATMATY2STA
// 		MATMATY3
// 			MATMATY3ALG
// 			MATMATY3MEA
// 			MATMATY3NUM
// 			MATMATY3PRO
// 			MATMATY3SPA
// 			MATMATY3STA
// 		MATMATY4
// 			MATMATY4ALG
// 			MATMATY4MEA
// 			MATMATY4NUM
// 			MATMATY4PRO
// 			MATMATY4SPA
// 			MATMATY4STA
// 		MATMATY5
// 			MATMATY5ALG
// 			MATMATY5MEA
// 			MATMATY5NUM
// 			MATMATY5PRO
// 			MATMATY5SPA
// 			MATMATY5STA
// 		MATMATY6
// 			MATMATY6ALG
// 			MATMATY6MEA
// 			MATMATY6NUM
// 			MATMATY6PRO
// 			MATMATY6SPA
// 			MATMATY6STA
// 		MATMATY7
// 			MATMATY7ALG
// 			MATMATY7MEA
// 			MATMATY7NUM
// 			MATMATY7PRO
// 			MATMATY7SPA
// 			MATMATY7STA
// 		MATMATY8
// 			MATMATY8ALG
// 			MATMATY8MEA
// 			MATMATY8NUM
// 			MATMATY8PRO
// 			MATMATY8SPA
// 			MATMATY8STA
// 		MATMATY9
// 			MATMATY9ALG
// 			MATMATY9MEA
// 			MATMATY9NUM
// 			MATMATY9PRO
// 			MATMATY9SPA
// 			MATMATY9STA
// 		MATMATY10
// 			MATMATY10ALG
// 			MATMATY10MEA
// 			MATMATY10NUM
// 			MATMATY10PRO
// 			MATMATY10SPA
// 			MATMATY10STA
