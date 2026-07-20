package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

// Post is one rendered markdown file. Title and Date come from convention
// rather than frontmatter: the date from the filename, the title from the
// leading level-1 heading.
type Post struct {
	Title string
	Slug  string
	Date  time.Time
	Body  template.HTML // template.HTML, or the page ships escaped markup
	Src   string        // source filename, for error messages and sort ties
}

func (p Post) DateLong() string { return p.Date.Format("January 2, 2006") }

// Month is a run of posts sharing a calendar month, for the index grouping.
type Month struct {
	Label string
	Posts []Post
}

var md = goldmark.New(
	// Unsafe is deliberately off: no post needs raw HTML, and loadPost turns
	// goldmark's silent omission of it into an error.
	goldmark.WithExtensions(extension.GFM, extension.Typographer),
	goldmark.WithParserOptions(parser.WithAutoHeadingID()),
)

var nameRE = regexp.MustCompile(`^(\d{8})_([a-z0-9_]+)\.md$`)

// parseFilename pulls the date and slug out of a name like
// 20260720_testing_strategy.md.
func parseFilename(name string) (time.Time, string, error) {
	m := nameRE.FindStringSubmatch(name)
	if m == nil {
		return time.Time{}, "", fmt.Errorf("bad post filename %q: want YYYYMMDD_slug_words.md", name)
	}
	date, err := time.Parse("20060102", m[1])
	if err != nil {
		return time.Time{}, "", fmt.Errorf("bad date in post filename %q: %w", name, err)
	}
	return date, strings.ReplaceAll(m[2], "_", "-"), nil
}

var h1RE = regexp.MustCompile(`^#[ \t]+(\S.*?)\s*$`)

// splitTitle takes the leading "# Title" heading off a post, returning it
// separately from the body so the title is not rendered twice.
func splitTitle(src []byte) (string, []byte, error) {
	rest := src
	for {
		line, tail, more := bytes.Cut(rest, []byte("\n"))
		if len(bytes.TrimSpace(line)) > 0 {
			m := h1RE.FindSubmatch(bytes.TrimRight(line, "\r"))
			if m == nil {
				return "", nil, fmt.Errorf("first line must be a \"# Title\" heading, got %q", line)
			}
			return string(m[1]), tail, nil
		}
		if !more {
			return "", nil, fmt.Errorf("post has no content")
		}
		rest = tail
	}
}

// loadPosts reads every post in dir, newest first. Anything that is not a .md
// file is skipped; a .md file that breaks the conventions is an error, because
// output is committed by hand and a silently dropped post would go unnoticed.
func loadPosts(dir string) ([]Post, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var posts []Post
	bySlug := make(map[string]string)
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.EqualFold(filepath.Ext(name), ".md") {
			continue
		}
		p, err := loadPost(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		if prev, dup := bySlug[p.Slug]; dup {
			return nil, fmt.Errorf("posts %s and %s both produce slug %q", prev, name, p.Slug)
		}
		bySlug[p.Slug] = name
		posts = append(posts, p)
	}

	// Newest first, ties broken by filename so the output is deterministic —
	// nondeterminism here would rewrite every page on every build.
	sort.Slice(posts, func(i, j int) bool {
		if !posts[i].Date.Equal(posts[j].Date) {
			return posts[i].Date.After(posts[j].Date)
		}
		return posts[i].Src > posts[j].Src
	})
	return posts, nil
}

func loadPost(path string) (Post, error) {
	name := filepath.Base(path)
	date, slug, err := parseFilename(name)
	if err != nil {
		return Post{}, err
	}

	src, err := os.ReadFile(path)
	if err != nil {
		return Post{}, err
	}

	title, body, err := splitTitle(src)
	if err != nil {
		return Post{}, fmt.Errorf("%s: %w", name, err)
	}

	var buf bytes.Buffer
	if err := md.Convert(body, &buf); err != nil {
		return Post{}, fmt.Errorf("%s: %w", name, err)
	}
	// The template already emits the title as the page's only <h1>.
	if bytes.Contains(buf.Bytes(), []byte("<h1")) {
		return Post{}, fmt.Errorf("%s: body has a second level-1 heading; use ## for sections", name)
	}
	// goldmark drops raw HTML quietly when unsafe rendering is off.
	if bytes.Contains(buf.Bytes(), []byte("raw HTML omitted")) {
		return Post{}, fmt.Errorf("%s: contains raw HTML, which is not enabled", name)
	}

	return Post{
		Title: title,
		Slug:  slug,
		Date:  date,
		Body:  template.HTML(buf.String()),
		Src:   name,
	}, nil
}

// groupByMonth collects an already-sorted post list into calendar months.
func groupByMonth(posts []Post) []Month {
	var months []Month
	for _, p := range posts {
		label := p.Date.Format("January 2006")
		if n := len(months); n > 0 && months[n-1].Label == label {
			months[n-1].Posts = append(months[n-1].Posts, p)
			continue
		}
		months = append(months, Month{Label: label, Posts: []Post{p}})
	}
	return months
}
