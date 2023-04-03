package main

var keywords = map[string]TokenKind{
  "else":   ELSE,
  "false":  FALSE,
  "fn":     FUNCTION,
  "if":     IF,
  "let":    LET,
  "return": RETURN,
  "true":   TRUE,
}

func LookupIdent(ident string) TokenKind {
  if kind, ok := keywords[ident]; ok {
    return kind
  }

  return IDENT
}

const (
  ASSIGN    = "="
  ASTERISK  = "*"
  BANG      = "!"
  COMMA     = ","
  ELSE      = "ELSE"
  EOF       = "EOF"
  FALSE     = "FALSE"
  FUNCTION  = "FUNCTION"
  GT        = ">"
  IDENT     = "IDENT"
  IF        = "IF"
  ILLEGAL   = "ILLEGAL"
  INT       = "INT"
  LBRACE    = "{"
  LET       = "LET"
  LPAREN    = "("
  LT        = "<"
  MINUS     = "-"
  PLUS      = "+"
  RBRACE    = "}"
  RETURN    = "RETURN"
  RPAREN    = ")"
  SEMICOLON = ";"
  SLASH     = "/"
  TRUE      = "TRUE"
)

type TokenKind string

type Token struct {
  Kind    TokenKind
  Literal string
}

func NewToken(kind TokenKind, ch byte) Token {
  return Token{Kind: kind, Literal: string(ch)}
}
