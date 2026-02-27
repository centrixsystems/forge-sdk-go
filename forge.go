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
	pdfBookmarks    *bool
	pdfPageNumbers  *bool
	pdfWatermarkText     *string
	pdfWatermarkImage    *string  // base64-encoded
	pdfWatermarkOpacity  *float64
	pdfWatermarkRotation *float64
	pdfWatermarkColor    *string
	pdfWatermarkFontSize *float64
	pdfWatermarkScale    *float64
	pdfWatermarkLayer    *string
	pdfWatermarkPages    *string
	pdfStandard          *PdfStandard
	pdfEmbeddedFiles     []EmbeddedFile
	pdfBarcodes          []BarcodeConfig
	pdfMode              *string
	pdfSignCertificate   *string
	pdfSignPassword      *string
	pdfSignName          *string
	pdfSignReason        *string
	pdfSignLocation      *string
	pdfSignTimestampUrl  *string
	pdfUserPassword      *string
	pdfOwnerPassword     *string
	pdfPermissions       *string
	pdfAccessibility     *string
	pdfLinearize         *bool
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

// PdfPageNumbers enables or disables "Page X of Y" footers on each page.
func (r *RenderRequest) PdfPageNumbers(enabled bool) *RenderRequest {
	r.pdfPageNumbers = &enabled
	return r
}

// PdfWatermarkText sets the watermark text overlay on each PDF page.
func (r *RenderRequest) PdfWatermarkText(text string) *RenderRequest {
	r.pdfWatermarkText = &text
	return r
}

// PdfWatermarkImage sets the watermark image (base64-encoded PNG/JPEG).
func (r *RenderRequest) PdfWatermarkImage(base64Data string) *RenderRequest {
	r.pdfWatermarkImage = &base64Data
	return r
}

// PdfWatermarkOpacity sets the watermark opacity (0.0-1.0, default 0.15).
func (r *RenderRequest) PdfWatermarkOpacity(opacity float64) *RenderRequest {
	r.pdfWatermarkOpacity = &opacity
	return r
}

// PdfWatermarkRotation sets the watermark rotation in degrees (default -45).
func (r *RenderRequest) PdfWatermarkRotation(degrees float64) *RenderRequest {
	r.pdfWatermarkRotation = &degrees
	return r
}

// PdfWatermarkColor sets the watermark text color as hex (default "#888888").
func (r *RenderRequest) PdfWatermarkColor(hex string) *RenderRequest {
	r.pdfWatermarkColor = &hex
	return r
}

// PdfWatermarkFontSize sets the watermark font size in PDF points.
func (r *RenderRequest) PdfWatermarkFontSize(size float64) *RenderRequest {
	r.pdfWatermarkFontSize = &size
	return r
}

// PdfWatermarkScale sets the watermark image scale (0.0-1.0, default 0.5).
func (r *RenderRequest) PdfWatermarkScale(scale float64) *RenderRequest {
	r.pdfWatermarkScale = &scale
	return r
}

// PdfWatermarkLayer sets the watermark layer position.
func (r *RenderRequest) PdfWatermarkLayer(layer WatermarkLayer) *RenderRequest {
	s := string(layer)
	r.pdfWatermarkLayer = &s
	return r
}

// PdfStandard sets the PDF standard compliance level.
func (r *RenderRequest) PdfStandard(standard PdfStandard) *RenderRequest {
	r.pdfStandard = &standard
	return r
}

// PdfAttach adds an embedded file to the PDF. Data must be base64-encoded.
func (r *RenderRequest) PdfAttach(path, data string, opts ...func(*EmbeddedFile)) *RenderRequest {
	ef := EmbeddedFile{Path: path, Data: data}
	for _, opt := range opts {
		opt(&ef)
	}
	r.pdfEmbeddedFiles = append(r.pdfEmbeddedFiles, ef)
	return r
}

// PdfWatermarkPages sets which pages the watermark applies to (e.g. "1,3-5").
func (r *RenderRequest) PdfWatermarkPages(pages string) *RenderRequest {
	r.pdfWatermarkPages = &pages
	return r
}

// PdfBarcode adds a barcode with the given type and data.
func (r *RenderRequest) PdfBarcode(barcodeType BarcodeType, data string) *RenderRequest {
	r.pdfBarcodes = append(r.pdfBarcodes, BarcodeConfig{Type: barcodeType, Data: data})
	return r
}

