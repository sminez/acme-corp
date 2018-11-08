acmfmt - auto format your source files
======================================

In the spirit of gofmt, running acmefmt in an acme session will watch for
windows being saved. If a window's name matches a known file type then the
appropriate formatters are run on the window contents.


### Supported file types
- Python
  - Black
  - isort
  - flake8

- Golang
  - gofmt
  - goimports
  - go vet

##### TODO
- Rust
- Shell
- Javascript
- Java
- Kotlin
- Haskell
