package main

import (
  "bufio"
  "fmt"
  "io"
)

const PROMPT = ">> "

func runRepl(in io.Reader, out io.Writer) {
  scanner := bufio.NewScanner(in)

  for {
    fmt.Printf(PROMPT)

    scanned := scanner.Scan()

    if !scanned {
      return
    }

    line := scanner.Text()

    l := NewLexer(line)

    for token := l.Advance(); token.Kind != EOF; token = l.Advance() {
      fmt.Printf("%+v\n", token)
    }
  }
}
