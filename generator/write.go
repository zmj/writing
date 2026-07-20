package main

import (
	"bytes"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
)

// writer emits the site under root, remembering every path it wrote so prune
// can clear out whatever a previous build left behind.
type writer struct {
	root    string
	written map[string]bool
}

func newWriter(root string) *writer {
	return &writer{root: root, written: make(map[string]bool)}
}

// write puts content at rel, touching the file only when the bytes differ.
// Unchanged pages keep their mtimes and stay out of the git diff.
func (w *writer) write(rel string, content []byte) error {
	rel = filepath.Clean(rel)
	w.written[rel] = true

	path := filepath.Join(w.root, rel)
	if old, err := os.ReadFile(path); err == nil && bytes.Equal(old, content) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0o644)
}

// prune deletes files under root that this build did not write — pages orphaned
// by a renamed or deleted post, which are committed and would otherwise stay
// live forever. Removals are logged rather than done silently.
func (w *writer) prune() error {
	if _, err := os.Stat(w.root); os.IsNotExist(err) {
		return nil
	}

	var stale []string
	err := filepath.WalkDir(w.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(w.root, path)
		if err != nil {
			return err
		}
		if !w.written[rel] {
			stale = append(stale, rel)
		}
		return nil
	})
	if err != nil {
		return err
	}

	sort.Strings(stale)
	for _, rel := range stale {
		log.Printf("removing stale %s", filepath.Join(w.root, rel))
		if err := os.Remove(filepath.Join(w.root, rel)); err != nil {
			return err
		}
	}
	return w.removeEmptyDirs()
}

// removeEmptyDirs clears directories left behind once their pages are gone.
func (w *writer) removeEmptyDirs() error {
	var dirs []string
	err := filepath.WalkDir(w.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != w.root {
			dirs = append(dirs, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Deepest first, so a directory holding only empty directories also goes.
	sort.Sort(sort.Reverse(sort.StringSlice(dirs)))
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		if len(entries) > 0 {
			continue
		}
		log.Printf("removing empty %s", dir)
		if err := os.Remove(dir); err != nil {
			return err
		}
	}
	return nil
}
