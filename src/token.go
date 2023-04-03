package main

var keywords = map[string]TokenKind{
  "fn":  FUNCTION,
  "let": LET,
}

func LookupIdent(ident string) TokenKind {
  if kind, ok := keywords[ident]; ok {
    return kind
  }

  return IDENT
}

const (
  ASSIGN    = "="
  COMMA     = ","
  EOF       = "EOF"
  FUNCTION  = "FUNCTION"
  IDENT     = "IDENT"
  ILLEGAL   = "ILLEGAL"
  INT       = "INT"
  LBRACE    = "{"
  LET       = "LET"
  LPAREN    = "("
  PLUS      = "+"
  RBRACE    = "}"
  RPAREN    = ")"
  SEMICOLON = ";"
)

type TokenKind string

type Token struct {
  Kind    TokenKind
  Literal string
}

func NewToken(kind TokenKind, ch byte) Token {
  return Token{Kind: kind, Literal: string(ch)}
}
