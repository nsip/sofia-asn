package tool

import (
	"fmt"
	"os"
	"testing"
)

func TestGetAllCode(t *testing.T) {
	data, err := os.ReadFile("../data/tree.pretty.json")
	if err != nil {
		panic(err)
	}
	mCodeParent := GetCodeParentMap(data)

	code := "ASARTDANY12"
	ancestors := GetCodeAncestors(mCodeParent, code)
	fmt.Println(ancestors)

	code = "AC9S1U02_E2"
	ancestors = GetCodeAncestors(mCodeParent, code)
	fmt.Println(ancestors)

	code = "PSCSEMC0_1"
	ancestors = GetCodeAncestors(mCodeParent, code)
	fmt.Println(ancestors)

	fmt.Println(GetAncestorTitle(mCodeParent, "ASMATY9L", ""))
	fmt.Println(GetAncestorTitle(mCodeParent, code, ""))

	fmt.Println(GetCodeAncestor(mCodeParent, code, 0))

	// for k := range mCodeTitle1 {
	// 	p := GetCodeAncestors(mCodeParent, k)[0]
	// 	if p == "GC" {
	// 		fmt.Println(p, k)
	// 	}
	// }

}
