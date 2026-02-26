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

### Render URL to PNG

```go
png, err := client.RenderURL("https://example.com").
	Format(forge.FormatPNG).
	Width(1280).
	Height(800).
	Send(ctx)
```

### Color Quantization

```go
eink, err := client.RenderHTML("<h1>Dashboard</h1>").
	Format(forge.FormatPNG).
	Palette(forge.PaletteEink).
	Dither(forge.DitherFloydSteinberg).
	Send(ctx)
```

### Custom Timeout

```go
client := forge.NewClient("http://forge:3000",
	forge.WithTimeout(120 * time.Second),
)
```

### Health Check

```go
ok, err := client.Health(ctx)
```

## API Reference

### Types

```go
type OutputFormat string  // FormatPDF, FormatPNG, FormatJPEG, FormatBMP, FormatTGA, FormatQOI, FormatSVG
type Orientation string   // Portrait, Landscape
type Flow string          // FlowAuto, FlowPaginate, FlowContinuous
type DitherMethod string  // DitherNone, DitherFloydSteinberg, DitherAtkinson, DitherOrdered
type Palette string       // PaletteAuto, PaletteBlackWhite, PaletteGrayscale, PaletteEink
```

### Errors

- `*ServerError` — 4xx/5xx server responses (has `StatusCode` and `Message`)
- `*ConnectionError` — network/connection failures (wraps cause)

## License

MIT
