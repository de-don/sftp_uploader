For working program need `xclip`. Install it
```bash
sudo apt-get install xclip
```

Fill credentials to sftp in `config.ini`

Compile go files `go build  -o sftp_uploader *.go` or use compiled binaries.

Run:

```bash
$ ./sftp_uploader # for upload image from clipboard buffer
$ ./sftp_uploader -r # for upload image from stdin
```

Example with [flameshot](https://github.com/lupoDharkael/flameshot):

```bash
$ flameshot gui -r  | ./sftp_uploader -r
```