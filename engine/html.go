package engine

import (
	"strings"
	"unicode"
	"unicode/utf8"
  "engine/dom"
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

func (p *Parser)consumeWhile(test func(rune)bool) string{
  var result strings.Builder
  for !p.eof() && test(p.nextChar()){
    result.WriteRune(p.consumeChar())
  }
  return result.String()
}

func (p *Parser)consumeWhiteSpace(){
  p.consumeWhile(unicode.IsSpace)
}

func (p *Parser)parseTagName() string {
  result := p.consumeWhile(func(c rune)bool{
    switch {
    case 'a' <= c && c <= 'z':
      return true
    case 'A' <= c && c <= 'Z':
      return true
    case '0' <= c && c <= '9':
      return true
    default:
      return false
    }
  })
  return result
}

//Parse functions

func (p *Parser)parseNode() dom.Node {
  switch p.nextChar() {
    case '<':
      return p.parseElement()
    default :
      return p.parseText()
  }
  
}

func (p *Parser)parseText() dom.Node {
  return dom.Text(p.consumeWhile(func(c rune)bool{return c != '<'})) 
}

func (p *Parser)parseElement() dom.Node {
  //TODO:implemt this 
}
