// Command generator turns posts/*.md into the static site under docs/, which
// GitHub Pages serves at words.zmj.dev. Run it from the repository root.
package main

import (
	"log"
	"path/filepath"
)

const (
	postsDir = "posts"
	outDir   = "docs"
)

func main() {
	log.SetFlags(0)

	posts, err := loadPosts(postsDir)
	if err != nil {
		log.Fatal(err)
	}

	w := newWriter(outDir)

	index, err := renderPage("index.html", indexData{Months: groupByMonth(posts)})
	if err != nil {
		log.Fatal(err)
	}
	if err := w.write("index.html", index); err != nil {
		log.Fatal(err)
	}

	for _, p := range posts {
		page, err := renderPage("post.html", p)
		if err != nil {
			log.Fatalf("%s: %v", p.Src, err)
		}
		if err := w.write(filepath.Join(p.Slug, "index.html"), page); err != nil {
			log.Fatal(err)
		}
	}

	for _, name := range []string{"style.css", "CNAME"} {
		content, err := staticFS.ReadFile("static/" + name)
		if err != nil {
			log.Fatal(err)
		}
		if err := w.write(name, content); err != nil {
			log.Fatal(err)
		}
	}

	// Written rather than embedded: //go:embed silently skips dotfiles, so an
	// embedded .nojekyll would go missing without any error.
	if err := w.write(".nojekyll", nil); err != nil {
		log.Fatal(err)
	}

	if err := w.prune(); err != nil {
		log.Fatal(err)
	}

	log.Printf("wrote %d post(s) to %s/", len(posts), outDir)
}
