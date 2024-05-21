package helpers

import (
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Args = struct {
	input string
}

type Test = struct {
	name string
	args Args
	want error
}

func runValidationTests(t *testing.T, validator func(string) error, tests []Test) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, validator(tt.args.input), tt.want)
		})
	}
}

func TestIsPathValidator(t *testing.T) {
	runValidationTests(t, IsPathValidator, []Test{
		{
			name: "Test path validator",
			args: Args{
				input: "./path/to/file",
			},
			want: nil,
		},
		{
			name: "Test path validator",
			args: Args{
				input: ".\\path\\to\\file",
			},
			want: nil,
		},
		{
			name: "Test path validator",
			args: Args{
				input: "./path/CON/file",
			},
			want: textinput.ErrInputValidation,
		},
		{
			name: "Test path validator with empty input",
			args: Args{
				input: "",
			},
			want: textinput.ErrInputValidation,
		},
	})
}
