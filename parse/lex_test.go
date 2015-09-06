package sumup

import (
	"strings"
	"testing"
)

type lexTest struct {
	name  string
	input string
	items []item
}

var (
	tEOF   = func(i Pos) item { return item{itemEOF, i, ""} }
	tPlus  = func(i Pos) item { return item{itemPlus, i, "+"} }
	tTimes = func(i Pos) item { return item{itemTimes, i, "*"} }
)

var lexTests = []lexTest{
	{"empty", "", []item{tEOF(0)}},
	{"plus", "+", []item{tPlus(0), tEOF(1)}},
	{"times", "*", []item{tTimes(0), tEOF(1)}},
	{"digit", "7", []item{item{itemNumber, 0, "7"}, tEOF(1)}},
	{"number", "1262", []item{item{itemNumber, 0, "1262"}, tEOF(4)}},
	{"add_two", "3+142", []item{
		item{itemNumber, 0, "3"},
		tPlus(1), item{itemNumber, 2, "142"}, tEOF(5),
	}},
	{"times_two", "3*142", []item{
		item{itemNumber, 0, "3"},
		tTimes(1), item{itemNumber, 2, "142"}, tEOF(5),
	}},
	{"add_four", "3+142+7+19874", []item{
		item{itemNumber, 0, "3"},
		tPlus(1), item{itemNumber, 2, "142"},
		tPlus(5), item{itemNumber, 6, "7"},
		tPlus(7), item{itemNumber, 8, "19874"}, tEOF(13),
	}},
	{"spaces", "3   +142 \t +\t  \t\t 7", []item{
		item{itemNumber, 0, "3"},
		item{itemSpace, 1, "   "},
		tPlus(4), item{itemNumber, 5, "142"},
		item{itemSpace, 8, " \t "},
		tPlus(11), item{itemSpace, 12, "\t  \t\t "},
		item{itemNumber, 18, "7"}, tEOF(19),
	}},
	{"error", "3+?", []item{
		item{itemNumber, 0, "3"},
		tPlus(1), item{itemError, 2, "unexpected rune: U+003F '?'"},
	}},
}

func collect(t *lexTest) (items []item) {
	l := lex(strings.NewReader(t.input))
	for {
		item := l.nextItem()
		items = append(items, item)
		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}
	return
}

func equal(i1, i2 []item, checkPos bool) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			return false
		}
		if i1[k].val != i2[k].val {
			return false
		}
		if checkPos && i1[k].pos != i2[k].pos {
			return false
		}
	}
	return true
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		items := collect(&test)
		if !equal(items, test.items, true) {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%v", test.name, items, test.items)
		}
	}
}
