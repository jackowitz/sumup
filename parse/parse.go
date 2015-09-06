package sumup

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

type SumNode []ProductNode
type ProductNode []NumberNode
type NumberNode int

type Tree struct {
	Root SumNode

	// Used during parsing only.
	acc   ProductNode
	lex   *lexer
	token item
}

func Parse(expr string) (*Tree, error) {
	tree := New()
	err := tree.Parse(expr)
	return tree, err
}

func New() *Tree {
	t := &Tree{
		Root: make([]ProductNode, 0),
		acc:  make([]NumberNode, 0),
	}
	return t
}

func (t *Tree) next() item {
	defer func() {
		t.token = t.lex.nextItem()
	}()
	return t.token
}

func (t *Tree) peek() item {
	return t.token
}

func (t *Tree) nextNonSpace() (token item) {
	for {
		token = t.next()
		if token.typ != itemSpace {
			break
		}
	}
	return token
}

func (t *Tree) peekNonSpace() (token item) {
	for {
		token = t.peek()
		if token.typ != itemSpace {
			break
		}
		t.next()
	}
	return token
}

func (t *Tree) expect(expected itemType) item {
	token := t.nextNonSpace()
	if token.typ != expected {
		t.unexpected(token)
	}
	return token
}

func (t *Tree) expectOneOf(expected1, expected2 itemType) item {
	token := t.nextNonSpace()
	if token.typ != expected1 && token.typ != expected2 {
		t.unexpected(token)
	}
	return token
}

func (t *Tree) unexpected(token item) {
	if token.typ == itemError {
		t.errorf("%s", token.val)
	} else {
		t.errorf("unexpected %s", token)
	}
}

func (t *Tree) errorf(format string, args ...interface{}) {
	t.Root = nil
	panic(fmt.Errorf(format, args...))
}

func (t *Tree) error(err error) {
	t.errorf("%s", err)
}

func (t *Tree) recover(errp *error) {
	e := recover()
	if e != nil {
		if _, ok := e.(runtime.Error); ok {
			panic(e)
		}
		if t != nil {
			t.lex.drain()
			t.stopParse()
		}
		*errp = e.(error)
	}
	return
}

func (t *Tree) Parse(expr string) (err error) {
	defer t.recover(&err)
	t.startParse(lex(strings.NewReader(expr)))
	t.parseOperand()
	for t.peekNonSpace().typ != itemEOF {
		t.parseOperator()
		t.parseOperand()
	}
	t.Root = append(t.Root, t.acc)
	t.stopParse()
	return nil
}

func (t *Tree) startParse(lex *lexer) {
	t.lex = lex
	t.token = t.lex.nextItem()
}

func (t *Tree) stopParse() {
	t.acc = nil
	t.lex = nil
}

func (t *Tree) parseOperand() {
	token := t.expect(itemNumber)
	value, err := strconv.Atoi(token.val)
	if err != nil {
		t.errorf("illegal number syntax: %q", token.val)
	}
	t.acc = append(t.acc, NumberNode(value))
}

func (t *Tree) parseOperator() {
	token := t.expectOneOf(itemPlus, itemTimes)
	switch token.typ {
	case itemPlus:
		t.Root = append(t.Root, t.acc)
		t.acc = make([]NumberNode, 0)
	case itemTimes:
	}
}

func (t *Tree) Execute() int {
	result := 0
	for i := range t.Root {
		product := 1
		for j := range t.Root[i] {
			product *= int(t.Root[i][j])
		}
		result += product
	}
	return result
}
