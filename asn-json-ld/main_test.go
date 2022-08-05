package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/digisan/gotk/filedir"
	gio "github.com/digisan/gotk/io"
)

func TestMain(t *testing.T) {
	main()
}

// ********** //
func TestRmDupl(t *testing.T) {

	// each json must be formatted (VSCode)
	in := "./out"
	files, _, err := filedir.WalkFileDir(in, false)
	if err != nil {
		panic(err)
	}

	out1 := "./out1"
	os.MkdirAll(out1, os.ModePerm)

	for _, f := range files {
		var (
			prevline  = ""
			prev2line = ""
		)
		gio.FileLineScan(f, func(line string) (bool, string) {
			if line == prevline {
				if strings.Contains(line, `"dc:description"`) {
					return false, ""
				}
			}

			{
				ln := strings.TrimSpace(line)
				pln := strings.TrimSpace(prevline)
				p2ln := strings.TrimSpace(prev2line)
				if strings.HasPrefix(p2ln, `"literal":`) && pln == "}," && strings.HasPrefix(ln, `"dc:description":`) {
					return false, ""
				}
			}

			///
			if len(prevline) > 0 {
				prev2line = prevline
			}
			if len(line) > 0 {
				prevline = line
			}
			return true, line
		}, filepath.Join(out1, filepath.Base(f)))
	}
}

func TestAddCtx(t *testing.T) {

	os.MkdirAll("./out", os.ModePerm)

	data, err := os.ReadFile("../asn-json/out/la-English.json")
	if err != nil {
		panic(err)
	}
	js := string(data)
	js = addContext(js, context)
	js = replace(js)

	os.WriteFile("./out/test-ld.json", []byte(js), os.ModePerm)
}
