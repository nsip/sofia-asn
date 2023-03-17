package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	. "github.com/digisan/go-generics/v2"
	jt "github.com/digisan/json-tool"
)

func TestTreeProc3(t *testing.T) {
	os.MkdirAll("./out/", os.ModePerm)

	file := "la-English.json"
	data, err := os.ReadFile(`../partition/out/` + file)
	if err != nil {
		log.Fatalln(err)
	}

	mData, _ := jt.Flatten(data)
	connSet := []string{}
	for k := range mData {
		if p := strings.Index(k, "connections."); p >= 0 { // strings.Contains(k, "connections.") {
			k = k[p:]
			if p := strings.LastIndex(k, "."); p >= 0 {
				k = k[:p]
			}
			connSet = append(connSet, k)
		}
	}
	connSet = Settify(connSet...)
	for _, conn := range connSet {
		fmt.Println(conn)
	}

	// s := treeProc3(data, "English")
	// os.WriteFile("./out/testout.json", []byte(s), os.ModePerm)
}

func TestSort(t *testing.T) {
	sortRule := func(s1, s2 string) bool {
		a1, a2 := []int{}, []int{}
		for i, s := range []string{s1, s2} {
			var a *[]int
			if i == 0 {
				a = &a1
			} else {
				a = &a2
			}
			for _, seg := range strings.Split(s, ".") {
				if IsNumeric(seg) {
					n, _ := strconv.Atoi(seg)
					*a = append(*a, n)
				}
			}
		}
		lmin := Min(len(a1), len(a2))
		for i := 0; i < lmin; i++ {
			n1, n2 := a1[i], a2[i]
			switch {
			case n1 < n2:
				return true
			case n1 > n2:
				return false
			default:
				continue
			}
		}
		return true
	}

	m := map[string]string{
		"a.2.b.1":         "",
		"c.2.y.2.r.9.q.1": "",
	}

	ks, vs := MapToKVs(m, sortRule, nil)
	fmt.Println(ks)
	fmt.Println(vs)
}
