# drawio-html-converter
Generates HTML representations for embedding drawio files in GitHub, GitLab, Gitea, etc.

[日本語版](README_ja.md)

## Background
This tool was created to easily view drawio files managed on [gitea](https://gitea.io/).  
While you can manage files as `.drawio.svg`, information is lost for multi-page diagrams.  
With this embedding method, you can embed multi-page drawio files as-is.  

Of course, it can also be used outside of gitea.

## Download
* [Linux](https://github.com/arika0093/drawio-html-converter/releases/latest/download/drawio-converter-linux-amd64)
* [macOS](https://github.com/arika0093/drawio-html-converter/releases/latest/download/drawio-converter-macos-amd64)
* [Windows](https://github.com/arika0093/drawio-html-converter/releases/latest/download/drawio-converter-windows-amd64.exe)

Alternatively, you can use the Docker image from [ghcr.io](https://ghcr.io/arika0093/drawio-html-converter).

## Usage
Pass a drawio file as an argument to generate HTML that can be displayed as-is. By default, the output is sent to standard output.

```
drawio-converter [options] <drawio-file>

  -d, --dark-mode       Display in dark mode
  -o, --output <file>   Specify output file (default: stdout)
  -t, --toolbar <items> Specify toolbar items as comma-separated list (default: "pages,zoom,layers,tags")
  --js <url>            Specify external JavaScript file URL (default: "https://viewer.diagrams.net/js/viewer-static.min.js")
                        Specify blank to suppress JavaScript tag output.

  Server options:
  --server              Start in API server mode
  --port <port>         Specify port to use (default: 8080)
```

### API Server Mode

Specify the `--server` option to start in API server mode. In this mode, HTTP requests are accepted and drawio files are converted to HTML.

### GET
`GET /convert?fileUri=<URL>`
This endpoint fetches the drawio file from the URL specified in the `fileUri` query parameter and converts it to HTML.

### POST
`POST /convert`
This endpoint converts the drawio file content included in the request body to HTML.

## Usage with gitea
### Using the CLI
Download the latest binary from the releases page and save it to any folder.

```
$ mkdir -p /data/gitea/bin
$ cd /data/gitea/bin
$ curl -sSL -o drawio-converter https://github.com/arika0093/drawio-html-converter/releases/latest/download/drawio-converter-linux-amd64
$ chmod 777 drawio-converter
```

Add the following to your `app.ini`:

```ini
[markup.drawio]
ENABLED         = true
FILE_EXTENSIONS = .drawio
RENDER_COMMAND  = /data/gitea/bin/drawio-converter
IS_INPUT_FILE   = false
RENDER_CONTENT_MODE = no-sanitizer
```
Finally, after restarting the gitea server, embedded display of drawio files will be enabled.

### Using the API

Start the API server. The quickest way is to use Docker:

```bash
$ docker run -d --name drawio-converter --port 8080:8080 ghcr.io/arika0093/drawio-html-converter
```

Next, add the following to your `app.ini`:

```ini
[markup.drawio]
ENABLED         = true
FILE_EXTENSIONS = .drawio
RENDER_COMMAND  = curl -sSL -X POST -d @- http://localhost:8080/convert
IS_INPUT_FILE   = true
RENDER_CONTENT_MODE = no-sanitizer
```

## Specification
drawio files are written in XML format.
For example, a simple drawio diagram is represented by the following XML:

![Hello, World](./assets/sample.svg)

```xml
<mxfile host="65bd71144e">
  <diagram id="dKW03aIZ6vnLPfy8lMd4" name="Page 1">
    <mxGraphModel dx="1723" dy="784" grid="1" gridSize="10" guides="1" tooltips="1" connect="1" arrows="1" fold="1" page="1" pageScale="1" pageWidth="1169" pageHeight="827" math="0" shadow="0">
      <root>
        <mxCell id="0"/>
        <mxCell id="1" parent="0"/>
        <mxCell id="2" value="Hello, World" style="rounded=1;whiteSpace=wrap;html=1;" vertex="1" parent="1">
          <mxGeometry x="160" y="90" width="120" height="60" as="geometry"/>
        </mxCell>
      </root>
    </mxGraphModel>
  </diagram>
</mxfile>
```

When this diagram is output in HTML format, it looks like this:

```html
<!-- draw.io diagram -->
<div class="mxgraph" style="max-width:100%;border:1px solid transparent;" data-mxgraph="{&quot;highlight&quot;:&quot;#0000ff&quot;,&quot;lightbox&quot;:false,&quot;nav&quot;:true,&quot;resize&quot;:true,&quot;toolbar&quot;:&quot;zoom layers tags&quot;,&quot;edit&quot;:&quot;_blank&quot;,&quot;xml&quot;:&quot;&lt;mxfile&gt;\n  &lt;diagram id=\&quot;dKW03aIZ6vnLPfy8lMd4\&quot; name=\&quot;Page 1\&quot;&gt;\n    &lt;mxGraphModel dx=\&quot;618\&quot; dy=\&quot;784\&quot; grid=\&quot;1\&quot; gridSize=\&quot;10\&quot; guides=\&quot;1\&quot; tooltips=\&quot;1\&quot; connect=\&quot;1\&quot; arrows=\&quot;1\&quot; fold=\&quot;1\&quot; page=\&quot;1\&quot; pageScale=\&quot;1\&quot; pageWidth=\&quot;1169\&quot; pageHeight=\&quot;827\&quot; math=\&quot;0\&quot; shadow=\&quot;0\&quot;&gt;\n      &lt;root&gt;\n        &lt;mxCell id=\&quot;0\&quot; /&gt;\n        &lt;mxCell id=\&quot;1\&quot; parent=\&quot;0\&quot; /&gt;\n        &lt;mxCell id=\&quot;2\&quot; value=\&quot;Hello, World\&quot; style=\&quot;rounded=1;whiteSpace=wrap;html=1;\&quot; parent=\&quot;1\&quot; vertex=\&quot;1\&quot;&gt;\n          &lt;mxGeometry x=\&quot;160\&quot; y=\&quot;90\&quot; width=\&quot;120\&quot; height=\&quot;60\&quot; as=\&quot;geometry\&quot; /&gt;\n        &lt;/mxCell&gt;\n      &lt;/root&gt;\n    &lt;/mxGraphModel&gt;\n  &lt;/diagram&gt;\n&lt;/mxfile&gt;\n&quot;}"></div>
<script type="text/javascript" src="https://viewer.diagrams.net/js/viewer-static.min.js"></script>
```

For readability, here is a formatted example:

```jsx
// Written in JSX format for explanation purposes.
export default function DrawioExample() {
  const xml = `
  <mxfile>
    <diagram id="dKW03aIZ6vnLPfy8lMd4" name="Page 1">
      <mxGraphModel dx="618" dy="784" grid="1" gridSize="10" guides="1" tooltips="1" connect="1" arrows="1" fold="1" page="1" pageScale="1" pageWidth="1169" pageHeight="827" math="0" shadow="0">
      <root>
        <mxCell id="0" />
        <mxCell id="1" parent="0" />
        <mxCell id="2" value="Hello, World" style="rounded=1;whiteSpace=wrap;html=1;" parent="1" vertex="1">
        <mxGeometry x="160" y="90" width="120" height="60" as="geometry" />
        </mxCell>
      </root>
      </mxGraphModel>
    </diagram>
  </mxfile>
  `;

  const drawio = {
    "highlight":"#0000ff",
    "lightbox":false,
    "nav":true,
    "resize":true,
    "page":0,
    "dark-mode":"auto",
    "toolbar":"pages zoom layers tags lightbox",
    "edit":"_blank",
    "xml":xml
  }

  return (
    <>
      <div class="mxgraph" style="max-width:100%;border:1px solid transparent;" data-mxgraph={drawio}></div>
      <script type="text/javascript" src="https://viewer.diagrams.net/js/viewer-static.min.js"></script>
    </>
  );
}
```
In short, the XML and display options are embedded directly into the HTML `data-mxgraph` attribute, and loading `viewer.diagrams.net/js/viewer-static.min.js` enables the diagram display.

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](./LICENSE) file for details.
