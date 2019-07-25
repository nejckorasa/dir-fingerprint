# dir-fingerprint

![GitHub release](https://img.shields.io/github/release/nejckorasa/dir-fingerprint)
[![Tweet](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://twitter.com/intent/tweet?url=https%3A%2F%2Fgithub.com%2Fnejckorasa%2Fdir-fingerprint&via=nejckorasa&text=Generate%20directory%20fingerprint%20with%20all%20its%20files%20in%20a%20tree%20to%20observe%20changes&hashtags=golang%2Cgo%2Cscript%2Cfingerprint%2Cgithub%2Cdirectory%2Chash)

Application written in [Go](https://golang.org/) that creates and stores directory fingerprint from all its files in a tree. Fingerprint is stored in `.fingerprint` file in JSON format, file name can be configured.

Fingerprint file also indicates if *the fingerprint has changed* = *something in the directory tree was modified / added / removed*.

# Usage

### Install with Homebrew

```bash
$ brew install nejckorasa/tap/dir-fingerprint
```
> Tap can be found [here](https://github.com/nejckorasa/homebrew-tap)

Creeate fingerprint for current directory with:

```bash
$ dir-fingerprint .
```

### Usage

```
Usage of ./dir-fingerprint:
  -d	debug, turn on debug logging
  -f string
    	fingerprint file name (default ".fingerprint")
  -files
    	files, include all files fingerprints in fingerprint file, mind that there might me a lot of them
  -h	help, display usage
  -q	quiet, turn off logging, only print result
```

# Output

#### Default logging

```
[0000]  INFO Root 	/Users/nejckorasa/dir/that/needs/fingerprint
[0000]  INFO File	.fingerprint
[0000]  INFO Done	[45f73b](0.0004) 	@ .DS_Store
[0000]  INFO Done	[e3b0c4](0.0002) 	@ some_file
[0000]  INFO Done	[d6d2fb](0.0005) 	@ some.txt
[0000]  INFO Done	[841a5c](0.0003) 	@ another.txt
[0000]  INFO Done	[d59967](0.0046) 	@ my.pdf
[0000]  INFO Took	0.00794 sec
[0000]  INFO For	5 files
[0000]  INFO Skip	0 files

Old		[8a7b73f9671004edd50500bc7d3f1837d841a5c086011207259eb2d183823adf]
New		[8a7b73f9671004edd50500bc7d3f1837d841a5c086011207259eb2d183823adf]
@		/Users/nejckorasa/dir/that/needs/fingerprint/.fingerprint
Diff		false
```

Where:

- **Old** represents old directory fingerprint
- **New** represents new directory fingerprint
- **@** represents fingerprint file location
- **Diff** indicates if directory fingerprint has changed


And every file processing info line has the following format:

```
[time_so_far] INFO Done  [fingerprint_suffix](seconds_took)  @ <location>
``` 

#### Quiet logging
 
```
Old		[8a7b73f9671004edd50500bc7d3f1837d841a5c086011207259eb2d183823adf]
New		[8a7b73f9671004edd50500bc7d3f1837d841a5c086011207259eb2d183823adf]
@		/Users/nejckorasa/dir/that/needs/fingerprint/.fingerprint
Diff		false
```

# Fingerprint file example

```json
{
  "Fingerprint": "8a7b73f9671004edd50500bc7d3f1837d841a5c086011207259eb2d183823adf",
  "Changed": false
}
```

Where:

- **Fingerprint** represents directory fingerprint
- **Changed** indicates if the fingerprint has changed


#### Include files fingerprints

Additionally fingerprint file can be configured to include all files fingerprint, passing `-files` as an argument, mind that there might me a lot of them. Default is false.

```json
{
  "Fingerprint": "8a7b73f9671004edd50500bc7d3f1837d841a5c086011207259eb2d183823adf",
   "FilesFingerprints": [
      "45f73be5a62b4cb92820280c15aca3b57e6ad5c910e030d337e76d4491f5f549",
      "d599676b7e0922ed0dc5f82b106c1d0f3d395bf085168b79abf67fb8208e5110",
      "d6d2fb210c23dc0c318bcfa5091366ede674052d1a46b406b2b04c46803764b0",
      "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    ],
  "Changed": false
}
```
Where:
 
- **FilesFingerprints** represents all files fingerprints.

# TODO

- Stop processing files when fingerprint change is found?
