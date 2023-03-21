package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	. "github.com/digisan/go-generics/v2"
	fd "github.com/digisan/gotk/file-dir"
	"github.com/digisan/gotk/strs"
)

func main() {

	var (
		todo = [5]bool{true, true, true, true, true}
	)

	inDir := "../out/json-ld/MRAC/2022/06/GC/CCP"  //"../out/url-ld/"
	outDir := "./out/json-ld/GC/CCP" // "./out/url-ld/"

	fd.MustCreateDir(outDir)

	files, err := os.ReadDir(inDir)
	if err != nil {
		log.Fatalln(err)
		return
	}

	for _, file := range files {

		fName := file.Name()

		fd.FileLineScanEx(
			filepath.Join(inDir, fName),
			3, 3, "", func(line string, cache []string) (bool, string) {

				// *** 1. date-time format ***
				//
				if todo[0] && strings.Contains(line, `"@value"`) {
					ln := strings.TrimSpace(line)
					above1, _, _ := strings.TrimSpace(cache[2]), strings.TrimSpace(cache[1]), strings.TrimSpace(cache[0])
					below1, _, _ := strings.TrimSpace(cache[4]), strings.TrimSpace(cache[5]), strings.TrimSpace(cache[6])
					if In(`"@type": "xsd:dateTime"`, above1, below1) {
						val := strings.TrimPrefix(ln, `"@value": `)
						val = strings.TrimSuffix(val, `,`)
						val = strings.Trim(val, `"`)
						dt, ok := TryToDateTime(val)
						if !ok {
							log.Println("cannot be date-time")
							os.Exit(-1)
						}
						return true, fmt.Sprintf(`%s"@value": "%s",`, strs.TrimTailFromLast(line, `"@value"`), dt.Format(`2006-01-02T15:04:05.000Z`))
					}
				}

				// *** 2. prefLabel on curriculum statements as concepts ***
				//
				if todo[1] && strings.Contains(line, `"asn:statementNotation"`) {
					ln := strings.TrimSpace(line)
					// above1, _, _ := strings.TrimSpace(cache[2]), strings.TrimSpace(cache[1]), strings.TrimSpace(cache[0])
					// below1, _, _ := strings.TrimSpace(cache[4]), strings.TrimSpace(cache[5]), strings.TrimSpace(cache[6])

					head := strs.TrimTailFromLast(line, `"asn:statementNotation"`)

					valStr := strings.TrimPrefix(ln, `"asn:statementNotation"`)
					valStr = strings.TrimPrefix(valStr, ":")
					valStr = strings.TrimSpace(valStr)
					valStr = strings.TrimSuffix(valStr, ",")
					valStr = strings.Trim(valStr, `"`)

					newFieldVal := fmt.Sprintf(`"skos:prefLabel": "%s",`, valStr)
					return true, line + "\n" + head + newFieldVal
				}

				// *** 3. skos shadowing of hasChild/isChildOf *** (NOT IMPLEMENT YET)
				//
				if todo[2] && strings.Contains(line, `"gem:hasChild"`) {
					// skos:narrower
				}
				if todo[2] && strings.Contains(line, `"gem:isChildOf"`) {
					// skos:broader
				}

				// *** 4. *** (No issue found as 'Remove the "connects" structure from JSON and JSON-LD output')
				//
				if todo[3] {

				}

				// *** 5. ***
				if todo[4] && strings.Contains(line, `"language"`) {
					return true, strings.Replace(line, `"language"`, `"@language"`, 1)
				}
				if todo[4] && strings.Contains(line, `"literal"`) {
					return true, strings.Replace(line, `"literal"`, `"@value"`, 1)
				}
				if todo[4] && strings.Contains(line, `"position"`) {
					return true, strings.Replace(line, `"position"`, `"@asn:listID"`, 1)
				}
				//

				return true, line
			},
			filepath.Join(outDir, fName),
		)
	}
}
