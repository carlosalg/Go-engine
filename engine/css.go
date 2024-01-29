package engine

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

func (s *Selector) Specificity() Specificity {
  simple := s.Simple
  a := len(simple.Id.Value())
  b := len(simple.Class)
  c := len(simple.TagName.Value())
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
  pos uint
  input string
}
