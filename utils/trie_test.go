package utils

import "testing"

func TestSearch(t *testing.T) {
	trie := NewTrie()

	trie.Insert("你好")
	trie.Insert("你好啊")
	trie.Insert("你在吗！")
	trie.Insert("are you ok?")

	tests := []struct {
		name string
		args string
		want bool
	}{{"不存在", "你", false}, {"存在", "你在吗！", true}, {"英文", "are you ok?", true}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if get := trie.Search(test.args); get != test.want {
				t.Errorf("Search()=%v,want %v", get, test.want)
			}
		})
	}

}
