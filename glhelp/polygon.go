package glhelp

import (
	"errors"
	"fmt"
	"log"
	"math"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Polygon struct {
	shader      *ShaderProgram // the shader program used to render this polygon - all polygons with the same number of sides can have the same shader program
	texture     *Texture       // texture image (if any)
	initialized bool           // have we sent all data to openGL
	indeces     []uint32       // array of indices that point into the verts array
	verts       []float32      // array of floats that constitute our verts
	texCoords   []float32      // texture coordinates
	colors      []float32      // color at each vert position
	vao         uint32         // Vector Attribute Object
	buffers     [4]uint32      // our 4 buffers

	colorMix  float32    // how to mix the colors of the vertices and the texture 0 = pure texture, 1 = pure color
	uniColor  mgl32.Vec4 // how to mix the colors of the vertices and the texture 0 = pure texture, 1 = pure color
	wPosX     float32    // The X position of the object - in world coordinates
	wPosY     float32    // The Y position of the object - in world coordinates
	wScaleX   float32    // The x-scaling of the object - in world coordinates
	wScaleY   float32    // The y-scaling of the object - in world coordinates
	wRotation float32    // the rotation, relative to the world's X-axis, of the object - in radians
	wireframe bool

	isCopy bool // is this a clone
}

type PolygonOptions struct {
	VertColors      []float32
	UniversalColor  *mgl32.Vec4
	TextureFilePath string
	TextureObj      *Texture
	ShaderObj       *ShaderProgram
	FragShaderFile  string
	VertShaderFile  string
}

func (O *PolygonOptions) Verify() (err []error) {
	textureConflict := O.TextureFilePath != "" && O.TextureObj != nil
	if textureConflict {
		err = append(err, errors.New("you cannot declare a texture object and a texture file path at the same time"))
	}

	fragShaderConflict := O.ShaderObj != nil && O.FragShaderFile != ""
	if fragShaderConflict {
		err = append(err, errors.New("you cannot declare frag shader file path as well as a shader object at the same time"))
	}

	vertShaderConflict := O.ShaderObj != nil && O.FragShaderFile != ""
	if vertShaderConflict {
		err = append(err, errors.New("you cannot declare a vert shader file path as well as a shader object at the same time"))
	}

	if (O.FragShaderFile == "" && O.VertShaderFile != "") || (O.FragShaderFile != "" && O.VertShaderFile == "") {
		err = append(err, errors.New("you must either declare both a vertex shader and a frag shader file path - or none of them"))
	}

	if O.UniversalColor != nil && len(O.VertColors) > 0 {
		err = append(err, errors.New("you must declare UniversalColor OR VertColors - not both"))
	}

	return err
}

// Create a 1x1 square that is axis-aligned and centered on the origin
func CreateSprite(texFilePath string) *Polygon {
	return CreatePolygonAdv(
		4,
		45*Degrees,
		math.Sqrt2,
		math.Sqrt2,
		[]float32{},
		texFilePath,
	)
}

func CreatePolygonAdv(sides int, angleOffset, width, height float64, colors []float32, texFilePath string) *Polygon {

	/// https://faun.pub/draw-circle-in-opengl-c-2da8d9c2c103

	anglePerSide := Tau / float64(sides)

	triangleCount := sides - 2

	var verts []float32
	var texCoords []float32
	var indeces []uint32
	var vertexColors []float32

	switch len(colors) {
	case sides * 4:
		// one color per vertex. OK
	case 0:
		// zero colors: OK
		for i := 0; i < sides; i++ {
			vertexColors = append(vertexColors, float32(1), float32(1), float32(1), float32(1))
		}
	case 4:
		// one single color: OK
		for i := 0; i < sides; i++ {
			vertexColors = append(vertexColors, colors[0], colors[1], colors[2], colors[3])
		}
	default:
		GlPanic(
			fmt.Errorf("to create a %d-sided polygon, you must supply 0, 1, or %d colors = 0, 4, or %d floats", sides, sides, sides*4),
		)
	}

	for i := 0; i < sides; i++ {
		angle := anglePerSide*float64(i) + angleOffset

		y, x := math.Sincos(angle)
		fx, fy := float32(x*width), float32(y*height)

		texCoords = append(texCoords, 0.5-0.5*fx, 0.5-0.5*fy)
		verts = append(verts, fx, fy, 1.0)

	}

	for i := 0; i < triangleCount; i++ {
		indeces = append(indeces,
			0,
			uint32(i+1),
			uint32(i+2),
		)
	}

	texture, err := CreateTextureFromFile(texFilePath, gl.CLAMP_TO_EDGE, gl.CLAMP_TO_EDGE)
	if err != nil {
		panic(err)
	}

	log.Printf("%+v\n", verts)

	return &Polygon{
		verts:     verts,
		indeces:   indeces,
		texCoords: texCoords,
		colors:    vertexColors,
		texture:   texture,
		colorMix:  0.5,
		uniColor: [4]float32{
			vertexColors[0],
			vertexColors[1],
			vertexColors[2],
			vertexColors[3],
		},
		wireframe: false,
	}
}

func (Pol *Polygon) GetPos() (float32, float32) {
	return Pol.wPosX, Pol.wPosY
}

func (Pol *Polygon) GetPosV() mgl32.Vec2 {

	return mgl32.Vec2{Pol.wPosX, Pol.wPosY}
}

func (Pol *Polygon) GetScale() (float32, float32) {
	return Pol.wScaleX, Pol.wScaleY
}

func (Pol *Polygon) GetScaleV() mgl32.Vec2 {

	return mgl32.Vec2{Pol.wScaleX, Pol.wScaleY}
}

func (Pol *Polygon) GetRotation() float32 {
	return Pol.wRotation
}

func (Pol *Polygon) GetColorMix() float32 {
	return Pol.colorMix
}

func (Pol *Polygon) GetTransformationMatrix() mgl32.Mat3 {
	translate := mgl32.Translate2D(Pol.wPosX, Pol.wPosY)
	rotate := mgl32.HomogRotate2D(Pol.wRotation)
	scale := mgl32.Scale2D(Pol.wScaleX, Pol.wScaleY)

	return MatMulMany(translate, rotate, scale)
}

func (Pol *Polygon) Draw(camMatrix mgl32.Mat3) {
	Pol.Initialize()
	Pol.shader.Use()
	gl.BindVertexArray(Pol.vao)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, Pol.texture.handle)

	trMatrix := MatMulMany(
		camMatrix,
		Pol.GetTransformationMatrix(),
	)

	Pol.shader.SetUniformAttr("uniWireframe", Pol.wireframe)
	err := Pol.shader.SetUniformAttr("uniColor", Pol.uniColor)
	GlPanicIfErrNotNil(err)
	Pol.shader.SetUniformAttr("uniColorMix", Pol.colorMix)
	Pol.shader.SetUniformAttr("uTransformation", trMatrix)

	gl.DrawElements(gl.TRIANGLES, int32(len(Pol.indeces)), gl.UNSIGNED_INT, nil)
	AssertGLOK("Draw", "D")
}

