package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fogleman/gg"
	"github.com/pravj/geopattern/pattern"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// Config stores values passed in as command-line flags.
type config struct {
	text        string
	note        string
	pattern     string
	width       int
	height      int
	showVersion bool
}

var (
	cfg config
	cmd *flag.FlagSet

	version = "devel"
)

func init() {
	cfg = config{}
	cmd = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	cmd.StringVar(&cfg.text, "text", "Hello, world.", "Text line")
	cmd.StringVar(&cfg.note, "note", "https://www.example.com", "Footnote text")
	cmd.StringVar(&cfg.pattern, "pattern", "nested-squares", "Geopattern for the background")
	cmd.IntVar(&cfg.width, "width", 1200, "output image width")
	cmd.IntVar(&cfg.height, "height", 628, "output image height")
	cmd.BoolVar(&cfg.showVersion, "version", false, "print version number")

	cmd.Usage = usage
}

func usage() {
	fmt.Fprintf(cmd.Output(), "ogi is a tool for generation social images with a geopattern background.\n\n")
	fmt.Fprintf(cmd.Output(), "USAGE\n\n")
	fmt.Fprintf(cmd.Output(), "  $ %s [OPTIONS]\n\n", os.Args[0])
	fmt.Fprintf(cmd.Output(), "OPTIONS\n\n")
	cmd.PrintDefaults()
	fmt.Fprintf(cmd.Output(), "\n")
	fmt.Fprintf(cmd.Output(), "EXAMPLES\n\n")
	fmt.Fprintf(cmd.Output(), "  # with default options\n")
	fmt.Fprintf(cmd.Output(), "  $ %s\n\n", os.Args[0])
	fmt.Fprintf(cmd.Output(), "  # with custom title\n")
	fmt.Fprintf(cmd.Output(), "  $ %s\n\n --title=\"Example\"", os.Args[0])
}

func main() {
	cmd.Parse(os.Args[1:])

	if cfg.showVersion {
		fmt.Fprintf(cmd.Output(), "%s %s (runtime: %s)\n", os.Args[0], version, runtime.Version())
		os.Exit(0)
	}

	if err := run(); err != nil {
		fmt.Fprintf(cmd.Output(), "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := gg.NewContext(cfg.width, cfg.height)

	if err := renderBackground(ctx); err != nil {
		return fmt.Errorf("Error rendering background: %w", err)
	}
	if err := renderOverlay(ctx); err != nil {
		return fmt.Errorf("Error rendering overlay: %w", err)
	}
	if err := renderText(ctx); err != nil {
		return fmt.Errorf("Error rendering text: %w", err)
	}
	if err := renderNote(ctx); err != nil {
		return fmt.Errorf("Error rendering note: %w", err)
	}
	if err := ctx.EncodePNG(os.Stdout); err != nil {
		return fmt.Errorf("Error encoding PNG: %w", err)
	}

	return nil
}

func renderBackground(ctx *gg.Context) error {
	tile, err := generateBackgroundTile()
	if err != nil {
		return fmt.Errorf("Error generating background tile: %w", err)
	}

	tileWidth := tile.Bounds().Size().X
	tileHeight := tile.Bounds().Size().Y

	// Calculate number of vertical (ny) and horizontal (nx) tiles
	nx := int(math.Ceil(float64(cfg.width) / float64(tileWidth)))
	ny := int(math.Ceil(float64(cfg.height) / float64(tileHeight)))

	// Put background tiles
	for y := 0; y < ny; y++ {
		for x := 0; x < nx; x++ {
			ctx.DrawImage(tile, x*tileWidth, y*tileHeight)
		}
	}

	return nil
}

func generateBackgroundTile() (image.Image, error) {
	p := pattern.New(map[string]string{
		"phrase":    cfg.text,
		"generator": cfg.pattern,
	})

	// Generate SVG code for given pattern.
	svg := p.SvgStr()

	// NOTE: These two lines MUST BE after a call to p.SvgStr() otherwise width
	//       and height of SVG is unknown.
	svgWidth := p.Svg.GetWidth()
	svgHeight := p.Svg.GetHeight()

	// HACK: Strings "width='100%'" and "height='100%'" are hardcoded in
	//       geopattern/pattern package. Some refactoring is required to make
	//       them flexible, but it'll take some time to implement.
	svg = strings.ReplaceAll(svg, "width='100%'", fmt.Sprintf("width='%d'", svgWidth))
	svg = strings.ReplaceAll(svg, "height='100%'", fmt.Sprintf("height='%d'", svgHeight))

	// Rasterize SVG.
	icon, err := oksvg.ReadIconStream(strings.NewReader(svg))
	if err != nil {
		return nil, fmt.Errorf("Error reading SVG: %w", err)
	}

	icon.SetTarget(0, 0, float64(svgWidth), float64(svgHeight))

	rgba := image.NewRGBA(image.Rect(0, 0, svgWidth, svgHeight))
	icon.Draw(rasterx.NewDasher(svgWidth, svgHeight, rasterx.NewScannerGV(svgWidth, svgHeight, rgba, rgba.Bounds())), 1)

	return image.Image(rgba), nil
}

// Renders semi-transparend overlay rectangle.
func renderOverlay(ctx *gg.Context) error {
	borderRadius := 8.0
	margin := 30.0
	x := margin
	y := margin
	width := float64(ctx.Width()) - (2.0 * margin)
	height := float64(ctx.Height()) - (2.0 * margin)

	ctx.SetColor(color.RGBA{0, 0, 0, 204})
	ctx.DrawRoundedRectangle(x, y, width, height, borderRadius)
	ctx.Fill()

	return nil
}

// Renders the given text line (with shadow).
func renderText(ctx *gg.Context) error {
	marginLeft := 60.0
	marginTop := 90.0

	x := marginLeft
	y := marginTop

	width := float64(ctx.Width()) - 2.0*marginLeft

	textColor := color.White
	shadowColor := color.Black

	fontSize := 90.0

	fontPath := filepath.Join("fonts", "JetBrainsMono-Bold.ttf")
	if err := ctx.LoadFontFace(fontPath, fontSize); err != nil {
		return fmt.Errorf("Error loading font: %w", err)
	}

	ctx.SetColor(shadowColor)
	ctx.DrawStringWrapped(cfg.text, x+1, y+1, 0, 0, width, 1.5, gg.AlignLeft)

	ctx.SetColor(textColor)
	ctx.DrawStringWrapped(cfg.text, x, y, 0, 0, width, 1.5, gg.AlignLeft)

	return nil
}

// Renders the given note.
func renderNote(ctx *gg.Context) error {
	fontSize := 40.0

	fontPath := filepath.Join("fonts", "JetBrainsMono-Regular.ttf")
	if err := ctx.LoadFontFace(fontPath, fontSize); err != nil {
		return fmt.Errorf("Error loading font: %w", err)
	}

	marginBottom := 30.0
	_, height := ctx.MeasureString(cfg.note)
	x := 70.0
	y := float64(ctx.Height()) - height - marginBottom

	ctx.DrawString(cfg.note, x, y)

	return nil
}
