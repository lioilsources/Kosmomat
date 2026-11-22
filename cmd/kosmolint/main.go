// Command kosmolint checks content/ against docs/EDITORIAL.md (rules L1–L11).
//
//	go run ./cmd/kosmolint [-content content] [-warnings] [-today YYYY-MM-DD]
//
// Exit code 1 when there is at least one error.
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/lioilsources/kosmomat/internal/content"
	"github.com/lioilsources/kosmomat/internal/lint"
)

func main() {
	root := flag.String("content", "content", "content root")
	showWarn := flag.Bool("warnings", false, "print warnings, not only errors")
	today := flag.String("today", "", "reference date for freshness rules (default: today)")
	flag.Parse()

	ref := time.Now()
	if *today != "" {
		t, err := time.Parse("2006-01-02", *today)
		if err != nil {
			slog.Error("bad -today", "err", err)
			os.Exit(2)
		}
		ref = t
	}
	c, err := content.Load(*root)
	if err != nil {
		slog.Error("load content", "err", err)
		os.Exit(2)
	}
	issues := lint.Run(c, lint.Options{Today: ref})
	for _, i := range issues {
		if i.Sev == lint.Error || *showWarn {
			fmt.Println(i)
		}
	}
	errs, warns := lint.Count(issues)
	fmt.Printf("kosmolint: %d nodes, %d works — %d errors, %d warnings\n", len(c.Nodes), len(c.Works), errs, warns)
	if errs > 0 {
		os.Exit(1)
	}
}
