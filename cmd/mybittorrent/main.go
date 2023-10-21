package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

func trimFromList(bencodedString string) string {
	l := len(bencodedString)
	i := 0
	for l > 1 {
		if bencodedString[i] == 'l' && bencodedString[l-1] == 'e' {
			i++
			l--
			continue
		}
		break
	}
	// fmt.Println("26: ", bencodedString[i:l])
	return bencodedString[i:l]
}

func printType(v interface{}) {
	fmt.Println(reflect.TypeOf(v))
}

func formatResult(s interface{}) interface{} {
	return []interface{}{s}
}

func decodeBencode(bencodedString string) (interface{}, error) {
	nest := len(bencodedString) - len(trimFromList(bencodedString))
	bencodedString = trimFromList(bencodedString)
	if len(bencodedString) == 0 {
		return []interface{}{}, nil
	}

	if bencodedString[0] == 'i' {
		if strings.Contains(bencodedString, ":") {
			var strIdx int
			for idx := 0; idx < len(bencodedString); idx++ {
				if bencodedString[idx] == ':' {
					strIdx = idx
					break
				}
			}
			strValue := bencodedString[strIdx+1:]
			intValue, _ := strconv.Atoi(bencodedString[1 : strIdx-2])
			var result interface{}
			for i := 0; i < nest; i++ {
				result = formatResult([]interface{}{
					intValue,
					strValue,
				})
			}
			return result, nil
		}
		intValue, _ := strconv.Atoi(bencodedString[1 : len(bencodedString)-1])
		return intValue, nil
	} else if unicode.IsDigit(rune(bencodedString[0])) {
		pattern := `i[0-9\-]+e`
		if contains, _ := regexp.MatchString(pattern, bencodedString); contains {
			var strIdx, lengthOfStrValue int
			var err error
			for idx := 0; idx < len(bencodedString); idx++ {
				if bencodedString[idx] == ':' {
					strIdx = idx
					lengthOfStrValue, err = strconv.Atoi(bencodedString[:strIdx])
					if err != nil {
						return nil, err
					}
					break
				}
			}
			strValue := bencodedString[strIdx+1 : lengthOfStrValue+strIdx+1]
			intValue, err := strconv.Atoi(bencodedString[lengthOfStrValue+strIdx+1+1 : len(bencodedString)-1])
			if err != nil {
				return nil, err
			}
			var result interface{}
			for i := 0; i < nest; i++ {
				result = formatResult([]interface{}{
					strValue,
					intValue,
				})
			}
			return result, nil
		}
	}
	var strIdx int
	for idx := 0; idx < len(bencodedString); idx++ {
		if bencodedString[idx] == ':' {
			strIdx = idx
			break
		}
	}
	return bencodedString[strIdx+1:], nil
}

func main() {
	// fmt.Println("Logs from your program will appear here!")

	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]

		decoded, err := decodeBencode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
