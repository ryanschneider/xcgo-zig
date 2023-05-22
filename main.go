package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

func main() {
	// check if go is 1.20.3 or higher (older versions broken)
	if err := checkGoVersion(); err != nil {
		panic(err)
	}

	goos, _ := os.LookupEnv("GOOS")
	goarch, _ := os.LookupEnv("GOARCH")
	zigTarget, err := lookupTarget(goos, goarch)

	env := append(os.Environ())
	args := append(os.Args)
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// replace our first argument w/ fully qualified path to zig
	arg0, err := exec.LookPath("zig")
	if err != nil {
		// We couldn't find zig, exit
		// TODO: download and install zig for this OS to .cache or other appropriate location
		panic(err)
	}

	if len(args) > 0 {
		args[0] = arg0
	}

	if len(zigTarget) > 0 {
		args = append(args, zigTarget...)
	}

	// don't write to stderr it seems to break go build on darwin
	// TODO: Setup logging to a temp file instead of disabling it here
	if false {
		log.Printf("exec: (pwd=%s) %v", pwd, args)
	}

	// Call exec to replace our process w/ zig
	if err := syscall.Exec(arg0, args, env); err != nil {
		panic(err)
	}
}

func checkGoVersion() error {
	version := runtime.Version()
	if !strings.HasPrefix(version, "go") {
		// assume it's a commit hash
		return nil
	}

	s := strings.Split(version, ".")
	if len(s) >= 2 {
		if s[0] != "go1" {
			// let's assume go2+ just work?
			return nil
		}
		minor, err := strconv.Atoi(s[1])
		if err != nil {
			return fmt.Errorf("unexpected go version format: %s", version)
		}
		if minor < 20 {
			return fmt.Errorf("go 1.20.3 or higher is necessary, not %s", version)
		}
		if minor == 20 {
			if len(s) >= 3 {
				// if bugfix .3 or higher we are ok
				bugfix, err := strconv.Atoi(s[2])
				if err != nil {
					return fmt.Errorf("unexpected go version format: %s", version)
				}
				if bugfix >= 3 {
					return nil
				}
			}
			return fmt.Errorf("go 1.20.3 or higher is necessary, not %s", version)
		}
		return nil
	} else {
		return fmt.Errorf("unexpected go version format: %s", version)
	}
}

// From https://dev.to/kristoff/zig-makes-go-cross-compilation-just-work-29ho
func lookupTarget(goos string, goarch string) ([]string, error) {
	if goos == "" && goarch == "" {
		// native, return nothing
		return nil, nil
	}

	if goos == runtime.GOOS && goarch == runtime.GOARCH {
		// native, return nothing
		return nil, nil
	}

	// TODO: go through and map all the other GOOS/GOARCH pairs to zig targets
	// For now just support linux/amd64 as a nice proof-of-concept
	if goos == "linux" && goarch == "amd64" {
		return []string{"-target", "x86_64-linux"}, nil
	}

	panic("unsupported GOOS/GOARCH pair")
}
