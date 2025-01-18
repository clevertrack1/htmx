package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/blevesearch/bleve/v2"
)

type SiteSearch struct {
	tmpl        *template.Template
	index       bleve.Index
	templateDir string
}

type HTMLFile struct {
	Path    string
	Content string
}

func InitSearchIndex(dir string, loadedTemplate *template.Template) *SiteSearch {

	s := SiteSearch{}
	s.tmpl = loadedTemplate
	s.templateDir = dir

	var err error
	indexMapping := bleve.NewIndexMapping()
	s.index, err = bleve.NewMemOnly(indexMapping)
	if err != nil {
		panic(err)
	}

	err = filepath.Walk(s.templateDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				panic(err)
			}

			cleanedContent := removeHTMLTags(string(content))
			file := strings.TrimSuffix(strings.TrimPrefix(path, s.templateDir), ".html")

			htmlFile := HTMLFile{
				Path:    file,
				Content: cleanedContent,
			}

			err = s.index.Index(path, htmlFile)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Indexed %s\n", path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return &s
}

func (s *SiteSearch) QueryIndex(search string) string {
	query := bleve.NewMatchQuery(search)
	request := bleve.NewSearchRequest(query)
	request.Highlight = bleve.NewHighlight()
	searchResults, err := s.index.Search(request)
	if err != nil {
		panic(err)
	}

	var matches []map[string]interface{}
	for _, hit := range searchResults.Hits {
		result := hit.ID
		trimmedFile := strings.TrimPrefix(strings.TrimSuffix(strings.TrimPrefix(result, s.templateDir), ".html"), "/views/")
		highlight := strings.Join(hit.Fragments["Content"], "")

		href := trimmedFile
		if href == "index" {
			href = "/"
		}

		matches = append(matches, map[string]interface{}{
			"File":      trimmedFile,
			"Highlight": template.HTML(highlight),
			"Href":      href,
		})
	}

	var searchResultsHTML bytes.Buffer
	err = s.tmpl.ExecuteTemplate(&searchResultsHTML, "searchresults", matches)
	if err != nil {
		panic(err)
	}

	return searchResultsHTML.String()
}

func removeHTMLTags(htmlContent string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(htmlContent, "")
}
