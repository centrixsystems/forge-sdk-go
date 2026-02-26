# forge-sdk-go

Go SDK for the [Forge](https://github.com/centrixsystems/forge) rendering engine. Converts HTML/CSS to PDF, PNG, and other formats via a running Forge server.

Zero dependencies beyond the standard library.

## Installation

```sh
go get github.com/centrixsystems/forge-sdk-go
```

## Quick Start

```go
package main

import (
	"context"
	"os"

	forge "github.com/centrixsystems/forge-sdk-go"
)

func main() {
	client := forge.NewClient("http://localhost:3000")

	pdf, err := client.RenderHTML("<h1>Invoice #1234</h1>").
		Format(forge.FormatPDF).
		Paper("a4").
		Send(context.Background())
	if err != nil {
		panic(err)
	}

	os.WriteFile("invoice.pdf", pdf, 0644)
}
```

## Usage

### Render HTML to PDF

```go
pdf, err := client.RenderHTML("<h1>Hello</h1>").
	Format(forge.FormatPDF).
	Paper("a4").
	Orientation(forge.Portrait).
	Margins("25.4,25.4,25.4,25.4").
	Flow(forge.FlowPaginate).
	Send(ctx)
```

### Render URL to PNG

```go
png, err := client.RenderURL("https://example.com").
	Format(forge.FormatPNG).
	Width(1280).
	Height(800).
	Density(2.0).
	Send(ctx)
```

### Color Quantization

Reduce colors for e-ink displays or limited-palette output.

```go
eink, err := client.RenderHTML("<h1>Dashboard</h1>").
	Format(forge.FormatPNG).
	Palette(forge.PaletteEink).
	Dither(forge.DitherFloydSteinberg).
	Send(ctx)
```

### Custom Palette

```go
img, err := client.RenderHTML("<h1>Brand</h1>").
	Format(forge.FormatPNG).
	CustomPalette([]string{"#000000", "#ffffff", "#ff0000"}).
	Dither(forge.DitherAtkinson).
	Send(ctx)
```

### PDF Metadata

Set PDF document properties and enable bookmarks.

```go
pdf, err := client.RenderHTML("<h1>Annual Report</h1><p>Contents...</p>").
	Format(forge.FormatPDF).
	Paper("a4").
	Flow(forge.FlowPaginate).
	PdfTitle("Annual Report 2026").
	PdfAuthor("Centrix Systems").
	PdfSubject("Financial Summary").
	PdfKeywords("finance,report,annual").
	PdfCreator("Forge SDK").
	PdfBookmarks(true).
	Send(ctx)
```

### PDF Watermarks

Add text or image watermarks to each page.

```go
pdf, err := client.RenderHTML("<h1>Draft Report</h1>").
	Format(forge.FormatPDF).
	PdfWatermarkText("DRAFT").
	PdfWatermarkOpacity(0.15).
	PdfWatermarkRotation(-45).
	PdfWatermarkColor("#888888").
	PdfWatermarkLayer(forge.WatermarkOver).
	Send(ctx)
```

### Custom Client Configuration

```go
import "time"

client := forge.NewClient("http://forge:3000",
	forge.WithTimeout(5 * time.Minute),
)

// Or bring your own http.Client
client := forge.NewClient("http://forge:3000",
	forge.WithHTTPClient(myHTTPClient),
)
```

### Health Check

```go
ok, err := client.Health(ctx)
```

Returns `(true, nil)` when the server is healthy, `(false, *ConnectionError)` when the server is unreachable.

## API Reference

### `Client`

| Function / Method | Description |
|--------------------|-------------|
| `NewClient(baseURL, ...Option)` | Create a client (default 120s timeout) |
| `client.RenderHTML(html)` | Start a render request from an HTML string |
| `client.RenderURL(url)` | Start a render request from a URL |
| `client.Health(ctx)` | Check server health |

### Options

| Function | Description |
|----------|-------------|
| `WithTimeout(d)` | Set HTTP request timeout (`time.Duration`) |
| `WithHTTPClient(hc)` | Use a custom `*http.Client` |

### `RenderRequest`

All methods return `*RenderRequest` for chaining. Call `.Send(ctx)` to execute.

| Method | Type | Description |
|--------|------|-------------|
| `Format` | `OutputFormat` | Output format (default: `FormatPDF`) |
| `Width` | `int` | Viewport width in CSS pixels |
| `Height` | `int` | Viewport height in CSS pixels |
| `Paper` | `string` | Paper size: a3, a4, a5, b4, b5, letter, legal, ledger |
| `Orientation` | `Orientation` | `Portrait` or `Landscape` |
| `Margins` | `string` | Preset (`default`, `none`, `narrow`) or `"T,R,B,L"` in mm |
| `Flow` | `Flow` | `FlowAuto`, `FlowPaginate`, or `FlowContinuous` |
| `Density` | `float64` | Output DPI (default: 96) |
| `Background` | `string` | CSS background color (e.g. `"#ffffff"`) |
| `Timeout` | `int` | Page load timeout in seconds |
| `Colors` | `int` | Quantization color count (2-256) |
| `Palette` | `Palette` | Built-in color palette preset |
| `CustomPalette` | `[]string` | Array of hex color strings |
| `Dither` | `DitherMethod` | Dithering algorithm |
| `PdfTitle` | `string` | PDF document title metadata |
| `PdfAuthor` | `string` | PDF document author metadata |
| `PdfSubject` | `string` | PDF document subject metadata |
| `PdfKeywords` | `string` | PDF keywords metadata (comma-separated) |
| `PdfCreator` | `string` | PDF creator application metadata |
| `PdfBookmarks` | `bool` | Enable PDF bookmarks from headings |
| `PdfWatermarkText` | `string` | Watermark text on each page |
| `PdfWatermarkImage` | `string` | Base64-encoded PNG/JPEG watermark image |
| `PdfWatermarkOpacity` | `float64` | Watermark opacity (0.0-1.0, default: 0.15) |
| `PdfWatermarkRotation` | `float64` | Watermark rotation in degrees (default: -45) |
| `PdfWatermarkColor` | `string` | Watermark text color as hex (default: #888888) |
| `PdfWatermarkFontSize` | `float64` | Watermark font size in PDF points (default: auto) |
| `PdfWatermarkScale` | `float64` | Watermark image scale (0.0-1.0, default: 0.5) |
| `PdfWatermarkLayer` | `WatermarkLayer` | Layer position: `WatermarkOver` or `WatermarkUnder` |
| `PdfStandard` | `PdfStandard` | PDF standard: `PdfStandardNone`, `PdfStandardA2B`, `PdfStandardA3B` |
| `PdfAttach` | `path, data string, opts...` | Embed file in PDF (base64 data) |

| Terminal Method | Returns | Description |
|-----------------|---------|-------------|
| `Send(ctx)` | `([]byte, error)` | Execute the render request |

### Type Constants

| Type | Constants |
|------|----------|
| `OutputFormat` | `FormatPDF`, `FormatPNG`, `FormatJPEG`, `FormatBMP`, `FormatTGA`, `FormatQOI`, `FormatSVG` |
| `Orientation` | `Portrait`, `Landscape` |
| `Flow` | `FlowAuto`, `FlowPaginate`, `FlowContinuous` |
| `DitherMethod` | `DitherNone`, `DitherFloydSteinberg`, `DitherAtkinson`, `DitherOrdered` |
| `Palette` | `PaletteAuto`, `PaletteBlackWhite`, `PaletteGrayscale`, `PaletteEink` |
| `WatermarkLayer` | `WatermarkOver`, `WatermarkUnder` |
| `PdfStandard` | `PdfStandardNone`, `PdfStandardA2B`, `PdfStandardA3B` |
| `EmbedRelationship` | `EmbedRelationshipAlternative`, `EmbedRelationshipSupplement`, `EmbedRelationshipData`, `EmbedRelationshipSource`, `EmbedRelationshipUnspecified` |

### Errors

| Type | Fields | Description |
|------|--------|-------------|
| `*ServerError` | `StatusCode int`, `Message string` | Server returned 4xx/5xx |
| `*ConnectionError` | `Cause error` | Network failure (implements `Unwrap()`) |

## Requirements

- Go 1.21+
- A running [Forge](https://github.com/centrixsystems/forge) server

## License

MIT
