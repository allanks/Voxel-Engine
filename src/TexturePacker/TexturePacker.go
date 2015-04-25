package TexturePacker

import (
	"image"
	"image/draw"
	"image/png"
	"os"
)

func PackTextures() {
	textureAtlas := image.NewRGBA(image.Rect(0, 0, 2048, 2048))

	mr := image.Rectangle{image.Point{0, 0}, image.Point{512, 512}}
	r := image.Rectangle{image.Point{0, 0}, image.Point{512, 512}}
	img := loadPNG("SkyBox/top.png")
	draw.DrawMask(textureAtlas, r, img, image.ZP, img, mr.Min, draw.Src)

	r = image.Rectangle{image.Point{512, 0}, image.Point{1024, 512}}
	img = loadPNG("SkyBox/bottom.png")
	draw.DrawMask(textureAtlas, r, img, image.ZP, img, mr.Min, draw.Src)

	r = image.Rectangle{image.Point{0, 512}, image.Point{512, 1024}}
	img = loadPNG("SkyBox/left.png")
	draw.DrawMask(textureAtlas, r, img, image.ZP, img, mr.Min, draw.Src)

	r = image.Rectangle{image.Point{512, 512}, image.Point{1024, 1024}}
	img = loadPNG("SkyBox/right.png")
	draw.DrawMask(textureAtlas, r, img, image.ZP, img, mr.Min, draw.Src)

	r = image.Rectangle{image.Point{0, 1024}, image.Point{512, 1536}}
	img = loadPNG("SkyBox/front.png")
	draw.DrawMask(textureAtlas, r, img, image.ZP, img, mr.Min, draw.Src)

	r = image.Rectangle{image.Point{512, 1024}, image.Point{1024, 1536}}
	img = loadPNG("SkyBox/back.png")
	draw.DrawMask(textureAtlas, r, img, image.ZP, img, mr.Min, draw.Src)

	r = image.Rectangle{image.Point{1024, 0}, image.Point{1536, 512}}
	img = loadPNG("CobbleStone/cobblestone.png")
	draw.DrawMask(textureAtlas, r, img, image.ZP, img, mr.Min, draw.Src)

	r = image.Rectangle{image.Point{1024, 512}, image.Point{1536, 1024}}
	img = loadPNG("Dirt/dirt.png")
	draw.DrawMask(textureAtlas, r, img, image.ZP, img, mr.Min, draw.Src)

	r = image.Rectangle{image.Point{1024, 1024}, image.Point{1536, 1536}}
	img = loadPNG("Grass/grass.png")
	draw.DrawMask(textureAtlas, r, img, image.ZP, img, mr.Min, draw.Src)

	r = image.Rectangle{image.Point{0, 1536}, image.Point{512, 2048}}
	img = loadPNG("Gravel/gravel.png")
	draw.DrawMask(textureAtlas, r, img, image.ZP, img, mr.Min, draw.Src)

	r = image.Rectangle{image.Point{512, 1536}, image.Point{1024, 2048}}
	img = loadPNG("Stone/stone.png")
	draw.DrawMask(textureAtlas, r, img, image.ZP, img, mr.Min, draw.Src)

	pngFile, err := os.Create("resource/texture/textureAtlas.png")
	if err != nil {
		panic(err)
	}
	png.Encode(pngFile, textureAtlas)
}

func loadPNG(filePath string) image.Image {

	var pngFile *os.File
	var err error
	var pngImg image.Image
	pngFile, err = os.Open("resource/texture/" + filePath)
	if err != nil {
		panic(err)
	}
	defer pngFile.Close()
	pngImg, err = png.Decode(pngFile)
	if err != nil {
		panic(err)
	}
	return pngImg
}
