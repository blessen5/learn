package export

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
)

//go:embed template.html
var htmlTemplate embed.FS

// NoteData holds the data for the HTML template.
type NoteData struct {
	Title    string
	Category string
	Date     string
	Tags     []string
	Content  template.HTML
}

// MarkdownToHTML converts markdown bytes to HTML bytes.
func MarkdownToHTML(md []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert(md, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ExtractFrontmatter parses YAML frontmatter and returns the body.
func ExtractFrontmatter(content string) (meta map[string]string, body string) {
	meta = make(map[string]string)
	body = content

	if !strings.HasPrefix(content, "---") {
		return
	}

	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return
	}

	front := strings.TrimSpace(parts[1])
	body = strings.TrimSpace(parts[2])

	for _, line := range strings.Split(front, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		kv := strings.SplitN(line, ":", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			val := strings.TrimSpace(kv[1])
			val = strings.Trim(val, "\"")
			meta[key] = val
		}
	}

	return
}

// GeneratePDF creates a PDF from a markdown file using weasyprint.
func GeneratePDF(mdPath, outputPath string) error {
	// Read markdown
	data, err := os.ReadFile(mdPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse frontmatter
	meta, body := ExtractFrontmatter(string(data))

	// Convert markdown body to HTML
	htmlBody, err := MarkdownToHTML([]byte(body))
	if err != nil {
		return fmt.Errorf("failed to convert markdown: %w", err)
	}

	// Parse tags
	var tags []string
	if t, ok := meta["tags"]; ok && t != "" && t != "[]" {
		t = strings.Trim(t, "[]")
		for _, tag := range strings.Split(t, ",") {
			tag = strings.TrimSpace(tag)
			tag = strings.Trim(tag, "\"")
			if tag != "" {
				tags = append(tags, tag)
			}
		}
	}

	// Fill template
	note := NoteData{
		Title:    meta["title"],
		Category: meta["category"],
		Date:     meta["date"],
		Tags:     tags,
		Content:  template.HTML(htmlBody),
	}

	if note.Title == "" {
		note.Title = strings.TrimSuffix(filepath.Base(mdPath), ".md")
	}

	// Render HTML
	tmplData, err := htmlTemplate.ReadFile("template.html")
	if err != nil {
		return fmt.Errorf("failed to read HTML template: %w", err)
	}

	tmpl, err := template.New("pdf").Parse(string(tmplData))
	if err != nil {
		return fmt.Errorf("failed to parse HTML template: %w", err)
	}

	var htmlBuf bytes.Buffer
	if err := tmpl.Execute(&htmlBuf, note); err != nil {
		return fmt.Errorf("failed to render HTML: %w", err)
	}

	// Write temp HTML
	tmpHTML, err := os.CreateTemp("", "learn-export-*.html")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpHTML.Name())

	if err := os.WriteFile(tmpHTML.Name(), htmlBuf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write temp HTML: %w", err)
	}
	tmpHTML.Close()

	// Convert to PDF with weasyprint
	cmd := exec.Command("weasyprint", tmpHTML.Name(), outputPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("weasyprint failed: %s\n%s", err, string(output))
	}

	return nil
}

// IsAvailable checks if weasyprint is installed.
func IsAvailable() bool {
	_, err := exec.LookPath("weasyprint")
	return err == nil
}
