# drawio-html-converter
drawioファイルをGitHub/GitLab/Gitea等に埋め込むためのHTML表現を生成します。

[English version](README.md)

## 背景
[gitea](https://gitea.io/)上で管理しているdrawioファイルを手軽に閲覧するために作成しました。  
`.drawio.svg`などの形でも管理は可能ですが、複数ページの場合その情報が失われてしまう課題がありました。  
この埋込み方法を使用すると、複数ページのdrawioファイルをそのまま埋め込むことができます。  
  
もちろん、gitea以外の用途でも使用可能です。

## ダウンロード
* [Linux](https://github.com/arika0093/drawio-html-converter/releases/latest/download/drawio-converter-linux-amd64)
* [macOS](https://github.com/arika0093/drawio-html-converter/releases/latest/download/drawio-converter-macos-amd64)
* [Windows](https://github.com/arika0093/drawio-html-converter/releases/latest/download/drawio-converter-windows-amd64.exe)

または、[ghcr.io](https://ghcr.io/arika0093/drawio-html-converter)のDockerイメージを使用することもできます。

## 使い方
drawioファイルを引数に渡すと、そのまま表示できるHTMLを生成します。特に指定をしない場合、標準出力にそのまま出力します。

```
drawio-converter [options] <drawio-file>

  -d, --dark-mode       ダークモードで表示する
  -o, --output <file>   出力先ファイルを指定 (省略時は標準出力)
  -t, --toolbar <items> ツールバーに表示するアイテムをカンマ区切りで指定 (デフォルト: "pages,zoom,layers,tags")
  --js <url>            外部JavaScriptファイルのURLを指定 (デフォルト: "https://viewer.diagrams.net/js/viewer-static.min.js")
                        空白を指定すると、JavaScriptタグの出力を抑制します。

  Server options:
  --server              APIサーバーモードで起動
  --port <port>         使用するポートを指定 (デフォルト: 8080)
```

### APIサーバーモード

`--server`オプションを指定すると、APIサーバーモードで起動します。このモードでは、HTTPリクエストを受け付けてdrawioファイルをHTMLに変換します。

### GET
`GET /convert?fileUri=<URL>`
このエンドポイントは、クエリパラメータ`fileUri`で指定されたURLからdrawioファイルを取得し、HTMLに変換します。

### POST
`POST /convert`
このエンドポイントは、リクエストボディに含まれるdrawioファイルの内容をHTMLに変換します。

## giteaでの使用方法
### CLIでの使用
リリースから最新のバイナリをダウンロードし、任意のフォルダに保存します。

```
$ mkdir -p /data/gitea/bin
$ cd /data/gitea/bin
$ curl -sSL -o drawio-converter https://github.com/arika0093/drawio-html-converter/releases/latest/download/drawio-converter-linux-amd64
$ chmod 777 drawio-converter
```

`app.ini` に以下を追加します。

```ini
[markup.drawio]
ENABLED         = true
FILE_EXTENSIONS = .drawio
RENDER_COMMAND  = /data/gitea/bin/drawio-converter
IS_INPUT_FILE   = false
RENDER_CONTENT_MODE = no-sanitizer
```

最後に、giteaサーバーを再起動するとdrawioファイルの埋め込み表示が有効になります。

### APIからの使用

APIサーバーを立ち上げます。最も手っ取り早い方法は、dockerを使用することです。

```bash
$ docker run -d --name drawio-converter --port 8080:8080 ghcr.io/arika0093/drawio-html-converter
```

次に、`app.ini` に以下を追加します。

```ini
[markup.drawio]
ENABLED         = true
FILE_EXTENSIONS = .drawio
RENDER_COMMAND  = curl -sSL -X POST -d @- http://localhost:8080/convert
IS_INPUT_FILE   = false
RENDER_CONTENT_MODE = no-sanitizer
```

## 仕様
drawioファイルはXML形式で記述されています。
例えば、以下の単純なdrawio画像はこのようなXMLで表現されます。

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

この画像をHTML形式で吐き出すと以下の形で出力されます。

```html
<!-- draw.io diagram -->
<div class="mxgraph" style="max-width:100%;border:1px solid transparent;" data-mxgraph="{&quot;highlight&quot;:&quot;#0000ff&quot;,&quot;lightbox&quot;:false,&quot;nav&quot;:true,&quot;resize&quot;:true,&quot;toolbar&quot;:&quot;zoom layers tags&quot;,&quot;edit&quot;:&quot;_blank&quot;,&quot;xml&quot;:&quot;&lt;mxfile&gt;\n  &lt;diagram id=\&quot;dKW03aIZ6vnLPfy8lMd4\&quot; name=\&quot;Page 1\&quot;&gt;\n    &lt;mxGraphModel dx=\&quot;618\&quot; dy=\&quot;784\&quot; grid=\&quot;1\&quot; gridSize=\&quot;10\&quot; guides=\&quot;1\&quot; tooltips=\&quot;1\&quot; connect=\&quot;1\&quot; arrows=\&quot;1\&quot; fold=\&quot;1\&quot; page=\&quot;1\&quot; pageScale=\&quot;1\&quot; pageWidth=\&quot;1169\&quot; pageHeight=\&quot;827\&quot; math=\&quot;0\&quot; shadow=\&quot;0\&quot;&gt;\n      &lt;root&gt;\n        &lt;mxCell id=\&quot;0\&quot; /&gt;\n        &lt;mxCell id=\&quot;1\&quot; parent=\&quot;0\&quot; /&gt;\n        &lt;mxCell id=\&quot;2\&quot; value=\&quot;Hello, World\&quot; style=\&quot;rounded=1;whiteSpace=wrap;html=1;\&quot; parent=\&quot;1\&quot; vertex=\&quot;1\&quot;&gt;\n          &lt;mxGeometry x=\&quot;160\&quot; y=\&quot;90\&quot; width=\&quot;120\&quot; height=\&quot;60\&quot; as=\&quot;geometry\&quot; /&gt;\n        &lt;/mxCell&gt;\n      &lt;/root&gt;\n    &lt;/mxGraphModel&gt;\n  &lt;/diagram&gt;\n&lt;/mxfile&gt;\n&quot;}"></div>
<script type="text/javascript" src="https://viewer.diagrams.net/js/viewer-static.min.js"></script>
```

読みにくいので整形すると、このようになります。

```jsx
// 説明のためにJSX形式で記述しています。
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
要するに、XML+表示オプションがそのままHTMLの`data-mxgraph`属性に埋め込まれ、`viewer.diagrams.net/js/viewer-static.min.js`を読み込むことで表示されます。

## ライセンス

このプロジェクトはApache License 2.0の下でライセンスされています。詳細は[LICENSE](./LICENSE)ファイルを参照してください。

