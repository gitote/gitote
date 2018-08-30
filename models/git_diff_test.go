package models

import (
	"html/template"
	"testing"

	dmp "github.com/sergi/go-diff/diffmatchpatch"
	"gitote.com/gitote/git-module"
)

func assertEqual(t *testing.T, s1 string, s2 template.HTML) {
	if s1 != string(s2) {
		t.Errorf("%s should be equal %s", s2, s1)
	}
}

func assertLineEqual(t *testing.T, d1 *git.DiffLine, d2 *git.DiffLine) {
	if d1 != d2 {
		t.Errorf("%v should be equal %v", d1, d2)
	}
}

func Test_diffToHTML(t *testing.T) {
	assertEqual(t, "+foo <span class=\"added-code\">bar</span> biz", diffToHTML([]dmp.Diff{
		dmp.Diff{dmp.DiffEqual, "foo "},
		dmp.Diff{dmp.DiffInsert, "bar"},
		dmp.Diff{dmp.DiffDelete, " baz"},
		dmp.Diff{dmp.DiffEqual, " biz"},
	}, git.DIFF_LINE_ADD))

	assertEqual(t, "-foo <span class=\"removed-code\">bar</span> biz", diffToHTML([]dmp.Diff{
		dmp.Diff{dmp.DiffEqual, "foo "},
		dmp.Diff{dmp.DiffDelete, "bar"},
		dmp.Diff{dmp.DiffInsert, " baz"},
		dmp.Diff{dmp.DiffEqual, " biz"},
	}, git.DIFF_LINE_DEL))
}
