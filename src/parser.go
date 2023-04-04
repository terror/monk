package main

import (
  "fmt"
)

type (
  prefixParseFn func() Expression
  infixParseFn  func(Expression) Expression
)

const (
  _ int = iota
  LOWEST
  EQUALS  // == LESSGREATER // > or <
  SUM     //+
  PRODUCT //*
  PREFIX  //-Xor!X
  CALL    // myFunction(X)
)

type Parser struct {
  curr           Token
  errors         []string
  lexer          *Lexer
  peek           Token
  prefixParseFns map[TokenKind]prefixParseFn
  infixParseFns  map[TokenKind]infixParseFn
}

func NewParser(lexer *Lexer) *Parser {
  p := &Parser{lexer: lexer, errors: []string{}}

  p.prefixParseFns = make(map[TokenKind]prefixParseFn)
  p.registerPrefix(IDENT, p.parseIdentifier)

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

func (p *Parser) registerPrefix(kind TokenKind, fn prefixParseFn) {
  p.prefixParseFns[kind] = fn
}

func (p *Parser) registerInfix(kind TokenKind, fn infixParseFn) {
  p.infixParseFns[kind] = fn
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
  p.errors = append(
    p.errors,
    fmt.Sprintf("Expected next token to be %s but got %s instead", kind, p.peek.Kind),
  )
}

func (p *Parser) parseStatement() Statement {
  switch p.curr.Kind {
  case LET:
    return p.parseLetStatement()
  case RETURN:
    return p.parseReturnStatement()
  default:
    return p.parseExpressionStatement()
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

func (p *Parser) parseReturnStatement() *ReturnStatement {
  statement := &ReturnStatement{Token: p.curr}

  p.advance()

  p.advanceUntil(SEMICOLON)

  return statement
}

func (p *Parser) parseExpression(precedence int) Expression {
  prefix := p.prefixParseFns[p.curr.Kind]

  if prefix == nil {
    return nil
  }

  return prefix()
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
  statement := &ExpressionStatement{Token: p.curr}

  statement.Expression = p.parseExpression(LOWEST)

  if p.peek.Kind == SEMICOLON {
    p.advance()
  }

  return statement
}

func (p *Parser) parseIdentifier() Expression {
  return &Identifier{Token: p.curr, Value: p.curr.Literal}
}
