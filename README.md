Fill `config.ini`

Compile go files `go build *.go`

Run:

```bash
$ ./clipboard_utils # for upload image from clipboard buffer
$ ./clipboard_utils -r # for upload image from stdin
```

Example with flameshot:

```bash
$ flameshot gui -r  | ./clipboard_utils -r
```