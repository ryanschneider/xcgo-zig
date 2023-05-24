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
	log.Default().SetFlags(0)
	log.Default().SetPrefix("xcgo-zig: ")
	if err := run(); err != nil {
		log.Fatalf("fatal error: %s", err.Error())
	}
}

func run() error {
	goCmd := os.Getenv("GO")
	if goCmd == "" {
		goCmd = "go"
	}

	goPath, err := exec.LookPath(goCmd)
	if err != nil {
		return fmt.Errorf("could not find go: %w", err)
	}

	goversion, err := getGoVersion(goPath)
	if err != nil {
		return fmt.Errorf("could not determine go version: %w", err)
	}

	zigPath, err := exec.LookPath("zig")
	if err != nil {
		return fmt.Errorf("could not find zig: %w", err)
	}

	goos, _ := os.LookupEnv("GOOS")
	goarch, _ := os.LookupEnv("GOARCH")
	zigTarget, err := lookupTarget(goos, goarch, goversion)
	if err != nil {
		return fmt.Errorf("incompatible GOOS/GOARCH pair: %w", err)
	}

	// Let's override CC and CXX w/ "/path/to/zig cc -target ..."
	if err := os.Setenv("CC", zigArgs(zigPath, "cc", zigTarget)); err != nil {
		return fmt.Errorf("cannot set CC: %w", err)
	}
	if err := os.Setenv("CXX", zigArgs(zigPath, "c++", zigTarget)); err != nil {
		return fmt.Errorf("cannot set CXX: %w", err)
	}
	if err := os.Setenv("CGO_ENABLED", "1"); err != nil {
		return fmt.Errorf("cannot set CGO_ENABLED: %w", err)
	}

	env := append(os.Environ())
	args := append(os.Args)

	// first log w/ the simple go command
	if len(args) > 0 {
		args[0] = goCmd
	}
	log.Printf("Running `%s` with:", strings.Join(args, " "))

	// then replace go with the fully qualified path
	if len(args) > 0 {
		args[0] = goPath
	}
	args0 := goPath

	log.Printf("GOOS=%s", goos)
	log.Printf("GOARCH=%s", goarch)
	log.Printf("CC=%s", os.Getenv("CC"))
	log.Printf("CXX=%s", os.Getenv("CXX"))

	// Call exec to replace our process w/ go
	if err := syscall.Exec(args0, args, env); err != nil {
		return err
	}

	return nil
}

func zigArgs(path string, cmd string, target []string) string {
	args := strings.Join(target, " ")
	if len(args) > 0 {
		args = " " + args
	}
	return fmt.Sprintf("%s %s%s", path, cmd, args)
}

type goVersion struct {
	Major  int
	Minor  int
	Bugfix int
	Hash   string
}

func getGoVersion(goPath string) (*goVersion, error) {
	// TODO: run `go version` instead of using the runtime.Version
	out, err := exec.Command(goPath, "version").Output()
	if err != nil {
		return nil, fmt.Errorf("could not run %s: %w", goPath, err)
	}
	s := strings.Split(string(out), " ")
	if len(s) < 3 || s[0] != "go" || s[1] != "version" {
		return nil, fmt.Errorf("unexpected go version output: %s", string(out))
	}
	version := s[2]
	if version == "devel" {
		if len(s) < 4 {
			return nil, fmt.Errorf("unexpected go development version output: %s", string(out))
		}
		version = s[3]
	}

	s = strings.Split(version, ".")
	if len(s) == 0 {
		return nil, fmt.Errorf("unexpected go version: %s", version)
	}

	v := &goVersion{}
	if s[0] == "go1" {
		v.Major = 1
	} else {
		// for now just bail until go2 is released
		return nil, fmt.Errorf("go2 not yet supported")
	}

	if len(s) >= 2 {
		// could be of the form 21-abcdefg
		if x := strings.Split(s[1], "-"); len(x) == 2 {
			minor, err := strconv.Atoi(x[0])
			if err != nil {
				return nil, fmt.Errorf("unexpected go minor version format: %s", version)
			}
			v.Minor = minor
			v.Hash = x[1]
		} else {
			minor, err := strconv.Atoi(s[1])
			if err != nil {
				return nil, fmt.Errorf("unexpected go minor version format: %s", version)
			}
			v.Minor = minor
		}
	}

	if len(s) >= 3 {
		bugfix, err := strconv.Atoi(s[2])
		if err != nil {
			return nil, fmt.Errorf("unexpected go bugfix version format: %s", version)
		}
		v.Bugfix = bugfix
	}
	return v, nil
}

func (v goVersion) String() string {
	var bugFixOrHash string
	if v.Hash != "" {
		bugFixOrHash = fmt.Sprintf("-%s", v.Hash)
	} else {
		bugFixOrHash = fmt.Sprintf(".%d", v.Bugfix)
	}
	return fmt.Sprintf("%d.%d%s", v.Major, v.Minor, bugFixOrHash)
}

func (v goVersion) AtLeast(major, minor, bugfix int) error {
	err := fmt.Errorf("go version %s less than %d.%d.%d", v.String(), major, minor, bugfix)
	if v.Major > major {
		return nil
	}
	if v.Major < major {
		return err
	}
	if v.Minor > minor {
		return nil
	}
	if v.Minor < minor {
		return err
	}
	if v.Bugfix < bugfix {
		return err
	}
	return nil
}

// From https://dev.to/kristoff/zig-makes-go-cross-compilation-just-work-29ho
func lookupTarget(goos string, goarch string, goversion *goVersion) ([]string, error) {
	if goos == "" && goarch == "" {
		// native, return nothing
		return nil, nil
	}

	if goos == runtime.GOOS && goarch == runtime.GOARCH {
		// native, return nothing
		return nil, nil
	}

	unsupportedTarget := fmt.Errorf("unsupported GOOS/GOARCH pair: %s/%s", goos, goarch)

	// For now let's be very explicit on GOOS/GOARCH pairs we support
	// eventually this could be made more dynamic, but I'm guessing there will
	// be lots of untested corner cases for various targets.
	switch goos {
	case "linux":
		switch goarch {
		case "amd64":
			return []string{"-target", "x86_64-linux-musl"}, nil
		case "arm64":
			return []string{"-target", "aarch64-linux-musl"}, nil
		default:
			return nil, unsupportedTarget
		}
	case "windows":
		switch goarch {
		case "amd64":
			return []string{"-target", "x86_64-windows"}, nil
		default:
			return nil, unsupportedTarget
		}
	case "darwin":
		switch goarch {
		case "arm64":
			if err := goversion.AtLeast(1, 20, 3); err != nil {
				return nil, fmt.Errorf("darwin/arm64 requires go 1.20.3 or higher")
			}
			return []string{"-target", "aarch64-macos"}, nil
		case "amd64":
			return []string{"-target", "x86_64-macos"}, nil
		default:
			return nil, unsupportedTarget
		}
	case "wasip1":
		switch goarch {
		case "wasm":
			// wasip1 fails with:
			//   cgo: unknown ptrSize for $GOARCH "wasm"
			// as of:
			//   go version devel go1.21-26f2569 Tue May 23 21:46:00 2023
			// but hopefully will be supported before release!
			if err := goversion.AtLeast(1, 21, 0); err != nil {
				return nil, fmt.Errorf("wasip1/wasm requires go 1.21.0 or higher")
			}
			return []string{"-target", "wasm32-wasi-musl"}, nil
		default:
			return nil, unsupportedTarget
		}
	default:
		return nil, unsupportedTarget
	}
}
