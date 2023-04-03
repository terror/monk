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
