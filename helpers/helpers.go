package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/textinput"
	"os"
	"regexp"
	"strings"
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

func Input(name string, value string, validator func(string) error, placeholder string, hide bool) string {
	if value == "" {
		input := textinput.New(name)
		input.Placeholder = placeholder
		input.Validate = validator
		input.Hidden = hide
		var e error
		value, e = input.RunPrompt()
		if e != nil {
			PrintErr(e.Error())
			os.Exit(1)
		}
	}
	if validator(value) != nil {
		PrintErr(fmt.Sprintf("Invalid %s: %s", name, value))
		os.Exit(1)
	}
	return value
}

func ParseBool(name string, value string, def bool) bool {
	if value == "" {
		return def
	}

	value = strings.ToLower(value)
	if value == "true" || value == "1" || value == "yes" {
		return true
	}
	if value == "false" || value == "0" || value == "no" {
		return false
	}

	PrintErr(fmt.Sprintf("Invalid %s: Can't parse \"%s\"", name, value))
	os.Exit(1)
	return false
}

func Confirmation(name string, value string, def *bool) bool {
	if value == "" {
		input := confirmation.New(name, def)
		decision, e := input.RunPrompt()
		if e != nil {
			PrintErr(e.Error())
			os.Exit(1)
		}
		return decision
	}
	return ParseBool(name, value, false)
}

func ParseDataSize(name string, size string, def uint64, max uint64) uint64 {
	if size == "" {
		return def
	}
	bytes, err := humanize.ParseBytes(size)
	if err != nil {
		PrintErr(fmt.Sprintf("Invalid %s: Can't parse \"%s\"", name, size))
		os.Exit(1)
	}
	if max != 0 && bytes > max {
		PrintErr(fmt.Sprintf("Invalid %s: Value \"%s\" is too large, max is \"%s\"", name, size, humanize.Bytes(max)))
		os.Exit(1)
	}
	return bytes
}

var separator = regexp.MustCompile(`,\s*`)

func ToJoinedValues(values []string) string {
	return strings.Join(values, ",")
}

func ParseValueArray(name string, value string, itemValidator func(string) error, def string) []string {
	if value == "" {
		value = def
	}
	var values []string
	for _, v := range separator.Split(value, -1) {
		if v == "" {
			continue
		}
		if itemValidator(v) != nil {
			PrintErr(fmt.Sprintf("Invalid %s: Can't parse \"%s\"", name, v))
			os.Exit(1)
		}
		values = append(values, v)
	}
	return values
}

func If[V any](condition bool, fn func(bool) V) V {
	return fn(condition)
}

func IfLazy[V any](condition bool, yes func() V, no func() V) V {
	if condition {
		return yes()
	} else {
		return no()
	}
}

func IsExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func Find[K comparable, V any](it map[K]V, fn func(V) bool) (*V, *K) {
	for k, v := range it {
		if fn(v) {
			return &v, &k
		}
	}
	return nil, nil
}
