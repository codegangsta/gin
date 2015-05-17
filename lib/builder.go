package gin

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/constabulary/gb"
	gbCmd "github.com/constabulary/gb/cmd"
)

type Builder interface {
	Build() error
	Binary() string
	Errors() string
}

type builder struct {
	dir      string
	binary   string
	errors   string
	useGodep bool
}

func NewBuilder(dir string, bin string, useGodep bool) Builder {
	if len(bin) == 0 {
		bin = "bin"
	}

	// does not work on Windows without the ".exe" extension
	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(bin, ".exe") { // check if it already has the .exe extension
			bin += ".exe"
		}
	}

	return &builder{dir: dir, binary: bin, useGodep: useGodep}
}

func (b *builder) Binary() string {
	return b.binary
}

func (b *builder) Errors() string {
	return b.errors
}

func (b *builder) Build() error {
	var command *exec.Cmd
	if b.useGodep {
		command = exec.Command("godep", "go", "build", "-o", b.binary)
	} else {
		command = exec.Command("go", "build", "-o", b.binary)
	}
	command.Dir = b.dir

	output, err := command.CombinedOutput()

	if command.ProcessState.Success() {
		b.errors = ""
	} else {
		b.errors = string(output)
	}

	if len(b.errors) > 0 {
		return fmt.Errorf(b.errors)
	}

	return err
}

type gbBuilder struct {
	dir  string
	root string
	proj *gb.Project

	pkg *gb.Package
	bin string
	err error
}

// NewGbBuilder creates a constabulary/gb builder using dir to find the $porject root and for the main pkg as well
// if dir == . it uses the current working directory of gin
// i.e dir == /home/meh/devel/proj/src/spcil/cmd/worker
// => $proj = /home/meh/devel/proj
// => main pkg = spcil/cmd/worker
// TODO(cryptix): maybe use --bin flag to sepcify main pkg?
func NewGbBuilder(dir string) Builder {
	b := new(gbBuilder)

	if dir == "." {
		dir = gbCmd.MustGetwd()
	}
	b.dir = dir

	var err error
	b.root, err = gbCmd.FindProjectroot(b.dir)
	if err != nil {
		b.err = fmt.Errorf("gb cmd.FindProjectroot(%q) failed:  %v", b.dir, err)
		return b
	}
	b.proj = gb.NewProject(b.root)

	if err := b.importResolveBuild(false); err != nil {
		b.err = fmt.Errorf("importAndResolve() failed: %v", err)
		return b
	}

	return b
}

func (b *gbBuilder) Binary() string {
	p := b.proj.Bindir()
	p += "/" + b.bin
	return p
}

func (b *gbBuilder) Errors() string {
	if b.err != nil {
		return b.err.Error()
	}
	return ""
}

func (b *gbBuilder) Build() error {
	if b.err != nil {
		return b.err
	}

	if err := b.importResolveBuild(true); err != nil {
		b.err = err
		return b.err
	}

	return nil
}

// importResolveBuild does everything because you can't reuse ResolvePackages return values
// if their imports changed and broke gb.Build will be confused what to do
func (b *gbBuilder) importResolveBuild(build bool) error {
	ctx, err := b.proj.NewContext()
	if err != nil {
		return fmt.Errorf("proj.NewContext() failed: %v", err)

	}

	args := gbCmd.ImportPaths(ctx, b.dir, []string{}) // args..?!
	if len(args) < 1 {
		return fmt.Errorf("No ImportPaths.")

	}
	b.bin = filepath.Base(args[0])

	pkgs, err := gbCmd.ResolvePackages(ctx, args[0])
	if err != nil {
		return fmt.Errorf("gb cmd.ResolvePackages(%q) failed: %v", args[0], err)
	}
	if len(pkgs) < 1 {
		return fmt.Errorf("No Pakages.")
	}

	if build {
		if err := gb.Build(pkgs[0]); err != nil {
			return fmt.Errorf("gb.Build(%s) failed: %v", pkgs[0], err)
		}
	}

	if err := ctx.Destroy(); err != nil {
		return fmt.Errorf("ctx.Destroy() failed: %v", err)

	}

	return nil
}
