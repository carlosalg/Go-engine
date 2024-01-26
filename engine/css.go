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
  TagName 
  Id
  Class []string
}
