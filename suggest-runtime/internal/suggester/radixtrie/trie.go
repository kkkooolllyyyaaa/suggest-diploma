package radixtrie

import (
	"sort"

	"suggest-runtime/internal/suggester"
)

const maxIndexSize = 500_000

// Node represents a node in the radix trie
type Node struct {
	Key      []rune         // The key segment stored at this node
	Children map[rune]*Node // Child nodes, indexed by their first rune
	Index    []int          // Indices of items stored in this node's subtree
}

// Trie is the radix trie data structure
type Trie struct {
	root          *Node                  // Root node of the trie
	index         []*suggester.IndexItem // All stored items
	indexSequence int                    // Sequence counter for items
	lastItemPtr   int                    // Pointer to the most recently added item
}

// NewTrie creates a new trie with the given configuration
func NewTrie() *Trie {
	t := &Trie{
		index: make([]*suggester.IndexItem, 0, 100_000),
	}
	return t
}

// Put adds a new item to the trie
func (t *Trie) Put(item *suggester.IndexItem) {
	// Store the item
	t.storeItem(item)
	query := t.getLastQuery()

	// Tokenize the query by spaces and insert each token
	tokens := t.tokenize(query)
	for _, token := range tokens {
		if len(token) > 0 {
			t.insert(token)
		}
	}
}

// tokenize splits a rune slice by space runes (U+0020)
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

// Get retrieves the node matching the given key
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

// IndexSize returns the number of items stored in the trie
func (t *Trie) IndexSize() int {
	return t.indexSequence
}

// findChild finds a child node matching the beginning of the query
func (t *Trie) findChild(node *Node, query []rune) (*Node, int) {
	if len(query) == 0 {
		return nil, 0
	}

	firstRune := query[0]
	child, exists := node.Children[firstRune]
	if !exists {
		return nil, 0
	}

	pos := commonPrefix(child.Key, query)
	return child, pos
}

// commonPrefix returns the length of the common prefix between two rune slices
func commonPrefix(a, b []rune) int {
	maxLen := min(len(a), len(b))
	for i := 0; i < maxLen; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return maxLen
}

// storeItem adds an item to the trie's index
func (t *Trie) storeItem(item *suggester.IndexItem) {
	t.index = append(t.index, item)
	t.lastItemPtr = t.indexSequence
	t.indexSequence++
}

// getLastQuery returns the query from the most recently added item
func (t *Trie) getLastQuery() []rune {
	return t.index[t.getLastIndexId()].Query
}

// getLastIndexId returns the index of the most recently added item
func (t *Trie) getLastIndexId() int {
	return t.lastItemPtr
}

// updateNodeIndex updates a node's index with the current item
func (t *Trie) updateNodeIndex(node *Node) {
	lastId := t.getLastIndexId()

	if false {
		// Insert maintaining sorted order by score
		insertPos := sort.Search(len(node.Index), func(i int) bool {
			return t.index[node.Index[i]].Score <= t.index[lastId].Score
		})

		// Insert at the found position
		node.Index = append(node.Index, 0)
		copy(node.Index[insertPos+1:], node.Index[insertPos:])
		node.Index[insertPos] = lastId

		// Trim if exceeds maximum size
		if len(node.Index) > maxIndexSize {
			node.Index = node.Index[:maxIndexSize]
		}
	} else {
		// Simple append
		node.Index = append(node.Index, lastId)
	}
}

// initRoot initializes the root node with the first key
func (t *Trie) initRoot(key []rune) {
	t.root = &Node{
		Children: make(map[rune]*Node),
	}

	newNode := t.newNode(key)
	t.root.Children[key[0]] = newNode
	t.updateNodeIndex(newNode)
}

// newNode creates a new node with the given key
func (t *Trie) newNode(key []rune) *Node {
	return &Node{
		Key:      key,
		Children: make(map[rune]*Node),
	}
}

// unpackNode splits a node to insert a new branch
func (t *Trie) unpackNode(node *Node, commonPrefix, nodeSuffix []rune) {
	// Create a bridge node with the suffix
	bridge := t.newNode(nodeSuffix)
	bridge.Children = node.Children
	bridge.Index = node.Index

	// Reset the original node
	node.Key = commonPrefix
	node.Children = make(map[rune]*Node)
	node.Children[nodeSuffix[0]] = bridge

	// Update the index
	t.updateNodeIndex(node)
}

// insert adds a query rune slice to the trie
func (t *Trie) insert(query []rune) {
	// Initialize root if needed
	if t.root == nil {
		t.initRoot(query)
		return
	}

	currentNode := t.root
	for currentNode != nil && len(query) > 0 {
		node, pos := t.findChild(currentNode, query)
		if node == nil {
			// No matching child - create a new node
			newNode := t.newNode(query)
			currentNode.Children[query[0]] = newNode
			t.updateNodeIndex(newNode)
			return
		}

		commonPrefix := query[:pos]
		querySuffix := query[pos:]
		nodeSuffix := node.Key[pos:]

		if len(nodeSuffix) == 0 {
			// Full match - move to child node
			currentNode = node
			t.updateNodeIndex(currentNode)
			query = querySuffix
			continue
		}

		// Partial match - split the node
		t.unpackNode(node, commonPrefix, nodeSuffix)
		currentNode = node
		query = querySuffix
	}
}
