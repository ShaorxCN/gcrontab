package utils

type Node struct {
	children map[rune]*Node
	isEnd    bool
}

type Trie struct {
	root *Node
}

func NewTrie() *Trie {
	return &Trie{root: &Node{children: make(map[rune]*Node)}}
}

// Insert 插入语句
func (t *Trie) Insert(word string) {
	node := t.root
	for _, v := range word {
		if _, exists := node.children[v]; !exists {
			node.children[v] = &Node{children: make(map[rune]*Node)}
		}

		node = node.children[v]
	}

	node.isEnd = true
}

func (t *Trie) Search(word string) bool {
	node := t.root

	for _, v := range word {
		if _, exists := node.children[v]; !exists {
			return false
		}

		node = node.children[v]
	}

	return node.isEnd
}
