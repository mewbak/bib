package main

import (
	"testing"

	"github.com/nickng/bibtex"
)

type TestEntry struct {
	Type   string
	Name   string
	Fields map[string]string
}

func (t TestEntry) Entry() *Entry {
	e := bibtex.NewBibEntry(t.Type, t.Name)
	for name, value := range t.Fields {
		e.AddField(name, bibtex.BibConst(value))
	}
	return &Entry{*e}
}

func TestFormat(t *testing.T) {
	cases := []struct {
		TestEntry
		Expect string
	}{
		{
			TestEntry: TestEntry{
				Name: "single_author",
				Type: "misc",
				Fields: map[string]string{
					"author": "First Author",
					"title":  "Title",
				},
			},
			Expect: "First Author. Title.",
		},
		{
			TestEntry: TestEntry{
				Name: "two_authors",
				Type: "misc",
				Fields: map[string]string{
					"author": "First Author and Second Author",
					"title":  "Title",
				},
			},
			Expect: "First Author and Second Author. Title.",
		},
		{
			TestEntry: TestEntry{
				Name: "multi_author",
				Type: "misc",
				Fields: map[string]string{
					"author": "First Author and Second Author and Third Author",
					"title":  "Title",
				},
			},
			Expect: "First Author, Second Author and Third Author. Title.",
		},
		{
			TestEntry: TestEntry{
				Name: "url_urldate",
				Type: "misc",
				Fields: map[string]string{
					"author":  "First Author",
					"title":   "Title",
					"url":     "https://golang.org",
					"urldate": "2020-02-06",
				},
			},
			Expect: "First Author. Title. https://golang.org (accessed February 6, 2020)",
		},
		{
			TestEntry: TestEntry{
				Name: "misc",
				Type: "misc",
				Fields: map[string]string{
					"author":       "First Author",
					"title":        "Title",
					"howpublished": "School Newsletter",
					"license":      "MIT License",
				},
			},
			Expect: "First Author. Title. School Newsletter. MIT License.",
		},
		{
			TestEntry: TestEntry{
				Name: "inproceedings",
				Type: "inproceedings",
				Fields: map[string]string{
					"author":    "First Author",
					"title":     "Title",
					"booktitle": "Handbook of Golang",
					"pages":     "42--78",
				},
			},
			Expect: "First Author. Title. In Handbook of Golang, pages 42--78.",
		},
	}
	for _, c := range cases {
		c := c // scopelint
		t.Run(c.Name, func(t *testing.T) {
			e := c.Entry()
			got, err := Format(e)
			if err != nil {
				t.Fatal(err)
			}
			if got != c.Expect {
				t.Logf("entry  = %s", e)
				t.Logf("got    = %s", got)
				t.Logf("expect = %s", c.Expect)
				t.Fail()
			}
		})
	}
}
