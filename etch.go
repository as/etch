// Package etch provides a simple facility to write graphical regression tests.
// The Assertf function provides the common case functionality. Provide it
// the test variable, along the image you have and want, and it will fail your
// case if want != have.
//
// Optionally, provide a filename to store the graphical difference as an
// uncompressed PNG if the test fails.
//
// The Extra data in the image (have but don't want) is represented in Red
// The Missing data (want, but dont have) is represented in Blue
// These can be changed by modifying Extra and Missing package variables
//
// To simplify the package, the alpha channel is ignored. A color triplet
// is equal to another if it's R,G,B values are identical.
//
// The foreground variable, FG, is what to paint on the delta image if two pixels match
// The background variable, BG, is the common background color between two images
//
// If two pixels at the same (x,y) coordinate don't match, the ambiguity is resolved
// by comparing the image you have's color value at that coordinate to the background
// color. If the color matches, the pixel you have is missing. Otherwise, it's extra.
//
// This may seem confusing, so a worked example is made available in the README
package etch

import (
	"github.com/as/frame/font"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"testing"
)

var (
	enc = png.Encoder{CompressionLevel: png.NoCompression}
	ft  = font.NewGoMono(20)
	sft = font.NewBasic(10)

	Red   = image.NewUniform(color.RGBA{255, 0, 0, 255})
	Blue  = image.NewUniform(color.RGBA{0, 0, 255, 255})
	Black = image.NewUniform(color.RGBA{0, 0, 0, 255})
	White = image.NewUniform(color.RGBA{255, 255, 255, 255})
	Gray  = image.NewUniform(color.RGBA{33, 33, 33, 255})
	Peach = image.NewUniform(color.RGBA{255, 248, 232, 255})

	FG      = White
	BG      = Peach
	Extra   = Red
	Missing = Blue
)

// Assertf fails a test if have != want, with an optional filename to write the delta
// image. By default, the pixels in blue represent the missing part of the image and
// the pixels in red represent unwanted additions.
func Assertf(t *testing.T, have, want image.Image, filename string, fm string, i ...interface{}) {
	delta, ok := Delta(have, want)
	if ok {
		return
	}
	t.Logf(fm, i...)
	if filename != "" {
		WriteFile(t, filename, Report(have, want, delta))
	}
	t.Fail()
}

// Report generates a visual summary of the actual (have) and
// expected (want) results, alongside the delta image. See
// Delta for details on the delta image format.
func Report(have, want, delta image.Image) image.Image {
	r := have.Bounds()
	r.Max.X = r.Min.X + r.Dx()*3 + 5*4
	r.Max.Y += 30
	rep := image.NewRGBA(r)
	draw.Draw(rep, r, Gray, rep.Bounds().Min, draw.Src)
	r.Min.X += 5
	s := []string{"Have", "Want", "Delta"}
	for i, src := range []image.Image{have, want, delta} {
		drawBorder(rep, r.Inset(-1), Black, image.ZP, 2)
		font.StringNBG(rep, image.Pt(r.Min.X+5, r.Max.Y-25), White, image.ZP, ft, []byte(s[i]))
		draw.Draw(rep, r, src, src.Bounds().Min, draw.Src)
		r.Min.X += want.Bounds().Dx() + 5
	}
	r.Min.X -= want.Bounds().Dx() + 5
	font.StringNBG(rep, image.Pt(r.Min.X+100-1, r.Max.Y-25-1), Black, image.ZP, sft, []byte("(Extra"))
	font.StringNBG(rep, image.Pt(r.Min.X+100, r.Max.Y-25), Extra, image.ZP, sft, []byte("(Extra"))

	font.StringNBG(rep, image.Pt(r.Min.X+100-1+45, r.Max.Y-25-1), Black, image.ZP, sft, []byte("/Missing)"))
	font.StringNBG(rep, image.Pt(r.Min.X+100+45, r.Max.Y-25), Missing, image.ZP, sft, []byte("/Missing)"))
	return rep
}

// Delta computes a difference between image a and b by
// comparing each pixel to the FG and BG colors. If a pixel
// in a and b are equal, the delta pixel is FG. Otherwise
// the pixel is either red or blue depending if its extra
// or missing respectively.
func Delta(a, b image.Image) (delta *image.RGBA, ok bool) {
	delta = image.NewRGBA(a.Bounds())
	dirty :=false 
	for y := a.Bounds().Min.Y; y < a.Bounds().Max.Y; y++ {
		for x := a.Bounds().Min.X; x < a.Bounds().Max.X; x++ {
			h := a.At(x, y)
			w := b.At(x, y)
			if EqualRGB(h, w) {
				delta.Set(x, y, FG)
				continue
			}
			dirty=true
			if EqualRGB(h, BG) {
				delta.Set(x, y, color.RGBA{0, 0, 255, 255})
			} else {
				delta.Set(x, y, color.RGBA{255, 0, 0, 255})
			}
		}
	}
	return delta, !dirty
}

// EqualRGB returns true if and only if
// the two colors share the same RGB triplets
func EqualRGB(c0, c1 color.Color) bool {
	r0, g0, b0, _ := c0.RGBA()
	r1, g1, b1, _ := c1.RGBA()
	return r0 == r1 && g0 == g1 && b0 == b1
}

func WriteFile(t *testing.T, file string, img image.Image) {
	fd, err := os.Create(file)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	defer fd.Close()
	err = enc.Encode(fd, img)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func drawBorder(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point, thick int) {
	draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+thick), src, sp, draw.Src)
	draw.Draw(dst, image.Rect(r.Min.X, r.Max.Y-thick, r.Max.X, r.Max.Y), src, sp, draw.Src)
	draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y, r.Min.X+thick, r.Max.Y), src, sp, draw.Src)
	draw.Draw(dst, image.Rect(r.Max.X-thick, r.Min.Y, r.Max.X, r.Max.Y), src, sp, draw.Src)
}
