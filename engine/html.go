package engine

import (
  "strings"
  "unicode/utf8"
)

type Parser struct {
  pos int
  input string
}

func (p *Parser)nextChar() rune{
  if p.pos < len(p.input){
    char := rune(p.input[p.pos])
    return char 
  }
  return 0
}

func (p *Parser)startsWith(s string) bool {
  return strings.HasPrefix(p.input[p.pos:], s)
}

func (p *Parser)eof() bool {
  return p.pos >= len(p.input)
}

func (p *Parser)consumeChar() rune {
  r, size := utf8.DecodeLastRuneInString(p.input[p.pos:])
  currChar := r
  nextPos := size
  p.pos += nextPos
  return currChar
}
