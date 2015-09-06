SumUp is a handwritten lexer and parser for a (rather trivial)
algebraic language supporting sum-of-product calculations.

The regular expression accepting the grammar:
```
[ \t]*[0-9]+([ \t]*[+*][ \t]*[0-9]+)+[ \t]*
```
Broken out:
```
[ \t]*   - optional leading whitespace
[0-9]+   - a mandatory number
(
  [ \t]* - optional whitespace
  [+*]   - an operator from the set {+,*}
  [ \t]* - whitespace, again optional
  [0-9]+ - another number
)*       - all repeated zero or more times
[ \t]*   - optional trailing whitespace
```

Notably, parentheses are not supported; precedence rules are
implicit and as expected. Namely, all multiplications are
performed before any additions.

`sumup.go` contains a basic REPL but all of the fun stuff
is under `parse`. Expressions are interpreted in three
steps. Lexing (`parse/lex.go`) transforms the raw input
into symbols then parsing (`parse/parse.go`) builds and
executes (although the latter should probably be
pulled out separately) the abstract syntax tree.

The general design and some of the utility functions for
text processing are modeled closely after the Go
`text/template`
[package](http://golang.org/src/text/template/).
