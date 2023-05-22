Simple CGO Example
==================

This folder includes a simple CGO example from: https://www.programming-books.io/essential/go/cgo-first-steps-tutorial-ed4cda13d7984045b85328a8d76211e5

(CC-BY-SA 3.0: https://creativecommons.org/licenses/by-sa/3.0)

To test `xgco-zig` try:

- Install `xgco-zig` via `github.com/ryanschneider/xcgo-zig`
- Build it natively and run it:

```
$ CGO_ENABLED=1 CC="xcgo-zig cc" go build main.go
$ ./main
Hello world
Sum of 5 + 4 is 9
```

- Now cross-compile:

```
$ CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC="xcgo-zig cc" go build -o main-linux-amd64 main.go
$ file ./main-linux-amd64
./main-linux-amd64: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), static-pie linked, Go BuildID=H0cG-tJrTdCRM2FMGtH8/lPdkEAqpFI_yWSiUaouT/JF49xhA9UvwOOdDFM7YI/BMGWJH47mI4YFLkXHId1, with debug_info, not stripped
```
