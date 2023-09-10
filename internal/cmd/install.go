package cmd

import (
	"embed"
	"errors"
	"io"
	"log"
	"os"
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

func installFiles(fs embed.FS, path, dest string) {
	if _, err := os.Stat(dest); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(dest, 0750)
	}

	dirs, err := fs.ReadDir(path)
	if err != nil {
		log.Fatalf("failed to list vim plugins: %q", err)
	}

	for _, f := range dirs {
		destFs, err := os.OpenFile(filepath.Join(dest, f.Name()), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0750)
		if err != nil {
			log.Fatalf("failed to create vim plugins: %q", err)
		}
		defer destFs.Close()

		srcFs, err := fs.Open(filepath.Join(path, f.Name()))
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
