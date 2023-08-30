package motor

import (
	"fmt"
	h "goat/glhelp"
	"path"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	Machine *MachineStruct
)

type MachineStruct struct {
	Shaders          map[string]*h.ShaderProgram   // A pointer to all the shader programs currently active in the world
	SubTextureDims   map[string]mgl32.Vec4         // stores subtextures as "sheet.png/image.png" => minX, minY, maxX, maxY
	AtlasDescriptors map[string]*h.AtlasDescriptor // stores atlasses as "sheet.png", not "sheet.xml"
	Textures         map[string]*h.Texture         // Pointers to all active textures
	Sprites          map[string]*SpriteGl          // Renderables are things that can be rendered (or that can render themselves)
	Cameras          map[string]*Camera            // Contains the projection matrices. You may want to render ceretain things with one cam, and other things with another cam
	AssetPath        string                        // Base path for all assets

	// Timing
	Now64     float64
	Prev64    float64
	Delta64   float64
	Delta     float32
	Now       float32
	Prev      float32
	TickCount uint64
}

func Start() {
	Machine = &MachineStruct{
		Shaders:          make(map[string]*h.ShaderProgram),
		SubTextureDims:   make(map[string]mgl32.Vec4),
		AtlasDescriptors: make(map[string]*h.AtlasDescriptor),
		Textures:         make(map[string]*h.Texture),
		Sprites:          make(map[string]*SpriteGl),
		Cameras:          make(map[string]*Camera),
		AssetPath:        "assets",
	}
}

func KeyPressed(key glfw.Key) bool {
	return glfw.GetCurrentContext().GetKey(key) != glfw.Release
}

func (W *MachineStruct) Tick() {
	W.TickCount += 1
	W.Prev64 = W.Now64
	W.Now64 = glfw.GetTime()
	W.Delta64 = W.Now64 - W.Prev64
	W.Delta = float32(W.Delta64)
	W.Now = float32(W.Now64)
	W.Prev = float32(W.Prev64)
}

func (W *MachineStruct) getPathForAsset(filePath string) string {
	return W.AssetPath + "/" + filePath
}

//	LOAD TEXTURE ATLAS
//
//	an atlas consists of two parts: a descriptor, and an image.
//	the descriptor contains info about the filename of the image
//	as well as info about all the subimages inside the main image.
//
// ///////////////////////////////////////////////////////////////////////////
func (W *MachineStruct) LoadTextureAtlas(filename string) *h.AtlasDescriptor {

	//
	// Success, the descriptor was found in the cache
	if descriptor, found := W.AtlasDescriptors[filename]; found {
		return descriptor
	}

	//
	// Load and unserialize the texture atlas lookup table
	assetPath := W.getPathForAsset(filename)
	descriptor, err := h.LoadTextureAtlas(assetPath)
	if err != nil {
		h.GlPanic(fmt.Errorf("cannot not load texture atlas '%s': %v", filename, err))
	}

	//
	// Store the lookup table for later use
	// for instance: W.TextureAtlasses["sheets/foo.xml"] = lookup
	W.AtlasDescriptors[filename] = descriptor

	//
	// The filename of the actual image.
	// The name of the image is located in the texture atlas lookup table
	// and the image itself must be located in the same directory as the lookup table
	iamgePath := path.Dir(filename) + "/" + descriptor.ImagePath
	// Load the texture
	descriptor.Texture, err = W.GetTexture(iamgePath)
	h.GlPanicIfErrNotNil(err)

	//
	// Populate the SubTextures table with subtexture dimensions
	// for use in the shader
	w, h := descriptor.Texture.GetSize()
	w_f32, h_f32 := float32(w), float32(h)
	for _, sub := range descriptor.SubTextures {
		W.SubTextureDims[filename+"/"+sub.Name] = sub.GetDims(w_f32, h_f32)
	}

	return descriptor
}

func (W *MachineStruct) GetTexture(filename string) (*h.Texture, error) {

	// Texture already loaded. Success.
	if tex, found := W.Textures[filename]; found {
		return tex, nil
	}

	assetPath := W.getPathForAsset(filename)

	tex, err := h.CreateTextureFromFile(assetPath, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE)

	if err != nil {
		return nil, err
	}

	W.Textures[filename] = tex

	return tex, nil
}

func (W *MachineStruct) GetShader(filename string) (*h.ShaderProgram, error) {

	vert := filename + ".vert"
	frag := filename + ".frag"

	// Do we already have this shader in the cache
	if prog, found := W.Shaders[filename]; found {
		return prog, nil
	}

	prog := h.CreateShaderProgramFromFiles(vert, frag)

	W.Shaders[filename] = prog

	return prog, nil
}

func (W *MachineStruct) GetCamera(name string) (cam *Camera, existsAlready bool) {

	if cam, found := W.Cameras[name]; found {
		return cam, found
	}

	cam = &Camera{}

	W.Cameras[name] = cam

	return
}

func (W *MachineStruct) GetDimsForSubtexture(atlasFilename, subTexFilename string) mgl32.Vec4 {
	key := atlasFilename + "/" + subTexFilename
	dims, found := W.SubTextureDims[key]
	if !found {
		h.GlPanic(fmt.Errorf("could not find subtexture '%s", key))
	}

	return dims
}

func (W *MachineStruct) GetAspectRatioForSubTexture(atlasFilename, subTexFilename string) float32 {
	dims := W.GetDimsForSubtexture(atlasFilename, subTexFilename)

	w := dims[2] - dims[0]
	h := dims[3] - dims[1]

	return w / h
}
