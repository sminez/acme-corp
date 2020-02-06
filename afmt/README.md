afmt - auto format your source files
====================================

In the spirit of `gofmt`, running `afmt` in an acme session will watch for
windows being saved. If a window's name matches a known file type then the
appropriate formatters are run on the window contents.


### Supported file types
- Python
  - [isort](https://github.com/timothycrosley/isort)
  - [Black](https://github.com/ambv/black)
  - [flake8](https://gitlab.com/pycqa/flake8)

- Golang
  - [gofmt](https://golang.org/cmd/gofmt/)
  - [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports)
  - [go vet](https://godoc.org/golang.org/x/tools/cmd/govet)

- Rust
  - [Rustfmt](https://github.com/rust-lang-nursery/rustfmt)

- Shell
  - Use [sed](https://en.wikipedia.org/wiki/Sed) to remove trailing whitespace.
  - [shellcheck](https://github.com/koalaman/shellcheck)

- Javascript
  - [js-beautify](https://github.com/beautify-web/js-beautify)
  - [jshint](https://github.com/jshint/jshint/)

- JSON
  - [python json.tool](https://docs.python.org/3.7/library/json.html#module-json.tool) [wrapper script](https://github.com/sminez/acme-corp/tree/master/scripts/json-format)

#### TODO
- Java
- Kotlin
- Haskell
- Switch JS formatting to prettier and also run for typescript files.

### Known bugs / current limitations
- Linter output ends up in the terminal that 'start-acme' was run from,
  not in the acme '+errors' window which makes more sense.
- Linter failures can sometimes crash the snooper it seems... (not sure
  if this really is the cause or if it is something else though)
