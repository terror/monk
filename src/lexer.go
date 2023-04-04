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

func (l *Lexer) Advance() Token {
  var token Token

  l.eat(isWhitespace)

  switch l.ch {
  case '!':
    if l.peek() == '=' {
      ch := l.ch
      l.read()
      token = Token{Kind: NOT_EQ, Literal: string(ch) + string(l.ch)}
    } else {
      token = NewToken(BANG, l.ch)
    }
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
    if l.peek() == '=' {
      ch := l.ch
      l.read()
      token = Token{Kind: EQ, Literal: string(ch) + string(l.ch)}
    } else {
      token = NewToken(ASSIGN, l.ch)
    }
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
      token.Literal = l.take(isLetter)
      token.Kind = LookupIdent(token.Literal)
      return token
    } else if isDigit(l.ch) {
      token.Kind = INT
      token.Literal = l.take(isDigit)
      return token
    } else {
      token = NewToken(ILLEGAL, l.ch)
    }
  }

  l.read()

  return token
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

func (l *Lexer) eat(pred func(byte) bool) {
  for pred(l.ch) {
    l.read()
  }
}

func (l *Lexer) take(pred func(byte) bool) string {
  position := l.position

  for pred(l.ch) {
    l.read()
  }

  return l.input[position:l.position]
}

func (l *Lexer) peek() byte {
  if l.readPosition >= len(l.input) {
    return 0
  } else {
    return l.input[l.readPosition]
  }
}
