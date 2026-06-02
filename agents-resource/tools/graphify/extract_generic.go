package main

// extract_generic.go — regex-based extractor for TS/JS/Python/Rust/Shell.

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type langRules struct {
	funcRe   []*regexp.Regexp
	classRe  []*regexp.Regexp
	importRe []*regexp.Regexp
	callRe   *regexp.Regexp
}

var rules map[string]*langRules

func init() {
	tsRules := &langRules{
		funcRe: []*regexp.Regexp{
			regexp.MustCompile(`(?:export\s+)?(?:async\s+)?function\s+(\w+)\s*[(<]`),
			regexp.MustCompile(`(?:export\s+)?(?:const|let|var)\s+(\w+)\s*=\s*(?:async\s+)?\(`),
			regexp.MustCompile(`(?:export\s+)?(?:const|let|var)\s+(\w+)\s*=\s*(?:async\s+)?function`),
			regexp.MustCompile(`^\s*(?:public|private|protected|static|async|override|\s)*\s*(\w+)\s*\([^)]*\)\s*(?::\s*\S+)?\s*\{`),
		},
		classRe: []*regexp.Regexp{
			regexp.MustCompile(`(?:export\s+)?(?:abstract\s+)?class\s+(\w+)`),
		},
		importRe: []*regexp.Regexp{
			regexp.MustCompile(`import\s+.*from\s+['"]([^'"]+)['"]`),
			regexp.MustCompile(`(?:import|require)\s*\(\s*['"]([^'"]+)['"]\s*\)`),
		},
		callRe: regexp.MustCompile(`(\w+)\s*\(`),
	}

	rules = map[string]*langRules{
		"ts": tsRules,
		"js": tsRules,
		"py": {
			funcRe: []*regexp.Regexp{
				regexp.MustCompile(`^\s*(?:async\s+)?def\s+(\w+)\s*\(`),
			},
			classRe: []*regexp.Regexp{
				regexp.MustCompile(`^\s*class\s+(\w+)\s*[:(]`),
			},
			importRe: []*regexp.Regexp{
				regexp.MustCompile(`^\s*from\s+(\S+)\s+import`),
				regexp.MustCompile(`^\s*import\s+(\S+)`),
			},
			callRe: regexp.MustCompile(`(\w+)\s*\(`),
		},
		"rust": {
			funcRe: []*regexp.Regexp{
				regexp.MustCompile(`(?:pub\s+)?(?:async\s+)?fn\s+(\w+)\s*[(<]`),
			},
			classRe: []*regexp.Regexp{
				regexp.MustCompile(`(?:pub\s+)?struct\s+(\w+)`),
				regexp.MustCompile(`(?:pub\s+)?enum\s+(\w+)`),
				regexp.MustCompile(`(?:pub\s+)?trait\s+(\w+)`),
			},
			importRe: []*regexp.Regexp{
				regexp.MustCompile(`use\s+([\w:]+)`),
			},
			callRe: regexp.MustCompile(`(\w+)\s*\(`),
		},
		"sh": {
			funcRe: []*regexp.Regexp{
				regexp.MustCompile(`^(?:function\s+)?(\w+)\s*\(\s*\)\s*\{`),
			},
		},
	}
}

func extractGeneric(path, rel, lang string) (ExtractedFile, error) {
	ef := ExtractedFile{RelPath: rel, Language: lang}

	r, ok := rules[lang]
	if !ok {
		ef.Nodes = append(ef.Nodes, RawNode{Label: lastPathSegment(rel), Location: "L1", Kind: "file"})
		return ef, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return ef, err
	}
	defer f.Close()

	fileLabel := lastPathSegment(rel)
	ef.Nodes = append(ef.Nodes, RawNode{Label: fileLabel, Location: "L1", Kind: "file"})

	var currentFunc string
	sc := bufio.NewScanner(f)
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := sc.Text()
		loc := fmt.Sprintf("L%d", lineNo)

		for _, re := range r.funcRe {
			if m := re.FindStringSubmatch(line); m != nil {
				name := m[1]
				if isKeyword(name) {
					continue
				}
				label := name + "()"
				ef.Nodes = append(ef.Nodes, RawNode{Label: label, Location: loc, Kind: "func"})
				ef.Edges = append(ef.Edges, RawEdge{FromLabel: fileLabel, ToLabel: label, Relation: "contains", Location: loc})
				currentFunc = label
				break
			}
		}

		for _, re := range r.classRe {
			if m := re.FindStringSubmatch(line); m != nil {
				name := m[1]
				if isKeyword(name) {
					continue
				}
				ef.Nodes = append(ef.Nodes, RawNode{Label: name, Location: loc, Kind: "class"})
				ef.Edges = append(ef.Edges, RawEdge{FromLabel: fileLabel, ToLabel: name, Relation: "contains", Location: loc})
				break
			}
		}

		for _, re := range r.importRe {
			if m := re.FindStringSubmatch(line); m != nil {
				pkg := lastImportSegment(m[1])
				ef.Edges = append(ef.Edges, RawEdge{FromLabel: fileLabel, ToLabel: pkg, Relation: "references", Location: loc})
				break
			}
		}

		if currentFunc != "" && r.callRe != nil {
			for _, m := range r.callRe.FindAllStringSubmatch(line, -1) {
				callee := m[1] + "()"
				if callee != currentFunc && !isKeyword(m[1]) {
					ef.Edges = append(ef.Edges, RawEdge{FromLabel: currentFunc, ToLabel: callee, Relation: "calls", Location: loc})
				}
			}
		}
	}

	return ef, sc.Err()
}

func lastPathSegment(rel string) string {
	parts := strings.Split(filepath.ToSlash(rel), "/")
	return parts[len(parts)-1]
}

func lastImportSegment(imp string) string {
	imp = strings.TrimPrefix(imp, "@")
	parts := strings.Split(imp, "/")
	return parts[len(parts)-1]
}

var keywords = map[string]bool{
	"if": true, "for": true, "while": true, "switch": true, "return": true,
	"new": true, "delete": true, "typeof": true, "instanceof": true, "void": true,
	"throw": true, "catch": true, "await": true, "yield": true, "async": true,
	"print": true, "len": true, "append": true, "make": true, "panic": true,
	"println": true, "printf": true, "sprintf": true, "fmt": true,
}

func isKeyword(s string) bool { return keywords[s] || len(s) <= 1 }

var mdHeadingRe = regexp.MustCompile(`^##\s+(.+)`)

// extractMarkdown parses H2 headings from a markdown file and emits them as "section" nodes.
// Used for raw-knowledge/*.md files so agents can discover knowledge by topic.
func extractMarkdown(path, rel string) (ExtractedFile, error) {
	ef := ExtractedFile{RelPath: rel, Language: "md"}
	fileLabel := lastPathSegment(rel)
	ef.Nodes = append(ef.Nodes, RawNode{Label: fileLabel, Location: "L1", Kind: "file"})

	f, err := os.Open(path)
	if err != nil {
		return ef, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := sc.Text()
		if m := mdHeadingRe.FindStringSubmatch(line); m != nil {
			heading := strings.TrimSpace(m[1])
			if heading == "" {
				continue
			}
			loc := fmt.Sprintf("L%d", lineNo)
			ef.Nodes = append(ef.Nodes, RawNode{Label: heading, Location: loc, Kind: "section"})
			ef.Edges = append(ef.Edges, RawEdge{FromLabel: fileLabel, ToLabel: heading, Relation: "contains", Location: loc})
		}
	}
	return ef, sc.Err()
}
