package main

import (
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// js must be The Arts Learning Area - Achievement Standard json
func reStructArt(js string) string {

	var (
		mLACode2ASCode = map[string]string{
			"ARTDANFY":   "",
			"ARTDANY12":  "ASARTDANY12",
			"ARTDANY34":  "ASARTDANY34",
			"ARTDANY56":  "ASARTDANY56",
			"ARTDANY78":  "ASARTDANY78",
			"ARTDANY910": "ASARTDANY910",
			"ARTDRAFY":   "",
			"ARTDRAY12":  "ASARTDRAY12",
			"ARTDRAY34":  "ASARTDRAY34",
			"ARTDRAY56":  "ASARTDRAY56",
			"ARTDRAY78":  "ASARTDRAY78",
			"ARTDRAY910": "ASARTDRAY910",
			"ARTMEDFY":   "",
			"ARTMEDY12":  "ASARTMEDY12",
			"ARTMEDY34":  "ASARTMEDY34",
			"ARTMEDY56":  "ASARTMEDY56",
			"ARTMEDY78":  "ASARTMEDY78",
			"ARTMEDY910": "ASARTMEDY910",
			"ARTMUSFY":   "",
			"ARTMUSY12":  "ASARTMUSY12",
			"ARTMUSY34":  "ASARTMUSY34",
			"ARTMUSY56":  "ASARTMUSY56",
			"ARTMUSY78":  "ASARTMUSY78",
			"ARTMUSY910": "ASARTMUSY910",
		}

		mLA2Path = map[string]string{
			"ARTDANFY":   "",
			"ARTDANY12":  "",
			"ARTDANY34":  "",
			"ARTDANY56":  "",
			"ARTDANY78":  "",
			"ARTDANY910": "",
			"ARTDRAFY":   "",
			"ARTDRAY12":  "",
			"ARTDRAY34":  "",
			"ARTDRAY56":  "",
			"ARTDRAY78":  "",
			"ARTDRAY910": "",
			"ARTMEDFY":   "",
			"ARTMEDY12":  "",
			"ARTMEDY34":  "",
			"ARTMEDY56":  "",
			"ARTMEDY78":  "",
			"ARTMEDY910": "",
			"ARTMUSFY":   "",
			"ARTMUSY12":  "",
			"ARTMUSY34":  "",
			"ARTMUSY56":  "",
			"ARTMUSY78":  "",
			"ARTMUSY910": "",
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

			// only deal with The Arts ******************************************
			if code != "ART" && code != "ASART" {
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
		if content != "" {
			js, _ = sjson.SetRaw(js, path, content)
		}
	}

	// remove AS part
	js, _ = sjson.Delete(js, "children.0")
	return js
}

// ASART
// 	ASARTDAN
// 		ASARTDANY12L
// 			ASARTDANY12
// 		ASARTDANY34L
// 			ASARTDANY34
// 		ASARTDANY56L
// 			ASARTDANY56
// 		ASARTDANY78L
// 			ASARTDANY78
// 		ASARTDANY910L
// 			ASARTDANY910
// 	ASARTDRA
// 		ASARTDRAY12L
// 			ASARTDRAY12
// 		ASARTDRAY34L
// 			ASARTDRAY34
// 		ASARTDRAY56L
// 			ASARTDRAY56
// 		ASARTDRAY78L
// 			ASARTDRAY78
// 		ASARTDRAY910L
// 			ASARTDRAY910
// 	ASARTFYL
// 		ASARTFY
// 	ASARTMED
// 		ASARTMEDY12L
// 			ASARTMEDY12
// 		ASARTMEDY34L
// 			ASARTMEDY34
// 		ASARTMEDY56L
// 			ASARTMEDY56
// 		ASARTMEDY78L
// 			ASARTMEDY78
// 		ASARTMEDY910L
// 			ASARTMEDY910
// 	ASARTMUS
// 		ASARTMUSY12L
// 			ASARTMUSY12
// 		ASARTMUSY34L
// 			ASARTMUSY34
// 		ASARTMUSY56L
// 			ASARTMUSY56
// 		ASARTMUSY78L
// 			ASARTMUSY78
// 		ASARTMUSY910L
// 			ASARTMUSY910
// 	ASARTVIS
// 		ASARTVISY12L
// 			ASARTVISY12
// 		ASARTVISY34L
// 			ASARTVISY34
// 		ASARTVISY56L
// 			ASARTVISY56
// 		ASARTVISY78L
// 			ASARTVISY78
// 		ASARTVISY910L
// 			ASARTVISY910
// 	ASARTY12L
// 		ASARTY12
// 	ASARTY34L
// 		ASARTY34
// 	ASRTY56L
// 		ASARTY56

// ART
// 	ARTDAN
// 		ARTDANFY
// 			ARTDANFYCRE
// 			ARTDANFYDEV
// 			ARTDANFYEXP
// 			ARTDANFYSHA
// 		ARTDANY12
// 			ARTDANY12CRE
// 			ARTDANY12DEV
// 			ARTDANY12EXP
// 			ARTDANY12SHA
// 		ARTDANY34
// 			ARTDANY34CRE
// 			ARTDANY34DEV
// 			ARTDANY34EXP
// 			ARTDANY34SHA
// 		ARTDANY56
// 			ARTDANY56CRE
// 			ARTDANY56DEV
// 			ARTDANY56EXP
// 			ARTDANY56SHA
// 		ARTDANY78
// 			ARTDANY78CRE
// 			ARTDANY78DEV
// 			ARTDANY78EXP
// 			ARTDANY78SHA
// 		ARTDANY910
// 			ARTDANY910CRE
// 			ARTDANY910DEV
// 			ARTDANY910EXP
// 			ARTDANY910SHA
// 	ARTDRA
// 		ARTDRAFY
// 			ARTDRAFYCRE
// 			ARTDRAFYDEV
// 			ARTDRAFYEXP
// 			ARTDRAFYSHA
// 		ARTDRAY12
// 			ARTDRAY12CRE
// 			ARTDRAY12DEV
// 			ARTDRAY12EXP
// 			ARTDRAY12SHA
// 		ARTDRAY34
// 			ARTDRAY34CRE
// 			ARTDRAY34DEV
// 			ARTDRAY34EXP
// 			ARTDRAY34SHA
// 		ARTDRAY56
// 			ARTDRAY56CRE
// 			ARTDRAY56DEV
// 			ARTDRAY56EXP
// 			ARTDRAY56SHA
// 		ARTDRAY78
// 			ARTDRAY78CRE
// 			ARTDRAY78DEV
// 			ARTDRAY78EXP
// 			ARTDRAY78SHA
// 		ARTDRAY910
// 			ARTDRAY910CRE
// 			ARTDRAY910DEV
// 			ARTDRAY910EXP
// 			ARTDRAY910SHA
// 	ARTMED
// 		ARTMEDFY
// 			ARTMEDFYCRE
// 			ARTMEDFYDEV
// 			ARTMEDFYEXP
// 			ARTMEDFYSHA
// 		ARTMEDY12
// 			ARTMEDY12CRE
// 			ARTMEDY12DEV
// 			ARTMEDY12EXP
// 			ARTMEDY12SHA
// 		ARTMEDY34
// 			ARTMEDY34CRE
// 			ARTMEDY34DEV
// 			ARTMEDY34EXP
// 			ARTMEDY34SHA
// 		ARTMEDY56
// 			ARTMEDY56CRE
// 			ARTMEDY56DEV
// 			ARTMEDY56EXP
// 			ARTMEDY56SHA
// 		ARTMEDY78
// 			ARTMEDY78CRE
// 			ARTMEDY78DEV
// 			ARTMEDY78EXP
// 			ARTMEDY78SHA
// 		ARTMEDY910
// 			ARTMEDY910CRE
// 			ARTMEDY910DEV
// 			ARTMEDY910EXP
// 			ARTMEDY910SHA
// 	ARTMUS
// 		ARTMUSFY
// 			ARTMUSFYCRE
// 			ARTMUSFYDEV
// 			ARTMUSFYEXP
// 			ARTMUSFYSHA
// 		ARTMUSY12
// 			ARTMUSY12CRE
// 			ARTMUSY12DEV
// 			ARTMUSY12EXP
// 			ARTMUSY12SHA
// 		ARTMUSY34
// 			ARTMUSY34CRE
// 			ARTMUSY34DEV
// 			ARTMUSY34EXP
// 			ARTMUSY34SHA
// 		ARTMUSY56
// 			ARTMUSY56CRE
// 			ARTMUSY56DEV
// 			ARTMUSY56EXP
// 			ARTMUSY56SHA
// 		ARTMUSY78
// 			ARTMUSY78CRE
// 			ARTMUSY78DEV
// 			ARTMUSY78EXP
// 			ARTMUSY78SHA
// 		ARTMUSY910
// 			ARTMUSY910CRE
// 			ARTMUSY910DEV
// 			ARTMUSY910EXP
// 			ARTMUSY910SHA
