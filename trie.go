package trie

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

var (
	ErrAlreadyExists = errors.New("val already exists in trie")
	ErrNotFound      = errors.New("key not found in trie")
)

type Trie[T any] struct {
	Root *Node[T]
}

type Node[T any] struct {
	Value    T
	Children []*Node[T]
	KeyRune  rune
	IsEnd    bool
}

func (n Node[T]) String() string {
	var s strings.Builder
	for i := range n.Children {
		s.WriteRune(n.Children[i].KeyRune)
		s.WriteString(", ")
	}
	return fmt.Sprintf("Node key=%s Children=%s", string(n.KeyRune), s.String())
}

func NewTrie[T any]() *Trie[T] {
	return &Trie[T]{
		Root: &Node[T]{},
	}
}

// Operations

func (t *Trie[T]) Insert(key string, value T) error {
	// for each rune in string
	// check if rune exists in current node,
	// if yes, current node = next.node,
	// if not create new node with value rune
	parent := t.Root
	numChars := len(key)
	for i, character := range key {
		found := false
		for j := range parent.Children {
			if parent.Children[j].KeyRune == character {
				// the word already exists in the trie!
				if (i + 1) == numChars {
					slog.Debug("word exists", "val", key)
					return ErrAlreadyExists
				}
				slog.Debug("char exists", "char", string(character))
				// continue to this node's children nodes
				parent = parent.Children[j]
				found = true
				break
			}
		}
		if found {
			continue
		}
		// did not find character, so create a new node
		newNode := &Node[T]{
			Children: []*Node[T]{},
			KeyRune:  character,
			Value:    *new(T),
			IsEnd:    (i + 1) == numChars, // if we're the last character, set the IsEnd to true
		}
		// add it to the parent's childrent
		slog.Debug("creating new node", "before", parent.Children, "after", append(parent.Children, newNode))
		parent.Children = append(parent.Children, newNode)
		slog.Debug("parent v new", "parent", parent, "new", newNode)
		parent = newNode
		// final node, should set to value
		if i+1 == numChars {
			newNode.Value = value
			return nil
		}
	}
	return nil
}

func (t *Trie[T]) Search(key string) (T, error) {
	current := t.Root
	for i, char := range key {
		for _, node := range current.Children {
			if node.KeyRune == char {
				if node.IsEnd && (i+1) == len(key) {
					return node.Value, nil
				}
				current = node
				break
			}
		}
	}
	return *new(T), ErrNotFound
}

func (t *Trie[T]) Delete(key string) (T, error) {
	return *new(T), nil
}

func (t *Trie[T]) GetAll() []string {
	// return findAllKeys(t.Root.Children, []rune{})
	// Create a function that will accumulate all words in node
	fun := func(node *Node[T], key string, accumulator []string) []string {
		accumulator = append(accumulator, key)
		return accumulator
	}
	return DepthFirstSearch(t.Root.Children, []rune{}, fun, []string{})
}

func (t *Trie[T]) Clear() {
	// DFS over every node and delete it (mark node nil for GC)
	fun := func(nodes []*Node[T], i int, isLeafNode bool, key string, accumulator []string) []string {
		nodes[i] = nil
		return nil
	}
	depthFirstSearchEveryNode(t.Root.Children, []rune{}, fun, nil)
	// reset the root node to
	t.Root = &Node[T]{}
}

// DepthFirstSearch() traverses every node in the trie and calls leafNodeFun() when it reaches a leaf node.
// leafNodeFun() parameters are the leaf *Node, the key for this Node, and the accumulator which is a value that is passed to every Leaf Node.
// Accumulator allows DepthFirstSearch to perform an operations and return some value, such as count keys in trie.
func DepthFirstSearch[T, A any](nodes []*Node[T], keys []rune, leafNodeFun func(*Node[T], string, A) A, accumulator A) A {
	if len(nodes) == 0 {
		slog.Error("no children nodes, this means that there is an unterminated node", "accumulator", accumulator)
		return accumulator
	}

	for _, node := range nodes {
		slog.Debug("node", "val", node)
		if node.IsEnd {
			keys = append(keys, node.KeyRune)
			accumulator = leafNodeFun(node, string(keys), accumulator)
			continue
		}
		// keys will start to diverge, for next DFS iteration create a copy of keys array
		keysForThisTreePath := append(keys, node.KeyRune)
		// not reached the end, continue DFS to nodes children
		accumulator = DepthFirstSearch(node.Children, keysForThisTreePath, leafNodeFun, accumulator)
	}
	return accumulator
}

// depthFirstSearchEveryNode like depthFirstSearch but will call nodeFun on every node, in DFS order:
//   - wiil call on first enountered leaf node
//   - then call on parent's of leaf node until another leaf node is found
func depthFirstSearchEveryNode[T, A any](nodes []*Node[T], keys []rune, nodeFun func([]*Node[T], int, bool, string, A) A, accumulator A) A {
	if len(nodes) == 0 {
		slog.Debug("no children nodes, this means that there is an unterminated node", "accumulator", accumulator)
		return accumulator
	}

	for i, node := range nodes {
		slog.Debug("node", "val", node)
		if node.IsEnd {
			keys = append(keys, node.KeyRune)
			accumulator = nodeFun(nodes, i, true, string(keys), accumulator)
			continue
		}

		keysForThisTreePath := append(keys, node.KeyRune)
		accumulator = depthFirstSearchEveryNode(node.Children, keysForThisTreePath, nodeFun, accumulator)
		nodeFun(nodes, i, false, string(keys), accumulator)
	}
	return accumulator
}
