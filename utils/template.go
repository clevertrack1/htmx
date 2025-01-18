package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
)

func RenderTemplate(tmpl *template.Template, view string, data interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, view, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type DocRender struct {
	DocTemplate template.HTML
}

func RenderDoc(tmpl *template.Template, data interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		doc := r.PathValue("doc")

		var docBody bytes.Buffer
		err := tmpl.ExecuteTemplate(&docBody, doc, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := DocRender{
			DocTemplate: template.HTML(docBody.String()),
		}

		if len(r.Header["Hx-Request"]) == 1 && len(r.Header["Hx-Boosted"]) == 0 {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(data.DocTemplate))
			return
		}

		err = tmpl.ExecuteTemplate(w, "docbase.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// dict creates a map from a list of key-value pairs.
// used for creating templates that conditionally render based on dicts defined inside the template
// useful for creating *.tmpl components that can be reused with different outputs without an imported package
func Dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}
