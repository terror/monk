package main

import (
  "fmt"
  "testing"
)

func validate(t *testing.T, p *Parser) {
  errors := p.Errors()

  if len(errors) == 0 {
    return
  }

  t.Errorf("Parser has %d errors", len(errors))

  for _, message := range errors {
    t.Errorf("Parser error: %q", message)
  }

  t.FailNow()
}

func setup(t *testing.T, input string) Program {
  parser := NewParser(NewLexer(input))
  program := parser.Parse()
  validate(t, parser)
  return *program
}

func testIntegerLiteral(t *testing.T, il Expression, value int64) bool {
  integ, ok := il.(*IntegerLiteral)

  if !ok {
    t.Errorf("il not *IntegerLiteral. got=%T", il)
    return false
  }

  if integ.Value != value {
    t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
    return false
  }

  if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
    t.Errorf("integ.TokenLiteral not %d. got=%s", value,
      integ.TokenLiteral())
    return false
  }

  return true
}

func TestLetStatements(t *testing.T) {
  input := `
    let x = 5;
    let y = 10;
    let foobar = 9000;
  `

  program := setup(t, input)

  if len(program.Statements) != 3 {
    t.Fatalf(
      "program.Statements does not contain 3 statements, got=%d",
      len(program.Statements),
    )
  }

  tests := []struct {
    expectedIdentifier string
  }{
    {"x"},
    {"y"},
    {"foobar"},
  }

  for i, tt := range tests {
    statement := program.Statements[i]

    if statement.TokenLiteral() != "let" {
      t.Errorf("statement.TokenLiteral not 'let', got=%q", statement.TokenLiteral())
    }

    letStatement, ok := statement.(*LetStatement)

    if letStatement.Name.Value != tt.expectedIdentifier {
      t.Errorf(
        "letStatement.Name.Value not '%s'. got=%s",
        tt.expectedIdentifier,
        letStatement.Name.Value,
      )
    }

    if letStatement.Name.TokenLiteral() != tt.expectedIdentifier {
      t.Errorf(
        "statement.Name not %s, got=%s",
        tt.expectedIdentifier,
        letStatement.Name,
      )
    }

    if !ok {
      t.Errorf("statement not *LetStatement, got=%T", statement)
    }
  }
}

func TestReturnStatement(t *testing.T) {
  input := `
    return 5;
    return 10;
    return 993322;
  `

  program := setup(t, input)

  if len(program.Statements) != 3 {
    t.Fatalf(
      "program.Statements does not contain 3 statements, got=%d",
      len(program.Statements),
    )
  }

  for _, statement := range program.Statements {
    returnStatement, ok := statement.(*ReturnStatement)

    if !ok {
      t.Errorf("Statement not ReturnStatement, got=%T", statement)
      continue
    }

    if returnStatement.TokenLiteral() != "return" {
      t.Errorf(
        "returnStatement.TokenLiteral not 'return', got %q",
        returnStatement.TokenLiteral(),
      )
    }
  }
}

func TestIdentifierExpression(t *testing.T) {
  input := "foobar;"

  program := setup(t, input)

  if len(program.Statements) != 1 {
    t.Fatalf(
      "Program doesn't have enough statements, got=%d",
      len(program.Statements),
    )
  }

  statement, ok := program.Statements[0].(*ExpressionStatement)

  if !ok {
    t.Fatalf(
      "program.Statements[0] is not an ExpressionStatement, got=%T",
      program.Statements[0],
    )
  }

  ident, ok := statement.Expression.(*Identifier)

  if !ok {
    t.Fatalf(
      "Expression is not a *Identifier, got=%T",
      statement.Expression,
    )
  }

  if ident.Value != "foobar" {
    t.Errorf(
      "ident.Value not %s, got=%s",
      "foobar",
      statement.Expression,
    )
  }

  if ident.TokenLiteral() != "foobar" {
    t.Errorf(
      "ident.TokenLiteral not %s, got=%s",
      "foobar",
      ident.TokenLiteral(),
    )
  }
}

