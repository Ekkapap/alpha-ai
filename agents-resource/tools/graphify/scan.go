package main

// scan.go — walk project directory, respect .graphifyignore, emit ExtractedFile list.

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// ExtractedFile holds raw extraction results for one file.
type ExtractedFile struct {
	RelPath  string // relative to scanRoot
	Language string // "go", "ts", "js", "py", "rust", "md", "generic"
	Nodes    []RawNode
	Edges    []RawEdge
}

// RawNode is a pre-community, pre-ID node from extraction.
type RawNode struct {
	Label    string
	Location string // "L<line>"
	Kind     string // "file", "func", "type", "class", "method", "var", "const"
}

// RawEdge is a pre-ID edge.
type RawEdge struct {
	FromLabel string
	ToLabel   string
	Relation  string // "contains", "calls", "references", "imports"
	Location  string
}

var langByExt = map[string]string{
	".go":   "go",
	".ts":   "ts",
	".tsx":  "ts",
	".js":   "js",
	".jsx":  "js",
	".mjs":  "js",
	".cjs":  "js",
	".py":   "py",
	".rs":   "rust",
	".md":   "md",
	".mdx":  "md",
	".txt":  "text",
	".yaml": "yaml",
	".yml":  "yaml",
	".toml": "toml",
	".json": "json",
	".sh":   "sh",
	".bash": "sh",
}

// scanProject walks root, applies .graphifyignore, and returns ExtractedFile list.
func scanProject(scanRoot, alphaDir string) ([]ExtractedFile, error) {
	patterns, _ := loadIgnorePatterns(alphaDir)

	var files []ExtractedFile
	err := filepath.WalkDir(scanRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(scanRoot, path)
		if rel == "." {
			return nil
		}
		if shouldIgnore(rel, d.IsDir(), patterns) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		lang, ok := langByExt[ext]
		if !ok {
			return nil
		}
		// skip non-code for graph nodes (md/yaml/etc have no AST edges worth extracting)
		if lang == "md" || lang == "text" || lang == "yaml" || lang == "toml" || lang == "json" {
			files = append(files, ExtractedFile{RelPath: rel, Language: lang})
			return nil
		}
		ef, err := extractFile(path, rel, lang)
		if err != nil {
			ef = ExtractedFile{RelPath: rel, Language: lang}
		}
		files = append(files, ef)
		return nil
	})
	return files, err
}

func loadIgnorePatterns(alphaDir string) ([]string, error) {
	// Look for .graphifyignore in alpha dir or its parent
	candidates := []string{
		filepath.Join(alphaDir, ".graphifyignore"),
		filepath.Join(filepath.Dir(alphaDir), ".graphifyignore"),
	}
	for _, p := range candidates {
		f, err := os.Open(p)
		if err != nil {
			continue
		}
		defer f.Close()
		var patterns []string
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			patterns = append(patterns, line)
		}
		return patterns, nil
	}
	return nil, nil
}

func shouldIgnore(rel string, _ bool, patterns []string) bool {
	relSlash := filepath.ToSlash(rel)
	parts := strings.Split(relSlash, "/")

	// Always ignore hidden files/dirs
	for _, p := range parts {
		if strings.HasPrefix(p, ".") {
			return true
		}
	}

	filename := parts[len(parts)-1]

	for _, pat := range patterns {
		pat = strings.TrimSpace(pat)
		if pat == "" || strings.HasPrefix(pat, "#") {
			continue
		}
		if globMatch(relSlash, parts, filename, pat) {
			return true
		}
	}
	return false
}

// globMatch handles gitignore-style patterns including **, *, ?.
func globMatch(relSlash string, parts []string, filename, pat string) bool {
	pat = filepath.ToSlash(pat)

	// Pattern like **/*.ext or *.ext → match filename only
	if strings.HasPrefix(pat, "**/") {
		suffix := pat[3:] // strip **/
		if !strings.Contains(suffix, "/") {
			// **/*.json → match filename
			if ok, _ := filepath.Match(suffix, filename); ok {
				return true
			}
			// also match any path segment
			for _, seg := range parts {
				if ok, _ := filepath.Match(suffix, seg); ok {
					return true
				}
			}
			return false
		}
		// **/dir/** → match if any part equals dir
		inner := strings.TrimSuffix(strings.TrimSuffix(suffix, "/**"), "/**/*")
		inner = strings.TrimSuffix(inner, "/")
		if !strings.Contains(inner, "/") {
			for _, seg := range parts {
				if seg == inner {
					return true
				}
			}
		}
		// fall through to full match
	}

	// Pattern without ** → direct match against relSlash or filename
	clean := strings.TrimSuffix(strings.TrimSuffix(pat, "/**/*"), "/**")
	clean = strings.TrimSuffix(clean, "/")
	// Exact segment match
	for _, seg := range parts {
		if seg == clean {
			return true
		}
	}
	// filepath.Match against full path
	if ok, _ := filepath.Match(pat, relSlash); ok {
		return true
	}
	// filepath.Match against filename
	if ok, _ := filepath.Match(pat, filename); ok {
		return true
	}
	return false
}

func extractFile(path, rel, lang string) (ExtractedFile, error) {
	switch lang {
	case "go":
		return extractGo(path, rel)
	default:
		return extractGeneric(path, rel, lang)
	}
}
