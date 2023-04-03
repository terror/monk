package main

func isDigit(ch byte) bool {
  return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
  return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isWhitespace(ch byte) bool {
  return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}
