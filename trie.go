package kin

import "strings"

type node struct {
	pattern	string
	part string
	children []*node
	isWild bool
}

func (n *node)matchChildByInsert(part string) *node  {
	for _, child := range n.children{
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}
func (n *node)matchChildrenByFind(part string) []*node {
	nodes := make([]*node,0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes,child)
		}
	}
	return nodes
}

func (n *node)matchChildByJudge(part string) ([]*node, bool) {
	nodes := make([]*node,0)
	for	_, child := range n.children{
		if	child.isWild {
			return nil, false
		}
		if child.part == part {
			nodes = append(nodes, child)
		}
	}
	return nodes, true
}

func (n *node)insert(patter string, parts []string, height int)  {
	if len(parts) == height {
		n.pattern = patter
		return
	}
	part := parts[height]
	child := n.matchChildByInsert(part)
	if child == nil{
		child = &node{part: part,isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children,child)
	}
	child.insert(patter,parts,height+1)
}

func (n *node)search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if	n.pattern == ""{
			return n
		}
	}

	part := parts[height]
	children :=	n.matchChildrenByFind(part)

	for _, child := range children{
		result := child.search(parts,height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

func (n *node)judge(parts []string, height int)(*node, bool){
	if len(parts) == height {
		return nil, false
	}
	part := parts[height]
	children, ok := n.matchChildByJudge(part)
	if !ok {
		return nil, false
	}
	for _, child := range children{
		result, ok := child.judge(parts,height+1)
		if !ok {
			return nil, false
		}
		if result != nil {
			return result,true
		}
	}
	return nil, false
}
