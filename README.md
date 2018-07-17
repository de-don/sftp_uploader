For working program need `xclip`. Install it
```bash
sudo apt-get install xclip
```

Fill credentials to sftp in `config.ini`

Compile go files `go build *.go` or use compiled binaries.

Run:

```bash
$ ./clipboard_utils # for upload image from clipboard buffer
$ ./clipboard_utils -r # for upload image from stdin
```

Example with [flameshot](https://github.com/lupoDharkael/flameshot):

```bash
$ flameshot gui -r  | ./clipboard_utils -r
```