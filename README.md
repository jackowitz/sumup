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