// PdfBarcodeWith adds a fully-configured barcode.
func (r *RenderRequest) PdfBarcodeWith(config BarcodeConfig) *RenderRequest {
	r.pdfBarcodes = append(r.pdfBarcodes, config)
	return r
}

// PdfMode sets the PDF rendering mode (auto, vector, or raster).
func (r *RenderRequest) PdfMode(mode PdfMode) *RenderRequest {
	s := string(mode)
	r.pdfMode = &s
	return r
}

// PdfSignCertificate sets the base64-encoded PKCS#12 certificate for PDF signing.
func (r *RenderRequest) PdfSignCertificate(data string) *RenderRequest {
	r.pdfSignCertificate = &data
	return r
}

// PdfSignPassword sets the password for the PKCS#12 certificate.
func (r *RenderRequest) PdfSignPassword(password string) *RenderRequest {
	r.pdfSignPassword = &password
	return r
}

// PdfSignName sets the signer name for the PDF signature.
func (r *RenderRequest) PdfSignName(name string) *RenderRequest {
	r.pdfSignName = &name
	return r
}

// PdfSignReason sets the reason for the PDF signature.
func (r *RenderRequest) PdfSignReason(reason string) *RenderRequest {
	r.pdfSignReason = &reason
	return r
}

// PdfSignLocation sets the location for the PDF signature.
func (r *RenderRequest) PdfSignLocation(location string) *RenderRequest {
	r.pdfSignLocation = &location
	return r
}

// PdfSignTimestampUrl sets the RFC 3161 timestamp server URL for the PDF signature.
func (r *RenderRequest) PdfSignTimestampUrl(url string) *RenderRequest {
	r.pdfSignTimestampUrl = &url
	return r
}

// PdfUserPassword sets the user password for PDF encryption (required to open).
func (r *RenderRequest) PdfUserPassword(password string) *RenderRequest {
	r.pdfUserPassword = &password
	return r
}

// PdfOwnerPassword sets the owner password for PDF encryption (required to edit).
func (r *RenderRequest) PdfOwnerPassword(password string) *RenderRequest {
	r.pdfOwnerPassword = &password
	return r
}

// PdfPermissions sets the PDF permission flags (comma-separated, e.g. "print,copy").
func (r *RenderRequest) PdfPermissions(permissions string) *RenderRequest {
	r.pdfPermissions = &permissions
	return r
}

// PdfAccessibility sets the PDF accessibility compliance level.
func (r *RenderRequest) PdfAccessibility(level AccessibilityLevel) *RenderRequest {
	s := string(level)
	r.pdfAccessibility = &s
	return r
}

