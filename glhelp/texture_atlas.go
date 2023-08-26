package glhelp

import (
	"encoding/xml"
	"io"
	"os"

	"github.com/go-gl/mathgl/mgl32"
)

// Struct for an XML texture atlas
type AtlasDescriptor struct {
	XMLName     xml.Name      `xml:"TextureAtlas"`
	ImagePath   string        `xml:"imagePath,attr"`
	SubTextures []*SubTexture `xml:"SubTexture"`
	Texture     *Texture
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

func (TA *AtlasDescriptor) GetSubTexture(filename string) *SubTexture {

	for _, st := range TA.SubTextures {
		if st.Name == filename {
			return st
		}
	}
	return nil
}

func (st *SubTexture) GetDims(sheetW, sheetH float32) mgl32.Vec4 {

	return mgl32.Vec4{
		float32(st.X) / sheetW,
		float32(st.Y) / sheetH,
		float32(st.X+st.Width) / sheetW,
		float32(st.Y+st.Height) / sheetH,
	}
}

func LoadTextureAtlas(filePath string) (*AtlasDescriptor, error) {
	// Open our xmlFile
	xmlFile, err := os.Open(filePath)
	// if os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()

	byteValue, _ := io.ReadAll(xmlFile)

	var textureAtlas AtlasDescriptor

	err = xml.Unmarshal(byteValue, &textureAtlas)

	if err != nil {
		return nil, err
	}

	return &textureAtlas, nil
}
