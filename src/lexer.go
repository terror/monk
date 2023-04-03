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
    token = NewToken(ASSIGN, l.ch)
  case ';':
    token = NewToken(SEMICOLON, l.ch)
  case '(':
    token = NewToken(LPAREN, l.ch)
  case ')':
    token = NewToken(RPAREN, l.ch)
  case ',':
    token = NewToken(COMMA, l.ch)
  case '+':
    token = NewToken(PLUS, l.ch)
  case '{':
    token = NewToken(LBRACE, l.ch)
  case '}':
    token = NewToken(RBRACE, l.ch)
  case 0:
    token.Literal = ""
    token.Kind = EOF
  }

  l.read()

  return token
}
