package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// DrawioFile represents the drawio file structure
type DrawioFile struct {
	XMLName xml.Name `xml:"mxfile"`
	Diagrams []Diagram `xml:"diagram"`
}

// Diagram represents a single diagram in the drawio file
type Diagram struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
	Content string `xml:",innerxml"`
}

func main() {
	// Parse command-line flags
	darkMode := flag.Bool("d", false, "Enable dark mode")
	darkModeLong := flag.Bool("dark-mode", false, "Enable dark mode")
	outputFile := flag.String("o", "", "Output file path")
	outputFileLong := flag.String("output", "", "Output file path")
	toolbar := flag.String("t", "pages,zoom,layers,tags", "Toolbar items to display (comma-separated)")
	toolbarLong := flag.String("toolbar", "pages,zoom,layers,tags", "Toolbar items to display (comma-separated)")
	jsURL := flag.String("js", "https://viewer.diagrams.net/js/viewer-static.min.js", "URL of external JavaScript file")
	serverMode := flag.Bool("server", false, "Run as an API server")
	port := flag.Int("port", 8080, "Port to run the API server on")

	flag.Parse()

	isDarkMode := *darkModeLong || *darkMode
	output := *outputFileLong
	if output == "" {
		output = *outputFile
	}
	
	toolbarItems := *toolbarLong
	if toolbarItems == "pages,zoom,layers,tags" && *toolbar != "pages,zoom,layers,tags" {
		toolbarItems = *toolbar
	}

	if *serverMode {
		// Start the server
		address := fmt.Sprintf(":%d", *port)
		fmt.Printf("Starting server on %s\n", address)
		http.HandleFunc("/convert", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				// Handle GET request with fileUri parameter
				fileUri := r.URL.Query().Get("fileUri")
				if fileUri == "" {
					http.Error(w, "fileUri parameter is required", http.StatusBadRequest)
					return
				}

				// Download the file from the provided URI
				resp, err := http.Get(fileUri)
				if err != nil {
					http.Error(w, "Failed to fetch file from URI", http.StatusInternalServerError)
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					http.Error(w, "Failed to fetch file: " + resp.Status, http.StatusInternalServerError)
					return
				}

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					http.Error(w, "Failed to read file content", http.StatusInternalServerError)
					return
				}

				// Parse and convert the file content
				toolbarItems := "pages,zoom,layers,tags" // Default toolbar items
				jsURL := "https://viewer.diagrams.net/js/viewer-static.min.js" // Default JS URL
				result := generateHTML(string(body), toolbarItems, jsURL, false)

				// Write the response
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(result))
				return
			}

			if r.Method == http.MethodPost {
				// Handle POST request (existing logic)
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Failed to read request body", http.StatusInternalServerError)
					return
				}

				// Parse the drawio XML content
				var drawioFile DrawioFile
				err = xml.Unmarshal(body, &drawioFile)
				if err != nil {
					http.Error(w, "Invalid drawio file format", http.StatusBadRequest)
					return
				}

				// Generate HTML content
				toolbarItems := "pages,zoom,layers,tags" // Default toolbar items
				jsURL := "https://viewer.diagrams.net/js/viewer-static.min.js" // Default JS URL
				result := generateHTML(string(body), toolbarItems, jsURL, false)

				// Write the response
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(result))
				return
			}

			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		})

		if err := http.ListenAndServe(address, nil); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Check if drawio file path is provided
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: drawio file path is required")
		fmt.Fprintln(os.Stderr, "Usage: drawio-converter [options] <drawio-file>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	drawioFilePath := args[0]
	
	// Read drawio file content
	content, err := ioutil.ReadFile(drawioFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Parse drawio file
	var drawioFile DrawioFile
	err = xml.Unmarshal(content, &drawioFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing drawio file: %v\n", err)
		os.Exit(1)
	}

	// Generate HTML content
	result := generateHTML(string(content), toolbarItems, *jsURL, isDarkMode)

	// Output result
	if output == "" {
		fmt.Print(result)
	} else {
		err = ioutil.WriteFile(output, []byte(result), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
			os.Exit(1)
		}
	}
}

func generateHTML(xmlContent string, toolbar string, jsURL string, darkMode bool) string {
	// Prepare data for mxgraph attribute
	mxgraph := map[string]interface{}{
		"highlight": "#0000ff",
		"lightbox":  false,
		"nav":       true,
		"resize":    true,
		"toolbar":   strings.Replace(toolbar, ",", " ", -1),
		"edit":      "_blank",
		"xml":       xmlContent,
	}

	if darkMode {
		mxgraph["dark-mode"] = true
	}

	// Convert mxgraph to JSON
	mxgraphJSON, err := json.Marshal(mxgraph)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating mxgraph JSON: %v\n", err)
		os.Exit(1)
	}

	// Escape special characters for HTML
	escapedMxgraphJSON := html.EscapeString(string(mxgraphJSON))

	// Generate HTML
	html := "<!-- draw.io diagram -->\n"
	html += "<div class=\"mxgraph\" style=\"max-width:100%;border:1px solid transparent;\" data-mxgraph='" + escapedMxgraphJSON + "'></div>\n"
	
	if jsURL != "" {
		html += "<script type=\"text/javascript\" src=\"" + jsURL + "\"></script>\n"
	}

	return html
}
