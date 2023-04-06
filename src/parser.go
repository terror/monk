package main

import (
  "fmt"
  "strconv"
)

type (
  prefixParseFn func() Expression
  infixParseFn  func(Expression) Expression
)

const (
  _ int = iota
  LOWEST
  EQUALS
  LESSGREATER
  SUM
  PRODUCT
  PREFIX
  CALL
)

var precedences = map[TokenKind]int{
  EQ:       EQUALS,
  NOT_EQ:   EQUALS,
  LT:       LESSGREATER,
  GT:       LESSGREATER,
  PLUS:     SUM,
  MINUS:    SUM,
  SLASH:    PRODUCT,
  ASTERISK: PRODUCT,
}

type Parser struct {
  curr   Token
  errors []string
  infix  map[TokenKind]infixParseFn
  lexer  *Lexer
  peek   Token
  prefix map[TokenKind]prefixParseFn
}

func NewParser(lexer *Lexer) *Parser {
  p := &Parser{lexer: lexer, errors: []string{}}

  p.prefix = make(map[TokenKind]prefixParseFn)
  p.registerPrefix(IDENT, p.parseIdentifier)
  p.registerPrefix(INT, p.parseIntegerLiteral)
  p.registerPrefix(BANG, p.parsePrefixExpression)
  p.registerPrefix(MINUS, p.parsePrefixExpression)

  p.infix = make(map[TokenKind]infixParseFn)
  p.registerInfix(ASTERISK, p.parseInfixExpression)
  p.registerInfix(EQ, p.parseInfixExpression)
  p.registerInfix(GT, p.parseInfixExpression)
  p.registerInfix(LT, p.parseInfixExpression)
  p.registerInfix(MINUS, p.parseInfixExpression)
  p.registerInfix(NOT_EQ, p.parseInfixExpression)
  p.registerInfix(PLUS, p.parseInfixExpression)
  p.registerInfix(SLASH, p.parseInfixExpression)

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
  p.prefix[kind] = fn
}

func (p *Parser) registerInfix(kind TokenKind, fn infixParseFn) {
  p.infix[kind] = fn
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

func (p *Parser) peekPrecedence() int {
  if p, ok := precedences[p.peek.Kind]; ok {
    return p
  }

  return LOWEST
}

func (p *Parser) currPrecedence() int {
  if p, ok := precedences[p.curr.Kind]; ok {
    return p
  }

  return LOWEST
}

func (p *Parser) peekError(kind TokenKind) {
  p.errors = append(
    p.errors,
    fmt.Sprintf(
      "Expected next token to be %s but got %s instead",
      kind,
      p.peek.Kind,
    ),
  )
}

func (p *Parser) missingPrefixError(kind TokenKind) {
  p.errors = append(
    p.errors,
    fmt.Sprintf("No prefix parse function for %s found", kind),
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
  prefix := p.prefix[p.curr.Kind]

  if prefix == nil {
    p.missingPrefixError(p.curr.Kind)
    return nil
  }

  left := prefix()

  for p.peek.Kind != SEMICOLON && precedence < p.peekPrecedence() {
    infix := p.infix[p.peek.Kind]

    if infix == nil {
      return left
    }

    p.advance()

    left = infix(left)
  }

  return left
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

func (p *Parser) parseIntegerLiteral() Expression {
  literal := &IntegerLiteral{Token: p.curr}

  value, err := strconv.ParseInt(p.curr.Literal, 0, 64)

  if err != nil {
    msg := fmt.Sprintf("Could not parse %q as integer", p.curr.Literal)
    p.errors = append(p.errors, msg)
    return nil
  }

  literal.Value = value

  return literal
}

func (p *Parser) parsePrefixExpression() Expression {
  expression := &PrefixExpression{
    Token:    p.curr,
    Operator: p.curr.Literal,
  }

  p.advance()

  expression.Right = p.parseExpression(PREFIX)

  return expression
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
  expression := &InfixExpression{
    Token:    p.curr,
    Operator: p.curr.Literal,
    Left:     left,
  }

  precedence := p.currPrecedence()

  p.advance()

  expression.Right = p.parseExpression(precedence)

  return expression
}
