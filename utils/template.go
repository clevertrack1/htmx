package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
)

// RenderTemplate renders a single HTML template and writes it to the HTTP response.
func RenderTemplate(tmpl *template.Template, view string, data interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.ExecuteTemplate(w, view, data); err != nil {
			sendErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
	}
}

// Helper data type for rendering templates with inline HTML content.
type DocRender struct {
	DocTemplate template.HTML
}

// RenderDoc handles rendering of a document (e.g., HTML with HTMX support).
func RenderDoc(tmpl *template.Template, data interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		doc := r.PathValue("doc") // Use r.URL.Path instead of r.PathValue.

		// Execute the specific document template into a buffer
		var docBody bytes.Buffer
		if err := tmpl.ExecuteTemplate(&docBody, doc, data); err != nil {
			sendErrorResponse(w, fmt.Errorf("failed to render document: %v", err), http.StatusInternalServerError)
			return
		}

		// Prepare the DocRender data structure
		renderData := DocRender{
			DocTemplate: template.HTML(docBody.String()),
		}

		// Handle HTMX-specific headers for partial updates
		if isHTMXRequest(r) {
			WriteHTMXResponse(w, renderData.DocTemplate)
			return
		}

		// Render the entire template (standard web response)
		if err := tmpl.ExecuteTemplate(w, "docbase.html", renderData); err != nil {
			sendErrorResponse(w, fmt.Errorf("failed to render base template: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

// WriteHTMXResponse handles HTMX-compatible responses with partial updates.
func WriteHTMXResponse(w http.ResponseWriter, content template.HTML) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(content)) // Directly write the HTML template content
}

// Send general HTML response with a standard structure
func WriteHTMLResponse(w http.ResponseWriter, content string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(content))
}

// Check if a request is an HTMX request by analyzing the headers.
func isHTMXRequest(r *http.Request) bool {
	// Check if the request is an HTMX request
	if r.Header.Get("Hx-Request") == "" {
		return false
	}

	// Optionally handle specific HTMX request types
	// Example: Separate based on "boosted" vs. non-boosted
	if r.Header.Get("Hx-Boosted") != "" {
		// Logic for boosted requests can go here (if needed)
		return true // or false if you only care about non-boosted HTMX requests
	}

	return true
}

// Helper function to send error responses consistently.
func sendErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	http.Error(w, fmt.Sprintf("Error: %v", err), statusCode)
}

// Dict creates a map from a list of key-value pairs.
// Used for creating templates that conditionally render based on dictionary values.
func Dict(values ...interface{}) (map[string]interface{}, error) {
	// Ensure even number of arguments
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("Dict: expected even number of arguments, got %d", len(values))
	}

	// Create the dictionary
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("Dict: key at index %d is not a string, value=%v", i, values[i])
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}
