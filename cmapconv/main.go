// Cmapconv is a tool for converting character map (cmap) files into Go code.
package main

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"fmt"
	"os"
	"slices"
	"strings"
	"text/template"
)

//go:embed cmap.go.tpl
var cmapTpl string

func main() {
	// Read source file
	gb1 := readCid2codeFile("./cmapconv/Adobe-GB1-6-cid2code.txt")

	cmap := func(ordering, from, to string) string {
		var records [][]string
		switch ordering {
		case "Adobe-GB1":
			records = gb1
		}
		i1 := slices.Index(records[0], from)
		i2 := slices.Index(records[0], to)

		var b strings.Builder
		fmt.Fprintf(&b, "// cmap from %s to %s\n", from, to)
		fmt.Fprintf(&b, "var cmap%s2%s = map[uint16]uint16{\n",
			strings.ReplaceAll(from, "-", ""),
			strings.ReplaceAll(to, "-", ""),
		)
		st := make(map[string]struct{})
		for i := 1; i < len(records); i++ {
			if records[i][i1] == "*" || records[i][i2] == "*" {
				continue
			}
			ss1 := strings.Split(records[i][i1], ",")
			ss2 := strings.Split(records[i][i2], ",")
			var s2 string
			for _, s := range ss2 {
				if len(s) == 5 {
					s = s[:4]
				}
				if len(s) != 4 {
					panic("invalid hex string length " + s2)
				}
				p1 := codePriority(s)
				p2 := codePriority(s2)
				if p1 == p2 {
					panic("code " + s2 + " and " + s + " have the same priority")
				}
				if p1 > p2 {
					s2 = s
				}
			}
			for _, s1 := range ss1 {
				if len(s1) == 5 {
					s1 = s1[:4]
				}
				if len(s1) != 4 {
					panic("invalid hex string length " + s1)
				}
				// a1ab -> ff5e
				// a1abv -> ff5ev (ignored)
				if _, ok := st[s1]; ok {
					continue
				}
				st[s1] = struct{}{}
				fmt.Fprintf(&b, "\t0x%s: 0x%s,\n", s1, s2)
			}
		}
		fmt.Fprintln(&b, "}")
		return b.String()
	}

	var buf bytes.Buffer
	t := template.New("file").Funcs(map[string]any{
		"cmap": cmap,
	})
	t = template.Must(t.Parse(cmapTpl))
	err := t.Execute(&buf, nil)
	if err != nil {
		panic(err)
	}
	os.WriteFile("cmap.go", buf.Bytes(), 0666)
}

func readCid2codeFile(path string) [][]string {
	file := must(os.ReadFile(path))
	cr := csv.NewReader(bytes.NewReader(file))
	cr.Comma = '\t'
	cr.Comment = '#'
	return must(cr.ReadAll())
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func codePriority(s string) int {
	s = strings.ToLower(s)
	switch {
	// CJK Unified Ideographs
	case s >= "4e00" && s <= "9fff":
		return 2
		// General Punctuation
	case s >= "2000" && s <= "206f":
		return 2
	// Mathematical Operators
	case s >= "2200" && s <= "22ff":
		return 1
	// Kangxi Radicals
	case s >= "2f00" && s <= "2fd5":
		return 1
	default:
		return 0
	case s == "":
		return -1
	}
}
