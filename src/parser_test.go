package main

import (
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
