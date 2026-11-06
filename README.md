# mdtoc

> **A table of contents for a Markdown file, generated at the speed of a pipe.**

`mdtoc` reads Markdown from a file or stdin, extracts the ATX headings, skips fenced code blocks and prints a nested table of contents with GitHub-style anchors. It is useful for READMEs, documentation sites and notes that are getting too long.

No dependencies outside the Go standard library.

## Why this exists

A long README without a table of contents is a wall of text. Writing the TOC by hand is boring and gets stale the moment a heading changes. `mdtoc` generates it from the source in one command, so the only thing to keep updated is the document itself.

## Install

```sh
go install github.com/melinoeluna4/mdtoc@latest
```

## Usage

```sh
mdtoc README.md
```

```text
- [mdtoc](#mdtoc)
  - [Why this exists](#why-this-exists)
  - [Install](#install)
  - [Usage](#usage)
```

Read from stdin, limit the depth or change the bullet:

```sh
cat README.md | mdtoc -max-depth 2
mdtoc -bullet "*" docs/guide.md
```

### Flags

| Flag | Default | Meaning |
| --- | --- | --- |
| `-max-depth` | `3` | maximum heading level to include (1–6) |
| `-bullet` | `-` | bullet marker |

### Exit codes

| Code | Meaning |
| --- | --- |
| `0` | the document was read |
| `2` | usage or I/O error |

## How it works

1. Scan the input line by line.
2. Toggle fenced-code state on ``` or ~~~ lines.
3. Parse `#` through `######` headings with a required space.
4. Slugify the heading text, adding `-0`, `-1` suffixes for duplicates.
5. Print a nested list up to `-max-depth`.

The slugger keeps Unicode letters and digits, lowercases the text and turns spaces, dashes, underscores, slashes and dots into a single hyphen.

## Development

```sh
go test ./...
go vet ./...
go build ./...
```

CI runs the same three commands on every push and pull request.

---

Part of the [Melinoe Luna portfolio](https://melinoeluna4.github.io/proyectos/mdtoc/) · Go 1.27 · standard library only.
