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

// Palette specifies a built-in color palette preset.
type Palette string

const (
	PaletteAuto       Palette = "auto"
	PaletteBlackWhite Palette = "bw"
	PaletteGrayscale  Palette = "grayscale"
	PaletteEink       Palette = "eink"
)
