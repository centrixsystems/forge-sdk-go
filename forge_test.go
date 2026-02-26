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

func TestBarcodeSimple(t *testing.T) {
	c := NewClient("http://localhost:3000")
	r := c.RenderHTML("<h1>Invoice</h1>").
		PdfBarcode(BarcodeQR, "https://example.com/invoice/123")

	p := r.buildPayload()
	pdf, ok := p["pdf"].(map[string]any)
	if !ok {
		t.Fatal("pdf not present")
	}
	barcodes, ok := pdf["barcodes"].([]map[string]interface{})
	if !ok {
		t.Fatal("barcodes not present or wrong type")
	}
	if len(barcodes) != 1 {
		t.Fatalf("barcodes len = %d, want 1", len(barcodes))
	}
	bc := barcodes[0]
	if bc["type"] != "qr" {
		t.Errorf("type = %v, want qr", bc["type"])
	}
	if bc["data"] != "https://example.com/invoice/123" {
		t.Errorf("data = %v", bc["data"])
	}
	// Optional fields should not be present
	for _, key := range []string{"x", "y", "width", "height", "anchor", "foreground", "background", "draw_background", "pages"} {
		if _, ok := bc[key]; ok {
			t.Errorf("%s should not be present", key)
		}
	}
}

func TestBarcodeWithFullConfig(t *testing.T) {
	c := NewClient("http://localhost:3000")
	x, y := 10.0, 20.0
	w, h := 100.0, 50.0
	anchor := AnchorBottomRight
	fg := "#000000"
	bg := "#ffffff"
	drawBg := true
	pages := "1,3-5"

	r := c.RenderHTML("<h1>Invoice</h1>").
		PdfBarcodeWith(BarcodeConfig{
			Type:       BarcodeCode128,
			Data:       "ABC-12345",
			X:          &x,
			Y:          &y,
			Width:      &w,
			Height:     &h,
			Anchor:     &anchor,
			Foreground: &fg,
			Background: &bg,
			DrawBg:     &drawBg,
			Pages:      &pages,
		})

	p := r.buildPayload()
	pdf, ok := p["pdf"].(map[string]any)
	if !ok {
		t.Fatal("pdf not present")
	}
	barcodes, ok := pdf["barcodes"].([]map[string]interface{})
	if !ok {
		t.Fatal("barcodes not present")
	}
	if len(barcodes) != 1 {
		t.Fatalf("barcodes len = %d, want 1", len(barcodes))
	}
	bc := barcodes[0]
	if bc["type"] != "code128" {
		t.Errorf("type = %v", bc["type"])
	}
	if bc["data"] != "ABC-12345" {
		t.Errorf("data = %v", bc["data"])
	}
	if bc["x"] != 10.0 {
		t.Errorf("x = %v", bc["x"])
	}
	if bc["y"] != 20.0 {
		t.Errorf("y = %v", bc["y"])
	}
	if bc["width"] != 100.0 {
		t.Errorf("width = %v", bc["width"])
	}
	if bc["height"] != 50.0 {
		t.Errorf("height = %v", bc["height"])
	}
	if bc["anchor"] != "bottom-right" {
		t.Errorf("anchor = %v", bc["anchor"])
	}
	if bc["foreground"] != "#000000" {
		t.Errorf("foreground = %v", bc["foreground"])
	}
	if bc["background"] != "#ffffff" {
		t.Errorf("background = %v", bc["background"])
	}
	if bc["draw_background"] != true {
		t.Errorf("draw_background = %v", bc["draw_background"])
	}
	if bc["pages"] != "1,3-5" {
		t.Errorf("pages = %v", bc["pages"])
	}
}

