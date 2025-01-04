package trie

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

var ErrAlreadyExists = errors.New("val already exists in trie")

type Trie[T any] struct {
	Root *Node[T]
}

type Node[T any] struct {
	Children []*Node[T]
	KeyRune  rune
	Value    T
	IsEnd    bool
}

func (n Node[T]) String() string {
	var s strings.Builder
	for i := range n.Children {
		s.WriteRune(n.Children[i].KeyRune)
		s.WriteString(", ")
	}
	return fmt.Sprintf("Node val=%s Children=%s", string(n.KeyRune), s.String())
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
			[]*Node[T]{},
			character,
			*new(T),
			(i + 1) == numChars, // if we're the last character, set the IsEnd to true
		}
		// add it to the parent's childrent
		slog.Debug("creating new node", "before", parent.Children, "after", append(parent.Children, newNode))
		parent.Children = append(parent.Children, newNode)
		slog.Debug("parent v new", "parent", parent, "new", newNode)
		parent = newNode
	}
	return nil
}

func (t *Trie[T]) Search(key string) error {
	return nil
}

func (t *Trie[T]) Delete(key string) error {
	return nil
}

func (t *Trie[T]) Dump() []string {
	// navigate DFS Trie until all words are found
	return DFS(t.Root.Children, []rune{})
}

func DFS[T any](nodes []*Node[T], elements []rune) []string {
	foundValues := []string{}
	if len(nodes) == 0 {
		slog.Error("no children nodes, this means that there is an unterminated node", "runes", string(elements))
		return []string{string(elements)}
	}

	for _, node := range nodes {
		slog.Debug("iterating over node", "val", node)
		if node.IsEnd {
			// create a new word
			foundValues = append(foundValues, string(elements)+string(node.KeyRune))
			slog.Debug("found end", "values", foundValues)
			continue
		}
		// not end, continue DFS to its children
		newValues := DFS(node.Children, append(elements, node.KeyRune))
		foundValues = append(foundValues, newValues...)
	}
	return foundValues
}
