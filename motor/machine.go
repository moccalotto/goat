package motor

import (
	"fmt"
	h "goat/glhelp"
	"path"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	Machine *MachineStruct
)

type ThingMap map[Thing]Thing

type Timing struct {
	Now64     float64
	Prev64    float64
	Delta64   float64
	Delta     float32
	Now       float32
	Prev      float32
	TickCount uint64
}

type MachineStruct struct {
	Shaders         map[string]*h.ShaderProgram  // A pointer to all the shader programs currently active in the world
	SubTextures     map[string]*h.SubTexture     // stores subtextures as "sheet.png/image.png"
	Named           map[string]Thing             // Things that aren't drawn and updated, but are kept as prototypes and templates, or just kept in reserve
	TextureAtlasses map[string]*h.TextureAtlas   // stores atlasses as "sheet.png", not "sheet.xml"
	Textures        map[string]*h.Texture        // Pointers to all active textures
	Renderables     map[string]*SpriteRenderable // Renderables are things that can be rendered (or that can render themselves)
	Cameras         map[string]*h.Camera         // Contains the projection matrices. You may want to render ceretain things with one cam, and other things with another cam
	Things          []Thing                      // Things that are dawn and updated every cycle
	Groups          map[string][]Thing           // Named arrays of things. So you can draw (physupdate) groups of things by themselves.
	Timing
}

func Start() {
	Machine = &MachineStruct{
		Shaders:         make(map[string]*h.ShaderProgram),
		SubTextures:     make(map[string]*h.SubTexture),
		Named:           make(map[string]Thing),
		TextureAtlasses: make(map[string]*h.TextureAtlas),
		Textures:        make(map[string]*h.Texture),
		Renderables:     make(map[string]*SpriteRenderable),
		Cameras:         make(map[string]*h.Camera),
	}
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

func (W *MachineStruct) UpdateThings() {
	for i := 0; i < len(W.Things); i++ {
		W.Things[i].Update()
	}
}
func (W *MachineStruct) DrawThings() {
	for i := 0; i < len(W.Things); i++ {
		W.Things[i].Draw()
	}
}

func (W *MachineStruct) LoadTextureAtlas(filePath string) {
	atlas, err := h.LoadTextureAtlas(filePath)

	if err != nil {
		h.GlPanic(err)
	}

	W.TextureAtlasses[atlas.ImagePath] = atlas

	mapPath := path.Dir(filePath) + "/" + atlas.ImagePath

	tex, _ := W.LoadTexture(atlas.ImagePath, mapPath)
	w, h := tex.GetSize()

	for _, sub := range atlas.SubTextures {
		sub.SheetW = uint(w)
		sub.SheetH = uint(h)
		W.SubTextures[atlas.ImagePath+"/"+sub.Name] = sub
	}
}

func (W *MachineStruct) LoadTexture(alias, filePath string) (tex *h.Texture, err error) {

	if _, found := W.Textures[alias]; found {
		err = fmt.Errorf("texture alias '%s' already in use", alias)
		return
	}

	tex, err = h.CreateTextureFromFile(filePath, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE)

	if err != nil {
		return
	}

	W.Textures[alias] = tex
	return
}

func (W *MachineStruct) LoadShader(alias, vert, frag string) (*h.ShaderProgram, error) {

	// Do we already have this shader in the cache
	if prog, found := W.Shaders[alias]; found {
		return prog, nil
	}

	prog := h.CreateShaderProgramFromFiles(vert, frag)

	W.Shaders[alias] = prog

	return prog, nil
}
