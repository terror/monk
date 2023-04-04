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

func TestLetStatements(t *testing.T) {
  input := `
    let x = 5;
    let y = 10;
    let foobar = 9000;
  `

  parser := NewParser(NewLexer(input))

  program := parser.Parse()

  validate(t, parser)

  if program == nil {
    t.Fatalf("Parse() returned nil")
  }

  if len(program.Statements) != 3 {
    t.Fatalf("program.Statements does not contain 3 statements, got=%d", len(program.Statements))
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
      t.Errorf("letStatement.Name.Value not '%s'. got=%s", tt.expectedIdentifier, letStatement.Name.Value)
    }

    if letStatement.Name.TokenLiteral() != tt.expectedIdentifier {
      t.Errorf("statement.Name not %s, got=%s", tt.expectedIdentifier, letStatement.Name)
    }

    if !ok {
      t.Errorf("statement not *LetStatement, got=%T", statement)
    }
  }
}
