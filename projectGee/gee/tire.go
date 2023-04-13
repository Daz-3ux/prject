package gee

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string  // 待匹配路由,例如: /p/:lang; 是否是一个完整的 URL,不是则返回空字符串
	part     string  // URL 块值,用 / 分割的部分,例如 :lang
	children []*node // 该节点下的子节点,例如 [doc, turtorial, intro]
	isWild   bool    // 是否是模糊匹配
}

/*
两个基础函数
第一个 用于查找第一个匹配的子节点
第二个 用于查找所有匹配的子节点
*/
// 找到第一个匹配的子节点
// 在插入时使用
func (n *node) matchChild(part string) *node {
	// 遍历 n 节点的所有子节点,查询是否有匹配的子节点,将其返回
	for _, child := range n.children {
		// 如果有模糊匹配的也会匹配上
		if child.part == part || child.isWild {
			return child
		}
	}

	return nil
}

// 找到所有可能匹配的子节点
// 在查找时使用:必须返回所有可能的子节点进行遍历查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}

	return nodes
}

// 一边匹配一边插入
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		// 如果匹配结束,就将 pattern 赋值给该 node, 表示这是一个完整的 URL
		// 递归终止条件
		n.pattern = pattern
		return
	}

	// part 是当前节点的
	part := parts[height]
	// child 是当前节点的 *node 列表
	child := n.matchChild(part)
	fmt.Println("height: ", height, "part: ", part, "child: ", child)
	if child == nil {
		// 没有匹配上,就进行生成,放到 n 节点的自列表
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	// 接着插入下一个part节点
	child.insert(pattern, parts, height+1)
}

// 在路由树中查找与输入路径最匹配的节点，并返回该节点
func (n *node) search(parts []string, height int) *node {
	// 递归终止条件:找到 末尾/通配符
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			// pattern表示这不是一个完整的URL,匹配失败
			return nil
		}
		return n
	}

	// 获取当前高度下的 part
	part := parts[height]
	// 获取所有可能的子路径
	children := n.matchChildren(part)

	// 对每个匹配的子节点进行递归查找
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			// 找到了就返回
			return result
		}
	}

	return nil
}

// // 查找所有完整的 URL, 保存到列表
// func (n *node) travel(list *([]*node)) {
// 	if n.pattern != "" {
// 		// 递归终止条件
// 		*list = append(*list, n)
// 	}

// 	for _, child := range n.children {
// 		// 一层一层的递归找 pattern 是非空的节点
// 		child.travel(list)
// 	}
// }
