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
 
