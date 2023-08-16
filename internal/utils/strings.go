package utils

import "strings"

func GetLastWord(input string) string {
	words := strings.Split(input, " ")
	return words[len(words)-1]
}
