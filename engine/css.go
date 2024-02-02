package engine

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Stylesheet struct {
  rules []Rule 
}

type Rule struct {
  selectors []Selector
  declarations []Declaration 
}

type Selector struct {
  Simple SimpleSelector
}

type SimpleSelector struct {
  TagName Optional[string] 
  Id Optional[string]
  Class []string
}

type Declaration struct {
  name string
  value Value 
}

type Value struct{
  Keyword string
  Length struct {
    Value float32
    Unit Unit 
  }
}

type Unit string

const (
  Px Unit = "px"
)

type Color struct {
  r uint8
  g uint8
  b uint8
  a uint8
}

type Specificity struct {
  vol1 int
  vol2 int
  vol3 int
}

func (s Specificity) Less(other Specificity) bool {
  if s.vol1 != other.vol1{
    return s.vol1 > other.vol1
  } else if s.vol2 != other.vol2 {
    return s.vol2 > other.vol2
  } else {
    return s.vol3 > other.vol3
  }
}

func (s *Selector) Specificity() Specificity {
  simple := s.Simple
  a := 0 
  if simple.Id.HasValue() {
    a = len(simple.Id.Value())
  }
  b := len(simple.Class)
  c := 0
  if simple.TagName.HasValue() {
    c = len(simple.TagName.Value())
  }
  return Specificity{vol1: a,vol2: b, vol3: c}
}

func (v *Value) ToPx() float32 {
  if v.Length.Unit == Px {
    return v.Length.Value
  }
  return 0.0
}

func ParseCss(source string) Stylesheet {
  parser := ParserCss {pos: 0, input: source}
  return Stylesheet {rules: parser.parseRules()}
}

type ParserCss struct {
  pos int
  input string
}

func (p *ParserCss)nextChar() rune {
  if p.pos < len(p.input){
    char := rune(p.input[p.pos])
    return char 
  }
  return 0
}

func (p *ParserCss)eof() bool {
  return p.pos >= len(p.input)
}

func (p *ParserCss)consumeChar() rune {
  r, size := utf8.DecodeRuneInString(p.input[p.pos:])
  currChar := r
  nextPos := size
  p.pos += nextPos
  return currChar
}

func (p *ParserCss)consumeWhile(test func(rune)bool) string {
  var result strings.Builder
  for !p.eof() && test(p.nextChar()){
    result.WriteRune(p.consumeChar())
  }
  return result.String()
}

func (p *ParserCss)consumeWhiteSpace() {
  p.consumeWhile(unicode.IsSpace)
}

//Parsing methods

func (p *ParserCss) parseRules() []Rule {
  var rules []Rule
  for{
    p.consumeWhiteSpace()
    if p.eof() {break}
    rule := p.parseRule()
    rules = append(rules, rule)
  }
  return rules
}

func (p *ParserCss) parseRule() Rule {
  return Rule{selectors: p.parseSelectors(), declarations: p.parseDeclarations()}
}

func (p *ParserCss) parseSelectors() []Selector {
  var selectors []Selector
  for{
    selectors = append(selectors, Selector{Simple: p.parseSimpleSelector()})
    p.consumeWhiteSpace()
    if p.nextChar() == ',' {
        p.consumeChar()
        p.consumeWhiteSpace()
    } else if p.nextChar() == '{'{
      break
    }
    
  }
  sort.Slice(selectors,func(i, j int) bool{
    return selectors[j].Specificity().Less(selectors[i].Specificity()) ||

           (selectors[j].Specificity() == selectors[i].Specificity() && 
         selectors[j].Simple.TagName.Value() < selectors[i].Simple.TagName.Value())
  })
  return selectors
}

func (p *ParserCss) parseSimpleSelector() SimpleSelector{
  selector := SimpleSelector{
    TagName: NoValue[string]{}, 
    Id:      NoValue[string]{}, 
    Class:   []string{},
  }
  for !p.eof() {
    switch p.nextChar() {
    case '#': 
      p.consumeChar()
      selector.Id =  SomeValue[string]{value: p.parseIdentifier()}
    case '.':
      p.consumeChar()
      selector.Class = append(selector.Class, p.parseIdentifier())
    case '*':
      p.consumeChar()
    default:
      if validIdentifierChar(p.nextChar()){
        selector.TagName = SomeValue[string]{value: p.parseIdentifier()}
      } else {
        break
      }
    }
  }
  return selector
}

func(p *ParserCss) parseDeclaration() Declaration {
  propertyName := p.parseIdentifier()
  p.consumeWhiteSpace()
  if p.consumeChar() != ':'{
    err := fmt.Errorf("expected ':' found: '%c' ", p.consumeChar())
    //TODO
  }
}

