package trie

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

var ErrAlreadyExists = errors.New("val already exists in trie")

type Trie struct {
	Root *Node
}

type Node struct {
	Children []*Node
	Val      rune
	IsEnd    bool
}

func (n Node) String() string {
	var s strings.Builder
	for i := range n.Children {
		s.WriteRune(n.Children[i].Val)
		s.WriteString(", ")
	}
	return fmt.Sprintf("Node val=%s Children=%s", string(n.Val), s.String())
}

func NewTrie() *Trie {
	return &Trie{
		Root: &Node{},
	}
}

// Operations

func (t *Trie) Insert(val string) error {
	// for each rune in string
	// check if rune exists in current node,
	// if yes, current node = next.node,
	// if not create new node with value rune
	parent := t.Root
	numChars := len(val)
	for i, character := range val {
		found := false
		for j := range parent.Children {
			if parent.Children[j].Val == character {
				// the word already exists in the trie!
				if (i + 1) == numChars {
					slog.Debug("word exists", "val", val)
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
		newNode := &Node{
			[]*Node{},
			character,
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

func (t *Trie) Search(val string) error {
	return nil
}

func (t *Trie) Delete(val string) error {
	return nil
}

func (t *Trie) Dump() []string {
	// navigate DFS Trie until all words are found
	return DFS(t.Root.Children, []rune{})
}

func DFS(nodes []*Node, elements []rune) []string {
	foundValues := []string{}
	if len(nodes) == 0 {
		slog.Error("no children nodes, this means that there is an unterminated node", "runes", string(elements))
		return []string{string(elements)}
	}

	for _, node := range nodes {
		slog.Debug("iterating over node", "val", node)
		if node.IsEnd {
			// create a new word
			foundValues = append(foundValues, string(elements)+string(node.Val))
			slog.Debug("found end", "values", foundValues)
			continue
		}
		// not end, continue DFS to its children
		newValues := DFS(node.Children, append(elements, node.Val))
		foundValues = append(foundValues, newValues...)
	}
	return foundValues
}
