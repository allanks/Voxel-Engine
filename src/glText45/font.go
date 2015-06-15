// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package glText45

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"os"

	"github.com/allanks/Voxel-Engine/src/Graphics"
	"github.com/go-gl/glh"
	"github.com/go-gl/glow/gl-core/4.5/gl"
)

// A Font allows rendering of text to an OpenGL context.
type Font struct {
	config         *FontConfig // Character set for this font.
	texture        uint32      // Holds the glyph texture id.
	maxGlyphWidth  int         // Largest glyph width.
	maxGlyphHeight int         // Largest glyph height.
}

// loadFont loads the given font data. This does not deal with font scaling.
// Scaling should be handled by the independent Bitmap/Truetype loaders.
// We therefore expect the supplied image and charset to already be adjusted
// to the correct font scale.
//
// The image should hold a sprite sheet, defining the graphical layout for
// every glyph. The config describes font metadata.
func loadFont(img *image.RGBA, config *FontConfig, vertexBuffer, textureDataStorageBlock uint32) (f *Font) {
	f = new(Font)
	f.config = config

	w, err := os.Create("resource/fonts/texture.png")
	if err != nil {
		panic(err)
	}
	png.Encode(w, img)

	// Resize image to next power-of-two.
	img = glh.Pow2Image(img).(*image.RGBA)
	ib := img.Bounds()

	f.texture, err = Graphics.NewTexture("resource/texture/textureAtlas.png")
	//f.texture, err = Graphics.NewTexture("resource/fonts/texture.png")
	if err != nil {
		panic(err)
	}

	// Create the texture itself. It will contain all glyphs.
	// Individual glyph-quads display a subset of this texture.
	gl.GenTextures(1, &f.texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.BindTexture(gl.TEXTURE_2D, f.texture)

	vertexData := []float32{}
	uvData := []float32{}

	texWidth := float32(ib.Dx())
	texHeight := float32(ib.Dy())

	for _, glyph := range config.Glyphs {
		// Update max glyph bounds.
		if glyph.Width > f.maxGlyphWidth {
			f.maxGlyphWidth = glyph.Width
		}

		if glyph.Height > f.maxGlyphHeight {
			f.maxGlyphHeight = glyph.Height
		}

		// Quad width/height
		vw := float32(glyph.Width)
		vh := float32(glyph.Height)

		// Texture coordinate offsets.
		tx1 := float32(glyph.X) / texWidth
		ty1 := float32(glyph.Y) / texHeight
		tx2 := (float32(glyph.X) + vw) / texWidth
		ty2 := (float32(glyph.Y) + vh) / texHeight

		uvData = append(uvData, tx1, ty2, tx2, ty2, tx2, ty1, tx1, ty1)
	}
	vertexData = []float32{0, 0, float32(f.maxGlyphWidth), 0, float32(f.maxGlyphWidth), float32(f.maxGlyphHeight), 0, float32(f.maxGlyphHeight)}

	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexData)*4, gl.Ptr(vertexData), gl.STATIC_DRAW)

	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, textureDataStorageBlock)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(uvData)*4, gl.Ptr(uvData), gl.STATIC_DRAW)
	gl.BufferSubData(gl.SHADER_STORAGE_BUFFER, 0, len(uvData)*4, gl.Ptr(uvData))

	return f
}

// Printf draws the given string at the specified coordinates.
// It expects the string to be a single line. Line breaks are not
// handled as line breaks and are rendered as glyphs.
//
// In order to render multi-line text, it is up to the caller to split
// the text up into individual lines of adequate length and then call
// this method for each line seperately.
func (f *Font) DisplayString(x, y float32, objectBuffer uint32, fs string, argv ...interface{}) {
	indices := []rune(fmt.Sprintf(fs, argv...))

	//fmt.Printf("Low:%v X:%v Y:%v\n", f.config.Low, x, y)
	//fmt.Printf("%v\n", indices)

	if len(indices) == 0 {
		return
	}

	instances := []float32{}

	// Runes form display list indices.
	// For this purpose, they need to be offset by -FontConfig.Low
	low := f.config.Low
	for i := range indices {
		indices[i] -= low
		instances = append(instances, x+float32(i*f.maxGlyphWidth), y, 0, float32(indices[i]))
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, objectBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(instances)*4, gl.Ptr(instances), gl.STATIC_DRAW)

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.BindTexture(gl.TEXTURE_2D, f.texture)

	gl.DrawArraysInstanced(gl.TRIANGLE_STRIP, 0, 4, int32(len(indices)))

	//fmt.Printf("Instances %v\n", instances)

}

