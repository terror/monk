package main

import ("fmt")

type Parser struct {
  curr  Token
  errors []string
  lexer *Lexer
  peek  Token
}

func NewParser(lexer *Lexer) *Parser {
  p := &Parser{lexer: lexer, errors: []string{}}

  p.advance()
  p.advance()

  return p
}

func (p *Parser) Errors() []string {
  return p.errors
}

func (p *Parser) Parse() *Program {
  program := &Program{}

  program.Statements = []Statement{}

  for p.curr.Kind != EOF {
    statement := p.parseStatement()

    if statement != nil {
      program.Statements = append(program.Statements, statement)
    }

    p.advance()
  }

  return program
}

func (p *Parser) advance() {
  p.curr = p.peek
  p.peek = p.lexer.Advance()
}

func (p *Parser) advanceUntil(kind TokenKind) {
  for p.curr.Kind != kind {
    p.advance()
  }
}

func (p *Parser) expectPeek(kind TokenKind) bool {
  if p.peek.Kind == kind {
    p.advance()
    return true
  } else {
    p.peekError(kind)
    return false
  }
}

func (p *Parser) peekError(kind TokenKind) {
  p.errors = append(p.errors, fmt.Sprintf("Expected next token to be %s but got %s instead", kind, p.peek.Kind))
}

func (p *Parser) parseStatement() Statement {
  switch p.curr.Kind {
  case LET:
    return p.parseLetStatement()
  default:
    return nil
  }
}

func (p *Parser) parseLetStatement() *LetStatement {
  statement := &LetStatement{Token: p.curr}

  if !p.expectPeek(IDENT) {
    return nil
  }

  statement.Name = &Identifier{Token: p.curr, Value: p.curr.Literal}

  if !p.expectPeek(ASSIGN) {
    return nil
  }

  p.advanceUntil(SEMICOLON)

  return statement
}