func TestIntegerLiteralExpression(t *testing.T) {
  input := "5;"

  program := setup(t, input)

  if len(program.Statements) != 1 {
    t.Fatalf(
      "Program doesn't have enough statements, got=%d",
      len(program.Statements),
    )
  }

  statement, ok := program.Statements[0].(*ExpressionStatement)

  if !ok {
    t.Fatalf(
      "program.statements[0] is not an *ExpressionStatement, got=%T",
      program.Statements[0],
    )
  }

  literal, ok := statement.Expression.(*IntegerLiteral)

  if !ok {
    t.Fatalf(
      "Expression is not a *IntegerLiteral, got=%T",
      statement.Expression,
    )
  }

  if literal.Value != 5 {
    t.Errorf(
      "literal.Value not %d, got=%d",
      5,
      literal.Value,
    )
  }

  if literal.TokenLiteral() != "5" {
    t.Errorf(
      "literal.TokenLiteral not %s, got=%s",
      "5",
      literal.TokenLiteral(),
    )
  }
}

func TestPrefixExpressions(t *testing.T) {
  prefixTests := []struct {
    input        string
    operator     string
    integerValue int64
  }{
    {"!5;", "!", 5},
    {"-15;", "-", 15},
  }

  for _, tt := range prefixTests {
    program := setup(t, tt.input)

    if len(program.Statements) != 1 {
      t.Fatalf(
        "program.Statements does not contain %d statement, got=%d",
        1,
        len(program.Statements),
      )
    }

    statement, ok := program.Statements[0].(*ExpressionStatement)

    if !ok {
      t.Fatalf(
        "program.Statements[0] is not an ExpressionStatement, got=%T",
        program.Statements[0],
      )
    }

    expression, ok := statement.Expression.(*PrefixExpression)

    if !ok {
      t.Fatalf(
        "statement is not a PrefixExpression, got=%T",
        statement.Expression,
      )
    }

    if expression.Operator != tt.operator {
      t.Fatalf(
        "expression.Operator is not '%s', got=%s",
        tt.operator,
        expression.Operator,
      )
    }

    if !testIntegerLiteral(t, expression.Right, tt.integerValue) {
      return
    }
  }
}

func TestInfixExpressions(t *testing.T) {
  infixTests := []struct {
    input      string
    leftValue  int64
    operator   string
    rightValue int64
  }{
    {"5 + 5;", 5, "+", 5},
    {"5 - 5;", 5, "-", 5},
    {"5 * 5;", 5, "*", 5},
    {"5 / 5;", 5, "/", 5},
    {"5 > 5;", 5, ">", 5},
    {"5 < 5;", 5, "<", 5},
    {"5 == 5;", 5, "==", 5},
    {"5 != 5;", 5, "!=", 5},
  }

  for _, tt := range infixTests {
    program := setup(t, tt.input)

    if len(program.Statements) != 1 {
      t.Fatalf(
        "program.Statements does not contain %d statements. got=%d\n",
        1, len(program.Statements))
    }

    statement, ok := program.Statements[0].(*ExpressionStatement)

    if !ok {
      t.Fatalf(
        "program.Statements[0] is not an ExpressionStatement. got=%T",
        program.Statements[0],
      )
    }

    expression, ok := statement.Expression.(*InfixExpression)

    if !ok {
      t.Fatalf(
        "expression is not an ast.InfixExpression. got=%T",
        statement.Expression,
      )
    }

    if !testIntegerLiteral(t, expression.Left, tt.leftValue) {
      return
    }

    if expression.Operator != tt.operator {
      t.Fatalf(
        "expression.Operator is not '%s'. got=%s",
        tt.operator,
        expression.Operator,
      )
    }

    if !testIntegerLiteral(t, expression.Right, tt.rightValue) {
      return
    }
  }
}
