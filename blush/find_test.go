package blush_test

import (
	"regexp"
	"testing"

	"github.com/arsham/blush/blush"
)

func TestNewLocatorExact(t *testing.T) {
	l := blush.NewLocator("aaa", false)
	if _, ok := l.(blush.Exact); !ok {
		t.Errorf("l = %T, want *blush.Exact", l)
	}
	l = blush.NewLocator("*aaa", false)
	if _, ok := l.(blush.Exact); !ok {
		t.Errorf("l = %T, want *blush.Exact", l)
	}
}

func TestNewLocatorIexact(t *testing.T) {
	l := blush.NewLocator("aaa", true)
	if _, ok := l.(blush.Iexact); !ok {
		t.Errorf("l = %T, want *blush.Iexact", l)
	}
	l = blush.NewLocator("*aaa", true)
	if _, ok := l.(blush.Iexact); !ok {
		t.Errorf("l = %T, want *blush.Iexact", l)
	}
}

func TestNewLocatorRx(t *testing.T) {
	tcs := []struct {
		name    string
		input   string
		matches []string
	}{
		{"empty", "^$", []string{""}},
		{"starts with", "^aaa", []string{"aaa", "aaa sss"}},
		{"ends with", "aaa$", []string{"aaa", "sss aaa"}},
		{"with star", "blah blah.*", []string{"blah blah", "aa blah blah aa"}},
		{"with curly brackets", "a{3}", []string{"aaa", "aa aaa aa"}},
		{"with brackets", "[ab]", []string{"kjhadf", "kjlrbrlkj", "sdbsdha"}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			l := blush.NewLocator(tc.input, false)
			if _, ok := l.(blush.Rx); !ok {
				t.Errorf("l = %T, want *blush.Rx", l)
			}
			l = blush.NewLocator(tc.input, false)
			if _, ok := l.(blush.Rx); !ok {
				t.Errorf("l = %T, want *blush.Rx", l)
			}
		})
	}
}

func TestExactFind(t *testing.T) {
	l := blush.Exact("nooooo")
	got, ok := l.Find("yessss", blush.NoColour)
	if got != "" {
		t.Errorf("got = %s, want `%s`", got, "")
	}
	if ok {
		t.Error("ok = true, want false")
	}

	tcs := []struct {
		name   string
		search string
		colour blush.Colour
		input  string
		want   string
		wantOk bool
	}{
		{"exact no colour", "aaa", blush.NoColour, "aaa", "aaa", true},
		{"exact not found", "aaaa", blush.NoColour, "aaa", "", false},
		{"some parts no colour", "aaa", blush.NoColour, "bb aaa bb", "bb aaa bb", true},
		{"exact blue", "aaa", blush.FgBlue, "aaa", blush.Colourise("aaa", blush.FgBlue), true},
		{"some parts blue", "aaa", blush.FgBlue, "bb aaa bb", "bb " + blush.Colourise("aaa", blush.FgBlue) + " bb", true},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			l := blush.Exact(tc.search)
			got, ok := l.Find(tc.input, tc.colour)
			if got != tc.want {
				t.Errorf("got = `%s`, want `%s`", got, tc.want)
			}
			if ok != tc.wantOk {
				t.Errorf("ok = %t, want %t", ok, tc.wantOk)
			}
		})
	}
}

func TestRxFind(t *testing.T) {
	l := blush.Rx{Regexp: regexp.MustCompile("nooooo")}
	got, ok := l.Find("yessss", blush.NoColour)
	if got != "" {
		t.Errorf("got = %s, want `%s`", got, "")
	}
	if ok {
		t.Error("ok = true, want false")
	}

	tcs := []struct {
		name   string
		search string
		colour blush.Colour
		input  string
		want   string
		wantOk bool
	}{
		{"exact no colour", "(^aaa$)", blush.NoColour, "aaa", "aaa", true},
		{"exact not found", "(^aa$)", blush.NoColour, "aaa", "", false},
		{"some parts no colour", "(aaa)", blush.NoColour, "bb aaa bb", "bb aaa bb", true},
		{"some parts not matched", "(Aaa)", blush.NoColour, "bb aaa bb", "", false},
		{"exact blue", "(aaa)", blush.FgBlue, "aaa", blush.Colourise("aaa", blush.FgBlue), true},
		{"some parts blue", "(aaa)", blush.FgBlue, "bb aaa bb", "bb " + blush.Colourise("aaa", blush.FgBlue) + " bb", true},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			l := blush.Rx{Regexp: regexp.MustCompile(tc.search)}
			got, ok := l.Find(tc.input, tc.colour)
			if got != tc.want {
				t.Errorf("got = `%s`, want `%s`", got, tc.want)
			}
			if ok != tc.wantOk {
				t.Errorf("ok = %t, want %t", ok, tc.wantOk)
			}
		})
	}

	rx := blush.NewLocator("a{3}", false)
	want := "this " + blush.Colourise("aaa", blush.FgBlue) + "meeting"
	got, ok = rx.Find("this aaameeting", blush.FgBlue)
	if got != want {
		t.Errorf("got = `%s`, want `%s`", got, want)
	}
	if !ok {
		t.Error("ok = false, want true")
	}
}

