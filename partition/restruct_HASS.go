package main

import (
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// js must be HASS Learning Area - Achievement Standard json
func reStructHASS(js string) string {

	var (
		mLACode2ASCode = map[string]string{
			"HASCIVY7":  "ASHASCIVY7",
			"HASCIVY8":  "ASHASCIVY8",
			"HASCIVY9":  "ASHASCIVY9",
			"HASCIVY10": "ASHASCIVY10",
			"HASECOY7":  "ASHASECOY7",
			"HASECOY8":  "ASHASECOY8",
			"HASECOY9":  "ASHASECOY9",
			"HASECOY10": "ASHASECOY10",
			"HASGEOY7":  "ASHASGEOY7",
			"HASGEOY8":  "ASHASGEOY8",
			"HASGEOY9":  "ASHASGEOY9",
			"HASGEOY10": "ASHASGEOY10",
			"HASHASFY":  "ASHASHASFY",
			"HASHASY1":  "ASHASHASY1",
			"HASHASY2":  "ASHASHASY2",
			"HASHASY3":  "ASHASHASY3",
			"HASHASY4":  "ASHASHASY4",
			"HASHASY5":  "ASHASHASY5",
			"HASHASY6":  "ASHASHASY6",
			"HASHISY7":  "ASHAHISY7",
			"HASHISY8":  "ASHAHISY8",
			"HASHISY9":  "ASHAHISY9",
			"HASHISY10": "ASHAHISY10",
		}

		mLA2Path = map[string]string{
			"HASCIVY7":  "",
			"HASCIVY8":  "",
			"HASCIVY9":  "",
			"HASCIVY10": "",
			"HASECOY7":  "",
			"HASECOY8":  "",
			"HASECOY9":  "",
			"HASECOY10": "",
			"HASGEOY7":  "",
			"HASGEOY8":  "",
			"HASGEOY9":  "",
			"HASGEOY10": "",
			"HASHASFY":  "",
			"HASHASY1":  "",
			"HASHASY2":  "",
			"HASHASY3":  "",
			"HASHASY4":  "",
			"HASHASY5":  "",
			"HASHASY6":  "",
			"HASHISY7":  "",
			"HASHISY8":  "",
			"HASHISY9":  "",
			"HASHISY10": "",
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

			// only deal with HASS ******************************************
			if code != "ASHAS" && code != "HAS" {
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
		// fmt.Println(laCode, path)
		// path += fmt.Sprintf(".%d", len(gjson.Get(js, path).Array())) // modify path, append to the last child
		content := mAS[mLACode2ASCode[laCode]]
		js, _ = sjson.SetRaw(js, path, content)
	}

	// remove AS part
	js, _ = sjson.Delete(js, "children.0")
	return js
}

// ASHAS
// 	ASHASCIV
// 		ASHASCIVY7L
// 			ASHASCIVY7
// 		ASHASCIVY8L
// 			ASHASCIVY8
// 		ASHASCIVY9L
// 			ASHASCIVY9
// 		ASHASCIVY10L
// 			ASHASCIVY10
// 	ASHASECO
// 		ASHASECOY7L
// 			ASHASECOY7
// 		ASHASECOY8L
// 			ASHASECOY8
// 		ASHASECOY9L
// 			ASHASECOY9
// 		ASHASECOY10L
// 			ASHASECOY10
// 	ASHASGEO
// 		ASHASGEOY7L
// 			ASHASGEOY7
// 		ASHASGEOY8L
// 			ASHASGEOY8
// 		ASHASGEOY9L
// 			ASHASGEOY9
// 		ASHASGEOY10L
// 			ASHASGEOY10
// 	ASHASHAS
// 		ASHASHASFYL
// 			ASHASHASFY
// 		ASHASHASY1L
// 			ASHASHASY1
// 		ASHASHASY2L
// 			ASHASHASY2
// 		ASHASHASY3L
// 			ASHASHASY3
// 		ASHASHASY4L
// 			ASHASHASY4
// 		ASHASHASY5L
// 			ASHASHASY5
// 		ASHASHASY6L
// 			ASHASHASY6
// 	ASHASHIS
// 		ASHAHISY7L
// 			ASHAHISY7
// 		ASHAHISY8L
// 			ASHAHISY8
// 		ASHAHISY9L
// 			ASHAHISY9
// 		ASHAHISY10L
// 			ASHAHISY10

// HAS
// 	HASCIV
// 		HASCIVY7
// 			HASCIVY7CI
// 			HASCIVY7CK
// 		HASCIVY8
// 			HASCIVY8CI
// 			HASCIVY8CK
// 		HASCIVY9
// 			HASCIVY9CI
// 			HASCIVY9CK
// 		HASCIVY10
// 			HASCIVY10CI
// 			HASCIVY10CK
// 	HASECO
// 		HASECOY7
// 			HASECOY7EI
// 			HASECOY7EK
// 		HASECOY8
// 			HASECOY8EI
// 			HASECOY8EK
// 		HASECOY9
// 			HASECOY9EI
// 			HASECOY9EK
// 		HASECOY10
// 			HASECOY10EI
// 			HASECOY10EK
// 	HASGEO
// 		HASGEOY7
// 			HASGEOY7GK
// 			HASGEOY7GS
// 		HASGEOY8
// 			HASGEOY8GK
// 			HASGEOY8GS
// 		HASGEOY9
// 			HASGEOY9GK
// 			HASGEOY9GS
// 		HASGEOY10
// 			HASGEOY10GK
// 			HASGEOY10GS
// 	HASHAS
// 		HASHASFY
// 			HASHASFYKNO
// 			HASHASFYSKI
// 		HASHASY1
// 			HASHASY1KNO
// 			HASHASY1SKI
// 		HASHASY2
// 			HASHASY2KNO
// 			HASHASY2SKI
// 		HASHASY3
// 			HASHASY3KNO
// 			HASHASY3SKI
// 		HASHASY4
// 			HASHASY4KNO
// 			HASHASY4SKI
// 		HASHASY5
// 			HASHASY5KNO
// 			HASHASY5SKI
// 		HASHASY6
// 			HASHASY6KNO
// 			HASHASY6SKI
// 	HASHIS
// 		HASHISY7
// 			HASHISY7HI
// 			HASHISY7HK
// 		HASHISY8
// 			HASHISY8HI
// 			HASHISY8HK
// 		HASHISY9
// 			HASHISY9HI
// 			HASHISY9HK
// 		HASHISY10
// 			HASHISY10HI
// 			HASHISY10HK
