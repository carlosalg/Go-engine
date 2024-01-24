package engine

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
  "log"
)

type Parser struct {
  pos int
  input string
}

type ParseError struct{
  error
  ExpectedChar rune
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

func (p *Parser)parseNode() (Node,error) {
  switch p.nextChar() {
    case '<':
      return p.parseElement()
    default :
      return p.parseText(), nil
  }
  
}

func (p *Parser)parseText() Node {
  return Text(p.consumeWhile(func(c rune)bool{return c != '<'})) 
}

func (p *Parser)parseElement() (Node, error) {
  var node Node 
  if char := p.consumeChar(); char != '<'{
    return node,&ParseError{
      error: fmt.Errorf("Expected '<' at the beggining, got '%c'", char),
      ExpectedChar: '<',
    }
  }
  
  tagName := p.parseTagName()
  attrs := p.parseAttributes()

  if char := p.consumeChar(); char != '>'{  
    return node, &ParseError{
      error: fmt.Errorf("Expected '<' at the beggining, got '%c'", char),
      ExpectedChar: '>',
    }
  }
  
  children := p.parseNodes()

 if char := p.consumeChar(); char != '<'{  
    return node,&ParseError{
      error: fmt.Errorf("Expected '<' at the beggining, got '%c'", char),
      ExpectedChar: '<',
    }
  }

 if char := p.consumeChar(); char != '/'{  
    return node, &ParseError{
      error: fmt.Errorf("Expected '/' at the beggining, got '%c'", char),
      ExpectedChar: '/',
    }
  }

 if tag := p.parseTagName(); tag != tagName {  
    log.Println("Tag mistmach:", tag)
  }

 if char := p.consumeChar(); char != '>'{  
    return node, &ParseError{
      error: fmt.Errorf("Expected '>' at the beggining, got '%c'", char),
      ExpectedChar: '>',
    }
  }
  node = Elem(tagName, attrs, children)
  return node, nil
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

  if openQuote != '"' {
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

func (p *Parser) parseAttributes() AttrMap {
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

func (p *Parser) parseNodes() []Node {
  var nodes []Node 
  for {
    p.consumeWhiteSpace()
    if p.eof() || p.startsWith("</") {
      break
    }
    node, _ := p.parseNode()
    nodes = append(nodes, node)
  }
  return nodes 
}