func (Pol *Polygon) Initialize() {
	if Pol.initialized {
		return
	}
	Pol.initShader()
	Pol.initBuffers()
	Pol.initTexture()

	AssertGLOK("Polygon.Initialize")
	log.Println("Shaders and Buffers are Ready")

	Pol.initialized = true
}

func (Pol *Polygon) initTexture() {
	Pol.texture.Initialize()
	// Pol.shader.SetUniformAttr("uniTexture", 0) // TODO: what ???
	Pol.shader.SetUniformAttr("uniColorMix", Pol.colorMix)
	Pol.shader.SetUniformAttr("uniColor", Pol.uniColor)
}

func (Pol *Polygon) initShader() {
	if Pol.initialized {
		return
	}
	if Pol.shader != nil {
		log.Printf("Shader was initialized before the buffers. %v\n", Pol)
		log.Printf("Possibly trying to init a previously destroyed polygon\n")
		panic("Trying to initialize a dirty object. Shader was already initialized")
	}

	Pol.shader = CreateShaderProgramFromFiles("shaders/polygon.vert", "shaders/polygon.frag")

	// tell GL to use the shader
	log.Println("Shaders compiled and loaded")
}

// Make a Vertex Array Object and return its ID
func (Pol *Polygon) initBuffers() {
	if Pol.initialized {
		return
	}
	if len(Pol.verts) < 3 {
		panic(fmt.Errorf("expected at least 3 elements in the verts array, only got %d", len(Pol.verts)))
	}

	Pol.shader.Use()

	// Create out buffers
	// VBO: 0	vertex
	// EBO: 1	element
	// CBO: 2	color
	// TBO: 3	texture
	gl.GenBuffers(int32(len(Pol.buffers)), &Pol.buffers[0])

	//
	// Vertex Array Object
	gl.GenVertexArrays(1, &Pol.vao)
	gl.BindVertexArray(Pol.vao)

	//
	// Vertex Buffer Object
	gl.BindBuffer(gl.ARRAY_BUFFER, Pol.buffers[0])
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(Pol.verts), gl.Ptr(Pol.verts), gl.STATIC_DRAW)

	Pol.shader.EnableVertexAttribArray("iVert")
	Pol.shader.VertexAttribPointer("iVert", 3, gl.FLOAT, false, 0, nil)
	defer Pol.shader.DisableVertexAttribArray("iVert")

	//
	// Element/Index Buffer Object
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, Pol.buffers[1])
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(Pol.indeces), gl.Ptr(Pol.indeces), gl.STATIC_DRAW)

	//
	// Color Buffer Object
	if Pol.shader.HasAttrib("iColor") {
		gl.BindBuffer(gl.ARRAY_BUFFER, Pol.buffers[2])
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(Pol.colors), gl.Ptr(Pol.colors), gl.STATIC_DRAW)

		Pol.shader.EnableVertexAttribArray("iColor")
		Pol.shader.VertexAttribPointer("iColor", 4, gl.FLOAT, false, 0, nil)
		defer Pol.shader.DisableVertexAttribArray("iColor")
	}

	//
	// Texture Buffers
	if Pol.shader.HasAttrib("iTexCoord") {
		gl.BindBuffer(gl.ARRAY_BUFFER, Pol.buffers[3])
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(Pol.texCoords), gl.Ptr(Pol.texCoords), gl.STATIC_DRAW)
		Pol.shader.EnableVertexAttribArray("iTexCoord")
		Pol.shader.VertexAttribPointer("iTexCoord", 2, gl.FLOAT, false, 0, nil)
		defer Pol.shader.DisableVertexAttribArray("iTexCoord")
	}

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	Pol.texture.Unbind()
}

