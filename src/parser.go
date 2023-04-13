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
  ASTERISK: PRODUCT,
  EQ:       EQUALS,
  GT:       LESSGREATER,
  LPAREN:   CALL,
  LT:       LESSGREATER,
  MINUS:    SUM,
  NOT_EQ:   EQUALS,
  PLUS:     SUM,
  SLASH:    PRODUCT,
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
  p.registerPrefix(BANG, p.parsePrefixExpression)
  p.registerPrefix(FALSE, p.parseBoolean)
  p.registerPrefix(FUNCTION, p.parseFunctionLiteral)
  p.registerPrefix(IDENT, p.parseIdentifier)
  p.registerPrefix(IF, p.parseIfExpression)
  p.registerPrefix(INT, p.parseIntegerLiteral)
  p.registerPrefix(LPAREN, p.parseGroupedExpression)
  p.registerPrefix(MINUS, p.parsePrefixExpression)
  p.registerPrefix(TRUE, p.parseBoolean)

  p.infix = make(map[TokenKind]infixParseFn)
  p.registerInfix(ASTERISK, p.parseInfixExpression)
  p.registerInfix(EQ, p.parseInfixExpression)
  p.registerInfix(GT, p.parseInfixExpression)
  p.registerInfix(LPAREN, p.parseCallExpression)
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

  p.advance()

  statement.Value = p.parseExpression(LOWEST)

  p.advanceUntil(SEMICOLON)

  return statement
}

func (p *Parser) parseReturnStatement() *ReturnStatement {
  statement := &ReturnStatement{Token: p.curr}

  p.advance()

  statement.ReturnValue = p.parseExpression(LOWEST)

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

func (p *Parser) parseFunctionLiteral() Expression {
  literal := &FunctionLiteral{Token: p.curr}

  if !p.expectPeek(LPAREN) {
    return nil
  }

  literal.Parameters = p.parseFunctionParameters()

  if !p.expectPeek(LBRACE) {
    return nil
  }

  literal.Body = p.parseBlockStatement()

  return literal
}

func (p *Parser) parseFunctionParameters() []*Identifier {
  identifers := []*Identifier{}

  if p.peek.Kind == RPAREN {
    p.advance()
    return identifers
  }

  p.advance()

  identifers = append(
    identifers,
    &Identifier{Token: p.curr, Value: p.curr.Literal},
  )

  for p.peek.Kind == COMMA {
    p.advance()
    p.advance()
    identifers = append(
      identifers,
      &Identifier{Token: p.curr, Value: p.curr.Literal},
    )
  }

  if !p.expectPeek(RPAREN) {
    return nil
  }

  return identifers
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

func (p *Parser) parseBoolean() Expression {
  return &Boolean{Token: p.curr, Value: p.curr.Kind == TRUE}
}

func (p *Parser) parseGroupedExpression() Expression {
  p.advance()

  expression := p.parseExpression(LOWEST)

  if !p.expectPeek(RPAREN) {
    return nil
  }

  return expression
}

func (p *Parser) parseIfExpression() Expression {
  expression := *&IfExpression{Token: p.curr}

  if !p.expectPeek(LPAREN) {
    return nil
  }

  p.advance()

  expression.Condition = p.parseExpression(LOWEST)

  if !p.expectPeek(RPAREN) {
    return nil
  }

  if !p.expectPeek(LBRACE) {
    return nil
  }

  expression.Consequence = p.parseBlockStatement()

  if p.peek.Kind == ELSE {
    p.advance()

    if !p.expectPeek(LBRACE) {
      return nil
    }

    expression.Alternative = p.parseBlockStatement()
  }

  return &expression
}

func (p *Parser) parseBlockStatement() *BlockStatement {
  block := &BlockStatement{Token: p.curr}

  block.Statements = []Statement{}

  p.advance()

  for p.curr.Kind != RBRACE && p.curr.Kind != EOF {
    statement := p.parseStatement()

    if statement != nil {
      block.Statements = append(block.Statements, statement)
    }

    p.advance()
  }

  return block
}

func (p *Parser) parseCallExpression(function Expression) Expression {
  exp := &CallExpression{Token: p.curr, Function: function}
  exp.Arguments = p.parseCallArguments()
  return exp
}

func (p *Parser) parseCallArguments() []Expression {
  args := []Expression{}

  if p.peek.Kind == RPAREN {
    p.advance()
    return args
  }

  p.advance()

  args = append(args, p.parseExpression(LOWEST))

  for p.peek.Kind == COMMA {
    p.advance()
    p.advance()
    args = append(args, p.parseExpression(LOWEST))
  }

  if !p.expectPeek(RPAREN) {
    return nil
  }

  return args
}
