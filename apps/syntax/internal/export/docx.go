package export

import (
	"archive/zip"
	"bytes"
	"fmt"
	"strings"

	"github.com/kyanite/syntax/internal/scene"
	"github.com/kyanite/syntax/internal/story"
)

// ExportDOCX exports the project to a DOCX file
func ExportDOCX(project *story.Project, scenes map[string]*scene.Scene) ([]byte, error) {
	// Sort scenes by chapter and scene number
	sortedScenes := scene.SortScenes(scenes)

	// Build document XML
	var content strings.Builder

	// Title
	content.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:jc w:val="center"/><w:rPr><w:b/><w:sz w:val="48"/></w:rPr></w:pPr><w:r><w:rPr><w:b/><w:sz w:val="48"/></w:rPr><w:t>%s</w:t></w:r></w:p>`, xmlEscape(project.Title)))

	// Author
	content.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:jc w:val="center"/><w:rPr><w:i/><w:sz w:val="32"/></w:rPr></w:pPr><w:r><w:rPr><w:i/><w:sz w:val="32"/></w:rPr><w:t>by %s</w:t></w:r></w:p>`, xmlEscape(project.Author)))

	// Empty paragraph
	content.WriteString(`<w:p><w:pPr/></w:p>`)

	// Scenes
	currentChapter := 0
	for _, sc := range sortedScenes {
		// Chapter heading
		if sc.Chapter != currentChapter {
			currentChapter = sc.Chapter
			// Page break
			content.WriteString(`<w:p><w:r><w:br w:type="page"/></w:r></w:p>`)

			// Chapter title
			content.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:rPr><w:b/><w:sz w:val="36"/></w:rPr></w:pPr><w:r><w:rPr><w:b/><w:sz w:val="36"/></w:rPr><w:t>Chapter %d</w:t></w:r></w:p>`, currentChapter))
			content.WriteString(`<w:p><w:pPr/></w:p>`)
		}

		// Scene name
		if sc.Name != "" {
			content.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:rPr><w:b/><w:sz w:val="28"/></w:rPr></w:pPr><w:r><w:rPr><w:b/><w:sz w:val="28"/></w:rPr><w:t>%s</w:t></w:r></w:p>`, xmlEscape(sc.Name)))
			content.WriteString(`<w:p><w:pPr/></w:p>`)
		}

		// Scene content - split into paragraphs
		paragraphs := strings.Split(sc.Content, "\n\n")
		for _, para := range paragraphs {
			if strings.TrimSpace(para) == "" {
				continue
			}
			content.WriteString(fmt.Sprintf(`<w:p><w:r><w:t xml:space="preserve">%s</w:t></w:r></w:p>`, xmlEscape(strings.TrimSpace(para))))
		}

		content.WriteString(`<w:p><w:pPr/></w:p>`)
	}

	// Statistics page
	content.WriteString(`<w:p><w:r><w:br w:type="page"/></w:r></w:p>`)
	content.WriteString(`<w:p><w:pPr><w:rPr><w:b/><w:sz w:val="32"/></w:rPr></w:pPr><w:r><w:rPr><w:b/><w:sz w:val="32"/></w:rPr><w:t>Story Statistics</w:t></w:r></w:p>`)
	content.WriteString(`<w:p><w:pPr/></w:p>`)

	totalWords := 0
	for _, sc := range sortedScenes {
		totalWords += sc.WordCount
	}

	content.WriteString(fmt.Sprintf(`<w:p><w:r><w:t>Total Scenes: %d</w:t></w:r></w:p>`, len(sortedScenes)))
	content.WriteString(fmt.Sprintf(`<w:p><w:r><w:t>Total Words: %d</w:t></w:r></w:p>`, totalWords))
	content.WriteString(fmt.Sprintf(`<w:p><w:r><w:t>Total Characters: %d</w:t></w:r></w:p>`, project.TotalCharacters))
	content.WriteString(fmt.Sprintf(`<w:p><w:r><w:t>Total Locations: %d</w:t></w:r></w:p>`, project.TotalLocations))

	// Create DOCX structure
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	// Add [Content_Types].xml
	if err := addZipFile(zipWriter, "[Content_Types].xml", contentTypesXML()); err != nil {
		return nil, err
	}

	// Add _rels/.rels
	if err := addZipFile(zipWriter, "_rels/.rels", relsXML()); err != nil {
		return nil, err
	}

	// Add word/_rels/document.xml.rels
	if err := addZipFile(zipWriter, "word/_rels/document.xml.rels", documentRelsXML()); err != nil {
		return nil, err
	}

	// Add word/document.xml
	if err := addZipFile(zipWriter, "word/document.xml", documentXML(content.String())); err != nil {
		return nil, err
	}

	// Close the zip writer
	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close zip: %w", err)
	}

	return buf.Bytes(), nil
}

func addZipFile(zw *zip.Writer, name, content string) error {
	f, err := zw.Create(name)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(content))
	return err
}

func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

func contentTypesXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`
}

func relsXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`
}

func documentRelsXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
</Relationships>`
}

func documentXML(content string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
%s
  </w:body>
</w:document>`, content)
}
