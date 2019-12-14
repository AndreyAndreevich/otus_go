package app

import (
	"sort"
	"strings"

	"github.com/jwangsadinata/go-multimap/slicemultimap"
)

func Top10(text string) []string {
	res := make([]string, 0, 10)

	if len(text) == 0 {
		return res
	}

	wordMap := make(map[string]int)
	top := slicemultimap.New()

	words := strings.Split(text, " ")

	for _, word := range words {
		count := wordMap[word] + 1
		wordMap[word] = count
		top.Remove(count-1, word)
		top.Put(count, word)
	}

	keys := top.KeySet()
	intKeys := make([]int, 0, len(keys))
	for _, key := range keys {
		intKeys = append(intKeys, key.(int))
	}

	sort.Sort(sort.Reverse(sort.IntSlice(intKeys)))

	count := 10

	for _, key := range intKeys {
		wordSlice, _ := top.Get(key)
		if len(wordSlice) > count {
			for i := 0; i < count; i++ {
				res = append(res, wordSlice[i].(string))
			}
			break
		}
		for _, word := range wordSlice {
			res = append(res, word.(string))
		}
		count -= len(wordSlice)
	}

	return res
}
