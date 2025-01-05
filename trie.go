package trie

import (
	"errors"
	"fmt"
	"log/slog"
	"slices"
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
	return fmt.Sprintf("key='%s' val='%v' Children='%s'", string(n.KeyRune), n.Value, s.String())
}

func (t Trie[T]) String() string {
	return fmt.Sprintf("\n%s\n", PrintTrie(t.Root, "", 0, true))
}

func NewTrie[T any]() *Trie[T] {
	return &Trie[T]{
		Root: &Node[T]{},
	}
}

// Operations

func (t *Trie[T]) Insert(key string, value T) error {
	return insert(t.Root, []rune(key), value)
}

func insert[T any](node *Node[T], key []rune, value T) error {
	if len(key) == 0 {
		if node.IsEnd {
			return ErrAlreadyExists
		}
		node.IsEnd = true
		node.Value = value
		return nil
	}

	if node == nil {
		return ErrAlreadyExists
	}
	for i := range node.Children {
		if key[0] == node.Children[i].KeyRune {
			err := insert(node.Children[i], key[1:], value)
			return err
		}
	}

	newNode := &Node[T]{
		Children: []*Node[T]{},
		KeyRune:  key[0],
		Value:    *new(T),
		IsEnd:    false,
	}
	// slog.Debug("insert new node", "node", newNode, "isEnd", isTerminal)
	// slog.Debug("node.children", "children", node.Children)

	node.Children = append(node.Children, newNode)
	// have we created a terminal node? (last char)
	if len(key) == 1 {
		newNode.Value = value
		newNode.IsEnd = true
		return nil
	}
	return insert(newNode, key[1:], value)
}

func (t *Trie[T]) Search(key string) (T, error) {
	// TODO: do recursively and return index path
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
	val, _, err := deleteNode(t.Root, []rune(key))
	return val, err
}

func deleteNode[T any](node *Node[T], key []rune) (T, bool, error) {
	// found key
	if len(key) == 0 {
		node.IsEnd = false // this removes the termination marker. Key will no longer be found
		// If node is Terminal, we can safely delete it, return true
		if len(node.Children) == 0 {
			return node.Value, true, nil
		} else {
			// Has other children, so  this is just a substring of another key. don't delete
			return node.Value, false, nil
		}
	}
	// not found key
	keyRune := key[0] // take first char
	for i := range node.Children {
		if node.Children[i].KeyRune == keyRune {
			// DFS into subsequent children that match the key chars
			val, safeToDelete, err := deleteNode(node.Children[i], key[1:])
			if err != nil { // did not find key
				return *new(T), false, err
			}
			// key has been found. Can we safely delete it?
			// Node is safe to delete if it the key has no children. which was already determined
			if safeToDelete {
				node.Children[i] = nil
				node.Children = slices.Delete(node.Children, i, i+1)
			}
			// also delete current node if it doesn't have any siblings. This will cleanup all unterminated leafs
			if len(node.Children) < 1 && !node.IsEnd {
				return val, true, nil
			}
			return val, false, nil
		}
	}
	return *new(T), false, ErrNotFound
}

// test if a deleting help when hello exists removes

func (t *Trie[T]) GetAll() []string {
	// Create a function that will accumulate all words in trie
	fun := func(nodes **Node[T], key string, accumulator []string) []string {
		if (*nodes).IsEnd {
			return append(accumulator, key)
		}
		return accumulator
	}
	return depthFirstSearchEveryNode(t.Root.Children, []rune{}, fun, []string{})
}

func (t *Trie[T]) Clear() {
	// DFS over every node and delete it (mark node nil for GC)
	fun := func(nodes **Node[T], key string, accumulator []string) []string {
		nodes = nil
		return nil
	}
	depthFirstSearchEveryNode(t.Root.Children, []rune{}, fun, nil)
	// reset the root node to
	t.Root = &Node[T]{}
}

func countNodesBelow[T any](node *Node[T], mapping map[*Node[T]]int) int {
	// do a look up
	if val, ok := mapping[node]; ok {
		return val
	}

	num := len(node.Children)
	for i := range node.Children {
		num += countNodesBelow(node.Children[i], mapping)
	}
	// update map
	mapping[node] = num
	return num
}