func (Pol *Polygon) Destroy() {

	if Pol.isCopy {
		return
	}

	Pol.initialized = false

	gl.DeleteBuffers(int32(len(Pol.buffers)), &Pol.buffers[0])
	AssertGLOK("Polygon.Destroy")

	Pol.shader.Destroy()
	Pol.shader = nil

	Pol.texture.Destroy()
	Pol.texture = nil

	AssertGLOK("Polygon.Destroy")
}

func (Pol *Polygon) SetWireframe(wf bool) {
	Pol.wireframe = wf
}

func (Pol *Polygon) GetWireframe() bool {
	return Pol.wireframe
}

func (Pol *Polygon) SetPosition(x, y float32) {
	Pol.wPosX = x
	Pol.wPosY = y
}

func (Pol *Polygon) SetColor(uniColor mgl32.Vec4) {
	Pol.uniColor = uniColor
}

func (Pol *Polygon) SetColorMix(uniColorMix float32) {
	Pol.colorMix = uniColorMix
}

func (Pol *Polygon) Move(dx, dy float32) {
	Pol.wPosX += dx
	Pol.wPosY += dy
}

func (Pol *Polygon) SetRotation(radians float32) {
	Pol.wRotation = radians
}

func (Pol *Polygon) Rotate(radians float32) {
	Pol.wRotation += radians
}

func (Pol *Polygon) SetScale(x, y float32) {
	Pol.wScaleX = x / 2
	Pol.wScaleY = y / 2
	// divide by 2 because the size is the "radius"
	// I want an object with scale = 4 to take up
	// the entire frame of a camera with the with = 4
}

// Scale(1.33) => increase size by 33%
func (Pol *Polygon) Scale(x, y float32) {
	Pol.wScaleX *= x
	Pol.wScaleY *= y
}

func (Pol *Polygon) Copy() *Polygon {
	return &Polygon{
		shader:      Pol.shader,
		texture:     Pol.texture,
		initialized: Pol.initialized,
		indeces:     Pol.indeces,
		verts:       Pol.verts,
		texCoords:   Pol.texCoords,
		colors:      Pol.colors,
		vao:         Pol.vao,
		buffers:     Pol.buffers,
		colorMix:    Pol.colorMix,
		uniColor:    Pol.uniColor,
		wPosX:       Pol.wPosX,
		wPosY:       Pol.wPosY,
		wScaleX:     Pol.wScaleX,
		wScaleY:     Pol.wScaleY,
		wRotation:   Pol.wRotation,
		isCopy:      true,
	}
}
