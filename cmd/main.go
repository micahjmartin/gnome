package main

import (
	"embed"
	"io/fs"
	"os"

	"github.com/nullmonk/gnome"
)

//go:embed example
var assets embed.FS

func main() {
	assets, _ := fs.Sub(assets, "example") // strip "example/" from the embedded asset names
	gnome.SetAssetLocker(assets)           // register the assets
	gnome.Run(os.Args)
}
