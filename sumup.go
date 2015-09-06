package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/jackowitzd2/sumup/parse"
)

type interpreter struct {
	r *bufio.Scanner
	w io.Writer
}

func NewInterpreter(input io.Reader, output io.Writer) *interpreter {
	return &interpreter{
		r: bufio.NewScanner(input),
		w: output,
	}
}

func (i *interpreter) loop() {
	for {
		io.WriteString(i.w, "> ")
		if success := i.r.Scan(); !success {
			panic("error reading expression from input")
		}
		expr := i.r.Text()
		tree, err := sumup.Parse(expr)
		if err != nil {
			fmt.Println(err)
		} else {
			io.WriteString(i.w, fmt.Sprintf("%v\n", tree.Execute()))
		}
	}
}

func main() {
	NewInterpreter(os.Stdin, os.Stdout).loop()
}
