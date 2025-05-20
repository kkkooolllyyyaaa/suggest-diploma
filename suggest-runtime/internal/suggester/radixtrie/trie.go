package radixtrie

import (
	"suggest-runtime/internal/suggester"
)

const maxIndexSize = 500_000
const initialCapacity = 1_000_000

type Node struct {
	Key      []rune
	Children map[rune]*Node
	Index    []int
}

type Trie struct {
	root          *Node
	index         []*suggester.IndexItem
	indexSequence int
	lastItemPtr   int
}

func NewTrie() *Trie {
	return &Trie{index: make([]*suggester.IndexItem, 0, initialCapacity)}
}

func (t *Trie) Put(item *suggester.IndexItem) {
	t.storeItem(item)
	tokens := t.tokenize(t.getLastQuery())
	for _, token := range tokens {
		if len(token) > 0 {
			t.insert(token)
		}
	}
}

func (t *Trie) tokenize(query []rune) [][]rune {
	var tokens [][]rune
	start := 0
	for i, r := range query {
		if r == ' ' {
			if i > start {
				tokens = append(tokens, query[start:i])
			}
			start = i + 1
		}
	}
	if start < len(query) {
		tokens = append(tokens, query[start:])
	}
	return tokens
}

func (t *Trie) Get(key []rune) *Node {
	currentNode := t.root
	query := key

	for currentNode != nil && len(query) > 0 {
		node, pos := t.findChild(currentNode, query)
		if node == nil {
			return nil
		}

		query = query[pos:]
		if len(query) == 0 {
			return node
		}
		currentNode = node
	}
	return nil
}

func (t *Trie) IndexSize() int {
	return t.indexSequence
}

func (t *Trie) findChild(node *Node, query []rune) (*Node, int) {
	if len(query) == 0 {
		return nil, 0
	}

	child, exists := node.Children[query[0]]
	if !exists {
		return nil, 0
	}

	pos := commonPrefix(child.Key, query)
	return child, pos
}

func commonPrefix(a, b []rune) int {
	maxLen := min(len(a), len(b))
	for i := 0; i < maxLen; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return maxLen
}

func (t *Trie) storeItem(item *suggester.IndexItem) {
	t.index = append(t.index, item)
	t.lastItemPtr = t.indexSequence
	t.indexSequence++
}

func (t *Trie) getLastQuery() []rune {
	return t.index[t.getLastIndexId()].Query
}

func (t *Trie) getLastIndexId() int {
	return t.lastItemPtr
}

func (t *Trie) updateNodeIndex(node *Node) {
	node.Index = append(node.Index, t.getLastIndexId())
}

func (t *Trie) initRoot(key []rune) {
	t.root = &Node{Children: make(map[rune]*Node)}
	newNode := t.newNode(key)
	t.root.Children[key[0]] = newNode
	t.updateNodeIndex(newNode)
}

func (t *Trie) newNode(key []rune) *Node {
	return &Node{
		Key:      key,
		Children: make(map[rune]*Node),
	}
}

func (t *Trie) unpackNode(node *Node, commonPrefix, nodeSuffix []rune) {
	bridge := t.newNode(nodeSuffix)
	bridge.Children = node.Children
	bridge.Index = node.Index

	node.Key = commonPrefix
	node.Children = make(map[rune]*Node)
	node.Children[nodeSuffix[0]] = bridge

	t.updateNodeIndex(node)
}

func (t *Trie) insert(query []rune) {
	if t.root == nil {
		t.initRoot(query)
		return
	}

	currentNode := t.root
	for currentNode != nil && len(query) > 0 {
		node, pos := t.findChild(currentNode, query)
		if node == nil {
			newNode := t.newNode(query)
			currentNode.Children[query[0]] = newNode
			t.updateNodeIndex(newNode)
			return
		}

		commonPrefix := query[:pos]
		querySuffix := query[pos:]
		nodeSuffix := node.Key[pos:]

		if len(nodeSuffix) == 0 {
			currentNode = node
			t.updateNodeIndex(currentNode)
			query = querySuffix
			continue
		}

		t.unpackNode(node, commonPrefix, nodeSuffix)
		currentNode = node
		query = querySuffix
	}
}
