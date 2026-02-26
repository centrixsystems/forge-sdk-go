package forge

import (
	"testing"
)

func TestMinimalHTMLPayload(t *testing.T) {
	c := NewClient("http://localhost:3000")
	r := c.RenderHTML("<h1>Hi</h1>")
	p := r.buildPayload()

	if p["html"] != "<h1>Hi</h1>" {
		t.Errorf("html = %v, want <h1>Hi</h1>", p["html"])
	}
	if p["format"] != "pdf" {
		t.Errorf("format = %v, want pdf", p["format"])
	}
	if _, ok := p["url"]; ok {
		t.Error("url should not be present")
	}
	if _, ok := p["quantize"]; ok {
		t.Error("quantize should not be present")
	}
}

func TestURLPayloadWithOptions(t *testing.T) {
	c := NewClient("http://localhost:3000")
	r := c.RenderURL("https://example.com").
		Format(FormatPNG).
		Width(1280).
		Height(800).
		Paper("letter").
		Orientation(Landscape).
		Margins("10,20,10,20").
		Flow(FlowPaginate).
		Density(300.0).
		Background("#ffffff").
		Timeout(60)

	p := r.buildPayload()

	if _, ok := p["html"]; ok {
		t.Error("html should not be present")
	}
	if p["url"] != "https://example.com" {
		t.Errorf("url = %v", p["url"])
	}
	if p["format"] != "png" {
		t.Errorf("format = %v", p["format"])
	}
	if p["width"] != 1280 {
		t.Errorf("width = %v", p["width"])
	}
	if p["height"] != 800 {
		t.Errorf("height = %v", p["height"])
	}
	if p["paper"] != "letter" {
		t.Errorf("paper = %v", p["paper"])
	}
	if p["orientation"] != "landscape" {
		t.Errorf("orientation = %v", p["orientation"])
	}
	if p["flow"] != "paginate" {
		t.Errorf("flow = %v", p["flow"])
	}
	if _, ok := p["quantize"]; ok {
		t.Error("quantize should not be present")
	}
}

func TestQuantizePayload(t *testing.T) {
	c := NewClient("http://localhost:3000")
	r := c.RenderHTML("<p>test</p>").
		Format(FormatPNG).
		Colors(16).
		Palette(PaletteAuto).
		Dither(DitherFloydSteinberg)

	p := r.buildPayload()
	q, ok := p["quantize"].(map[string]any)
	if !ok {
		t.Fatal("quantize not present")
	}
	if q["colors"] != 16 {
		t.Errorf("colors = %v", q["colors"])
	}
	if q["palette"] != "auto" {
		t.Errorf("palette = %v", q["palette"])
	}
	if q["dither"] != "floyd-steinberg" {
		t.Errorf("dither = %v", q["dither"])
	}
}

func TestCustomPalette(t *testing.T) {
	c := NewClient("http://localhost:3000")
	r := c.RenderHTML("<p>test</p>").
		CustomPalette([]string{"#000000", "#ffffff", "#ff0000"}).
		Dither(DitherAtkinson)

	p := r.buildPayload()
	q, ok := p["quantize"].(map[string]any)
	if !ok {
		t.Fatal("quantize not present")
	}
	palette, ok := q["palette"].([]string)
	if !ok {
		t.Fatal("palette not a string slice")
	}
	if len(palette) != 3 {
		t.Errorf("palette len = %d, want 3", len(palette))
	}
	if q["dither"] != "atkinson" {
		t.Errorf("dither = %v", q["dither"])
	}
}

func TestNoQuantize(t *testing.T) {
	c := NewClient("http://localhost:3000")
	r := c.RenderHTML("<p>test</p>").Format(FormatPNG)
	p := r.buildPayload()
	if _, ok := p["quantize"]; ok {
		t.Error("quantize should not be present")
	}
}

func TestPdfOptionsPayload(t *testing.T) {
	c := NewClient("http://localhost:3000")
	r := c.RenderHTML("<h1>Report</h1>").
		PdfTitle("Annual Report").
		PdfAuthor("Centrix Systems").
		PdfSubject("Financial Summary").
		PdfKeywords("finance,report,2026").
		PdfCreator("Forge SDK").
		PdfBookmarks(true)

	p := r.buildPayload()
	pdf, ok := p["pdf"].(map[string]any)
	if !ok {
		t.Fatal("pdf not present")
	}
	if pdf["title"] != "Annual Report" {
		t.Errorf("title = %v", pdf["title"])
	}
	if pdf["author"] != "Centrix Systems" {
		t.Errorf("author = %v", pdf["author"])
	}
	if pdf["subject"] != "Financial Summary" {
		t.Errorf("subject = %v", pdf["subject"])
	}
	if pdf["keywords"] != "finance,report,2026" {
		t.Errorf("keywords = %v", pdf["keywords"])
	}
	if pdf["creator"] != "Forge SDK" {
		t.Errorf("creator = %v", pdf["creator"])
	}
	if pdf["bookmarks"] != true {
		t.Errorf("bookmarks = %v", pdf["bookmarks"])
	}
}

func TestPdfPartialOptions(t *testing.T) {
	c := NewClient("http://localhost:3000")
	r := c.RenderHTML("<h1>Test</h1>").
		PdfTitle("My Title").
		PdfBookmarks(false)

	p := r.buildPayload()
	pdf, ok := p["pdf"].(map[string]any)
	if !ok {
		t.Fatal("pdf not present")
	}
	if pdf["title"] != "My Title" {
		t.Errorf("title = %v", pdf["title"])
	}
	if pdf["bookmarks"] != false {
		t.Errorf("bookmarks = %v", pdf["bookmarks"])
	}
	if _, ok := pdf["author"]; ok {
		t.Error("author should not be present")
	}
	if _, ok := pdf["subject"]; ok {
		t.Error("subject should not be present")
	}
	if _, ok := pdf["keywords"]; ok {
		t.Error("keywords should not be present")
	}
	if _, ok := pdf["creator"]; ok {
		t.Error("creator should not be present")
	}
}

func TestNoPdf(t *testing.T) {
	c := NewClient("http://localhost:3000")
	r := c.RenderHTML("<p>test</p>").Format(FormatPDF)
	p := r.buildPayload()
	if _, ok := p["pdf"]; ok {
		t.Error("pdf should not be present when no pdf options set")
	}
}

func TestTrailingSlash(t *testing.T) {
	c := NewClient("http://localhost:3000/")
	if c.baseURL != "http://localhost:3000" {
		t.Errorf("baseURL = %v", c.baseURL)
	}
}
