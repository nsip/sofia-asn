package main

import (
	"fmt"
	"log"
	"strings"

	. "github.com/digisan/go-generics/v2"
	lk "github.com/digisan/logkit"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func init() {
	lk.WarnDetail(false)
}

// js must be The Arts Learning Area - Achievement Standard json
func reStructArt(js string) string {

	var (
		mLA2Path = map[string]string{}
		mAS      = make(map[string]string)
		asCodes  = []string{}
		laCodes  = []string{}
	)

	fSf := fmt.Sprintf

	for I := 0; I < 2; I++ {
		path := fSf("children.%d.code", I)
		code1 := gjson.Get(js, path).String()
		if code1 == "" {
			break
		}
		fmt.Println(code1)

		for i := 0; i < 100; i++ {
			path := fSf("children.%d.children.%d.code", I, i)
			code2 := gjson.Get(js, path).String()
			if code2 == "" {
				break
			}
			fmt.Printf("\t%s\n", code2)

			if code2 != "ART" && code2 != "ASART" {
				log.Fatalln(code2)
			}

			for j := 0; j < 100; j++ {
				path := fSf("children.%d.children.%d.children.%d.code", I, i, j)
				code3 := gjson.Get(js, path).String()
				if code3 == "" {
					break
				}
				// fmt.Printf("\t\t%s\n", code3)

				path = fSf("children.%d.children.%d.children.%d.doc.typeName", I, i, j)
				if typeName := gjson.Get(js, path).String(); typeName == "Level" {
					fmt.Printf("\t\t%s - ok\n", code3)

					if code1 == "AS" {
						asCodes = append(asCodes, code3)
						mAS[code3] = gjson.Get(js, fSf("children.%d.children.%d.children.%d.children.0", I, i, j)).String()
					}
					if code1 == "LA" {
						laCodes = append(laCodes, code3)
						mLA2Path[code3] = fSf("children.%d.children.%d.children.%d.asn_hasLevel", I, i, j)
					}
				}

				for k := 0; k < 100; k++ {
					path := fSf("children.%d.children.%d.children.%d.children.%d.code", I, i, j, k)
					code4 := gjson.Get(js, path).String()
					if code4 == "" {
						break
					}
					// fmt.Printf("\t\t\t%s", code4)

					path = fSf("children.%d.children.%d.children.%d.children.%d.doc.typeName", I, i, j, k)
					if typeName := gjson.Get(js, path).String(); typeName == "Level" {
						fmt.Printf("\t\t\t%s - ok\n", code4)

						if code1 == "AS" {
							asCodes = append(asCodes, code4)
							mAS[code4] = gjson.Get(js, fSf("children.%d.children.%d.children.%d.children.%d.children.0", I, i, j, k)).String()
						}
						if code1 == "LA" {
							laCodes = append(laCodes, code4)
							mLA2Path[code4] = fSf("children.%d.children.%d.children.%d.children.%d.asn_hasLevel", I, i, j, k)
						}
					}
				}
			}
		}
	}

	for _, as := range asCodes {
		la := strings.TrimPrefix(as, "AS")
		la = strings.TrimSuffix(la, "L")
		if NotIn(la, laCodes...) {
			lk.Warn("AS has [%s], BUT LA has no [%s]", as, la)
		}
	}

	for _, la := range laCodes {
		as := "AS" + la + "L"
		if NotIn(as, asCodes...) {
			lk.Warn("LA has [%s], BUT AS has no [%s]", la, as)
		}
	}

	for laCode, path := range mLA2Path {
		// path += fmt.Sprintf(".%d", len(gjson.Get(js, path).Array())) // modify path, append to the last child
		if content, ok := mAS["AS"+laCode+"L"]; ok {
			if content != "" {
				js, _ = sjson.SetRaw(js, path, content)
			}
		} else {
			js, _ = sjson.Set(js, path, "")
		}
	}

	// remove AS part
	js, _ = sjson.Delete(js, "children.0")
	return js
}
