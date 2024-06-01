package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	ConsoleReset = 0

	ConsoleFgBlack        = 30
	ConsoleFgDarkRed      = 31
	ConsoleFgDarkGreen    = 32
	ConsoleFgDarkYellow   = 33
	ConsoleFgDarkBlue     = 34
	ConsoleFgDarkMagenta  = 35
	ConsoleFgDarkCyan     = 36
	ConsoleFgLightGray    = 37
	ConsoleFgDarkGray     = 90
	ConsoleFgLightRed     = 91
	ConsoleFgLightGreen   = 92
	ConsoleFgLightYellow  = 93
	ConsoleFgLightBlue    = 94
	ConsoleFgLightMagenta = 95
	ConsoleFgLightCyan    = 96
	ConsoleFgWhite        = 97

	ConsoleBgBlack        = 40
	ConsoleBgDarkRed      = 41
	ConsoleBgDarkGreen    = 42
	ConsoleBgDarkYellow   = 43
	ConsoleBgDarkBlue     = 44
	ConsoleBgDarkMagenta  = 45
	ConsoleBgDarkCyan     = 46
	ConsoleBgLightGray    = 47
	ConsoleBgDarkGray     = 100
	ConsoleBgLightRed     = 101
	ConsoleBgLightGreen   = 102
	ConsoleBgLightYellow  = 103
	ConsoleBgLightBlue    = 104
	ConsoleBgLightMagenta = 105
	ConsoleBgLightCyan    = 106
	ConsoleBgWhite        = 107

	ConsoleBold      = 1
	ConsoleUnderline = 4
	ConsoleFramed    = 54
)

func GetConsoleColor(effects []int) string {
	var codes []string
	for _, i := range effects {
		codes = append(codes, strconv.Itoa(i))
	}

	return fmt.Sprintf("\033[%sm", strings.Join(codes, ";"))
}

func WithColor(effects []int, text string) string {
	return fmt.Sprintf("%s%s%s", GetConsoleColor(effects), text, GetConsoleColor([]int{ConsoleReset}))
}

func ToTitle(v string) string {
	return cases.Title(language.AmericanEnglish).String(v)
}

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

func PrintMap(data interface{}, name string, omit []string) {
	var jsonMap map[string]interface{}
	v, e := json.Marshal(data)
	if e != nil {
		panic(e)
	}
	err := json.Unmarshal(v, &jsonMap)
	if err != nil {
		panic(e)
	}

	fmt.Println()
	fmt.Println(fmt.Sprintf("%s:", name))
	for k, v := range jsonMap {
		if Includes(omit, k) {
			continue
		}
		fmt.Println(fmt.Sprintf("  %s: %s", ToTitle(k), v))
	}
	fmt.Println()
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

func Select(name string, value string, options []string) string {
	if value == "" {
		selector := selection.New(name, options)
		selector.PageSize = len(options)
		var e error
		value, e = selector.RunPrompt()
		if e != nil {
			PrintErr(e.Error())
			os.Exit(1)
		}
	}
	if !Includes(options, value) {
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
	values := []string{}
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

func If[V any](condition bool, yes func() V, no func() V) V {
	if condition {
		return yes()
	} else {
		return no()
	}
}

func EmptyIfNull[T any](arr []T) []T {
	if len(arr) == 0 {
		return []T{}
	}
	return arr
}

func IsExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func Find[T any](ts []T, fn func(T) bool) *T {
	for _, t := range ts {
		if fn(t) {
			return &t
		}
	}
	return nil
}

func Includes[T comparable](ts []T, item T) bool {
	for _, t := range ts {
		if t == item {
			return true
		}
	}
	return false
}

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func NormalizePath(path string) string {
	return strings.ReplaceAll(filepath.Clean(path), "\\", "/")
}

func PreparePath(path string) string {
	abs, e := filepath.Abs(path)
	if e != nil {
		PrintErr(fmt.Sprintf("Invalid path: %s: %s", e.Error(), path))
	}
	return NormalizePath(abs)
}

func GetPathParts(path string) []string {
	return strings.Split(NormalizePath(path), "/")
}

// IsChildPath accepts two normalized paths
func IsChildPath(parent, child string) (bool, error) {
	if parent == child {
		return true, nil
	}

	// path-comparisons using filepath.Abs don't work reliably according to docs (no unique representation).
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false, err
	}
	if !strings.HasPrefix(rel, ".."+string(os.PathSeparator)) && rel != ".." {
		return true, nil
	}
	return false, nil
}

func Filter[T any](ts []T, fn func(T) bool) []T {
	var res []T
	for _, v := range ts {
		if fn(v) {
			res = append(res, v)
		}
	}
	return res
}
