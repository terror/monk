package main

import (
  "testing"
)

func TestAdvance(t *testing.T) {
  input := `
    let five = 5;

    let ten = 10;

    let add = fn(x, y) {
      x + y;
    };

    let result = add(five, ten);
  `

  tests := []struct {
    expectedKind    TokenKind
    expectedLiteral string
  }{
    {LET, "let"},
    {IDENT, "five"},
    {ASSIGN, "="},
    {INT, "5"},
    {SEMICOLON, ";"},
    {LET, "let"},
    {IDENT, "ten"},
    {ASSIGN, "="},
    {INT, "10"},
    {SEMICOLON, ";"},
    {LET, "let"},
    {IDENT, "add"},
    {ASSIGN, "="},
    {FUNCTION, "fn"},
    {LPAREN, "("},
    {IDENT, "x"},
    {COMMA, ","},
    {IDENT, "y"},
    {RPAREN, ")"},
    {LBRACE, "{"},
    {IDENT, "x"},
    {PLUS, "+"},
    {IDENT, "y"},
    {SEMICOLON, ";"},
    {RBRACE, "}"},
    {SEMICOLON, ";"},
    {LET, "let"},
    {IDENT, "result"},
    {ASSIGN, "="},
    {IDENT, "add"},
    {LPAREN, "("},
    {IDENT, "five"},
    {COMMA, ","},
    {IDENT, "ten"},
    {RPAREN, ")"},
    {SEMICOLON, ";"},
    {EOF, ""},
  }

  l := NewLexer(input)

  for i, tt := range tests {
    token := l.Advance()

    if token.Kind != tt.expectedKind {
      t.Fatalf("tests[%d] - Wrong token kind: expected=%q, got=%q", i, tt.expectedKind, token.Kind)
    }

    if token.Literal != tt.expectedLiteral {
      t.Fatalf("test[%d] - Wrong literal: expected=%q, got=%q", i, tt.expectedLiteral, token.Literal)
    }
  }
}
