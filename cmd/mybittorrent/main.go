package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/bittorrent-starter-go/cmd/bencode"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

func main() {
	// fmt.Println("Logs from your program will appear here!")

	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]

		sr := strings.NewReader(bencodedValue)

		decoded, err := bencode.Unmarshall(bufio.NewReader(sr))
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
