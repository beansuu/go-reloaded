package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func transformWords(words []string, transformFunc func(string) string) string {
	for i := range words {
		words[i] = transformFunc(words[i])
	}
	return strings.Join(words, " ")
}

func transformAndReplace(match string, submatches []string, transformFunc func(string) string) string {
	count, _ := strconv.Atoi(submatches[2])
	if count == 0 {
		count = 1
	}
	words := strings.Fields(submatches[1])
	wordsTransformed := transformWords(words[len(words)-count:], transformFunc)
	return strings.Join(append(words[:len(words)-count], wordsTransformed), " ")
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run . <input_file> <output_file>")
		return
	}

	inputFile := os.Args[1]
	inputBytes, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	input := string(inputBytes)

	transformFunctions := []struct {
		pattern *regexp.Regexp
		replace func(string) string
	}{
		{regexp.MustCompile(`([0-9A-Fa-f]+) \(hex\)`), func(match string) string {
			hexStr := regexp.MustCompile(`([0-9A-Fa-f]+) \(hex\)`).FindStringSubmatch(match)[1]
			decimalVal, _ := strconv.ParseInt(hexStr, 16, 64)
			return strconv.Itoa(int(decimalVal))
		}},
		{regexp.MustCompile(`([01]+) \(bin\)`), func(match string) string {
			binStr := regexp.MustCompile(`([01]+) \(bin\)`).FindStringSubmatch(match)[1]
			decimalVal, _ := strconv.ParseInt(binStr, 2, 64)
			return strconv.Itoa(int(decimalVal))
		}},
		{regexp.MustCompile(`(\w+(?:\s+\w+)*) \(up(?:, (\d+))?\)`), func(match string) string {
			submatches := regexp.MustCompile(`(\w+(?:\s+\w+)*) \(up(?:, (\d+))?\)`).FindStringSubmatch(match)
			return transformAndReplace(match, submatches, strings.ToUpper)
		}},
		{regexp.MustCompile(`(\w+(?:\s+\w+)*) \(low(?:, (\d+))?\)`), func(match string) string {
			submatches := regexp.MustCompile(`(\w+(?:\s+\w+)*) \(low(?:, (\d+))?\)`).FindStringSubmatch(match)
			return transformAndReplace(match, submatches, strings.ToLower)
		}},
		{regexp.MustCompile(`(\w+(?:\s+\w+)*) \(cap(?:, (\d+))?\)`), func(match string) string {
			submatches := regexp.MustCompile(`(\w+(?:\s+\w+)*) \(cap(?:, (\d+))?\)`).FindStringSubmatch(match)
			return transformAndReplace(match, submatches, strings.Title)
		}},
		{regexp.MustCompile(`'\s*([^']*)\s*'`), func(match string) string {
			trimmed := strings.TrimSpace(regexp.MustCompile(`'\s*([^']*)\s*'`).FindStringSubmatch(match)[1])
			return "'" + trimmed + "'"
		}},
		{regexp.MustCompile(`\ba\b\s+([aeiouhAEIOUH])`), func(match string) string {
			return "an " + strings.TrimSpace(match[2:])
		}},
	}

	for _, tf := range transformFunctions {
		input = tf.pattern.ReplaceAllStringFunc(input, tf.replace)
	}

	// Additional cleanup for punctuation
	cleanupPatterns := []string{` :`, ` \.`, ` \?`, ` \!`}
	for _, pattern := range cleanupPatterns {
		re := regexp.MustCompile(pattern)
		input = re.ReplaceAllStringFunc(input, func(match string) string {
			return strings.TrimSpace(match)
		})
	}

	// Remove space before comma
	input = regexp.MustCompile(` ,`).ReplaceAllString(input, ",")

	// Add space after each comma if it's not already there
	input = regexp.MustCompile(`,([^ ])`).ReplaceAllString(input, ", $1")

	// Write to output file
	outputFile := os.Args[2]
	err = ioutil.WriteFile(outputFile, []byte(strings.TrimSpace(input)), 0644)
	if err != nil {
		fmt.Println("Error writing output file:", err)
		return
	}
}
