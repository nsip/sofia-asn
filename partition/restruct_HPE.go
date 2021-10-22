package main

import (
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// js must be HPE Learning Area - Achievement Standard json
func reStructHPE(js string) string {

	var (
		mLACode2ASCode = map[string]string{
			"HPEHPEFY":   "ASHPEFY",
			"HPEHPEY12":  "ASHPEY1Y2",
			"HPEHPEY34":  "ASHPEY34",
			"HPEHPEY56":  "ASHPEY56",
			"HPEHPEY78":  "ASHPEY78",
			"HPEHPEY910": "ASHPEY910",
		}

		mLA2Path = map[string]string{
			"HPEHPEFY":   "",
			"HPEHPEY12":  "",
			"HPEHPEY34":  "",
			"HPEHPEY56":  "",
			"HPEHPEY78":  "",
			"HPEHPEY910": "",
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
			if code != "HPE" && code != "ASHPE" {
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
						mLA2Path[code] = fSf("children.%d.children.%d.children.%d.children.%d.children", I, i, j, k)
					}
				}
			}
		}
	}

	for laCode, path := range mLA2Path {
		path += fmt.Sprintf(".%d", len(gjson.Get(js, path).Array())) // modify path, append to the last child
		content := mAS[mLACode2ASCode[laCode]]
		js, _ = sjson.SetRaw(js, path, content)
	}

	// remove AS part
	js, _ = sjson.Delete(js, "children.0")
	return js
}

// ASHPE
// 	ASHPEHPE
// 		ASHPEFYL
// 			ASHPEFY
// 		ASHPEY1Y2L
// 			ASHPEY1Y2
// 		ASHPEY34L
// 			ASHPEY34
// 		ASHPEY56L
// 			ASHPEY56
// 		ASHPEY78L
// 			ASHPEY78
// 		ASHPEY910L
// 			ASHPEY910

// HPE
// 	HPEHPE
// 		HPEHPEFY
// 			HPEHPEFYMOV
// 			HPEHPEFYPER
// 		HPEHPEY12
// 			HPEHPEY12MOV
// 			HPEHPEY12PER
// 		HPEHPEY34
// 			HPEHPEY34MOV
// 			HPEHPEY34PER
// 		HPEHPEY56
// 			HPEHPEY56MOV
// 			HPEHPEY56PER
// 		HPEHPEY78
// 			HPEHPEY78MOV
// 			HPEHPEY78PER
// 		HPEHPEY910
// 			HPEHPEY910MOV
// 			HPEHPEY910PER
