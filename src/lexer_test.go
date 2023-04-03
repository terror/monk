package main

import (
  "testing"
)

func TestAdvance(t *testing.T) {
  input := `=+(){},;`

  tests := []struct {
    expectedKind    TokenKind
    expectedLiteral string
  }{
    {ASSIGN, "="},
    {PLUS, "+"},
    {LPAREN, "("},
    {RPAREN, ")"},
    {LBRACE, "{"},
    {RBRACE, "}"},
    {COMMA, ","},
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