// GlyphBounds returns the largest width and height for any of the glyphs
// in the font. This constitutes the largest possible bounding box
// a single glyph will have.
func (f *Font) GlyphBounds() (int, int) {
	return f.maxGlyphWidth, f.maxGlyphHeight
}

// Direction represents the direction in which strings should be rendered.
type Direction uint8

// Known directions.
const (
	LeftToRight Direction = iota // E.g.: Latin
	RightToLeft                  // E.g.: Arabic
	TopToBottom                  // E.g.: Chinese
)

// FontConfig describes raster font metadata.
//
// It can be loaded from, or saved to a JSON encoded file,
// which should come with any bitmap font image.
type FontConfig struct {
	// The direction determines the orientation of rendered strings and should
	// hold any of the pre-defined Direction constants.
	Dir Direction `json:"direction"`

	// Lower rune boundary
	Low rune `json:"rune_low"`

	// Upper rune boundary.
	High rune `json:"rune_high"`

	// Glyphs holds a set of glyph descriptors, defining the location,
	// size and advance of each glyph in the sprite sheet.
	Glyphs Charset `json:"glyphs"`
}

// Load reads font configuration data from the given JSON encoded stream.
func (fc *FontConfig) Load(r io.Reader) (err error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	return json.Unmarshal(data, fc)
}

// Save writes font configuration data to the given stream as JSON data.
func (fc *FontConfig) Save(w io.Writer) (err error) {
	data, err := json.MarshalIndent(fc, "", "  ")
	if err != nil {
		return
	}
	_, err = io.Copy(w, bytes.NewBuffer(data))
	return
}

// A Glyph describes metrics for a single font glyph.
// These indicate which area of a given image contains the
// glyph data and how the glyph should be spaced in a rendered string.
type Glyph struct {
	X      int `json:"x"`      // The x location of the glyph on a sprite sheet.
	Y      int `json:"y"`      // The y location of the glyph on a sprite sheet.
	Width  int `json:"width"`  // The width of the glyph on a sprite sheet.
	Height int `json:"height"` // The height of the glyph on a sprite sheet.

	// Advance determines the distance to the next glyph.
	// This is used to properly align non-monospaced fonts.
	Advance int `json:"advance"`
}

// A Charset represents a set of glyph descriptors for a font.
// Each glyph descriptor holds glyph metrics which are used to
// properly align the given glyph in the resulting rendered string.
type Charset []Glyph

// Scale scales all glyphs by the given factor and repositions them
// appropriately. A scale of 1 retains the original size. A scale of 2
// doubles the size of each glyph, etc.
//
// This is useful when the accompanying sprite sheet is scaled by the
// same factor. In this case, we want the glyph data to match up with the
// new image.
func (c Charset) Scale(factor int) {
	if factor <= 1 {
		// A factor of zero results in zero-sized glyphs and
		// is therefore not valid. A factor of 1 does not change
		// the glyphs, so we can ignore it.
		return
	}

	// Multiply each glyph field by the given factor
	// to scale them up to the new size.
	for i := range c {
		c[i].X *= factor
		c[i].Y *= factor
		c[i].Width *= factor
		c[i].Height *= factor
		c[i].Advance *= factor
	}
}

// Dir returns the font's rendering orientation.
func (f *Font) Dir() Direction { return f.config.Dir }

// Low returns the font's lower rune bound.
func (f *Font) Low() rune { return f.config.Low }

// High returns the font's upper rune bound.
func (f *Font) High() rune { return f.config.High }

// Glyphs returns the font's glyph descriptors.
func (f *Font) Glyphs() Charset { return f.config.Glyphs }

// Metrics returns the pixel width and height for the given string.
// This takes the scale and rendering direction of the font into account.
//
// Unknown runes will be counted as having the maximum glyph bounds as
// defined by Font.GlyphBounds().
func (f *Font) Metrics(text string) (int, int) {
	if len(text) == 0 {
		return 0, 0
	}

	gw, gh := f.GlyphBounds()

	if f.config.Dir == TopToBottom {
		return gw, f.advanceSize(text)
	}

	return f.advanceSize(text), gh
}

// advanceSize computes the pixel width or height for the given single-line
// input string. This iterates over all of its runes, finds the matching
// Charset entry and adds up the Advance values.
//
// Unknown runes will be counted as having the maximum glyph bounds as
// defined by Font.GlyphBounds().
func (f *Font) advanceSize(line string) int {
	gw, gh := f.GlyphBounds()
	glyphs := f.config.Glyphs
	low := f.config.Low
	indices := []rune(line)

	var size int
	for _, r := range indices {
		r -= low

		if r >= 0 && int(r) < len(glyphs) {
			size += glyphs[r].Advance
			continue
		}

		if f.config.Dir == TopToBottom {
			size += gh
		} else {
			size += gw
		}
	}

	return size
}
