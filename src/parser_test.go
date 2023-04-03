package main

import (
  "testing"
)

func TestLetStatements(t *testing.T) {
  input := `
    let x = 5;
    let y = 10;
    let foobar = 9000;
  `

  lexer := NewLexer(input)
  parser := NewParser(lexer)

  program := parser.Parse()

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
      t.Errorf("statement not *ast.LetStatement, got=%T", statement)
    }
  }
}
