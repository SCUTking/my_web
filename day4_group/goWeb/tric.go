package goWeb

import "strings"

type node struct {
	pattern  string  //该节点匹配的路由，查询时如果不为空，说明路由存在匹配
	part     string  //路由的当前部分
	children []*node //前缀树的子节点
	isWild   bool    //用于表示是否是精确匹配，动态路由的关键  有*的话，后续加入的
}

// 针对当前节点的插入
// pattern表示匹配的完整路径
// paths 表示分解的路径
// height 表示当前的长度 用于判断是否结束获取当前插入的part
func (n *node) insert(pattern string, paths []string, height int) {
	if len(paths) == height {
		n.pattern = pattern
		return
	}

	part := paths[height]
	child := n.matchChild(part)

	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}

	child.insert(pattern, paths, height+1)

}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "{" {
			return nil
		}
		return n
	}

	part := parts[height]
	childs := n.matchChilds(part)
	//如果有多条路由满足规则，按照插入的顺序选择第一个匹配的
	for _, child := range childs {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

// 用于插入 找到指定node的第一个匹配的节点————
// 如果不存在，返回空，根据是否为空判断是否创建新的；如果存在，直接返回
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		//或的使用：如果是动态的，走到这里，也表示匹配了。
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 用于查找  所有匹配成功的节点，用于查找
func (n *node) matchChilds(part string) []*node {
	res := make([]*node, 0)
	for _, child := range n.children {
		//或的使用：如果是动态的，走到这里，也表示匹配了。
		if child.part == part || child.isWild {
			res = append(res, child)
		}
	}
	return res
}
