package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/digisan/go-generics/v2"
	gio "github.com/digisan/gotk/io"
	"github.com/digisan/gotk/misc"
	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
	"github.com/tidwall/gjson"
)

func init() {
	lk.WarnDetail(false)
	lk.Log2F(true, "./out.log")
}

// compute weight by nested array element sequence.
// in each level's 'children' count less than 100
func weight(path string, lvl int) int64 {
	ss := strings.Split(path, ".")
	ns := ""
	for i := 0; i < 2*lvl && i < len(ss); i++ {
		if i%2 == 1 {
			ns += fmt.Sprintf("%02s", ss[i])
		}
	}
	n, err := strconv.ParseInt(ns, 10, 64)
	if err != nil {
		log.Fatalln(err)
	}
	return n
}

// get certain level's field('code') & its value,
// return map for path('code')-value
func getCodes(fpath string, level int) map[string]string {

	jsBytes, err := os.ReadFile(fpath)
	if err != nil {
		log.Fatalln(err)
	}
	js := string(jsBytes)

	paths, _ := jt.GetLeavesPathOrderly(js)
	fmt.Println(len(paths), "total paths")

	fpaths := jt.GetLeafPathsOrderly("code", paths)
	fmt.Println(len(fpaths), "code paths")

	mPathCode := make(map[string]string, len(fpaths))

	if level < 0 {
		for _, path := range fpaths {
			mPathCode[path] = gjson.Get(js, path).String()
		}
	} else {
		for _, path := range fpaths {
			if lvl := strings.Count(path, "children"); lvl == level {
				v := gjson.Get(js, path).String()
				// fmt.Printf("%05d - %d - %s - %s\n", i, lvl, path, v)
				mPathCode[path] = v
			}
		}
	}
	return mPathCode
}

func extractNodeCodes(fpath, opath string) int {
	mCodesNode := getCodes("../data/node.pretty.json", -1)
	_, codesNode := Map2KVs(mCodesNode, nil, nil)
	for _, code := range codesNode {
		gio.MustAppendFile(opath, []byte(code), true)
	}
	return len(mCodesNode)
}

func loadNodeCodes(fpath string) (codes []string) {
	gio.FileLineScan(fpath, func(line string) (bool, string) {
		codes = append(codes, line)
		return false, ""
	}, "")
	return Settify(codes...)
}

func loadScotTxtMapCodeUrls(fpath string) map[string][]string {
	mCodeUrls := make(map[string][]string)
	gio.FileLineScan(fpath, func(line string) (bool, string) {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			ss := strings.Split(line, "\t")
			lk.FailOnErrWhen(len(ss) != 2, "%v", fmt.Errorf("SCOT TXT line is NOT %d parts", 2))
			mCodeUrls[ss[0]] = append(mCodeUrls[ss[0]], ss[1])
		}
		return false, ""
	}, "")
	return mCodeUrls
}

func loadScotTxtMapUrlCodes(fpath string) map[string][]string {
	mUrlCodes := make(map[string][]string)
	gio.FileLineScan(fpath, func(line string) (bool, string) {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			ss := strings.Split(line, "\t")
			lk.FailOnErrWhen(len(ss) != 2, "%v", fmt.Errorf("SCOT TXT line is NOT %d parts", 2))

			// extra process...
			url := strings.TrimSuffix(ss[1], ".html")
			if strings.Count(url, "//") > 1 {
				p := strings.LastIndex(url, "//")
				url = url[:p] + url[p+1:]
			}

			mUrlCodes[url] = append(mUrlCodes[url], ss[0])
		}
		return false, ""
	}, "")
	return mUrlCodes
}

// id is url style
func loadScotJsonldIDs(fpath string) (ids []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r, err := os.Open(fpath)
	lk.FailOnErr("%v", err)

	cObj, cErr, err := jt.ScanObjectInArray(ctx, r, true)
	lk.FailOnErr("%v", err)

	for obj := range cObj {
		err := <-cErr
		lk.FailOnErr("%v", err)

		id := gjson.Get(obj, "@id").String()
		ids = append(ids, id)
	}
	return Settify(ids...)
}