func TestMultipleBarcodes(t *testing.T) {
	c := NewClient("http://localhost:3000")
	r := c.RenderHTML("<h1>Product</h1>").
		PdfBarcode(BarcodeQR, "https://example.com").
		PdfBarcode(BarcodeEAN13, "4006381333931")

	p := r.buildPayload()
	pdf, ok := p["pdf"].(map[string]any)
	if !ok {
		t.Fatal("pdf not present")
	}
	barcodes, ok := pdf["barcodes"].([]map[string]interface{})
	if !ok {
		t.Fatal("barcodes not present")
	}
	if len(barcodes) != 2 {
		t.Fatalf("barcodes len = %d, want 2", len(barcodes))
	}
	if barcodes[0]["type"] != "qr" {
		t.Errorf("first type = %v", barcodes[0]["type"])
	}
	if barcodes[1]["type"] != "ean13" {
		t.Errorf("second type = %v", barcodes[1]["type"])
	}
}

func TestBarcodeOnlyTriggersPdf(t *testing.T) {
	c := NewClient("http://localhost:3000")
	r := c.RenderHTML("<p>test</p>").
		PdfBarcode(BarcodeQR, "test-data")

	p := r.buildPayload()
	pdf, ok := p["pdf"].(map[string]any)
	if !ok {
		t.Fatal("pdf should be present when barcodes are set")
	}
	if _, ok := pdf["barcodes"]; !ok {
		t.Error("barcodes should be present in pdf")
	}
	// No other pdf keys should be present
	if _, ok := pdf["title"]; ok {
		t.Error("title should not be present")
	}
	if _, ok := pdf["watermark"]; ok {
		t.Error("watermark should not be present")
	}
}

func TestWatermarkPages(t *testing.T) {
	c := NewClient("http://localhost:3000")
	r := c.RenderHTML("<h1>Draft</h1>").
		PdfWatermarkText("DRAFT").
		PdfWatermarkPages("1,3-5")

	p := r.buildPayload()
	pdf, ok := p["pdf"].(map[string]any)
	if !ok {
		t.Fatal("pdf not present")
	}
	wm, ok := pdf["watermark"].(map[string]any)
	if !ok {
		t.Fatal("watermark not present")
	}
	if wm["text"] != "DRAFT" {
		t.Errorf("text = %v", wm["text"])
	}
	if wm["pages"] != "1,3-5" {
		t.Errorf("pages = %v", wm["pages"])
	}
}

func TestWatermarkPagesOnlyTriggers(t *testing.T) {
	c := NewClient("http://localhost:3000")
	r := c.RenderHTML("<h1>Test</h1>").
		PdfWatermarkPages("2-4")

	p := r.buildPayload()
	pdf, ok := p["pdf"].(map[string]any)
	if !ok {
		t.Fatal("pdf should be present when watermark pages set")
	}
	wm, ok := pdf["watermark"].(map[string]any)
	if !ok {
		t.Fatal("watermark should be present")
	}
	if wm["pages"] != "2-4" {
		t.Errorf("pages = %v", wm["pages"])
	}
}

func TestBarcodeTypeConstants(t *testing.T) {
	tests := []struct {
		bt   BarcodeType
		want string
	}{
		{BarcodeQR, "qr"},
		{BarcodeCode128, "code128"},
		{BarcodeEAN13, "ean13"},
		{BarcodeUPCA, "upca"},
		{BarcodeCode39, "code39"},
	}
	for _, tt := range tests {
		if string(tt.bt) != tt.want {
			t.Errorf("BarcodeType %v = %q, want %q", tt.bt, string(tt.bt), tt.want)
		}
	}
}

func TestBarcodeAnchorConstants(t *testing.T) {
	tests := []struct {
		a    BarcodeAnchor
		want string
	}{
		{AnchorTopLeft, "top-left"},
		{AnchorTopRight, "top-right"},
		{AnchorBottomLeft, "bottom-left"},
		{AnchorBottomRight, "bottom-right"},
	}
	for _, tt := range tests {
		if string(tt.a) != tt.want {
			t.Errorf("BarcodeAnchor %v = %q, want %q", tt.a, string(tt.a), tt.want)
		}
	}
}

func TestTrailingSlash(t *testing.T) {
	c := NewClient("http://localhost:3000/")
	if c.baseURL != "http://localhost:3000" {
		t.Errorf("baseURL = %v", c.baseURL)
	}
}
