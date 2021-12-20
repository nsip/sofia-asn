package main

import (
	"fmt"
	"log"
	"os"
)

var (
	fSf = fmt.Sprintf

	uri4id = "http://uat.vocabulary.curriculum.edu.au/" // "http://rdf.curriculum.edu.au/202110/"
)

func main() {

	outdir := "out"
	os.MkdirAll(fSf("./%s/", outdir), os.ModePerm)

	//////////////////////////////////////////////////////////////

	data, err := os.ReadFile("../data/node.pretty.json")
	if err != nil {
		panic(err)
	}

	bytesMeta, err := os.ReadFile("../data/metadata.pretty.json")
	if err != nil {
		panic(err)
	}
	mMeta := parseMeta(string(bytesMeta))
	nodeProcess(data, uri4id, mMeta, outdir)

	//////////////////////////////////////////////////////////////

	data, err = os.ReadFile("../data/tree.pretty.json")
	if err != nil {
		panic(err)
	}
	js := string(data)

	fileContent := ccp(js, outdir)
	err = os.WriteFile(fSf("./%s/ccp-%s.json", outdir, "Cross-curriculum Priorities"), []byte(fileContent), os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	for gc, fileContent := range gc(js) {
		err = os.WriteFile(fSf("./%s/gc-%s.json", outdir, gc), []byte(fileContent), os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}

	for la, fileContent := range la(js) {
		err := os.WriteFile(fSf("./%s/la-%s.json", outdir, la), []byte(fileContent), os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}

	//////////////////////////////////////////////////////////////

	func() {
		file := "./out/la-English.json"
		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatalln(err)
		}
		js := reStructEng(string(data))
		js = ConnFieldMapping(js, uri4id, mMeta)
		os.WriteFile(file, []byte(js), os.ModePerm)
	}()

	func() {
		file := "./out/la-HASS.json"
		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatalln(err)
		}
		js := reStructHASS(string(data))
		js = ConnFieldMapping(js, uri4id, mMeta)
		os.WriteFile(file, []byte(js), os.ModePerm)
	}()

	func() {
		file := "./out/la-HPE.json"
		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatalln(err)
		}
		js := reStructHPE(string(data))
		js = ConnFieldMapping(js, uri4id, mMeta)
		os.WriteFile(file, []byte(js), os.ModePerm)
	}()

	func() {
		file := "./out/la-Languages.json"
		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatalln(err)
		}
		js := reStructLang(string(data))
		js = ConnFieldMapping(js, uri4id, mMeta)
		os.WriteFile(file, []byte(js), os.ModePerm)
	}()

	func() {
		file := "./out/la-Mathematics.json"
		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatalln(err)
		}
		js := reStructMath(string(data))
		js = ConnFieldMapping(js, uri4id, mMeta)
		os.WriteFile(file, []byte(js), os.ModePerm)
	}()

	func() {
		file := "./out/la-Science.json"
		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatalln(err)
		}
		js := reStructSci(string(data))
		js = ConnFieldMapping(js, uri4id, mMeta)
		os.WriteFile(file, []byte(js), os.ModePerm)
	}()

	func() {
		file := "./out/la-Technologies.json"
		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatalln(err)
		}
		js := reStructTech(string(data))
		js = ConnFieldMapping(js, uri4id, mMeta)
		os.WriteFile(file, []byte(js), os.ModePerm)
	}()

	func() {
		file := "./out/la-The Arts.json"
		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatalln(err)
		}
		js := reStructArt(string(data))
		js = ConnFieldMapping(js, uri4id, mMeta)
		os.WriteFile(file, []byte(js), os.ModePerm)
	}()
}
