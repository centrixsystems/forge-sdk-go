// Package forge provides a client for the Forge rendering engine.
//
// Forge converts HTML/CSS to PDF, PNG, and other formats via an HTTP API.
package forge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client communicates with a Forge rendering server.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// Option configures a Client.
type Option func(*Client)

// WithTimeout sets the HTTP request timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = d
	}
}

// WithHTTPClient sets a custom http.Client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// NewClient creates a Forge client.
func NewClient(baseURL string, opts ...Option) *Client {
	// Strip trailing slashes.
	for len(baseURL) > 0 && baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}

	c := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// RenderHTML starts a render request from an HTML string.
func (c *Client) RenderHTML(html string) *RenderRequest {
	return &RenderRequest{client: c, html: &html}
}

// RenderURL starts a render request from a URL.
func (c *Client) RenderURL(url string) *RenderRequest {
	return &RenderRequest{client: c, url: &url}
}

// Health checks if the server is healthy.
func (c *Client) Health(ctx context.Context) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/health", nil)
	if err != nil {
		return false, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, &ConnectionError{Cause: err}
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK, nil
}

// RenderRequest builds a render request.
type RenderRequest struct {
	client      *Client
	html        *string
	url         *string
	format      string
	width       *int
	height      *int
	paper       *string
	orientation *string
	margins     *string
	flow        *string
	density     *float64
	background  *string
	timeout     *int
	colors      *int
	palette     any
	dither      *string
	pdfTitle    *string
	pdfAuthor   *string
	pdfSubject  *string
	pdfKeywords *string
	pdfCreator  *string
	pdfBookmarks *bool
}

// Format sets the output format (default: "pdf").
func (r *RenderRequest) Format(f OutputFormat) *RenderRequest {
	s := string(f)
	r.format = s
	return r
}

// Width sets the viewport width in CSS pixels.
func (r *RenderRequest) Width(px int) *RenderRequest {
	r.width = &px
	return r
}

// Height sets the viewport height in CSS pixels.
func (r *RenderRequest) Height(px int) *RenderRequest {
	r.height = &px
	return r
}

// Paper sets the paper size.
func (r *RenderRequest) Paper(size string) *RenderRequest {
	r.paper = &size
	return r
}

// Orientation sets the page orientation.
func (r *RenderRequest) Orientation(o Orientation) *RenderRequest {
	s := string(o)
	r.orientation = &s
	return r
}

// Margins sets page margins.
func (r *RenderRequest) Margins(m string) *RenderRequest {
	r.margins = &m
	return r
}

// Flow sets the document flow mode.
func (r *RenderRequest) Flow(f Flow) *RenderRequest {
	s := string(f)
	r.flow = &s
	return r
}

// Density sets the output DPI.
func (r *RenderRequest) Density(dpi float64) *RenderRequest {
	r.density = &dpi
	return r
}

// Background sets the CSS background color.
func (r *RenderRequest) Background(color string) *RenderRequest {
	r.background = &color
	return r
}

// Timeout sets the page load timeout in seconds.
func (r *RenderRequest) Timeout(seconds int) *RenderRequest {
	r.timeout = &seconds
	return r
}

// Colors sets the number of colors for quantization (2-256).
func (r *RenderRequest) Colors(n int) *RenderRequest {
	r.colors = &n
	return r
}

// Palette sets a built-in palette preset.
func (r *RenderRequest) Palette(p Palette) *RenderRequest {
	r.palette = string(p)
	return r
}

// CustomPalette sets a custom palette of hex color strings.
func (r *RenderRequest) CustomPalette(colors []string) *RenderRequest {
	r.palette = colors
	return r
}

// Dither sets the dithering algorithm.
func (r *RenderRequest) Dither(method DitherMethod) *RenderRequest {
	s := string(method)
	r.dither = &s
	return r
}

// PdfTitle sets the PDF document title metadata.
func (r *RenderRequest) PdfTitle(title string) *RenderRequest {
	r.pdfTitle = &title
	return r
}

