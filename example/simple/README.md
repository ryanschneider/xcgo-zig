Simple CGO Example
==================

This folder includes a simple CGO example from: https://www.programming-books.io/essential/go/cgo-first-steps-tutorial-ed4cda13d7984045b85328a8d76211e5

(CC-BY-SA 3.0: https://creativecommons.org/licenses/by-sa/3.0)

To test `xgco-zig` try:

- Install `xgco-zig` using `go install github.com/ryanschneider/xcgo-zig`
- Build it natively and run it:

```
$ xcgo-zig build main.go
xcgo-zig: Running `go build main.go` with:
xcgo-zig: GOOS=
xcgo-zig: GOARCH=
xcgo-zig: CC=/opt/homebrew/bin/zig cc
xcgo-zig: CXX=/opt/homebrew/bin/zig c++

$ ./main
Hello world
Sum of 5 + 4 is 9
```

- Now cross-compile:

```
$ GOOS=linux GOARCH=amd64 xcgo-zig build -o main-linux-amd64 main.go
xcgo-zig: Running `go build -o main-linux-arm64 main.go` with:
xcgo-zig: GOOS=linux
xcgo-zig: GOARCH=amd64
xcgo-zig: CC=/opt/homebrew/bin/zig cc -target x86_64-linux-musl
xcgo-zig: CXX=/opt/homebrew/bin/zig c++ -target x86_64-linux-musl

$ file main-linux-amd64
main-linux-amd64: ELF 64-bit LSB executable, x86-64, version 1 (SYSV), static-pie linked, Go BuildID=l5KDublXNsr4_nsYvjoo/SF0v_NlOt4ZT55dIfo_f/JF49xhA9UvwOOdDFM7YI/BMGWJH47mI4YFLkXHId1, with debug_info, not stripped
```
