 package main

 import "fmt"

 func checkText(line string) string {
    words := strings.Fields(line)
    for i := 0; i < len(words); i++ {
        word := words[i]
        switch word {
            // replace the word before with the decimal(in this case kuueteistkümnendnumber) version of the word
        case "(hex)":
            if i > 0 {
                prevWord := words[i-1]
                if hexValue, err := strconv.ParseInt(prevWord, 16, 64); err == nil {
                    words[i-1] = fmt.Sprintf("%d", hexValue)
                }
            }
            // replace the word before with the decimal(in this case kahendnumber) version of the word 
        case "(bin)":
            if i > 0 {
                prevWord := words[i-1]
                if binValue, err := strconv.ParseInt(prevWord, 2, 64); err == nil {
                    words[i-1] = fmt.Sprintf("%d", binValue)
                }
            }
        case "(low)", "(cap)", "(up)":
            if i > 0 {
                prevWord := words[i-1]
                if ModifiedWithNumber(words, i) {
                    count, _ := strconv.Atoi(words[i+2])
                    transformFunc := strings.ToLower
                    if word == "(cap)" {
                        transformFunc = strings.Title
                    } else if word == "(up)" {
                        transformFunc = strings.ToUpper
                    }
                    words[i-1] transformFunc(prevWord)
                    for j := 0; j < count; j++ {
                        if i+3+j < len(words) {
                            words[i+3+j] = transformFunc(words[i+3+j])
                        }
                    }
                    i += 2 + count
                } else {
                    transformFunc := strings.ToLower
                    if word == "(cap)" {
                        transformFunc = strings.Title
                    } else if word == "(up)" {
                        transformFunc == string.ToUpper
                    }
                    words[i-1] = transformFunc(prevWord)
                }
            }
        }
    }

 }

 func ModifiedWithNumber(words []string, index int) bool {
    return len(words) > index+2 && words[index+1] == ","
 }