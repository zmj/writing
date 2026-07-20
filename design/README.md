# Blog template handoff

Static templates for a markdown → HTML pipeline.

- `style.css` — the whole design. One shared stylesheet for both page types. Light/dark follows `prefers-color-scheme` automatically via CSS variables at the top; that's also the only place to touch colors.
- `post.html` — post template. Your converter fills `<title>`, the `<h1>`, `.post-date`, and the body between the marked comments. Standard markdown output (p, h2, em/strong, a, code, pre>code, ul/ol) is already styled; no per-element classes needed.
- `index.html` — index template. One `.month-group` per month, newest first; each post is a bare `<a>` inside `.month-posts`.

Notes:
- Fonts load from Google Fonts (Source Serif 4 + IBM Plex Mono). For fully self-contained hosting, download the woff2 files and swap the `<link>`s for `@font-face` rules in style.css.
- `<meta name="color-scheme" content="light dark">` keeps browser UI (scrollbars, form controls) matching the theme — keep it in every page.
- Line length is capped at 620px; below 560px viewport width the index's month labels stack above their titles.
