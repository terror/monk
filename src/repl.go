package main

import (
  "bufio"
  "fmt"
  "io"
)

const PROMPT = ">> "

func repl(in io.Reader, out io.Writer) {
  scanner := bufio.NewScanner(in)

  for {
    fmt.Print(PROMPT)

    scanned := scanner.Scan()

    if !scanned {
      return
    }

    line := scanner.Text()

    parser := NewParser(NewLexer(line))

    program := parser.Parse()

    if len(parser.Errors()) != 0 {
      for _, message := range parser.Errors() {
        fmt.Println(message)
      }
    } else {
      fmt.Println(program.String())
    }
  }
}
