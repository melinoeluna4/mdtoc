package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseHeadingsIgnoresCodeFences(t *testing.T) {
	t.Parallel()

	input := `# Title

## Install

` + "```sh" + `
# not a heading
` + "```" + `

### Details

#### Too deep
`
	headings, err := parseHeadings(strings.NewReader(input))
	if err != nil {
		t.Fatalf("parseHeadings: %v", err)
	}

	if len(headings) != 4 {
		t.Fatalf("headings = %d, want 4", len(headings))
	}
	if headings[0].text != "Title" || headings[0].level != 1 {
		t.Fatalf("first heading = %#v", headings[0])
	}
	if headings[2].text != "Details" || headings[2].level != 3 {
		t.Fatalf("third heading = %#v", headings[2])
	}
}

func TestParseHeadingRejectsInvalidATX(t *testing.T) {
	t.Parallel()

	for _, line := range []string{"#NoSpace", "#######", "title"} {
		if _, _, ok := parseHeading(line); ok {
			t.Fatalf("parseHeading(%q) = true, want false", line)
		}
	}
}

func TestSlugifyHandlesAccentsAndDuplicates(t *testing.T) {
	t.Parallel()

	if got := slugify("Caché en Cloudflare"); got != "caché-en-cloudflare" {
		t.Fatalf("slugify = %q, want caché-en-cloudflare", got)
	}
	if got := slugify("Go / Rust"); got != "go-rust" {
		t.Fatalf("slugify = %q, want go-rust", got)
	}
	if got := slugify("!!!"); got != "section" {
		t.Fatalf("slugify = %q, want section", got)
	}
}

func TestParseHeadingsAddsDuplicateSuffixes(t *testing.T) {
	t.Parallel()

	headings, err := parseHeadings(strings.NewReader("# Notes\n\n# Notes\n"))
	if err != nil {
		t.Fatalf("parseHeadings: %v", err)
	}

	if headings[0].slug != "notes" || headings[1].slug != "notes-0" {
		t.Fatalf("slugs = %q, %q", headings[0].slug, headings[1].slug)
	}
}

func TestRunFromStdin(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	code := run(
		[]string{"-max-depth", "2"},
		strings.NewReader("# Title\n\n## Install\n\n### Details\n"),
		&stdout,
		&stderr,
	)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0; stderr = %s", code, stderr.String())
	}

	want := "- [Title](#title)\n  - [Install](#install)\n"
	if stdout.String() != want {
		t.Fatalf("stdout = %q, want %q", stdout.String(), want)
	}
}

func TestRunFromFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "README.md")
	if err := os.WriteFile(path, []byte("# Title\n\n## Usage\n"), 0o600); err != nil {
		t.Fatalf("os.WriteFile: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{path}, strings.NewReader(""), &stdout, &stderr)

	if code != 0 {
		t.Fatalf("run exit code = %d, want 0; stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "[Usage](#usage)") {
		t.Fatalf("stdout = %q, want Usage", stdout.String())
	}
}

func TestRunRejectsInvalidDepth(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	code := run([]string{"-max-depth", "7"}, strings.NewReader("# Title\n"), &stdout, &stderr)

	if code != 2 {
		t.Fatalf("run exit code = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "max-depth") {
		t.Fatalf("stderr = %q, want max-depth", stderr.String())
	}
}
