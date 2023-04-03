package main

const (
  ILLEGAL   = "ILLEGAL"
  EOF       = "EOF"
  IDENT     = "IDENT"
  INT       = "INT"
  ASSIGN    = "="
  PLUS      = "+"
  COMMA     = ","
  SEMICOLON = ";"
  LPAREN    = "("
  RPAREN    = ")"
  LBRACE    = "{"
  RBRACE    = "}"
  FUNCTION  = "FUNCTION"
  LET       = "LET"
)

type TokenKind string

type Token struct {
  Kind    TokenKind
  Literal string
}

func NewToken(kind TokenKind, ch byte) Token {
  return Token{Kind: kind, Literal: string(ch)}
}
