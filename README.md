# xcgo-zig

Cross-compiler for CGO-based Golang projects using zig cc

## What is `xcgo-zig`?

`xcgo-zig` is, at the moment, a very basic proof-of-concept of a wrapper to make cross-compiling CGO binaries using `zig cc` easier.

It mainly eases this by inspecting `GOOS` and `GOARCH` and setting the appropriate `zig cc` targets.

So instead of running:

```bash
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC="zig cc -target x86_64-linux" CXX="zig c++ -target x86_64-linux" go build main.go
```

One can use `xcgo-zig` as a wrapper for `go` and simply run:

```bash
GOOS=linux GOARCH=amd64 xcgo-zig build main.go
```

Behind the scenes, `xcgo-zig` will build and inject the correct `CC` and `CXX` environment variables:

```
$ GOOS=linux GOARCH=amd64 xcgo-zig build main.go
xcgo-zig: Running `go build main.go` with:
xcgo-zig: GOOS=linux
xcgo-zig: GOARCH=amd64
xcgo-zig: CC=/opt/homebrew/bin/zig cc -target x86_64-linux-musl
xcgo-zig: CXX=/opt/homebrew/bin/zig c++ -target x86_64-linux-musl
```

By default, `xcgo-zig` will find and use the first `go` available in your `$PATH`.  However, this behavior can be changed with `GO=` environment variable:

```
$ GO=gotip xcgo-zig version
xcgo-zig: Running `gotip version` with:
xcgo-zig: GOOS=
xcgo-zig: GOARCH=
xcgo-zig: CC=/opt/homebrew/bin/zig cc
xcgo-zig: CXX=/opt/homebrew/bin/zig c++
go version devel go1.21-26f2569 Tue May 23 21:46:00 2023 +0000 darwin/arm64
```

## Requirements

So far this has only been tested cross-compiling from an M1 macOS machine to linux, which required: 

- Golang 1.20.3 or higher (due to https://github.com/golang/go/issues/58935)
- Zig 0.11 or higher (I used `0.11.0-dev.3220+447a30299`)
  - `zig` must be found in your `$PATH` 
  - In the future, I'd like to support downloading and extracting a known-good version of `zig` to use instead.


