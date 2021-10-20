package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/digisan/gotk/filedir"
)

func TestMain(t *testing.T) {
	main()
}

func TestRestructure(t *testing.T) {

	fs, _, err := filedir.WalkFileDir("./out/", false)
	if err != nil {
		log.Fatalln(err)
	}

	for _, f := range fs {
		if strings.Contains(f, "la-") {
			data, err := os.ReadFile(f)
			if err != nil {
				log.Fatalln(err)
			}
			laRestructure(string(data), 1)
			fmt.Println("-----------------------", f)
		}
	}
}
