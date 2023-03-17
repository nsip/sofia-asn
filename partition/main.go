package main

import (
	"fmt"
	"os"

	lk "github.com/digisan/logkit"
)

var (
	fSf    = fmt.Sprintf
	uri4id = "http://vocabulary.curriculum.edu.au/" // "http://rdf.curriculum.edu.au/202110/"
)

func init() {
	lk.Log2F(true, "partition.log")
}

func main() {

	outDir := "out"
	os.MkdirAll(fSf("./%s/", outDir), os.ModePerm)

	//////////////////////////////////////////////////////////////

	bytesMeta, err := os.ReadFile("../data/Sofia-API-Meta-Data-09062022.json")
	lk.FailOnErr("%v", err)
	mMeta := parseMeta(string(bytesMeta))

	bytesNodes, err := os.ReadFile("../data/Sofia-API-Nodes-Data-09062022.json")
	lk.FailOnErr("%v", err)
	nodeProcess(bytesNodes, uri4id, mMeta, outDir)

	//////////////////////////////////////////////////////////////

	data, err := os.ReadFile("../data/Sofia-API-Tree-Data-09062022.json") // tree.pretty.json
	lk.FailOnErr("%v", err)
	js := string(data)

	fileContent := ccp(js, outDir)
	err = os.WriteFile(fSf("./%s/ccp-%s.json", outDir, "Cross-curriculum Priorities"), []byte(fileContent), os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	for gc, fileContent := range gc(js) {
		err = os.WriteFile(fSf("./%s/gc-%s.json", outDir, gc), []byte(fileContent), os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}

	for la, fileContent := range la(js) {
		err := os.WriteFile(fSf("./%s/la-%s.json", outDir, la), []byte(fileContent), os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}

	//////////////////////////////////////////////////////////////

	func() {
		uri4id := "http://vocabulary.curriculum.edu.au/MRAC/LA/ENG/"
		out := "./out/la-English.json"
		data, err := os.ReadFile(out)
		lk.WarnOnErr("%v", err)
		if err != nil {
			return
		}
		js := reStruct(string(data))
		js = ConnFieldMapping(js, uri4id, mMeta)
		if len(js) > 0 {
			os.WriteFile(out, []byte(js), os.ModePerm)
		}
	}()

	func() {
		uri4id := "http://vocabulary.curriculum.edu.au/MRAC/LA/HASS/"
		out := "./out/la-Humanities and Social Sciences.json" // Humanities and Social Sciences.json // HASS.json
		data, err := os.ReadFile(out)
		lk.WarnOnErr("%v", err)
		if err != nil {
			return
		}
		js := reStruct(string(data))
		js = ConnFieldMapping(js, uri4id, mMeta)
		if len(js) > 0 {
			os.WriteFile(out, []byte(js), os.ModePerm)
		}
	}()

	func() {
		uri4id := "http://vocabulary.curriculum.edu.au/MRAC/LA/HPE/"
		out := "./out/la-Health and Physical Education.json" // Health and Physical Education.json // HPE.json
		data, err := os.ReadFile(out)
		lk.WarnOnErr("%v", err)
		if err != nil {
			return
		}
		js := reStruct(string(data))
		js = ConnFieldMapping(js, uri4id, mMeta)
		if len(js) > 0 {
			os.WriteFile(out, []byte(js), os.ModePerm)
		}
	}()

	func() {
		uri4id := "http://vocabulary.curriculum.edu.au/MRAC/LA/LAN/"
		out := "./out/la-Languages.json"
		data, err := os.ReadFile(out)
		lk.WarnOnErr("%v", err)
		if err != nil {
			return
		}
		js := reStruct(string(data))
		js = ConnFieldMapping(js, uri4id, mMeta)
		if len(js) > 0 {
			os.WriteFile(out, []byte(js), os.ModePerm)
		}
	}()

	func() {
		uri4id := "http://vocabulary.curriculum.edu.au/MRAC/LA/MAT/"
		out := "./out/la-Mathematics.json"
		data, err := os.ReadFile(out)
		lk.WarnOnErr("%v", err)
		if err != nil {
			return
		}
		js := reStruct(string(data))
		js = ConnFieldMapping(js, uri4id, mMeta)
		if len(js) > 0 {
			os.WriteFile(out, []byte(js), os.ModePerm)
		}
	}()

	func() {
		uri4id := "http://vocabulary.curriculum.edu.au/MRAC/LA/SCI/"
		out := "./out/la-Science.json"
		data, err := os.ReadFile(out)
		lk.WarnOnErr("%v", err)
		if err != nil {
			return
		}
		js := reStruct(string(data))
		js = ConnFieldMapping(js, uri4id, mMeta)
		if len(js) > 0 {
			os.WriteFile(out, []byte(js), os.ModePerm)
		}
	}()

	func() {
		uri4id := "http://vocabulary.curriculum.edu.au/MRAC/LA/ART/"
		out := "./out/la-The Arts.json"
		data, err := os.ReadFile(out)
		lk.WarnOnErr("%v", err)
		if err != nil {
			return
		}
		js := reStruct(string(data))
		js = ConnFieldMapping(js, uri4id, mMeta)
		if len(js) > 0 {
			os.WriteFile(out, []byte(js), os.ModePerm)
		}
	}()

	func() {
		uri4id := "http://vocabulary.curriculum.edu.au/MRAC/LA/TEC/"
		out := "./out/la-Technologies.json"
		data, err := os.ReadFile(out)
		lk.WarnOnErr("%v", err)
		if err != nil {
			return
		}
		js := reStruct(string(data))
		js = ConnFieldMapping(js, uri4id, mMeta)
		if len(js) > 0 {
			os.WriteFile(out, []byte(js), os.ModePerm)
		}
	}()
}
