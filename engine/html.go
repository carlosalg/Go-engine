package engine

import (
	"engine/dom"
	"fmt"
	"strings"
	"unicode"
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
  if char := p.consumeChar(); char != '<'{
    return fmt.Errorf("Expected '<' at the beggining, got '%c'", char)
  }
  tagName := p.parseTagName()
  attrs := p.parseAttributes()
  if char := p.consumeChar(); char != '>'{  
    return fmt.Errorf("Expected '<' at the beggining, got '%c'", char)
  }
  
  children := p.parseNodes()
if char := p.consumeChar(); char != '<'{  
    return fmt.Errorf("Expected '<' at the beggining, got '%c'", char)
  }
if char := p.consumeChar(); char != '/'{  
    return fmt.Errorf("Expected '/' at the beggining, got '%c'", char)
  }
if tag := p.parseTagName(); tag != tagName {  
    return fmt.Errorf("Expected  tag, got '%s'", tag)
  }
if char := p.consumeChar(); char != '>'{  
    return fmt.Errorf("Expected '>' at the beggining, got '%c'", char)
  }
  
  return dom.Elem(tagName,attrs,children)
}

func (p *Parser) parseAttr() (string,string){
  name := p.parseTagName()
  if char := p.consumeChar(); char != '='{  
    err := fmt.Errorf("Expected '=' at the beggining, got '%c'", char)
    return "", err.Error()
  }
  value := p.parseAttrValue()
  return name,value
}

func (p *Parser)parseAttrValue() string {
  openQuote := p.consumeChar()

  if openQuote != '"' || openQuote != '\'' {
    err := fmt.Errorf(`Expected " or ', got '%c' `, openQuote)
    return err.Error()
  }

  value := p.consumeWhile(func(c rune) bool{return c != openQuote})

  if char := p.consumeChar(); char != openQuote {
    err := fmt.Errorf(`Expected " or ' found '%c'`, openQuote)
    return err.Error()
  }

  return value

}

func (p *Parser) parseAttributes() dom.AttrMap {
  attributes := make(map[string]string)
  for {
    p.consumeWhiteSpace()

    if p.nextChar() == '>'{
      break
    }

    name, value := p.parseAttr()

    attributes[name] = value
  }
  return attributes
}

func (p *Parser) parseNodes() []dom.Node {
  var nodes []dom.Node 
  for {
    p.consumeWhiteSpace()
    if p.eof() || p.startsWith("</") {
      break
    }

    nodes = append(nodes, p.parseNode())
  }
  return nodes 
}


