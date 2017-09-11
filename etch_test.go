package etch

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
)

func TestDifferent(t *testing.T) {
	r := image.Rect(0, 0, 200, 200)
	have, want := image.NewRGBA(r), image.NewRGBA(r)
	draw.Draw(have, r, Peach, r.Min, draw.Src)
	draw.Draw(want, r, Peach, r.Min, draw.Src)
	have.Set(25, 25, color.RGBA{33, 33, 33, 255})
	want.Set(55, 55, color.RGBA{55, 55, 55, 255})
	delta, ok := Delta(have, want)
	if !ok {
		return
	}
	t.Logf("failed: see TestDifferent.png")
	WriteFile(t, "TestDifferent.png", Report(have, want, delta))
	t.Fail()
}

func TestIdentical(t *testing.T) {
	r := image.Rect(0, 0, 200, 200)
	have, want := image.NewRGBA(r), image.NewRGBA(r)
	draw.Draw(have, r, Peach, r.Min, draw.Src)
	draw.Draw(want, r, Peach, r.Min, draw.Src)

	//	have.Set(25,25,	color.RGBA{33, 33, 33, 255})
	//	want.Set(55,55,	color.RGBA{55, 55, 55, 255})
	Assertf(t, have, want, "TestIdentical.png", "TestIdentical failed") //
}
