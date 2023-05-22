# xcgo-zig

Cross-compiler for CGO-based Golang projects using zig cc

## What is `xcgo-zig`?

`xcgo-zig` is, at the moment, a very basic proof-of-concept of a wrapper to make cross-compiling CGO binaries using `zig cc` easier.

It mainly eases this by inspecting `GOOS` and `GOARCH` and setting the appropriate `zig cc` targets.

So, instead of:

```bash
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC="zig cc -target x86_64-linux" CXX="zig c++ -target x86_64-linux" go build main.go
```

One can forgo the `-target` arguments and simply run:

```bash
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC="xcgo-zig cc" CXX="xcgo-zig c++" go build main.go
```

## Requirements

So far this has only been tested cross-compiling from an M1 macOS machine to linux, which required: 

- Golang 1.20.3 or higher (due to https://github.com/golang/go/issues/58935)
- Zig 0.11 or higher (I used `0.11.0-dev.3220+447a30299`)
  - `zig` must be found in your `$PATH` 
  - In the future, I'd like to support downloading and extracting a known-good version of `zig` to use instead.