// PdfAuthor sets the PDF document author metadata.
func (r *RenderRequest) PdfAuthor(author string) *RenderRequest {
	r.pdfAuthor = &author
	return r
}

// PdfSubject sets the PDF document subject metadata.
func (r *RenderRequest) PdfSubject(subject string) *RenderRequest {
	r.pdfSubject = &subject
	return r
}

// PdfKeywords sets the PDF document keywords metadata (comma-separated).
func (r *RenderRequest) PdfKeywords(keywords string) *RenderRequest {
	r.pdfKeywords = &keywords
	return r
}

// PdfCreator sets the PDF document creator metadata.
func (r *RenderRequest) PdfCreator(creator string) *RenderRequest {
	r.pdfCreator = &creator
	return r
}

// PdfBookmarks enables or disables PDF bookmarks from headings.
func (r *RenderRequest) PdfBookmarks(enabled bool) *RenderRequest {
	r.pdfBookmarks = &enabled
	return r
}

// buildPayload builds the JSON payload map.
func (r *RenderRequest) buildPayload() map[string]any {
	p := map[string]any{}

	if r.html != nil {
		p["html"] = *r.html
	}
	if r.url != nil {
		p["url"] = *r.url
	}

	format := r.format
	if format == "" {
		format = "pdf"
	}
	p["format"] = format

	if r.width != nil {
		p["width"] = *r.width
	}
	if r.height != nil {
		p["height"] = *r.height
	}
	if r.paper != nil {
		p["paper"] = *r.paper
	}
	if r.orientation != nil {
		p["orientation"] = *r.orientation
	}
	if r.margins != nil {
		p["margins"] = *r.margins
	}
	if r.flow != nil {
		p["flow"] = *r.flow
	}
	if r.density != nil {
		p["density"] = *r.density
	}
	if r.background != nil {
		p["background"] = *r.background
	}
	if r.timeout != nil {
		p["timeout"] = *r.timeout
	}

	if r.colors != nil || r.palette != nil || r.dither != nil {
		q := map[string]any{}
		if r.colors != nil {
			q["colors"] = *r.colors
		}
		if r.palette != nil {
			q["palette"] = r.palette
		}
		if r.dither != nil {
			q["dither"] = *r.dither
		}
		p["quantize"] = q
	}

	if r.pdfTitle != nil || r.pdfAuthor != nil || r.pdfSubject != nil ||
		r.pdfKeywords != nil || r.pdfCreator != nil || r.pdfBookmarks != nil {
		pdf := map[string]any{}
		if r.pdfTitle != nil {
			pdf["title"] = *r.pdfTitle
		}
		if r.pdfAuthor != nil {
			pdf["author"] = *r.pdfAuthor
		}
		if r.pdfSubject != nil {
			pdf["subject"] = *r.pdfSubject
		}
		if r.pdfKeywords != nil {
			pdf["keywords"] = *r.pdfKeywords
		}
		if r.pdfCreator != nil {
			pdf["creator"] = *r.pdfCreator
		}
		if r.pdfBookmarks != nil {
			pdf["bookmarks"] = *r.pdfBookmarks
		}
		p["pdf"] = pdf
	}

	return p
}

// Send executes the render request and returns the raw output bytes.
func (r *RenderRequest) Send(ctx context.Context) ([]byte, error) {
	payload := r.buildPayload()

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("forge: marshal error: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost,
		r.client.baseURL+"/render",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("forge: request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.httpClient.Do(req)
	if err != nil {
		return nil, &ConnectionError{Cause: err}
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("forge: read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error string `json:"error"`
		}
		msg := fmt.Sprintf("HTTP %d", resp.StatusCode)
		if json.Unmarshal(data, &errResp) == nil && errResp.Error != "" {
			msg = errResp.Error
		}
		return nil, &ServerError{
			StatusCode: resp.StatusCode,
			Message:    msg,
		}
	}

	return data, nil
}