func main() {

	defer misc.TrackTime(time.Now())

	/// 1) create "node.code.txt"
	//
	// {
	// 	n := extractNodeCodes("../data/node.pretty.json", "../data/tmp/node.code.txt")
	// 	fmt.Println(n)
	// 	return
	// }

	// Tree.json /////////////////////////////////////////////////////////////////////
	var codesTree []string
	{
		lk.Log(" Tree.json ...")

		// fpath := "../data/tree.pretty.json"
		fpath := "../data/Sofia_API_Data_25052022.json"

		lvl := 4
		mCodesTree := getCodes(fpath, lvl)
		// for k, v := range codes {
		// 	fmt.Println(k, v)
		// }

		_, codesTree = Map2KVs(mCodesTree, func(i, j string) bool {
			return weight(i, lvl) < weight(j, lvl)
		}, nil)
		codesTree = Settify(codesTree...)

		// for i := 0; i < len(pathsTree); i++ {
		// 	fmt.Printf("%05d - %s - %s\n", i, pathsTree[i], codesTree[i])
		// }
		fmt.Println("codesTree count:", len(codesTree))
	}

	// Node.json /////////////////////////////////////////////////////////////////////
	var codesNode []string
	{
		lk.Log(" Node.json ...")

		/// slow if from original json
		//
		// mCodesNode := getCodes("../data/node.pretty.json", -1)
		// _, codesNode := Map2KVs(mCodesNode, nil, nil)

		/// fast if already did 1)
		//
		codesNode := loadNodeCodes("../data/tmp/node.code.txt")
		lk.Warn("codesNode count: %v", len(codesNode))
	}

	// SCOT.txt /////////////////////////////////////////////////////////////////////
	var (
		mCodeUrls map[string][]string
		mUrlCodes map[string][]string
	)
	{
		lk.Log(" SCOT.txt ...")

		mCodeUrls = loadScotTxtMapCodeUrls("../data/ACv9_ScOT_BC_20220422.txt")
		lk.Warn("mCodeUrls count: %v", len(mCodeUrls))

		mUrlCodes = loadScotTxtMapUrlCodes("../data/ACv9_ScOT_BC_20220422.txt")
		lk.Warn("mUrlCodes count: %v", len(mUrlCodes))
	}

	//
	// SCOT.jsonld /////////////////////////////////////////////////////////////////
	var ids []string
	{
		lk.Log(" SCOT.jsonld ...")

		ids = loadScotJsonldIDs("../data/scot.jsonld")
		lk.Warn("scot-jsonld ids count: %v", len(ids))
	}

	//--------------------------------------------------------------------------------------------------------//

	lk.Log("Checking------------------------------------------------------------------------------")

	// Check Tree <---> Node ///////////////////////////////////////////////////////
	{
		lk.Log(" Check which is in tree but not in node ...")

		missing := Minus(codesTree, codesNode)
		lk.Warn(" missing count: %v", len(missing))
		for _, m := range missing {
			lk.Log("%v", m)
		}
	}

	// Check Tree <---> Scot Txt ///////////////////////////////////////////////////////
	{
		lk.Log(" Check which is in tree but not in Scot Txt ...")

		missing := []string{}
		for _, code := range codesTree {
			if _, ok := mCodeUrls[code]; !ok {
				missing = append(missing, code)
			}
		}
		lk.Warn(" missing count: %v", len(missing))
		for _, m := range missing {
			lk.Log("%v", m)
		}
	}

	// Check Scot Txt <---> Scot jsonld ///////////////////////////////////////////////////////
	{
		lk.Log(" Check which is in Scot Txt but not in Scot jsonld ...")

		missing := []string{}
		for url := range mUrlCodes {
			if NotIn(url, ids...) {
				missing = append(missing, url)
			}
		}
		lk.Warn(" missing count: %d", len(missing))
		for _, m := range missing {
			lk.Log("%v", m)
		}
	}

}
