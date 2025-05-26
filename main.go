package main

import (
  "encoding/json"
  "encoding/xml"
  "flag"
  "fmt"
  "html"
  "io/ioutil"
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
