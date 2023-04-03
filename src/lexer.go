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

func (l *Lexer) consume(pred func(byte) bool) string {
  position := l.position

  for pred(l.ch) {
    l.read()
  }

  return l.input[position:l.position]
}

func (l *Lexer) Advance() Token {
  var token Token

  l.consume(isWhitespace)

  switch l.ch {
  case '!':
    token = NewToken(BANG, l.ch)
  case '(':
    token = NewToken(LPAREN, l.ch)
  case ')':
    token = NewToken(RPAREN, l.ch)
  case '*':
    token = NewToken(ASTERISK, l.ch)
  case '+':
    token = NewToken(PLUS, l.ch)
  case ',':
    token = NewToken(COMMA, l.ch)
  case '-':
    token = NewToken(MINUS, l.ch)
  case '/':
    token = NewToken(SLASH, l.ch)
  case ';':
    token = NewToken(SEMICOLON, l.ch)
  case '<':
    token = NewToken(LT, l.ch)
  case '=':
    token = NewToken(ASSIGN, l.ch)
  case '>':
    token = NewToken(GT, l.ch)
  case '{':
    token = NewToken(LBRACE, l.ch)
  case '}':
    token = NewToken(RBRACE, l.ch)
  case 0:
    token.Literal = ""
    token.Kind = EOF
  default:
    if isLetter(l.ch) {
      token.Literal = l.consume(isLetter)
      token.Kind = LookupIdent(token.Literal)
      return token
    } else if isDigit(l.ch) {
      token.Kind = INT
      token.Literal = l.consume(isDigit)
      return token
    } else {
      token = NewToken(ILLEGAL, l.ch)
    }
  }

  l.read()

  return token
}
