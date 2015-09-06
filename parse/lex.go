package sumup

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Pos int

type item struct {
	typ itemType
	pos Pos
	val string
}

type itemType int

const (
	itemError itemType = iota
	itemEOF
	itemSpace
	itemNumber
	itemPlus
	itemTimes
)

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return i.val
	case len(i.val) > 10:
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

type lexer struct {
	name    string
	input   *bufio.Reader
	pos     Pos
	state   stateFn
	lastPos Pos
	items   chan item
}

const eof = -1

func (l *lexer) read() rune {
	r, _, err := l.input.ReadRune()
	if err != nil {
		return eof
	}
	return r
}

func (l *lexer) unread() {
	l.input.UnreadRune()
}

func (l *lexer) emit(t itemType, v string) {
	l.items <- item{t, l.pos, v}
}

func lex(input io.Reader) *lexer {
	l := &lexer{
		input: bufio.NewReader(input), // XXX check first
		items: make(chan item),
	}
	go l.run()
	return l
}

func (l *lexer) run() {
	for l.state = lexExpr; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}

func (l *lexer) nextItem() item {
	item := <-l.items
	l.lastPos = item.pos
	return item
}

func (l *lexer) drain() {
	for range l.items {
	}
}

type stateFn func(*lexer) stateFn

func lexExpr(l *lexer) stateFn {
	switch r := l.read(); {
	case r == eof || isEndOfLine(r):
		l.emit(itemEOF, "")
		return nil
	case isSpace(r):
		l.unread()
		return lexSpace
	case unicode.IsDigit(r):
		l.unread()
		return lexNumber
	case r == '+':
		l.emit(itemPlus, "+")
		l.pos++
		return lexExpr
	case r == '*':
		l.emit(itemTimes, "*")
		l.pos++
		return lexExpr
	default:
		l.emit(itemError, fmt.Sprintf("unexpected rune: %#U", r))
		return nil
	}
}

func lexNumber(l *lexer) stateFn {
	return lexRun(l, itemNumber, "0123456789")
}

func lexSpace(l *lexer) stateFn {
	return lexRun(l, itemSpace, " \t")
}

func lexRun(l *lexer, typ itemType, valid string) stateFn {
	var buf bytes.Buffer
	buf.WriteRune(l.read())
Loop:
	for {
		switch r := l.read(); {
		case strings.IndexRune(valid, r) >= 0:
			buf.WriteRune(r)
		default:
			l.unread()
			break Loop
		}
	}
	token := buf.String()
	l.emit(typ, token)
	l.pos += Pos(len(token))
	return lexExpr
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}
