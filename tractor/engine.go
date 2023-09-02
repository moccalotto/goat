package tractor

import (
	"fmt"
	shed "goat/shed"
	"path"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	Engine *EngineType
)

// =============================================================================================
// ||
// || Engine.
// ||
// || Windowing
// || Rendering
// || Input
// || Behaviors
// ||
// =============================================================================================
// The engine.
// It drives most things.
type EngineType struct {
	shaders          map[string]*shed.ShaderProgram   // A pointer to all the shader programs currently active in the world
	subTextureDims   map[string]shed.V4               // stores subtextures as "sheet.png/image.png" => minX, minY, maxX, maxY
	atlasDescriptors map[string]*shed.AtlasDescriptor // stores atlasses as "sheet.png", not "sheet.xml"
	textures         map[string]*shed.TextureWrapper  // Pointers to all active textures
	cameras          map[string]*Camera               // Contains the projection matrices. You may want to render ceretain things with one cam, and other things with another cam
	AssetPath        string                           // Base path for all assets
	MainCamera       *Camera
	Controls         *ControlsType

	// Timing
	Now64     float64
	Prev64    float64
	Delta64   float64
	Delta     float32
	Now       float32
	Prev      float32
	TickCount uint64
	Window    *glfw.Window
	Dispose   func()
}

// Start the goat Motor and assign it to the global variable Motor
func StartMain(o *WindowOptions) {
	Engine = StartCustom(o)
}

// Spin up a goat Motor and return it.
func StartCustom(o *WindowOptions) *EngineType {
	M := &EngineType{
		shaders:          make(map[string]*shed.ShaderProgram),
		subTextureDims:   make(map[string]shed.V4),
		atlasDescriptors: make(map[string]*shed.AtlasDescriptor),
		textures:         make(map[string]*shed.TextureWrapper),
		cameras:          make(map[string]*Camera),
		AssetPath:        "assets",
		Window:           nil,
	}

	M.Controls = &ControlsType{E: M}

	var err error

	M.Dispose, M.Window, err = glfwCreateWin(o)
	shed.GlPanicIfErrNotNil(err)

	M.GetCamera("main")

	return M
}

// ============================================
// ||
// || Update all timers
// ||
// ============================================
func (W *EngineType) Tick() {
	W.TickCount += 1
	W.Prev64 = W.Now64
	W.Now64 = glfw.GetTime()
	W.Delta64 = W.Now64 - W.Prev64
	W.Delta = float32(W.Delta64)
	W.Now = float32(W.Now64)
	W.Prev = float32(W.Prev64)
}

// ============================================
// || LOOP:
// ||
// || Clear screen
// || Call fn(),
// || Update screen
// ============================================
func (W *EngineType) Loop(fn func()) {

	for !W.Window.ShouldClose() {
		shed.Clear()

		W.Tick()

		fn()

		W.Window.SwapBuffers()

		shed.AssertGLOK("End Of Loop")

		glfw.PollEvents()
	}
}

// Append AssetPath to a file path
func (W *EngineType) getPathForAsset(filePath string) string {
	return W.AssetPath + "/" + filePath
}

//	LOAD TEXTURE ATLAS
//
//	an atlas consists of two parts: a descriptor, and an image.
//	the descriptor contains info about the filename of the image
//	as well as info about all the subimages inside the main image.
//
// ///////////////////////////////////////////////////////////////////////////
func (W *EngineType) LoadTextureAtlas(filename string) *shed.AtlasDescriptor {

	//
	// Success, the descriptor was found in the cache
	if descriptor, found := W.atlasDescriptors[filename]; found {
		return descriptor
	}

	//
	// Load and unserialize the texture atlas lookup table
	assetPath := W.getPathForAsset(filename)
	descriptor, err := shed.LoadTextureAtlasFile(assetPath)
	if err != nil {
		shed.GlPanic(fmt.Errorf("cannot not load texture atlas '%s': %v", filename, err))
	}

	//
	// Store the lookup table for later use
	// for instance: W.TextureAtlasses["sheets/foo.xml"] = lookup
	W.atlasDescriptors[filename] = descriptor

	//
	// The filename of the actual image.
	// The name of the image is located in the texture atlas lookup table
	// and the image itself must be located in the same directory as the lookup table
	iamgePath := path.Dir(filename) + "/" + descriptor.ImagePath
	// Load the texture
	descriptor.Texture, err = W.GetTexture(iamgePath)
	shed.GlPanicIfErrNotNil(err)

	//
	// Populate the SubTextures table with subtexture dimensions
	// for use in the shader
	w, h := descriptor.Texture.GetSize()
	w_f32, h_f32 := float32(w), float32(h)
	for _, sub := range descriptor.SubTextures {
		W.subTextureDims[filename+"/"+sub.Name] = sub.GetDims(w_f32, h_f32)
	}

	return descriptor
}

// Load a texture from a file, or retrieve it from the cache if it had previously been loaded
func (W *EngineType) GetTexture(filename string) (*shed.TextureWrapper, error) {

	// Texture already loaded. Success.
	if tex, found := W.textures[filename]; found {
		return tex, nil
	}

	assetPath := W.getPathForAsset(filename)

	tex, err := shed.CreateTextureFromFile(assetPath, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE)

	if err != nil {
		return nil, err
	}

	W.textures[filename] = tex

	return tex, nil
}

// Load a shader program from the BASENAME of a file.
// The vert shader must have the .vert extension
// The frag shader must have the .frag extension
func (W *EngineType) GetShader(filename string) (*shed.ShaderProgram, error) {

	vert := filename + ".vert"
	frag := filename + ".frag"

	// Do we already have this shader in the cache
	if prog, found := W.shaders[filename]; found {
		return prog, nil
	}

	prog := shed.CreateShaderProgramFromFiles(vert, frag)

	W.shaders[filename] = prog

	return prog, nil
}

// Acquire a camera by the given name.
func (W *EngineType) GetCamera(name string) (cam *Camera, existsAlready bool) {

	if cam, found := W.cameras[name]; found {
		return cam, found
	}

	cam = &Camera{}

	if W.MainCamera == nil {
		W.MainCamera = cam
	}

	if name == "main" {
		W.MainCamera = cam
	}

	W.cameras[name] = cam

	return
}

// Get the location and size of a given subtexture
func (W *EngineType) GetDimsForSubtexture(atlasFilename, subTexFilename string) shed.V4 {
	key := atlasFilename + "/" + subTexFilename
	dims, found := W.subTextureDims[key]
	if !found {
		shed.GlPanic(fmt.Errorf("could not find subtexture '%s", key))
	}

	return dims
}

func (W *EngineType) GetAspectRatioForSubTexture(atlasFilename, subTexFilename string) float32 {
	dims := W.GetDimsForSubtexture(atlasFilename, subTexFilename)

	w := dims.C3 - dims.C1
	h := dims.C4 - dims.C2

	return w / h
}

func (W *EngineType) GracefulShutdown() {
	W.Window.SetShouldClose(true)
}
