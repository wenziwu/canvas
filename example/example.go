package main

import (
	"fmt"
	"image/color"
	"image/png"
	"os"

	"github.com/tdewolff/canvas"
)

var dejaVuSerif canvas.Font

func main() {
	var err error
	dejaVuSerif, err = canvas.LoadFontFile("DejaVuSerif", canvas.Regular, "DejaVuSerif.woff")
	if err != nil {
		panic(err)
	}
	dejaVuSerif.Use(canvas.CommonLigatures)

	c := canvas.New(200, 80)
	Draw(c)

	////////////////

	svgFile, err := os.Create("example.svg")
	if err != nil {
		panic(err)
	}
	defer svgFile.Close()
	c.WriteSVG(svgFile)

	////////////////

	// SLOW
	pngFile, err := os.Create("example.png")
	if err != nil {
		panic(err)
	}
	defer pngFile.Close()

	img := c.WriteImage(144.0)
	err = png.Encode(pngFile, img)
	if err != nil {
		panic(err)
	}

	////////////////

	pdfFile, err := os.Create("example.pdf")
	if err != nil {
		panic(err)
	}
	defer pdfFile.Close()

	err = c.WritePDF(pdfFile)
	if err != nil {
		panic(err)
	}

	////////////////

	epsFile, err := os.Create("example.eps")
	if err != nil {
		panic(err)
	}
	defer epsFile.Close()
	c.WriteEPS(epsFile)
}

func drawText(c *canvas.C, x, y float64, face canvas.FontFace, rich *canvas.RichText) {
	metrics := face.Metrics()
	width, height := 80.0, 30.0
	text := rich.ToText(width, height, canvas.Justify, canvas.Top, 0.0)

	c.SetColor(canvas.Gainsboro)
	c.DrawPath(x, y, 0.0, text.Bounds().ToPath())
	c.SetColor(color.RGBA{0, 0, 0, 50})
	c.DrawPath(x, y, 0.0, canvas.Rectangle(0, 0, width, -metrics.LineHeight))
	c.DrawPath(x, y, 0.0, canvas.Rectangle(0, metrics.CapHeight-metrics.Ascent, width, -metrics.CapHeight-metrics.Descent))
	c.DrawPath(x, y, 0.0, canvas.Rectangle(0, metrics.XHeight-metrics.Ascent, width, -metrics.XHeight))

	c.SetColor(canvas.Black)
	c.DrawPath(x, y, 0.0, canvas.Rectangle(0.0, 0.0, width, -height).Stroke(0.2, canvas.RoundCapper, canvas.RoundJoiner))
	c.DrawText(x, y, 0.0, text)
}

func Draw(c *canvas.C) {
	face := dejaVuSerif.Face(12.0)
	rich := canvas.NewRichText()
	rich.Add(face, "Lorem ")
	rich.Add(face.FauxBold(), "ipsum ")
	rich.Add(face.FauxItalic(), "dolor ")
	rich.Add(face, "sit amet, ")
	rich.Add(face.Decoration(canvas.DoubleUnderline), "confiscatur")
	rich.Add(face, " ")
	rich.Add(face.Decoration(canvas.SineUnderline), "patria")
	rich.Add(face, " ")
	rich.Add(face.Decoration(canvas.SawtoothUnderline), "est")
	rich.Add(face, " ")
	rich.Add(face.Decoration(canvas.DottedUnderline), "gravus")
	rich.Add(face, " ")
	rich.Add(face.Decoration(canvas.DashedUnderline), "instantum")
	rich.Add(face, " ")
	rich.Add(face.Decoration(canvas.Underline), "norpe")
	rich.Add(face, " ")
	rich.Add(face.Decoration(canvas.Strikethrough), "targe")
	rich.Add(face, " ")
	rich.Add(face.Decoration(canvas.Overline), "yatum")
	drawText(c, 10, 70, face, rich)

	face = dejaVuSerif.Face(80.0)
	p := canvas.NewText(face, "Stroke").ToPath()
	c.DrawPath(5, 10, 0.0, p.Stroke(0.75, canvas.RoundCapper, canvas.RoundJoiner))

	latex, err := canvas.ParseLaTeX(`$y = \sin\left(\frac{x}{180}\pi\right)$`)
	if err != nil {
		panic(err)
	}
	latex.Rotate(-30, 0, 0)
	c.SetColor(canvas.Black)
	c.DrawPath(140, 65, 0.0, latex)

	ellipse, err := canvas.ParseSVG(fmt.Sprintf("A10 20 30 1 0 20 0z"))
	if err != nil {
		panic(err)
	}
	c.SetColor(canvas.WhiteSmoke)
	c.DrawPath(110, 40, 0.0, ellipse)
	ellipse = ellipse.Dash(2.0, 4.0, 2.0).Stroke(0.5, canvas.RoundCapper, canvas.RoundJoiner)
	c.SetColor(canvas.Black)
	c.DrawPath(110, 40, 0.0, ellipse)

	p = &canvas.Path{}
	p.LineTo(20.0, 0.0)
	p.LineTo(20.0, 10.0)
	p.LineTo(0.0, 20.0)
	p.Close()
	q := p.Smoothen()
	c.DrawPath(170, 10, 0.0, q)
	c.SetColor(canvas.Grey)
	c.DrawPath(170, 10, 0.0, p)
}
