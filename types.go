package forge

// OutputFormat specifies the rendered output format.
type OutputFormat string

const (
	FormatPDF  OutputFormat = "pdf"
	FormatPNG  OutputFormat = "png"
	FormatJPEG OutputFormat = "jpeg"
	FormatBMP  OutputFormat = "bmp"
	FormatTGA  OutputFormat = "tga"
	FormatQOI  OutputFormat = "qoi"
	FormatSVG  OutputFormat = "svg"
)

// Orientation specifies page orientation.
type Orientation string

const (
	Portrait  Orientation = "portrait"
	Landscape Orientation = "landscape"
)

// Flow specifies the document flow mode.
type Flow string

const (
	FlowAuto       Flow = "auto"
	FlowPaginate   Flow = "paginate"
	FlowContinuous Flow = "continuous"
)

// DitherMethod specifies the dithering algorithm.
type DitherMethod string

const (
	DitherNone          DitherMethod = "none"
	DitherFloydSteinberg DitherMethod = "floyd-steinberg"
	DitherAtkinson      DitherMethod = "atkinson"
	DitherOrdered       DitherMethod = "ordered"
)

// WatermarkLayer specifies whether the watermark renders over or under content.
type WatermarkLayer string

const (
	WatermarkOver  WatermarkLayer = "over"
	WatermarkUnder WatermarkLayer = "under"
)

// PdfStandard represents a PDF standard compliance level.
type PdfStandard string

const (
	PdfStandardNone PdfStandard = "none"
	PdfStandardA2B  PdfStandard = "pdf/a-2b"
	PdfStandardA3B  PdfStandard = "pdf/a-3b"
)

// EmbedRelationship represents the relationship of an embedded file to the PDF.
type EmbedRelationship string

const (
	EmbedRelationshipAlternative EmbedRelationship = "alternative"
	EmbedRelationshipSupplement  EmbedRelationship = "supplement"
	EmbedRelationshipData        EmbedRelationship = "data"
	EmbedRelationshipSource      EmbedRelationship = "source"
	EmbedRelationshipUnspecified EmbedRelationship = "unspecified"
)

// EmbeddedFile represents a file to embed in the PDF.
type EmbeddedFile struct {
	Path         string
	Data         string // base64-encoded
	MimeType     string
	Description  string
	Relationship EmbedRelationship
}

// Palette specifies a built-in color palette preset.
type Palette string

const (
	PaletteAuto       Palette = "auto"
	PaletteBlackWhite Palette = "bw"
	PaletteGrayscale  Palette = "grayscale"
	PaletteEink       Palette = "eink"
)
