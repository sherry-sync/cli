package helpers

import (
	"encoding/json"
	"fmt"
	flag "github.com/jessevdk/go-flags"
	"os"
)

func PrintErr(err string) {
	_, _ = fmt.Fprintf(os.Stderr, err)
	_, _ = fmt.Fprintf(os.Stderr, "\n")
}

func PrintMessage(msg string) {
	fmt.Println(msg)
}

func PrintJson(data interface{}) {
	dataJson, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(dataJson))
}

func AddCommand(cmd *flag.Command, command string, shortDescription string, longDescription string, data interface{}) {
	if _, err := cmd.AddCommand(
		command,
		shortDescription,
		longDescription,
		data,
	); err != nil {
		fmt.Println(err)
		return
	}
}
