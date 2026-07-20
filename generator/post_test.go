package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFilename(t *testing.T) {
	tests := []struct {
		name      string
		wantDate  string
		wantSlug  string
		wantError bool
	}{
		{name: "20260720_testing_strategy.md", wantDate: "2026-07-20", wantSlug: "testing-strategy"},
		{name: "20260401_one.md", wantDate: "2026-04-01", wantSlug: "one"},
		{name: "20260401_go_1_22_notes.md", wantDate: "2026-04-01", wantSlug: "go-1-22-notes"},

		{name: "20261340_bad_month.md", wantError: true},
		{name: "20260231_bad_day.md", wantError: true},
		{name: "testing_strategy.md", wantError: true},
		{name: "2026720_short_date.md", wantError: true},
		{name: "20260720-dashes-not-underscores.md", wantError: true},
		{name: "20260720_Mixed_Case.md", wantError: true},
		{name: "20260720_testing_strategy.MD", wantError: true},
		{name: "20260720_.md", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date, slug, err := parseFilename(tt.name)
			if tt.wantError {
				if err == nil {
					t.Fatalf("parseFilename(%q) = %v, %q; want error", tt.name, date, slug)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseFilename(%q): %v", tt.name, err)
			}
			if got := date.Format("2006-01-02"); got != tt.wantDate {
				t.Errorf("date = %q, want %q", got, tt.wantDate)
			}
			if slug != tt.wantSlug {
				t.Errorf("slug = %q, want %q", slug, tt.wantSlug)
			}
		})
	}
}

func TestSplitTitle(t *testing.T) {
	tests := []struct {
		desc      string
		src       string
		wantTitle string
		wantBody  string
		wantError bool
	}{
		{
			desc:      "title then body",
			src:       "# How to Test\n\nFirst paragraph.\n",
			wantTitle: "How to Test",
			wantBody:  "\nFirst paragraph.\n",
		},
		{
			desc:      "leading blank lines are skipped",
			src:       "\n\n# Title\nbody\n",
			wantTitle: "Title",
			wantBody:  "body\n",
		},
		{
			desc:      "trailing whitespace is trimmed",
			src:       "#   Spaced Out   \n\nbody\n",
			wantTitle: "Spaced Out",
			wantBody:  "\nbody\n",
		},
		{
			desc:      "title only, no body",
			src:       "# Just A Title\n",
			wantTitle: "Just A Title",
			wantBody:  "",
		},
		{
			desc:      "carriage returns are tolerated",
			src:       "# Title\r\nbody\r\n",
			wantTitle: "Title",
			wantBody:  "body\r\n",
		},

		{desc: "body before any heading", src: "Some prose.\n\n# Late Title\n", wantError: true},
		{desc: "level-2 heading first", src: "## Section\n\nbody\n", wantError: true},
		{desc: "no space after hash", src: "#Title\n\nbody\n", wantError: true},
		{desc: "empty heading", src: "#\n\nbody\n", wantError: true},
		{desc: "empty file", src: "", wantError: true},
		{desc: "blank lines only", src: "\n\n\n", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			title, body, err := splitTitle([]byte(tt.src))
			if tt.wantError {
				if err == nil {
					t.Fatalf("splitTitle(%q) = %q, %q; want error", tt.src, title, body)
				}
				return
			}
			if err != nil {
				t.Fatalf("splitTitle(%q): %v", tt.src, err)
			}
			if title != tt.wantTitle {
				t.Errorf("title = %q, want %q", title, tt.wantTitle)
			}
			if string(body) != tt.wantBody {
				t.Errorf("body = %q, want %q", body, tt.wantBody)
			}
		})
	}
}

func TestGroupByMonth(t *testing.T) {
	posts, err := loadPostsFrom(t, map[string]string{
		"20260720_late_july.md":  "# Late July\n\nbody\n",
		"20260701_early_july.md": "# Early July\n\nbody\n",
		"20260615_june.md":       "# June\n\nbody\n",
	})
	if err != nil {
		t.Fatal(err)
	}

	months := groupByMonth(posts)
	if len(months) != 2 {
		t.Fatalf("got %d months, want 2", len(months))
	}
	if months[0].Label != "July 2026" || len(months[0].Posts) != 2 {
		t.Errorf("months[0] = %q with %d posts, want July 2026 with 2", months[0].Label, len(months[0].Posts))
	}
	if months[1].Label != "June 2026" || len(months[1].Posts) != 1 {
		t.Errorf("months[1] = %q with %d posts, want June 2026 with 1", months[1].Label, len(months[1].Posts))
	}
	// Newest first within a month.
	if months[0].Posts[0].Title != "Late July" {
		t.Errorf("months[0].Posts[0] = %q, want %q", months[0].Posts[0].Title, "Late July")
	}
}

func TestLoadPostsRejectsDuplicateSlugs(t *testing.T) {
	// "a_b" and "a-b" both slugify to "a-b"... but dashes are rejected by the
	// filename pattern, so collide two dates onto the same word instead.
	_, err := loadPostsFrom(t, map[string]string{
		"20260720_same_name.md": "# One\n\nbody\n",
		"20260721_same_name.md": "# Two\n\nbody\n",
	})
	if err == nil {
		t.Fatal("loadPosts accepted two posts with the same slug")
	}
}

func TestLoadPostsRejectsSecondH1(t *testing.T) {
	_, err := loadPostsFrom(t, map[string]string{
		"20260720_two_titles.md": "# Title\n\nbody\n\n# Another Title\n\nmore\n",
	})
	if err == nil {
		t.Fatal("loadPosts accepted a post with a second level-1 heading")
	}
}

func TestLoadPostsSkipsNonMarkdown(t *testing.T) {
	posts, err := loadPostsFrom(t, map[string]string{
		"20260720_real_post.md": "# Real\n\nbody\n",
		".DS_Store":             "junk",
		"notes.txt":             "junk",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != 1 {
		t.Fatalf("got %d posts, want 1", len(posts))
	}
}

func loadPostsFrom(t *testing.T, files map[string]string) ([]Post, error) {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return loadPosts(dir)
}
