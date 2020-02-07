package snoop

// TODO: Parse file shebangs to determine filetype as a fallback.

import (
	"regexp"
	"strings"
)

var formatableTypes = []FileType{
	golang, python, shell, c, rust, javascript, json,
}

var golang = FileType{
	extensions: []string{"go"},
	Tools: []Tool{
		Tool{cmd: "goimports", args: []string{"-w"}},
		Tool{cmd: "go", args: []string{"vet"}},
	},
}

var python = FileType{
	extensions:   []string{"py", "pyw"},
	shebangProgs: []string{"python"},
	Tools: []Tool{
		Tool{
			cmd:          "isort",
			args:         []string{"-m", "5"},
			ignoreOutput: true,
		},
		Tool{cmd: "black", args: []string{"-q", "--line-length", "100"}},
		// Black is pep8 compliant but flake8 is not...
		Tool{cmd: "flake8", args: []string{"--ignore=E203,W503"}},
	},
}

var rust = FileType{
	extensions: []string{"rs"},
	Tools:      []Tool{Tool{cmd: "rustfmt"}},
}

var shell = FileType{
	extensions: []string{"sh", "bash", "zsh"},
	Tools: []Tool{
		// Remove trailing whitespace and whitespace only lines
		Tool{cmd: "sed", args: []string{"-i", "'s/[[:blank:]]*$//g'"}},
		Tool{cmd: "shellcheck", args: []string{"--color=never"}},
	},
}

var javascript = FileType{
	extensions: []string{"js"},
	Tools: []Tool{
		Tool{cmd: "js-beautify", args: []string{"-r"}},
		Tool{
			cmd: "jshint",
			outputFixer: func(s string) string {
				// Convert to button3 friendly output
				s = strings.Replace(s, " line ", "", -1)
				// Remove the column number as it just adds line noise
				re := regexp.MustCompile(", col .*,")
				return re.ReplaceAllString(s, " ->")
			},
		},
	},
}

var json = FileType{
	extensions: []string{"json"},
	Tools:      []Tool{Tool{cmd: "json-format"}},
}

var c = FileType{
	extensions: []string{"c", "h"},
	Tools: []Tool{
		Tool{cmd: "c-astyle"},
		Tool{
			cmd: "splint",
			args: []string{
				"+charintliteral", "+charint", "-exportlocal", "-compdef",
				"-usedef", "-retvalint", "+relaxtypes",
			},
		},
	},
}
