package snoop

import (
	"regexp"
	"strings"
)

var formatableTypes = []FileType{
	golang, python, shell, rust, c, javascript, json,
}

var golang = FileType{
	extensions: []string{"go"},
	Tools: []Tool{
		Tool{cmd: "goimports", args: []string{"-w"}, appendFilePath: true},
		Tool{cmd: "golint", appendDirPath: true},
		Tool{cmd: "go", args: []string{"vet"}, appendDirPath: true},
	},
}

var python = FileType{
	extensions:   []string{"py", "pyw"},
	shebangProgs: []string{"python", "python3"},
	Tools: []Tool{
		Tool{
			cmd:            "isort",
			args:           []string{"-m", "5"},
			ignoreOutput:   true,
			appendFilePath: true,
		},
		Tool{cmd: "black", args: []string{"-q", "--line-length", "100"}, appendFilePath: true},
		// Black is pep8 compliant but flake8 is not...
		Tool{cmd: "flake8", args: []string{"--ignore=E203,E501,W503"}, appendFilePath: true},
	},
}

var rust = FileType{
	extensions: []string{"rs"},
	Tools: []Tool{
		Tool{cmd: "cargo", args: []string{"fmt", "--"}, appendFilePath: true},
		// Tool{cmd: "cargo", args: []string{"check"}},
		// Tool{cmd: "cargo", args: []string{"clippy"}},
	},
}

var shell = FileType{
	extensions:   []string{"sh", "bash", "zsh"},
	shebangProgs: []string{"bash", "sh", "zsh"},
	Tools: []Tool{
		// Remove trailing whitespace and whitespace only lines
		//Tool{cmd: "sed", args: []string{"-i", "s/[[:blank:]]*$//g"}, appendFilePath: true},
		Tool{cmd: "shellcheck", args: []string{"-f", "gcc"}, appendFilePath: true},
	},
}

var javascript = FileType{
	extensions:   []string{"js"},
	shebangProgs: []string{"node"},
	Tools: []Tool{
		Tool{cmd: "js-beautify", args: []string{"-r"}, appendFilePath: true},
		Tool{
			cmd:            "jshint",
			appendFilePath: true,
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
	Tools:      []Tool{Tool{cmd: "json-format", appendFilePath: true}},
}

var c = FileType{
	extensions: []string{"c", "h"},
	Tools: []Tool{
		Tool{cmd: "c-astyle", appendFilePath: true},
		Tool{
			cmd: "splint",
			args: []string{
				"+charintliteral", "+charint", "-exportlocal", "-compdef",
				"-usedef", "-retvalint", "+relaxtypes",
			},
			appendFilePath: true,
		},
	},
}
