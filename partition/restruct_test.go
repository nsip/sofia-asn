package main

import (
	"log"
	"os"
	"testing"
)

func TestRestructEng(t *testing.T) {
	data, err := os.ReadFile("./out/la-English.json")
	if err != nil {
		log.Fatalln(err)
	}
	js := reStructEng(string(data))
	os.WriteFile("./out/la-English(restruct).json", []byte(js), os.ModePerm)
}
 
func TestRestructHASS(t *testing.T) {
	data, err := os.ReadFile("./out/la-HASS.json")
	if err != nil {
		log.Fatalln(err)
	}
	js := reStructHASS(string(data))
	os.WriteFile("./out/la-HASS(restruct).json", []byte(js), os.ModePerm)
}

func TestRestructHPE(t *testing.T) {
	data, err := os.ReadFile("./out/la-HPE.json")
	if err != nil {
		log.Fatalln(err)
	}
	js := reStructHPE(string(data))
	os.WriteFile("./out/la-HPE(restruct).json", []byte(js), os.ModePerm)
}

func TestRestructLang(t *testing.T) {
	data, err := os.ReadFile("./out/la-Languages.json")
	if err != nil {
		log.Fatalln(err)
	}
	js := reStructLang(string(data))
	os.WriteFile("./out/la-Languages(restruct).json", []byte(js), os.ModePerm)
}

func TestRestructMath(t *testing.T) {
	data, err := os.ReadFile("./out/la-Mathematics.json")
	if err != nil {
		log.Fatalln(err)
	}
	js := reStructMath(string(data))
	os.WriteFile("./out/la-Mathematics(restruct).json", []byte(js), os.ModePerm)
}

func TestRestructSci(t *testing.T) {
	data, err := os.ReadFile("./out/la-Science.json")
	if err != nil {
		log.Fatalln(err)
	}
	js := reStructSci(string(data))
	os.WriteFile("./out/la-Science(restruct).json", []byte(js), os.ModePerm)
}

func TestRestructTech(t *testing.T) {
	data, err := os.ReadFile("./out/la-Technologies.json")
	if err != nil {
		log.Fatalln(err)
	}
	js := reStructTech(string(data))
	os.WriteFile("./out/la-Technologies(restruct).json", []byte(js), os.ModePerm)
}

func TestRestructArts(t *testing.T) {
	data, err := os.ReadFile("./out/la-The Arts.json")
	if err != nil {
		log.Fatalln(err)
	}
	js := reStructArt(string(data))
	os.WriteFile("./out/la-The Arts(restruct).json", []byte(js), os.ModePerm)
}
