package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	gio "github.com/digisan/gotk/io"
)

func LoadUrl(fpath string) map[string]string {
	m := make(map[string]string)
	gio.FileLineScan(fpath, func(line string) (bool, string) {
		ss := strings.Split(line, "\t")
		m[ss[0]] = ss[1]
		return true, ""
	}, "")
	return m
}

func main() {

	// rUUID := regexp.MustCompile(`[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`)

	mCodeUrl := LoadUrl("../data/code-url.txt")
	mIdUrl := LoadUrl("../data/id-url.txt")

	fmt.Println(len(mCodeUrl))
	fmt.Println(len(mIdUrl))

	///////////////////////////////////////

	r := regexp.MustCompile(`http://vocabulary.curriculum.edu.au/+[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`)

	///////////////////////////////////////

	dirIn := "../asn-json/out"
	dirOut := "../asn-json/out/url"
	gio.MustCreateDir(dirOut)

	fs, err := os.ReadDir(dirIn)
	if err != nil {
		panic(err)
	}
	for _, f := range fs {

		if f.IsDir() {
			continue
		}

		fpath := filepath.Join(dirIn, f.Name())
		// fmt.Println(fpath)

		data, err := os.ReadFile(fpath)
		if err != nil {
			panic(err)
		}
		str := string(data)

		fd := r.FindAllString(str, -1)
		fmt.Println(fpath, len(fd))

		str = r.ReplaceAllStringFunc(str, func(s string) string {
			id := s[len(s)-36:]
			if url, ok := mIdUrl[id]; ok {
				return url + id
			}
			return s
		})

		fpath = filepath.Join(dirOut, f.Name())
		if err := os.WriteFile(fpath, []byte(str), os.ModePerm); err != nil {
			panic(err)
		}
	}

	///////////////////////////////////////

	fmt.Println("-------------------------------------------------------")

	dirIn = "../asn-json-ld/out1"
	dirOut = "../asn-json-ld/out1/url"
	gio.MustCreateDir(dirOut)

	fs, err = os.ReadDir(dirIn)
	if err != nil {
		panic(err)
	}
	for _, f := range fs {

		if f.IsDir() {
			continue
		}

		fpath := filepath.Join(dirIn, f.Name())
		// fmt.Println(fpath)

		data, err := os.ReadFile(fpath)
		if err != nil {
			panic(err)
		}
		str := string(data)

		fd := r.FindAllString(str, -1)
		fmt.Println(fpath, len(fd))

		str = r.ReplaceAllStringFunc(str, func(s string) string {
			id := s[len(s)-36:]
			if url, ok := mIdUrl[id]; ok {
				return url + id
			}
			return s
		})

		fpath = filepath.Join(dirOut, f.Name())
		if err := os.WriteFile(fpath, []byte(str), os.ModePerm); err != nil {
			panic(err)
		}
	}
}
