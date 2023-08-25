package glhelp

import (
	"encoding/xml"
	"io"
	"os"

	"github.com/go-gl/mathgl/mgl32"
)

// Struct for an XML texture atlas
type TextureAtlas struct {
	XMLName     xml.Name      `xml:"TextureAtlas"`
	ImagePath   string        `xml:"imagePath,attr"`
	SubTextures []*SubTexture `xml:"SubTexture"`
}

type SubTexture struct {
	Name   string `xml:"name,attr"`
	X      uint   `xml:"x,attr"`
	Y      uint   `xml:"y,attr"`
	Width  uint   `xml:"width,attr"`
	Height uint   `xml:"height,attr"`
	SheetW uint
	SheetH uint
}

func (TA *TextureAtlas) GetSubTexture(filename string) (*SubTexture, int) {

	for i, st := range TA.SubTextures {
		if st.Name == filename {
			return st, i
		}
	}
	return nil, -1
}

func (st *SubTexture) GetDims() mgl32.Vec4 {

	fsw, fsh := float32(st.SheetW), float32(st.SheetH)

	spriteW := float32(st.Width)
	spriteH := float32(st.Height)

	spriteX := float32(st.X)
	spriteY := float32(st.Y)

	return mgl32.Vec4{
		spriteX / fsw,
		spriteY / fsh,
		(spriteX + spriteW) / fsw,
		(spriteY + spriteH) / fsh,
	}
}

func LoadTextureAtlas(filePath string) (*TextureAtlas, error) {
	// Open our xmlFile
	xmlFile, err := os.Open(filePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()

	byteValue, _ := io.ReadAll(xmlFile)

	var textureAtlas TextureAtlas

	err = xml.Unmarshal(byteValue, &textureAtlas)

	if err != nil {
		return nil, err
	}

	return &textureAtlas, nil
}
