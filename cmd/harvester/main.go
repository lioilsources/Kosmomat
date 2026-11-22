// Command harvester fetches bibliographic metadata and checks content/works.
//
//	harvester verify [-write] [-json report.json] [work-key …]
//	harvester arxiv    [-from D] [-until D] [-sets physics:astro-ph,…]
//	harvester crossref -issn 0004-637X,… [-from D] [-until D]
//	harvester scix     [-from D] [-until D]          (needs ADS_API_TOKEN)
//	harvester stats
//
// Dates are YYYY-MM-DD and default to yesterday (UTC). Harvested records go
// to the SQLite file given by -db (default out/harvest.sqlite).
// Environment: CROSSREF_MAILTO (optional contact for the Crossref polite pool),
// ADS_API_TOKEN (SciX / NASA ADS).
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/lioilsources/kosmomat/internal/content"
	"github.com/lioilsources/kosmomat/internal/harvest"
)

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	cmd, args := os.Args[1], os.Args[2:]
	var err error
	switch cmd {
	case "verify":
		err = verify(ctx, args)
	case "arxiv", "crossref", "scix":
		err = harvestCmd(ctx, cmd, args)
	case "stats":
		err = stats(args)
	default:
		usage()
	}
	if err != nil {
		slog.Error(cmd, "err", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: harvester verify|arxiv|crossref|scix|stats [flags]")
	os.Exit(2)
}

func clients() (*harvest.Crossref, *harvest.Arxiv, *harvest.Scix) {
	cr := harvest.NewClient(time.Second)
	cr.Mailto = os.Getenv("CROSSREF_MAILTO")
	ax := harvest.NewClient(3 * time.Second) // arXiv: max 1 request / 3 s
	var sx *harvest.Scix
	if tok := os.Getenv("ADS_API_TOKEN"); tok != "" {
		c := harvest.NewClient(time.Second)
		c.Header = map[string]string{"Authorization": "Bearer " + tok}
		sx = &harvest.Scix{C: c}
	}
	return &harvest.Crossref{C: cr}, &harvest.Arxiv{C: ax}, sx
}

func verify(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	root := fs.String("content", "content", "content root")
	write := fs.Bool("write", false, "set verified: true on works whose metadata fully agrees")
	jsonOut := fs.String("json", "", "write the full report as JSON to this file")
	all := fs.Bool("all", false, "also re-check works already verified (never writes them)")
	fs.Parse(args)

	c, err := content.Load(*root)
	if err != nil {
		return err
	}
	keys := fs.Args()
	if *all && len(keys) == 0 {
		keys = c.WorkKeys()
	}
	cr, ax, sx := clients()
	if sx == nil {
		slog.Info("ADS_API_TOKEN not set: bibcode-only works are skipped")
	}
	v := &harvest.Verifier{Crossref: cr, Arxiv: ax, Scix: sx,
		Log: func(f string, a ...any) { fmt.Fprintf(os.Stderr, f+"\n", a...) }}
	checks, err := v.Run(ctx, c, keys, *write)
	if *jsonOut != "" && len(checks) > 0 {
		b, _ := json.MarshalIndent(checks, "", "  ")
		if werr := os.WriteFile(*jsonOut, b, 0o644); werr != nil {
			return werr
		}
	}
	if err != nil {
		return err
	}
	printReport(checks)
	return nil
}

func printReport(checks []harvest.Check) {
	order := []harvest.Status{harvest.StatusRetracted, harvest.StatusMismatch, harvest.StatusNotFound,
		harvest.StatusError, harvest.StatusSkipped, harvest.StatusOK}
	for _, st := range order {
		for _, ch := range checks {
			if ch.Status != st {
				continue
			}
			mark := ""
			if ch.Written {
				mark = " (written)"
			}
			fmt.Printf("%-10s %-30s %-8s%s\n", ch.Status, ch.Key, ch.Source, mark)
			for _, n := range ch.Notes {
				fmt.Printf("           · %s\n", n)
			}
			for _, p := range ch.Problems {
				fmt.Printf("           - %s\n", p)
			}
			if ch.Status == harvest.StatusMismatch && ch.Venue != "" {
				fmt.Printf("           · source venue: %s\n", ch.Venue)
			}
		}
	}
	n := harvest.Counts(checks)
	fmt.Printf("\n%d checked: %d ok, %d mismatch, %d not found, %d retracted, %d error, %d skipped\n",
		len(checks), n[harvest.StatusOK], n[harvest.StatusMismatch], n[harvest.StatusNotFound],
		n[harvest.StatusRetracted], n[harvest.StatusError], n[harvest.StatusSkipped])
}

func harvestCmd(ctx context.Context, source string, args []string) error {
	fs := flag.NewFlagSet(source, flag.ExitOnError)
	yesterday := time.Now().UTC().AddDate(0, 0, -1).Format(time.DateOnly)
	db := fs.String("db", filepath.Join("out", "harvest.sqlite"), "SQLite database")
	from := fs.String("from", yesterday, "first day (YYYY-MM-DD)")
	until := fs.String("until", "", "last day (YYYY-MM-DD, default = from)")
	sets := fs.String("sets", strings.Join(harvest.DefaultSets, ","), "arXiv OAI sets")
	issn := fs.String("issn", "", "Crossref: comma-separated journal ISSNs")
	fs.Parse(args)
	if *until == "" {
		*until = *from
	}
	for _, d := range []string{*from, *until} {
		if _, err := time.Parse(time.DateOnly, d); err != nil {
			return fmt.Errorf("bad date %q", d)
		}
	}
	if err := os.MkdirAll(filepath.Dir(*db), 0o755); err != nil {
		return err
	}
	st, err := harvest.OpenStore(*db)
	if err != nil {
		return err
	}
	defer st.Close()

	now := func() string { return time.Now().UTC().Format(time.RFC3339) }
	cr, ax, sx := clients()
	run, err := st.StartRun(source, fmt.Sprintf("from=%s until=%s sets=%s issn=%s", *from, *until, *sets, *issn), now())
	if err != nil {
		return err
	}
	put := func(r harvest.Record) error {
		run.Fetched++
		isNew, err := st.Upsert(r, now())
		if isNew {
			run.Inserted++
		}
		return err
	}
	switch source {
	case "arxiv":
		for _, set := range strings.Split(*sets, ",") {
			slog.Info("arxiv", "set", set, "from", *from, "until", *until)
			if err = ax.Harvest(ctx, set, *from, *until, harvest.DefaultCategories, put); err != nil {
				break
			}
		}
	case "crossref":
		if *issn == "" {
			err = fmt.Errorf("crossref needs -issn")
			break
		}
		for _, is := range strings.Split(*issn, ",") {
			slog.Info("crossref", "issn", is, "from", *from, "until", *until)
			if err = cr.ByISSN(ctx, strings.TrimSpace(is), *from, *until, put); err != nil {
				break
			}
		}
	case "scix":
		if sx == nil {
			err = fmt.Errorf("scix needs ADS_API_TOKEN")
			break
		}
		err = sx.Refereed(ctx, *from, *until, put)
	}
	if ferr := st.FinishRun(run, now(), err); ferr != nil && err == nil {
		err = ferr
	}
	slog.Info("done", "source", source, "fetched", run.Fetched, "new", run.Inserted)
	return err
}

func stats(args []string) error {
	fs := flag.NewFlagSet("stats", flag.ExitOnError)
	db := fs.String("db", filepath.Join("out", "harvest.sqlite"), "SQLite database")
	fs.Parse(args)
	st, err := harvest.OpenStore(*db)
	if err != nil {
		return err
	}
	defer st.Close()
	m, err := st.Stats()
	if err != nil {
		return err
	}
	b, _ := json.MarshalIndent(m, "", "  ")
	fmt.Println(string(b))
	return nil
}
