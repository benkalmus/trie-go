# Trie

A trie library implementation in Go using generics.

## TODO

- thread safe
- parallelised
- nil and error checking
- `[]byte` instead of `string` for keys
- benchmark
  - compare perf against a hashset
  - compare scalability

## Usage

Visualize a trie:

```go
trie := NewTrie[string]()
... // insert keys
fmt.Println(PrintTrie(trie))


├── c
|   └── a
|       └── a
|           ├── t*
|           ├── l
|           |   ├── m*
|           |   └── c*
|           |       ├── u*
|           |       └── r*
|           └── b*
|               └── l
|                   └── e*
└── a
    ├── s*
    |   └── k*
    └── t*

```