// PrintTrie recursively prints the prefix tree in a structured format
// //TODO: tidy
func PrintTrie[T any](node *Node[T], prefix string, offset int, isLast bool) string {
	if node == nil {
		return ""
	}
	str := ""

	// Print the current node
	if node.KeyRune != 0 {
		offset += 4
		if isLast {

			str += fmt.Sprintf("%s└── %c", prefix, node.KeyRune)
			for i := 0; i < 4; i++ {
				if i%offset == 4 {
					prefix += "|"
					continue
				}
				prefix += " "
			}
		} else {
			str += fmt.Sprintf("%s├── %c", prefix, node.KeyRune)
			// prefix = ""
			for i := 0; i < 4; i++ {
				if i%offset == 0 {
					prefix += "|"
					continue
				}
				prefix += " "
			}
			// prefix += "│   "
		}
	}
	if node.IsEnd {
		str += "*"
	}
	str += "\n"

	// Recursively print the children
	for i, child := range node.Children {
		str += PrintTrie(child, prefix, offset, i == len(node.Children)-1)
	}
	return str
}

func leftPad(amount int, char rune) string {
	var s strings.Builder
	for i := 0; i < amount; i++ {
		s.WriteRune(char)
	}
	return s.String()
}

// DepthFirstSearchWord() traverses every node in the trie and calls endNodeFun() when it reaches a end node, that is a key.
// endNodeFun() parameters are the end *Node, the key for this Node, and the accumulator which is a value that is passed to every end Node.
// Accumulator allows DepthFirstSearchWord to perform an operations and return some value, such as count keys in trie.
func DepthFirstSearchWord[T, A any](nodes []*Node[T], keys []rune, endNodeFun func(*Node[T], string, A) A, accumulator A) A {
	if len(nodes) == 0 {
		// slog.Debug("no children nodes, reached end of subtree", "accumulator", accumulator)
		return accumulator
	}

	for _, node := range nodes {
		slog.Debug("node", "val", node)
		keys := append(keys, node.KeyRune)

		if node.IsEnd {
			// keys = append(keys, node.KeyRune)
			accumulator = endNodeFun(node, string(keys), accumulator)
		}
		// continue DFS to this node's children
		accumulator = DepthFirstSearchWord(node.Children, keys, endNodeFun, accumulator)
	}
	return accumulator
}

// depthFirstSearchEveryNode like depthFirstSearch but will call nodeFun on every node, in DFS order:
//   - wiil call on first enountered end node
//   - then call on parent's of leaf node until another leaf node is found
//     NOTE: to be able to edit the node passed into nodeFun, a pointer is passed to the original pointer.
//     This is because in Go, pointers are passed by value (creates another pointer that points to the original object)
//     This will allow the caller to modify the node in any way, even setting it to nil
//     TODO: should add some testcoverage around this
//   - Alternate solution: return a new node from nodeFun which can be set back to the slice
func depthFirstSearchEveryNode[T, A any](nodes []*Node[T], keys []rune, nodeFun func(**Node[T], string, A) A, accumulator A) A {
	if len(nodes) == 0 {
		return accumulator
	}

	for i := range nodes {
		slog.Debug("node", "val", nodes[i])

		keys := append(keys, nodes[i].KeyRune)
		accumulator = depthFirstSearchEveryNode(nodes[i].Children, keys, nodeFun, accumulator)
		accumulator = nodeFun(&nodes[i], string(keys), accumulator)
	}
	return accumulator
}

func breadthFirstSearch[T, A any](queue []*Node[T], keys []rune, nodeFun func([]*Node[T], int, int, string, A) A, accumulator A) ([]*Node[T], A) {
	currentQueueLen := len(queue)
	currentLevel := 0
	for i := 0; i < len(queue); i++ {
		if i == currentQueueLen {
			currentQueueLen = len(queue)
			currentLevel++
		}
		acc := nodeFun(queue, i, currentLevel, string(keys), accumulator)
		accumulator = acc
		queue = append(queue, queue[i].Children...)
	}

	return queue, accumulator
}
