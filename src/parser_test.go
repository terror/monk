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

func testIntegerLiteral(t *testing.T, expression Expression, value int64) bool {
  integ, ok := expression.(*IntegerLiteral)

  if !ok {
    t.Errorf("expression not *IntegerLiteral. got=%T", expression)
    return false
  }

  if integ.Value != value {
    t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
    return false
  }

  if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
    t.Errorf(
      "integ.TokenLiteral not %d. got=%s",
      value,
      integ.TokenLiteral(),
    )
    return false
  }

  return true
}

func testBooleanLiteral(t *testing.T, expression Expression, value bool) bool {
  boolean, ok := expression.(*Boolean)

  if !ok {
    t.Errorf("Expression is not a *Boolean, got=%T", expression)
  }

  if boolean.Value != value {
    t.Errorf("boolean.Value not %t, got=%t", value, boolean.Value)
    return false
  }

  if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
    t.Errorf(
      "boolean.TokenLiteral() not %t, got=%s",
      value,
      boolean.TokenLiteral(),
    )
    return false
  }

  return true
}

func testIdentifier(t *testing.T, expression Expression, value string) bool {
  identifier, ok := expression.(*Identifier)

  if !ok {
    t.Errorf("Expression not an *Identifier, got=%T", expression)
    return false
  }

  if identifier.Value != value {
    t.Errorf("identifier.Value not %s, got=%s", value, identifier.Value)
    return false
  }

  if identifier.TokenLiteral() != value {
    t.Errorf(
      "identifier.TokenLiteral() not %s, got=%s",
      value,
      identifier.TokenLiteral(),
    )
    return false
  }

  return true
}

func testLiteralExpression(
  t *testing.T,
  expression Expression,
  expected interface{},
) bool {
  switch v := expected.(type) {
  case int:
    return testIntegerLiteral(t, expression, int64(v))
  case int64:
    return testIntegerLiteral(t, expression, v)
  case string:
    return testIdentifier(t, expression, v)
  case bool:
    return testBooleanLiteral(t, expression, v)
  }

  t.Errorf("Type of expression not handled, got=%T", expression)

  return false
}

func testInfixExpression(
  t *testing.T,
  expression Expression,
  left interface{},
  operator string,
  right interface{},
) bool {
  op, ok := expression.(*InfixExpression)

  if !ok {
    t.Errorf(
      "Expression is not an *InfixExpression, got=%T(%s)",
      expression,
      expression,
    )
    return false
  }

  if !testLiteralExpression(t, op.Left, left) {
    return false
  }

  if op.Operator != operator {
    t.Errorf("Expression operator is not %s, got=%q", operator, op.Operator)
    return false
  }

  return testLiteralExpression(t, op.Right, right)
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
    expectedExpression interface{}
  }{
    {"x", 5},
    {"y", 10},
    {"foobar", 9000},
  }

  for i, tt := range tests {
    statement := program.Statements[i]

    if statement.TokenLiteral() != "let" {
      t.Errorf(
        "statement.TokenLiteral not 'let', got=%q",
        statement.TokenLiteral(),
      )
    }

    letStatement, ok := statement.(*LetStatement)

    if !ok {
      t.Errorf("statement not *LetStatement, got=%T", statement)
    }

    if letStatement.Name.Value != tt.expectedIdentifier {
      t.Errorf(
        "letStatement.Name.Value not '%s'. got=%s",
        tt.expectedIdentifier,
        letStatement.Name.Value,
      )
    }

    if letStatement.Name.TokenLiteral() != tt.expectedIdentifier {
      t.Errorf(
        "letStatement.Name not %s, got=%s",
        tt.expectedIdentifier,
        letStatement.Name,
      )
    }

    if !testLiteralExpression(
      t,
      letStatement.Value,
      tt.expectedExpression,
    ) {
      t.Errorf(
        "letStatement.Value not %T, got=%T",
        letStatement.Value,
        tt.expectedExpression,
      )
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
    input    string
    operator string
    expected interface{}
  }{
    {"!5;", "!", 5},
    {"-15;", "-", 15},
    {"!true;", "!", true},
    {"!false;", "!", false},
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

    if !testLiteralExpression(t, expression.Right, tt.expected) {
      return
    }
  }
}

