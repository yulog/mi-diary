package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"

	"github.com/yulog/mi-diary/app"
	"github.com/yulog/mi-diary/infra"
	"github.com/yulog/mi-diary/logic"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

var (
	BIN                 string = "mi-diary"
	VERSION             string = getVersion()
	CURRENT_REVISION, _        = sh.Output("git", "rev-parse", "--short", "HEAD")
	BUILD_LDFLAGS       string = "-s -w -X main.revision=" + CURRENT_REVISION
)

// func init() {
// 	VERSION = getVersion()
// 	CURRENT_REVISION, _ = sh.Output("git", "rev-parse", "--short", "HEAD")
// }

func getVersion() string {
	_, err := exec.LookPath("gobump")
	if err != nil {
		fmt.Println("installing gobump")
		sh.Run("go", "install", "github.com/x-motemen/gobump/cmd/gobump@latest")
	}
	v, _ := sh.Output("gobump", "show", "-r", ".")
	return v
}

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	// mg.Deps(InstallDeps)
	fmt.Println("Building...")
	bin := BIN
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	cmd := exec.Command("go", "build", "-trimpath", "-ldflags="+BUILD_LDFLAGS, "-o", bin, ".")
	return cmd.Run()
}

// A custom install step if you need your bin someplace other than go/bin
func Install() error {
	mg.Deps(Build)
	fmt.Println("Installing...")
	cmd := exec.Command("go", "install", "-ldflags="+BUILD_LDFLAGS, ".")
	return cmd.Run()
}

// Manage your deps, or running package managers.
// func InstallDeps() error {
// 	fmt.Println("Installing Deps...")
// 	cmd := exec.Command("go", "get", "github.com/stretchr/piglatin")
// 	return cmd.Run()
// }

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("goxz")
	bin := BIN
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	os.RemoveAll(bin)
}

func ShowVersion() {
	fmt.Println(getVersion())
}

func Credits() {
	_, err := exec.LookPath("gocredits")
	if err != nil {
		fmt.Println("installing gocredits")
		sh.Run("go", "install", "github.com/Songmu/gocredits")
	}
	s, _ := sh.Output("gocredits", ".")
	f, _ := os.Create("CREDITS")
	f.WriteString(s)
	defer f.Close()
}

func Cross() {
	_, err := exec.LookPath("goxz")
	if err != nil {
		fmt.Println("installing goxz")
		sh.Run("go", "install", "github.com/Songmu/goxz/cmd/goxz@latest")
	}
	sh.Run("goxz", "-n", BIN, "-pv=v"+VERSION, ".")
}

func Bump() {
	_, err := exec.LookPath("gobump")
	if err != nil {
		fmt.Println("installing gobump")
		sh.Run("go", "install", "github.com/x-motemen/gobump/cmd/gobump@latest")
	}
	sh.Run("gobump", "up", "-w", ".")
}

func Upload() {
	_, err := exec.LookPath("ghr")
	if err != nil {
		fmt.Println("installing ghr")
		sh.Run("go", "install", "github.com/tcnksm/ghr@latest")
	}
	sh.Run("ghr", "-draft", "v"+VERSION, "goxz")
}

func GenSchema() {
	fmt.Println("Generate schema...")
	app := app.New()
	l := logic.New().
		WithRepo(infra.New()).
		WithMigrationServiceUsingRepo().
		WithConfigRepo(infra.NewConfigInfra(app)).
		Build()
	l.GenerateSchema()
}

func GenMigrate() {
	GenSchema()
	fmt.Println("Generate diff...")
	sh.Run("atlas", "migrate", "diff", "migration", "--dir", "file://migrate/migrations?format=golang-migrate", "--dev-url", "sqlite://file?mode=memory", "--to", "file://migrate/schema.sql")
}

func Migrate() {
	fmt.Println("Migration...")
	app := app.New()
	l := logic.New().
		WithRepo(infra.New()).
		WithMigrationServiceUsingRepo().
		WithConfigRepo(infra.NewConfigInfra(app)).
		Build()
	l.Migrate()
}
