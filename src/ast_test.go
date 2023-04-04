package main

import (
  "testing"
)

func TestString(t *testing.T) {
  program := &Program{
    Statements: []Statement{
      &LetStatement{
        Token: Token{Kind: LET, Literal: "let"},
        Name: &Identifier{
          Token: Token{Kind: IDENT, Literal: "myVar"},
          Value: "myVar"},
        Value: &Identifier{
          Token: Token{Kind: IDENT, Literal: "anotherVar"},
          Value: "anotherVar",
        },
      },
    },
  }

  if program.String() != "let myVar = anotherVar;" {
    t.Errorf(
      "program.String() wrong. got=%q",
      program.String(),
    )
  }
}
