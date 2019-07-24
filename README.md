# dir-fingerprint

Application written in [Go](https://golang.org/) that creates and stores directory fingerprint from all its files in a tree. Fingerprint is stored in `.fingerprint` file in JSON format, file name can be configured.

Fingerprint file also indicates if *the fingerprint has changed* = *something in the directory tree was modified/added/removed*.

# Usage

```bash
Usage of ./dir-fingerprint:
  -d	debug, turn on debug logging
  -f string
    	fingerprint file name (default ".fingerprint")
  -files
    	files, include all files fingerprints in fingerprint file, mind that there might me a lot of them
  -h	help, display usage
  -q	quiet, turn off logging, only print result
```
# TODO

## Impl
- Stop processing files when fingerprint change is found?
- Add to homebrew tap

## Docs
- Add examples
- Explain usage
- Add output example
