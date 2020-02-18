package main

import (
	"golang.org/x/tour/wc"
	"strings"
)

func WordCount(s string) map[string]int {

	dic := make(map[string]int)
	for _,w := range strings.Fields(s){
		dic[w]++
	}

	return dic
}

func main() {
	wc.Test(WordCount)
}

