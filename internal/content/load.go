package content

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Langs are the languages required at launch.
var Langs = []string{"cs", "en"}

// Levels are the reading levels; "lay" is mandatory.
var Levels = []string{"kid", "lay", "student"}

// Load reads the content tree rooted at root (the folder that contains nodes/, works/ …).
// Parse failures are collected in Corpus.LoadErrors instead of aborting.
func Load(root string) (*Corpus, error) {
	c := &Corpus{
		Root:      root,
		Nodes:     map[string]*Node{},
		Works:     map[string]*Work{},
		Redirects: map[string]string{},
	}
	if _, err := os.Stat(filepath.Join(root, "nodes")); err != nil {
		return nil, fmt.Errorf("content root %q: %w", root, err)
	}
	if err := c.loadWorks(filepath.Join(root, "works")); err != nil {
		return nil, err
	}
	if err := c.loadNodes(filepath.Join(root, "nodes")); err != nil {
		return nil, err
	}
	c.loadMedia(filepath.Join(root, "media", "media.yaml"))
	c.loadPaths(filepath.Join(root, "paths"))
	c.loadRedirects(filepath.Join(root, "redirects.yaml"))
	return c, nil
}

func (c *Corpus) fail(file string, err error) {
	c.LoadErrors = append(c.LoadErrors, LoadError{File: c.rel(file), Err: err})
}

func (c *Corpus) rel(p string) string {
	if r, err := filepath.Rel(c.Root, p); err == nil {
		return filepath.ToSlash(r)
	}
	return p
}

func readYAML(path string, v any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	dec := yaml.NewDecoder(strings.NewReader(string(b)))
	dec.KnownFields(true)
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("yaml: %w", err)
	}
	return nil
}

func (c *Corpus) loadWorks(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		return nil
	}
	return filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(p, ".yaml") {
			return nil
		}
		var w Work
		if err := readYAML(p, &w); err != nil {
			c.fail(p, err)
			return nil
		}
		w.File = c.rel(p)
		if want := strings.TrimSuffix(filepath.Base(p), ".yaml"); w.Key != want {
			c.fail(p, fmt.Errorf("key %q does not match file name %q", w.Key, want))
		}
		if prev, ok := c.Works[w.Key]; ok {
			c.fail(p, fmt.Errorf("duplicate work key %q (also in %s)", w.Key, prev.File))
			return nil
		}
		c.Works[w.Key] = &w
		return nil
	})
}

func (c *Corpus) loadNodes(dir string) error {
	return filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Name() != "node.yaml" {
			return nil
		}
		var n Node
		if err := readYAML(p, &n); err != nil {
			c.fail(p, err)
			return nil
		}
		n.Dir = c.rel(filepath.Dir(p))
		n.Docs = map[string]*Doc{}
		for i := range n.Predictions {
			n.Predictions[i].Inline = true
			if n.Predictions[i].Theory == "" {
				n.Predictions[i].Theory = n.ID
			}
		}
		for _, lang := range Langs {
			for _, lvl := range Levels {
				name := lang + "." + lvl + ".md"
				b, err := os.ReadFile(filepath.Join(filepath.Dir(p), name))
				if err != nil {
					continue
				}
				n.Docs[lang+"."+lvl] = &Doc{File: n.Dir + "/" + name, Lang: lang, Level: lvl, Blocks: ParseMarkdown(string(b))}
			}
		}
		if prev, ok := c.Nodes[n.ID]; ok {
			c.fail(p, fmt.Errorf("duplicate node id %q (also in %s)", n.ID, prev.Dir))
			return nil
		}
		c.Nodes[n.ID] = &n
		return nil
	})
}

func (c *Corpus) loadMedia(path string) {
	if _, err := os.Stat(path); err != nil {
		return
	}
	var m struct {
		Items []Media `yaml:"items"`
	}
	if err := readYAML(path, &m); err != nil {
		c.fail(path, err)
		return
	}
	c.Media = m.Items
}

func (c *Corpus) loadPaths(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		var p Path
		if err := readYAML(filepath.Join(dir, e.Name()), &p); err != nil {
			c.fail(filepath.Join(dir, e.Name()), err)
			continue
		}
		c.Paths = append(c.Paths, p)
	}
}

func (c *Corpus) loadRedirects(path string) {
	if _, err := os.Stat(path); err != nil {
		return
	}
	var m map[string]string
	if err := readYAML(path, &m); err != nil {
		c.fail(path, err)
		return
	}
	c.Redirects = m
}

// NodeIDs returns node ids sorted, for deterministic output.
func (c *Corpus) NodeIDs() []string {
	ids := make([]string, 0, len(c.Nodes))
	for id := range c.Nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// WorkKeys returns work keys sorted.
func (c *Corpus) WorkKeys() []string {
	keys := make([]string, 0, len(c.Works))
	for k := range c.Works {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// AllPredictions returns inline predictions and prediction nodes, sorted by id.
func (c *Corpus) AllPredictions() []Prediction {
	var out []Prediction
	for _, id := range c.NodeIDs() {
		n := c.Nodes[id]
		out = append(out, n.Predictions...)
		if n.Type == "prediction" {
			out = append(out, n.AsPrediction())
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// Title returns the node title in lang, falling back to the id.
func (c *Corpus) Title(id, lang string) string {
	if n, ok := c.Nodes[id]; ok {
		if t := n.Title[lang]; t != "" {
			return t
		}
	}
	return id
}
