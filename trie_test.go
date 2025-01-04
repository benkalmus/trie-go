package trie

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrieInsert(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger) // Set this logger as the default for tests

	t.Run("insert one word into trie", func(t *testing.T) {
		trie := NewTrie[string]()
		word := "hello"
		err := trie.Insert(word, "")
		assert.Equal(t, nil, err, "expected no errors on insert")

		values := trie.Dump()
		expected := []string{word}
		assert.Equal(t, expected, values, "expected 'hello' in trie")
	})
	t.Run("insert multiple but different words", func(t *testing.T) {
		trie := NewTrie[string]()
		word := "hello"
		word2 := "world"
		err := trie.Insert(word, "")
		assert.Equal(t, nil, err, "expected no errors on insert")
		err = trie.Insert(word2, "")
		assert.Equal(t, nil, err, "expected no errors on insert")

		values := trie.Dump()
		expected := []string{word, word2}
		assert.Equal(t, expected, values)
	})
	t.Run("insert multiple but similar words", func(t *testing.T) {
		trie := NewTrie[string]()
		word := "hello"
		word2 := "help"
		err := trie.Insert(word, "")
		assert.Equal(t, nil, err, "expected no errors on insert")
		err = trie.Insert(word2, "")
		assert.Equal(t, nil, err, "expected no errors on insert")

		values := trie.Dump()
		expected := []string{word, word2}
		assert.Equal(t, expected, values)
	})
	t.Run("insert the same word returns error", func(t *testing.T) {
		trie := NewTrie[string]()
		word := "hello"
		word2 := word
		err := trie.Insert(word, "")
		assert.Equal(t, nil, err, "expected no errors on insert")
		err = trie.Insert(word2, "")
		assert.Equal(t, ErrAlreadyExists, err, "expected an error inserting the same word")

		values := trie.Dump()
		expected := []string{word}
		assert.Equal(t, expected, values)
	})
}

func TestTrieSearch(t *testing.T) {
	t.Run("find key and fetch value", func(t *testing.T) {
		trie := NewTrie[string]()
		word := "hello"
		val := "ok"
		err := trie.Insert(word, val)
		assert.Equal(t, nil, err, "expected no errors on insert")

		got, err := trie.Search(word)
		assert.Equal(t, nil, err)
		assert.Equal(t, val, got)
	})
	t.Run("if key does not exist, return err", func(t *testing.T) {
		trie := NewTrie[string]()
		word := "hello"
		search := "hello2"
		val := "ok"
		err := trie.Insert(word, val)
		assert.Equal(t, nil, err, "expected no errors on insert")

		got, err := trie.Search(search)
		assert.Equal(t, ErrNotFound, err)
		assert.Equal(t, "", got)
	})
}
