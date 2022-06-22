package main

import (
	"fmt"
	"strings"

	. "github.com/digisan/go-generics/v2"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// js must be Languages Learning Area - Achievement Standard json
func reStructLang(js string) string {

	var (
		//  mLACode2ASCode = map[string]string{
		// 	"LANCHISEC7-1Y78":  "",
		// 	"LANCHISEC7-1Y910": "",
		// 	"LANCHISECF-1FY":   "",
		// 	"LANCHISECF-1Y12":  "",
		// 	"LANCHISECF-1Y34":  "",
		// 	"LANCHISECF-1Y56":  "",
		// 	"LANCHISECF-1Y78":  "",
		// 	"LANCHISECF-1Y910": "",
		// 	"LANFRE7-1Y78":     "",
		// 	"LANFRE7-1Y910":    "",
		// 	"LANFREF-1FY":      "",
		// 	"LANFREF-1Y12":     "",
		// 	"LANFREF-1Y34":     "",
		// 	"LANFREF-1Y56":     "",
		// 	"LANFREF-1Y78":     "",
		// 	"LANFREF-1Y910":    "",
		// 	"LANITA7-1Y78":     "",
		// 	"LANITA7-1Y910":    "",
		// 	"LANITAF-1FY":      "",
		// 	"LANITAF-1Y12":     "",
		// 	"LANITAF-1Y34":     "",
		// 	"LANITAF-1Y56":     "",
		// 	"LANITAF-1Y78":     "",
		// 	"LANITAF-1Y910":    "",
		// 	"LANJAP7-1Y78":     "",
		// 	"LANJAP7-1Y910":    "",
		// 	"LANJAPF-1FY":      "",
		// 	"LANJAPF-1Y12":     "",
		// 	"LANJAPF-1Y34":     "",
		// 	"LANJAPF-1Y56":     "",
		// 	"LANJAPF-1Y78":     "",
		// 	"LANJAPF-1Y910":    "",
		// }

		mLA2Path = map[string]string{
			"LANCHISEC7-1Y78":  "",
			"LANCHISEC7-1Y910": "",
			"LANCHISECF-1FY":   "",
			"LANCHISECF-1Y12":  "",
			"LANCHISECF-1Y34":  "",
			"LANCHISECF-1Y56":  "",
			"LANCHISECF-1Y78":  "",
			"LANCHISECF-1Y910": "",
			"LANFRE7-1Y78":     "",
			"LANFRE7-1Y910":    "",
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
			"LANJAP7-1Y78":     "",
			"LANJAP7-1Y910":    "",
			"LANJAPF-1FY":      "",
			"LANJAPF-1Y12":     "",
			"LANJAPF-1Y34":     "",
			"LANJAPF-1Y56":     "",
			"LANJAPF-1Y78":     "",
			"LANJAPF-1Y910":    "",
		}

		mAS = make(map[string]string)
	)

	fSf := fmt.Sprintf

	for I := 0; I < 2; I++ {

		CODE1 := ""

		path := fSf("children.%d.code", I)
		code := gjson.Get(js, path).String() // AS LA
		if code == "" {
			break
		}
		fmt.Println(code)

		CODE1 = code

		for i := 0; i < 100; i++ {

			path := fSf("children.%d.children.%d.code", I, i)
			code := gjson.Get(js, path).String() //
			if code == "" {
				break
			}
			fmt.Printf("\t%s\n", code)

			// only deal with Languages ******************************************
			if code != "ASLAN" && code != "LAN" {
				continue
			}

			for j := 0; j < 100; j++ {

				CODE3 := ""

				path := fSf("children.%d.children.%d.children.%d.code", I, i, j)
				code := gjson.Get(js, path).String()
				if code == "" {
					break
				}
				fmt.Printf("\t\t%s\n", code)

				CODE3 = code

				for k := 0; k < 100; k++ {
					path := fSf("children.%d.children.%d.children.%d.children.%d.code", I, i, j, k)
					code := gjson.Get(js, path).String()
					if code == "" {
						break
					}
					fmt.Printf("\t\t\t%s - %s\n", code, path)

					for l := 0; l < 100; l++ {
						path := fSf("children.%d.children.%d.children.%d.children.%d.children.%d.code", I, i, j, k, l)
						code := gjson.Get(js, path).String()
						if code == "" {
							break
						}
						fmt.Printf("\t\t\t\t%s - %s\n", code, path)

						////////////////////////////////////////////
						// here "ASLANFRE", "ASLANITA", "ASLANJAP"

						// path = fSf("children.%d.children.%d.children.%d.children.%d.children.%d.children", I, i, j, k, l)
						// children := gjson.Get(js, path).Array()
						// fmt.Println("\t\t\t\tchildren count:", len(children))

						if strings.HasPrefix(code, "AS") && strings.HasSuffix(code, "L") {
							mAS[code] = gjson.Get(js, fSf("children.%d.children.%d.children.%d.children.%d.children.%d.children.0", I, i, j, k, l)).String()
						}
						if CODE1 == "LA" && In(CODE3, "LANFRE", "LANITA", "LANJAP") {
							mLA2Path[code] = fSf("children.%d.children.%d.children.%d.children.%d.children.%d.asn_hasLevel", I, i, j, k, l)
						}

						for m := 0; m < 100; m++ {
							path := fSf("children.%d.children.%d.children.%d.children.%d.children.%d.children.%d.code", I, i, j, k, l, m)
							code := gjson.Get(js, path).String()
							if code == "" {
								break
							}
							fmt.Printf("\t\t\t\t\t%s\n", code)

							// here "ASLANCHI"
							if strings.HasPrefix(code, "AS") && strings.HasSuffix(code, "L") {
								mAS[code] = gjson.Get(js, fSf("children.%d.children.%d.children.%d.children.%d.children.%d.children.%d.children.0", I, i, j, k, l, m)).String()
							}
							if CODE1 == "LA" && In(CODE3, "LANCHI") {
								mLA2Path[code] = fSf("children.%d.children.%d.children.%d.children.%d.children.%d.children.%d.asn_hasLevel", I, i, j, k, l, m)
							}
						}
					}
				}
			}
		}
	}

	for laCode, path := range mLA2Path {
		// fmt.Println(laCode, path)
		// path += fmt.Sprintf(".%d", len(gjson.Get(js, path).Array())) // modify path, append to the last child
		content := mAS["AS"+laCode+"L"]
		js, _ = sjson.SetRaw(js, path, content)
	}

	// remove AS part
	js, _ = sjson.Delete(js, "children.0")
	return js

	////////////////////////////////////////

	// las, _ := Map2KVs(mLA2Path, nil, nil)
	// ass, _ := Map2KVs(mAS, nil, nil)

	// for _, la := range las {
	// 	as := "AS" + la + "L"
	// 	if NotIn(as, ass...) {
	// 		panic("CODE ERROR 1")
	// 	}
	// }

	// for _, as := range ass {
	// 	la := strings.TrimPrefix(as, "AS")
	// 	la = strings.TrimSuffix(la, "L")
	// 	if NotIn(la, las...) {
	// 		panic("CODE ERROR 2")
	// 	}
	// }

	// return ""
}

func chkCodeLang(js string) bool {

	fnAsCode := func(LaCode string) string {
		return "AS" + LaCode + "L"
	}

	mLACode2ASCode := map[string]string{
		"LANCHISEC7-1Y78":  "",
		"LANCHISEC7-1Y910": "",
		"LANCHISECF-1FY":   "",
		"LANCHISECF-1Y12":  "",
		"LANCHISECF-1Y34":  "",
		"LANCHISECF-1Y56":  "",
		"LANCHISECF-1Y78":  "",
		"LANCHISECF-1Y910": "",
		"LANFRE7-1Y78":     "",
		"LANFRE7-1Y910":    "",
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
		"LANJAP7-1Y78":     "",
		"LANJAP7-1Y910":    "",
		"LANJAPF-1FY":      "",
		"LANJAPF-1Y12":     "",
		"LANJAPF-1Y34":     "",
		"LANJAPF-1Y56":     "",
		"LANJAPF-1Y78":     "",
		"LANJAPF-1Y910":    "",
	}

	for la := range mLACode2ASCode {
		{
			search := fmt.Sprintf(`"code": "%s",`, la)
			if !strings.Contains(js, search) {
				fmt.Println("la:", search)
				return false
			}
		}
		{
			as := fnAsCode(la)
			search := fmt.Sprintf(`"code": "%s",`, as)
			if !strings.Contains(js, search) {
				fmt.Println("as", search)
				return false
			}
		}
	}
	return true
}

// AS                                       LA
//   ASLAN									  LAN
//     ASLANCHI									LANCHI
//       ASLANCHISEC							  LANCHISEC
//         ASLANCHISEC7-1							LANCHISEC7-1
//           ASLANCHISEC7-1Y78L						  LANCHISEC7-1Y78
//           ASLANCHISEC7-1Y910L                      LANCHISEC7-1Y910
//         ASLANCHISECF-1                           LANCHISECF-1
//           ASLANCHISECF-1FYL                        LANCHISECF-1FY
//           ASLANCHISECF-1Y12L                       LANCHISECF-1Y12 *
//           ASLANCHISECF-1Y34L                       LANCHISECF-1Y34
//           ASLANCHISECF-1Y56L                       LANCHISECF-1Y56
//           ASLANCHISECF-1Y78L                       LANCHISECF-1Y78
//           ASLANCHISECF-1Y910L                      LANCHISECF-1Y910
//     ASLANFRE
//       ASLANFRE7-1
//         ASLANFRE7-1Y78L                          LANFRE7-1Y78
//         ASLANFRE7-1Y910L                         LANFRE7-1Y910
//       ASLANFREF-1
//         ASLANFREF-1FYL                           LANFREF-1FY
//         ASLANFREF-1Y12L                          LANFREF-1Y12
//         ASLANFREF-1Y34L                          LANFREF-1Y34
//         ASLANFREF-1Y56L                          LANFREF-1Y56
//         ASLANFREF-1Y78L                          LANFREF-1Y78
//         ASLANFREF-1Y910L                         LANFREF-1Y910
//     ASLANITA
//       ASLANITA7-1
//         ASLANITA7-1Y78L                          LANITA7-1Y78
//         ASLANITA7-1Y910L                         LANITA7-1Y910
//       ASLANITAF-1
//         ASLANITAF-1FYL                           LANITAF-1FY
//         ASLANITAF-1Y12L                          LANITAF-1Y12
//         ASLANITAF-1Y34L                          LANITAF-1Y34
//         ASLANITAF-1Y56L                          LANITAF-1Y56
//         ASLANITAF-1Y78L                          LANITAF-1Y78
//         ASLANITAF-1Y910L                         LANITAF-1Y910
//     ASLANJAP
//       ASLANJAP7-1
//         ASLANJAP7-1Y78L                          LANJAP7-1Y78
//         ASLANJAP7-1Y910L                         LANJAP7-1Y910
//       ASLANJAPF-1
//         ASLANJAPF-1FYL                           LANJAPF-1FY
//         ASLANJAPF-1Y12L                          LANJAPF-1Y12
//         ASLANJAPF-1Y34L                          LANJAPF-1Y34
//         ASLANJAPF-1Y56L                          LANJAPF-1Y56
//         ASLANJAPF-1Y78L                          LANJAPF-1Y78
//         ASLANJAPF-1Y910L                         LANJAPF-1Y910
