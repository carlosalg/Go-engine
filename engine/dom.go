package engine

type Node struct {
  //data common to all nodes:
  children []Node
  //data specific to each node type:
  node_type NodeType
}

type ElementData struct{
  tag_name string
  attributes AttrMap
}

type AttrMap map[string]string

//Equivalent to a enum:
type NodeType struct {
  Text string
  Element ElementData
}
func Text(data string) Node {
  return Node{
    children: []Node{},
    node_type: NodeType{Text: data},
  }
}

func Elem(name string, attrs AttrMap, children []Node) Node {
  return Node{
    children: children,
    node_type: NodeType{Element: ElementData{tag_name: name, attributes: attrs}},
  }
}



