package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	wordArray []string
	flag      bool
)

func CheckBracket(pos int) int {
	switch wordArray[pos][:4] {
	case "(bin":
		ConvertNumber(pos-1, 2)
	case "(hex":
		ConvertNumber(pos-1, 16)
	case "(cap":
		ConvertString(pos, strings.Title)
	case "(up,", "(up)":
		ConvertString(pos, strings.ToUpper)
	case "(low":
		ConvertString(pos, strings.ToLower)
	default:
		pos++
	}
	return pos - 1
}

func ConvertNumber(pos, base int) {
	decimalValue, err := strconv.ParseInt(wordArray[pos], base, 64)
	if err != nil {
		fmt.Println("Error: Unable to convert to decimal.")
		os.Exit(1)
	}
	wordArray[pos] = strconv.Itoa(int(decimalValue))
	removeItems(pos+1, 1)
}

func ConvertString(pos int, transformFunc func(string) string) {
	if len(wordArray)-1 == pos {
		wordArray[pos-1] = transformFunc(wordArray[pos-1])
		removeItems(pos, 1)
		return
	}

	count, err := strconv.Atoi(wordArray[pos+1][:len(wordArray[pos+1])-1])
	if err == nil {
		for index := pos - 1; index >= pos-int(count); index-- {
			wordArray[index] = transformFunc(wordArray[index])
		}
		removeItems(pos, count)
	} else {
		wordArray[pos-1] = transformFunc(wordArray[pos-1])
		removeItems(pos, 1)
	}
}

func CorrectPunctuation(pos int) int {
	if strings.ContainsRune(".,!?:;", rune(wordArray[pos][0])) {
		wordArray[pos-1] += string(wordArray[pos][0])
		wordArray[pos] = wordArray[pos][1:]
		if len(wordArray[pos]) == 0 {
			removeItems(pos, 1)
			pos--
		}
		pos = CorrectPunctuation(pos)
	}
	return pos
}

func HandleArticles(pos int) {
	firstLetter := string(wordArray[pos+1][0])
	if strings.ContainsAny(firstLetter, "aeiouh") && (wordArray[pos] == "a" || wordArray[pos] == "A") {
		wordArray[pos] += "n"
	} else if !strings.ContainsAny(firstLetter, "aeiouh") && (wordArray[pos] == "an" || wordArray[pos] == "An") {
		wordArray[pos] = string(wordArray[pos][0])
	}
}

func AdjustCommas(pos int) {
	initial, size1 := utf8.DecodeRuneInString(wordArray[pos])
	final, size2 := utf8.DecodeLastRune([]byte(wordArray[pos]))

	if size1 != 3 || size2 != 3 {
		return
	}

	if flag {
		if string(final) == "‘" {
			wordArray[pos] = wordArray[pos][size2:] + "’"
		}
		if wordArray[pos] == "’" {
			wordArray[pos-1] += "’"
			removeItems(pos, 1)
		}
		flag = false
	} else {
		if string(initial) == "’" {
			wordArray[pos] = "‘" + wordArray[pos][size1:]
		}
		if wordArray[pos] == "‘" {
			wordArray[pos+1] = "‘" + wordArray[pos+1]
			removeItems(pos, 1)
		}
		flag = true
	}
}

func ParseContent() {
	fileContent, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error: File not found.")
		os.Exit(1)
	}
	wordArray = strings.Fields(string(fileContent))
}

func removeItems(index, count int) {
	wordArray = append(wordArray[:index], wordArray[index+count:]...)
}

func main() {
	ParseContent()
	for i := 0; i < len(wordArray); i++ {
		if len(wordArray[i]) > 3 {
			i = CheckBracket(i)
		}
		i = CorrectPunctuation(i)
		if i != len(wordArray)-1 {
			HandleArticles(i)
		}
		AdjustCommas(i)
	}
	os.WriteFile("result.txt", []byte(strings.Join(wordArray, " ")), 0644)
}
