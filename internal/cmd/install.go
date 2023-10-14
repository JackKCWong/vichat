package cmd

import (
	"embed"
	"errors"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
)

//go:embed vim
var vimPlugins embed.FS

var InstallCmd = &cobra.Command{
	Use:     "install-vim-plugin",
	Short:   "install the vim plugin",
	Aliases: []string{"i"},
	Run: func(cmd *cobra.Command, args []string) {
		installVim()
	},
}

func installFiles(efs embed.FS, efsPath, dest string) {
	if _, err := os.Stat(dest); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(dest, 0750)
	}

	dirs, err := efs.ReadDir(efsPath)
	if err != nil {
		log.Fatalf("failed to list vim plugins: %q", err)
	}

	for _, f := range dirs {
		destFs, err := os.OpenFile(filepath.Join(dest, f.Name()), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0750)
		if err != nil {
			log.Fatalf("failed to create vim plugins: %q", err)
		}
		defer destFs.Close()

		srcFs, err := efs.Open(path.Join(efsPath, f.Name())) // embed.FS use / regardless OS
		if err != nil {
			log.Fatalf("failed to read vim plugins: %q", err)
		}
		defer srcFs.Close()

		if _, err := io.Copy(destFs, srcFs); err != nil {
			log.Printf("failed to install %s", destFs.Name())
		} else {
			log.Printf("installed %s", destFs.Name())
		}
	}
}

func installVim() {
	installFiles(vimPlugins, "vim/ftdetect", os.ExpandEnv("$HOME/.vim/ftdetect"))
	installFiles(vimPlugins, "vim/ftplugin", os.ExpandEnv("$HOME/.vim/ftplugin"))
	installFiles(vimPlugins, "vim/syntax", os.ExpandEnv("$HOME/.vim/syntax"))
}

func hasDifference(efs embed.FS, fsPath, dest string) bool {
	if _, err := os.Stat(dest); errors.Is(err, os.ErrNotExist) {
		return false
	}

	dirs, err := efs.ReadDir(fsPath)
	if err != nil {
		log.Fatalf("failed to list vim plugins: %q", err)
	}

	for _, f := range dirs {
		destFs, err := os.Open(filepath.Join(dest, f.Name()))
		if err != nil {
			return true
		}
		defer destFs.Close()

		srcFs, err := efs.Open(path.Join(fsPath, f.Name()))
		if err != nil {
			log.Fatalf("failed to read vim plugins: %q", err)
		}
		defer srcFs.Close()

		destFi, err := destFs.Stat()
		if err != nil {
			log.Printf("failed to read vim plugins: %q", err)
		}

		srcFi, err := srcFs.Stat()
		if err != nil {
			log.Printf("failed to read vim plugins: %q", err)
		}

		if destFi.Size() != srcFi.Size() {
			return true
		}
	}

	return false
}