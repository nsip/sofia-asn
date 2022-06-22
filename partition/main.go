package main

import (
	"fmt"
	"os"

	lk "github.com/digisan/logkit"
)

var (
	fSf    = fmt.Sprintf
	uri4id = "http://uat.vocabulary.curriculum.edu.au/" // "http://rdf.curriculum.edu.au/202110/"
)

func main() {

	outdir := "out"
	os.MkdirAll(fSf("./%s/", outdir), os.ModePerm)

	//////////////////////////////////////////////////////////////

	bytesMeta, err := os.ReadFile("../data/Sofia-API-Meta-Data-09062022.json")
	lk.FailOnErr("%v", err)
	mMeta := parseMeta(string(bytesMeta))

	bytesNodes, err := os.ReadFile("../data/Sofia-API-Nodes-Data-09062022.json")
	lk.FailOnErr("%v", err)
	nodeProcess(bytesNodes, uri4id, mMeta, outdir)

	//////////////////////////////////////////////////////////////

	data, err := os.ReadFile("../data/Sofia-API-Tree-Data-09062022.json") // tree.pretty.json
	lk.FailOnErr("%v", err)
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
		in := "./out/la-English.json"
		// out := "./out/la-English-out.json"
		// if fd.FileExists(out) {
		// 	return
		// }
		out := in

		data, err := os.ReadFile(in)
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
		in := "./out/la-Humanities and Social Sciences.json" // Humanities and Social Sciences.json // HASS.json
		// out := "./out/la-Humanities and Social Sciences-out.json"
		// if fd.FileExists(out) {
		// 	return
		// }
		out := in

		data, err := os.ReadFile(in)
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
		in := "./out/la-Health and Physical Education.json" // Health and Physical Education.json // HPE.json
		// out := "./out/la-Health and Physical Education-out.json"
		// if fd.FileExists(out) {
		// 	return
		// }
		out := in

		data, err := os.ReadFile(in)
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
		in := "./out/la-Languages.json"
		// out := "./out/la-Languages-out.json"
		// if fd.FileExists(out) {
		// 	return
		// }
		out := in

		data, err := os.ReadFile(in)
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
		in := "./out/la-Mathematics.json"
		// out := "./out/la-Mathematics-out.json"
		// if fd.FileExists(out) {
		// 	return
		// }
		out := in

		data, err := os.ReadFile(in)
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
		in := "./out/la-Science.json"
		// out := "./out/la-Science-out.json"
		// if fd.FileExists(out) {
		// 	return
		// }
		out := in

		data, err := os.ReadFile(in)
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
		in := "./out/la-The Arts.json"
		// out := "./out/la-The Arts-out.json"
		// if fd.FileExists(out) {
		// 	return
		// }
		out := in

		data, err := os.ReadFile(in)
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
		in := "./out/la-Technologies.json"
		// out := "./out/la-Technologies-out.json"
		// if fd.FileExists(out) {
		// 	return
		// }
		out := in

		data, err := os.ReadFile(in)
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