// PdfLinearize enables or disables PDF linearization (fast web view).
func (r *RenderRequest) PdfLinearize(enabled bool) *RenderRequest {
	r.pdfLinearize = &enabled
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

	hasWatermark := r.pdfWatermarkText != nil || r.pdfWatermarkImage != nil ||
		r.pdfWatermarkOpacity != nil || r.pdfWatermarkRotation != nil ||
		r.pdfWatermarkColor != nil || r.pdfWatermarkFontSize != nil ||
		r.pdfWatermarkScale != nil || r.pdfWatermarkLayer != nil ||
		r.pdfWatermarkPages != nil

	hasSignature := r.pdfSignCertificate != nil || r.pdfSignPassword != nil ||
		r.pdfSignName != nil || r.pdfSignReason != nil || r.pdfSignLocation != nil ||
		r.pdfSignTimestampUrl != nil

	hasEncryption := r.pdfUserPassword != nil || r.pdfOwnerPassword != nil ||
		r.pdfPermissions != nil

	if r.pdfTitle != nil || r.pdfAuthor != nil || r.pdfSubject != nil ||
		r.pdfKeywords != nil || r.pdfCreator != nil || r.pdfBookmarks != nil ||
		r.pdfPageNumbers != nil || hasWatermark ||
		r.pdfStandard != nil || len(r.pdfEmbeddedFiles) > 0 || len(r.pdfBarcodes) > 0 ||
		r.pdfMode != nil || hasSignature || hasEncryption || r.pdfAccessibility != nil ||
		r.pdfLinearize != nil {
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
		if r.pdfPageNumbers != nil {
			pdf["page_numbers"] = *r.pdfPageNumbers
		}
		if hasWatermark {
			wm := map[string]any{}
			if r.pdfWatermarkText != nil {
				wm["text"] = *r.pdfWatermarkText
			}
			if r.pdfWatermarkImage != nil {
				wm["image_data"] = *r.pdfWatermarkImage
			}
			if r.pdfWatermarkOpacity != nil {
				wm["opacity"] = *r.pdfWatermarkOpacity
			}
			if r.pdfWatermarkRotation != nil {
				wm["rotation"] = *r.pdfWatermarkRotation
			}
			if r.pdfWatermarkColor != nil {
				wm["color"] = *r.pdfWatermarkColor
			}
			if r.pdfWatermarkFontSize != nil {
				wm["font_size"] = *r.pdfWatermarkFontSize
			}
			if r.pdfWatermarkScale != nil {
				wm["scale"] = *r.pdfWatermarkScale
			}
			if r.pdfWatermarkLayer != nil {
				wm["layer"] = *r.pdfWatermarkLayer
			}
			if r.pdfWatermarkPages != nil {
				wm["pages"] = *r.pdfWatermarkPages
			}
			pdf["watermark"] = wm
		}
		if r.pdfStandard != nil {
			pdf["standard"] = string(*r.pdfStandard)
		}
		if len(r.pdfEmbeddedFiles) > 0 {
			files := make([]map[string]interface{}, len(r.pdfEmbeddedFiles))
			for i, ef := range r.pdfEmbeddedFiles {
				f := map[string]interface{}{
					"path": ef.Path,
					"data": ef.Data,
				}
				if ef.MimeType != "" {
					f["mime_type"] = ef.MimeType
				}
				if ef.Description != "" {
					f["description"] = ef.Description
				}
				if ef.Relationship != "" {
					f["relationship"] = string(ef.Relationship)
				}
				files[i] = f
			}
			pdf["embedded_files"] = files
		}
		if len(r.pdfBarcodes) > 0 {
			barcodes := make([]map[string]interface{}, len(r.pdfBarcodes))
			for i, bc := range r.pdfBarcodes {
				b := map[string]interface{}{
					"type": string(bc.Type),
					"data": bc.Data,
				}
				if bc.X != nil {
					b["x"] = *bc.X
				}
				if bc.Y != nil {
					b["y"] = *bc.Y
				}
				if bc.Width != nil {
					b["width"] = *bc.Width
				}
				if bc.Height != nil {
					b["height"] = *bc.Height
				}
				if bc.Anchor != nil {
					b["anchor"] = string(*bc.Anchor)
				}
				if bc.Foreground != nil {
					b["foreground"] = *bc.Foreground
				}
				if bc.Background != nil {
					b["background"] = *bc.Background
				}
				if bc.DrawBg != nil {
					b["draw_background"] = *bc.DrawBg
				}
				if bc.Pages != nil {
					b["pages"] = *bc.Pages
				}
				barcodes[i] = b
			}
			pdf["barcodes"] = barcodes
		}
		if r.pdfMode != nil {
			pdf["mode"] = *r.pdfMode
		}
		if hasSignature {
			sig := map[string]any{}
			if r.pdfSignCertificate != nil {
				sig["certificate_data"] = *r.pdfSignCertificate
			}
			if r.pdfSignPassword != nil {
				sig["password"] = *r.pdfSignPassword
			}
			if r.pdfSignName != nil {
				sig["signer_name"] = *r.pdfSignName
			}
			if r.pdfSignReason != nil {
				sig["reason"] = *r.pdfSignReason
			}
			if r.pdfSignLocation != nil {
				sig["location"] = *r.pdfSignLocation
			}
			if r.pdfSignTimestampUrl != nil {
				sig["timestamp_url"] = *r.pdfSignTimestampUrl
			}
			pdf["signature"] = sig
		}
		if hasEncryption {
			enc := map[string]any{}
			if r.pdfUserPassword != nil {
				enc["user_password"] = *r.pdfUserPassword
			}
			if r.pdfOwnerPassword != nil {
				enc["owner_password"] = *r.pdfOwnerPassword
			}
			if r.pdfPermissions != nil {
				enc["permissions"] = *r.pdfPermissions
			}
			pdf["encryption"] = enc
		}
		if r.pdfAccessibility != nil {
			pdf["accessibility"] = *r.pdfAccessibility
		}
		if r.pdfLinearize != nil {
			pdf["linearize"] = *r.pdfLinearize
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
