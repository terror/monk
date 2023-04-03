package main

type Lexer struct {
  ch           byte
  input        string
  position     int
  readPosition int
}

func NewLexer(input string) *Lexer {
  l := &Lexer{input: input}
  l.read()
  return l
}

func (l *Lexer) read() {
  if l.readPosition >= len(l.input) {
    l.ch = 0
  } else {
    l.ch = l.input[l.readPosition]
  }
  l.position = l.readPosition
  l.readPosition += 1
}

func (l *Lexer) Advance() Token {
  var token Token

  switch l.ch {
  case '=':
    token = newToken(ASSIGN, l.ch)
  case ';':
    token = newToken(SEMICOLON, l.ch)
  case '(':
    token = newToken(LPAREN, l.ch)
  case ')':
    token = newToken(RPAREN, l.ch)
  case ',':
    token = newToken(COMMA, l.ch)
  case '+':
    token = newToken(PLUS, l.ch)
  case '{':
    token = newToken(LBRACE, l.ch)
  case '}':
    token = newToken(RBRACE, l.ch)
  case 0:
    token.Literal = ""
    token.Kind = EOF
  }

  l.read()

  return token
}

func newToken(tokenKind TokenKind, ch byte) Token { return Token{Kind: tokenKind, Literal: string(ch)} }
