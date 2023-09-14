package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "strconv"
    "unicode"
)

func checkText(line string) string {
    words := strings.Fields(line)

    for i := 0; i < len(words); i++ {
        word := words[i]
        switch word {
        case "(hex)":
            if i > 0 {
                prevWord := words[i-1]
                if hexValue, err := strconv.ParseInt(prevWord, 16, 64); err == nil {
                    words[i-1] = fmt.Sprintf("%d", hexValue)
                    // Remove the (hex) marker
                    words = append(words[:i], words[i+1:]...)
                }
            }
        case "(bin)":
            if i > 0 {
                prevWord := words[i-1]
                if binValue, err := strconv.ParseInt(prevWord, 2, 64); err == nil {
                    words[i-1] = fmt.Sprintf("%d", binValue)
                    // Remove the (bin) marker
                    words = append(words[:i], words[i+1:]...)
                }
            }
        case "(low)", "(cap)", "(up)":
            if i > 0 {
                prevWord := words[i-1]
                if modifiedWithNumber(words, i) {
                    count, _ := strconv.Atoi(words[i+2])
                    transformFunc := strings.ToLower
                    if word == "(cap)" {
                        transformFunc = strings.Title
                    } else if word == "(up)" {
                        transformFunc = strings.ToUpper
                    }
                    words[i-1] = transformFunc(prevWord)
                    for j := 0; j < count; j++ {
                        if i+3+j < len(words) {
                            words[i+3+j] = transformFunc(words[i+3+j])
                        }
                    }
                    // Remove the marker and count
                    words = append(words[:i], words[i+3+count:]...)
                } else {
                    transformFunc := strings.ToLower
                    if word == "(cap)" {
                        transformFunc = strings.Title
                    } else if word == "(up)" {
                        transformFunc = strings.ToUpper
                    }
                    words[i-1] = transformFunc(prevWord)
                    // Remove the marker
                    words = append(words[:i], words[i+1:]...)
                }
            }
        }
    }
    return strings.Join(words, " ")
}

func modifiedWithNumber(words []string, index int) bool {
    return len(words) > index+2 && words[index+1] == ","
}

func formatPunctuation(text string) string {
    text = strings.ReplaceAll(text, " ,", ",")
    text = strings.ReplaceAll(text, " .", ".")
    text = strings.ReplaceAll(text, " !", "!")
    text = strings.ReplaceAll(text, " ?", "?")
    text = strings.ReplaceAll(text, " :", ":")
    text = strings.ReplaceAll(text, " ;", ";")

    text = strings.ReplaceAll(text, "...", "...")
    text = strings.ReplaceAll(text, "!?", "!?")
    return text
}

func handleAAn(text string) string {
    words := strings.Fields(text)
    for i := 0; i < len(words)-1; i++ {
        word := words[i]
        nextWord := words[i+1]
        if word == "a" && (startsWithVowel(nextWord) || strings.HasPrefix(nextWord, "h")) {
            words[i] = "an"
        }
    }
    return strings.Join(words, " ")
}

// vowel = täishäälik
func startsWithVowel(s string) bool {
    firstChar := []rune(s)[0]
    return unicode.Is(unicode.Latin, firstChar) &&
        (firstChar == 'a' || firstChar == 'e' || firstChar == 'i' || firstChar == 'o' || firstChar == 'u')
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: go run main.go sample.txt result.txt")
        os.Exit(1)
    }
    sampleTxt := os.Args[1]
    resultTxt := os.Args[2]

    inputFile, err := os.Open(sampleTxt)
    if err != nil {
        fmt.Println("Error opening input file:", err)
        os.Exit(1)
    }
    defer inputFile.Close()

    outputFile, err := os.Create(resultTxt)
    if err != nil {
        fmt.Println("Error creating output file:", err)
        os.Exit(1)
    }
    defer outputFile.Close()

    scanner := bufio.NewScanner(inputFile)

    for scanner.Scan() {
        line := scanner.Text()
        modifiedLine := checkText(line)
        fmt.Fprintln(outputFile, modifiedLine)
    }
    fmt.Println("Editing complete. Output written to", resultTxt)
}
