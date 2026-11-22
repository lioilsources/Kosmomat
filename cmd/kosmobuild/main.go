// Command kosmobuild lints content/ and builds the content pack:
// out/pack.sqlite(.zst), out/manifest.json and out/web/ for the Astro site.
//
//	go run ./cmd/kosmobuild [-content content] [-out out] [-version v] [-base-url URL]
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/lioilsources/kosmomat/internal/content"
	"github.com/lioilsources/kosmomat/internal/lint"
	"github.com/lioilsources/kosmomat/internal/pack"
)

func main() {
	root := flag.String("content", "content", "content root")
	out := flag.String("out", "out", "output directory")
	version := flag.String("version", "", "pack version (default dev-<content hash>)")
	minApp := flag.Int("min-app", 1, "minimum app build able to read this pack")
	baseURL := flag.String("base-url", "https://kosmomat.example/packs", "public URL prefix of the pack")
	appAssets := flag.String("app", "", "also copy pack.sqlite + manifest.json into this app asset dir (e.g. app/assets/pack)")
	editorial := flag.String("editorial", "docs/EDITORIAL.md", "editorial constitution shipped in the pack (cs); empty to skip")
	force := flag.Bool("force", false, "build even when lint has errors (never for release)")
	flag.Parse()

	c, err := content.Load(*root)
	if err != nil {
		slog.Error("load content", "err", err)
		os.Exit(2)
	}
	issues := lint.Run(c, lint.Options{Today: time.Now()})
	if errs, _ := lint.Count(issues); errs > 0 {
		for _, i := range issues {
			if i.Sev == lint.Error {
				fmt.Println(i)
			}
		}
		if !*force {
			slog.Error("lint failed; fix errors or pass -force", "errors", errs)
			os.Exit(1)
		}
	}
	ed := map[string]string{}
	if *editorial != "" {
		b, err := os.ReadFile(*editorial)
		if err != nil {
			slog.Error("read editorial", "err", err)
			os.Exit(1)
		}
		ed["cs"] = string(b)
	}
	m, err := pack.Build(c, pack.Options{OutDir: *out, Version: *version, MinApp: *minApp, BaseURL: *baseURL, Editorial: ed})
	if err != nil {
		slog.Error("build", "err", err)
		os.Exit(1)
	}
	if *appAssets != "" {
		if err := pack.CopyToApp(*out, *appAssets); err != nil {
			slog.Error("copy to app", "err", err)
			os.Exit(1)
		}
	}
	fmt.Printf("kosmobuild: pack %s, %d nodes, %d bytes (zst), sha256 %s\n", m.Version, m.Nodes, m.Size, m.SHA256[:16])
}