func TestIexact(t *testing.T) {
	l := blush.Iexact("nooooo")
	got, ok := l.Find("yessss", blush.NoColour)
	if got != "" {
		t.Errorf("got = %s, want `%s`", got, "")
	}
	if ok {
		t.Error("ok = true, want false")
	}

	tcs := []struct {
		name   string
		search string
		colour blush.Colour
		input  string
		want   string
		wantOk bool
	}{
		{"exact no colour", "aaa", blush.NoColour, "aaa", "aaa", true},
		{"exact not found", "aaaa", blush.NoColour, "aaa", "", false},
		{"i exact no colour", "AAA", blush.NoColour, "aaa", "aaa", true},
		{"some parts no colour", "aaa", blush.NoColour, "bb aaa bb", "bb aaa bb", true},
		{"i some parts no colour", "AAA", blush.NoColour, "bb aaa bb", "bb aaa bb", true},
		{"exact blue", "aaa", blush.FgBlue, "aaa", blush.Colourise("aaa", blush.FgBlue), true},
		{"i exact blue", "AAA", blush.FgBlue, "aaa", blush.Colourise("aaa", blush.FgBlue), true},
		{"some parts blue", "aaa", blush.FgBlue, "bb aaa bb", "bb " + blush.Colourise("aaa", blush.FgBlue) + " bb", true},
		{"i some parts blue", "AAA", blush.FgBlue, "bb aaa bb", "bb " + blush.Colourise("aaa", blush.FgBlue) + " bb", true},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			l := blush.Iexact(tc.search)
			got, ok := l.Find(tc.input, tc.colour)
			if got != tc.want {
				t.Errorf("got = `%s`, want `%s`", got, tc.want)
			}
			if ok != tc.wantOk {
				t.Errorf("ok = %t, want %t", ok, tc.wantOk)
			}
		})
	}

}

func TestRxInsensitiveFind(t *testing.T) {
	tcs := []struct {
		name   string
		search string
		colour blush.Colour
		input  string
		want   string
		wantOk bool
	}{
		{"exact no colour", "^AAA$", blush.NoColour, "aaa", "aaa", true},
		{"exact not found", "^AA$", blush.NoColour, "aaa", "", false},
		{"some words no colour", `AAA*`, blush.NoColour, "bb aaa bb", "bb aaa bb", true},
		{"exact blue", "^AAA$", blush.FgBlue, "aaa", blush.Colourise("aaa", blush.FgBlue), true},
		{"some words blue", "AAA?", blush.FgBlue, "bb aaa bb", "bb " + blush.Colourise("aaa", blush.FgBlue) + " bb", true},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			l := blush.NewLocator(tc.search, true)
			if _, ok := l.(blush.Rx); !ok {
				t.Fatalf("l = %T, want blush.Rx", l)
			}
			got, ok := l.Find(tc.input, tc.colour)
			if got != tc.want {
				t.Errorf("got = `%s`, want `%s`", got, tc.want)
			}
			if ok != tc.wantOk {
				t.Errorf("ok = %t, want %t", ok, tc.wantOk)
			}
		})
	}

	rx := blush.NewLocator("A{3}", true)
	want := "this " + blush.Colourise("aaa", blush.FgBlue) + "meeting"
	got, ok := rx.Find("this aaameeting", blush.FgBlue)
	if got != want {
		t.Errorf("got = `%s`, want `%s`", got, want)
	}
	if !ok {
		t.Error("ok = false, want true")
	}
}