func TestInfixExpressions(t *testing.T) {
  infixTests := []struct {
    input      string
    leftValue  interface{}
    operator   string
    rightValue interface{}
  }{
    {"5 + 5;", 5, "+", 5},
    {"5 - 5;", 5, "-", 5},
    {"5 * 5;", 5, "*", 5},
    {"5 / 5;", 5, "/", 5},
    {"5 > 5;", 5, ">", 5},
    {"5 < 5;", 5, "<", 5},
    {"5 == 5;", 5, "==", 5},
    {"5 != 5;", 5, "!=", 5},
    {"true == true", true, "==", true},
    {"true != false", true, "!=", false},
    {"false == false", false, "==", false},
  }

  for _, tt := range infixTests {
    program := setup(t, tt.input)

    if len(program.Statements) != 1 {
      t.Fatalf(
        "program.Statements does not contain %d statements. got=%d\n",
        1,
        len(program.Statements),
      )
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
        "expression is not an *InfixExpression. got=%T",
        statement.Expression,
      )
    }

    if !testInfixExpression(
      t,
      expression,
      tt.leftValue,
      tt.operator,
      tt.rightValue,
    ) {
      return
    }
  }
}

func TestOperatorPrecdence(t *testing.T) {
  tests := []struct {
    input    string
    expected string
  }{
    {"-a * b", "((-a) * b)"},
    {"!-a", "(!(-a))"},
    {"a + b + c", "((a + b) + c)"},
    {"a + b - c", "((a + b) - c)"},
    {"a * b * c", "((a * b) * c)"},
    {"a * b / c", "((a * b) / c)"},
    {"a + b / c", "(a + (b / c))"},
    {"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
    {"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
    {"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
    {"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
    {
      "3 + 4 * 5 == 3 * 1 + 4 * 5",
      "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
    },
    {
      "3 + 4 * 5 == 3 * 1 + 4 * 5",
      "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
    },
    {"true", "true"},
    {"false", "false"},
    {"3 > 5 == false", "((3 > 5) == false)"},
    {"3 < 5 == true", "((3 < 5) == true)"},
    {
      "1 + (2 + 3) + 4",
      "((1 + (2 + 3)) + 4)",
    },
    {
      "(5 + 5) * 2",
      "((5 + 5) * 2)",
    },
    {
      "2 / (5 + 5)",
      "(2 / (5 + 5))",
    },
    {
      "-(5 + 5)",
      "(-(5 + 5))",
    },
    {
      "!(true == true)",
      "(!(true == true))",
    },
  }

  for _, tt := range tests {
    program := setup(t, tt.input)

    actual := program.String()

    if actual != tt.expected {
      t.Errorf("Expected=%q, got=%q", tt.expected, actual)
    }
  }
}

func TestIfExpression(t *testing.T) {
  program := setup(t, `if (x < y) { x } else { y }`)

  if len(program.Statements) != 1 {
    t.Fatalf("program.Body does not contain %d statements. got=%d\n",
      1, len(program.Statements))
  }

  stmt, ok := program.Statements[0].(*ExpressionStatement)

  if !ok {
    t.Fatalf("program.Statements[0] is not an *ExpressionStatement. got=%T",
      program.Statements[0])
  }

  exp, ok := stmt.Expression.(*IfExpression)

  if !ok {
    t.Fatalf("stmt.Expression is not an *IfExpression. got=%T",
      stmt.Expression)
  }

  if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
    return
  }

  if len(exp.Consequence.Statements) != 1 {
    t.Errorf("consequence is not 1 statements. got=%d\n",
      len(exp.Consequence.Statements))
  }

  consequence, ok := exp.Consequence.Statements[0].(*ExpressionStatement)

  if !ok {
    t.Fatalf("Statements[0] is not an *ExpressionStatement. got=%T",
      exp.Consequence.Statements[0])
  }

  if !testIdentifier(t, consequence.Expression, "x") {
    return
  }

  if exp.Alternative == nil {
    t.Errorf(
      "exp.Alternative.Statements was nil.",
    )
  }
}
