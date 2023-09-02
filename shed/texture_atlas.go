package shed

import (
	"encoding/xml"
	"io"
	"os"
)

// Struct for an XML texture atlas
type AtlasDescriptor struct {
	XMLName     xml.Name      `xml:"TextureAtlas"`   // Formality
	ImagePath   string        `xml:"imagePath,attr"` // Path of image file, relative to the sheet descriptor file
	SubTextures []*SubTexture `xml:"SubTexture"`     // Array of all subtextures in the atlas

	Texture *TextureWrapper // The GOAT texture object that contains the entire texture image
}

type SubTexture struct {
	Name   string `xml:"name,attr"`   // Name of the subtex
	X      uint   `xml:"x,attr"`      // x coordinate (in pixels, relative to (0, 0) in the image file
	Y      uint   `xml:"y,attr"`      // y coordinate (in pixels, relative to (0, 0) in the image file
	Width  uint   `xml:"width,attr"`  // width of the subtex
	Height uint   `xml:"height,attr"` // height of the subtex
}

// Laod a file containing a texture atlas.
// The texture image itself will not be loaded.
func LoadTextureAtlasFile(filePath string) (*AtlasDescriptor, error) {

	// Open the xml file
	xmlFile, err := os.Open(filePath)
	if err != nil { // if os.Open returns an error then handle it
		return nil, err
	}

	defer xmlFile.Close()

	byteValue, _ := io.ReadAll(xmlFile)

	var descriptor AtlasDescriptor

	err = xml.Unmarshal(byteValue, &descriptor)

	if err != nil {
		return nil, err
	}

	return &descriptor, nil
}

func (TA *AtlasDescriptor) GetSubTexture(filename string) *SubTexture {

	for _, st := range TA.SubTextures {
		if st.Name == filename {
			return st
		}
	}
	return nil
}

func (st *SubTexture) GetDims(sheetW, sheetH float32) V4 {

	return V4{
		float32(st.X) / sheetW,
		float32(st.Y) / sheetH,
		float32(st.X+st.Width) / sheetW,
		float32(st.Y+st.Height) / sheetH,
	}
}
