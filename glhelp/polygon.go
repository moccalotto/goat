package glhelp

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Polygon struct {
	shader      *ShaderProgram
	texture     *Texture
	indeces     []uint32  // array of indices that point into the verts array
	verts       []float32 // array of floats that constitute our verts
	texCoords   []float32 // texture coordinates
	colors      []float32 // color at each vert position
	ready       bool      // have the buffers been initialized
	vao         uint32    // Vector Attribute Object
	vbo         uint32    // Vector Buffer Object
	indexBuffer uint32    // Element/Index Buffer Object
	cbo         uint32    // Color Buffer Object
	tbo         uint32    // Texture Buffer Object
}

func CreatePolygon(sides int, colors []float32, texFilePath string) *Polygon {

	/// https://faun.pub/draw-circle-in-opengl-c-2da8d9c2c103

	tau := math.Pi * 2

	anglePerSide := tau / float64(sides)

	triangleCount := sides - 2

	var verts []float32
	var texCoords []float32
	var indeces []uint32
	var actualColors []float32

	switch len(colors) {
	case sides * 4:
		// one color per vertex. OK
	case 0:
		// zero colors: OK
		for i := 0; i < sides; i++ {
			actualColors = append(actualColors, float32(1), float32(1), float32(1), float32(1))
		}
	case 4:
		// one single color: OK
		actualColors = append(actualColors, colors[0], colors[1], colors[2], colors[3])
	default:
		log.Printf("To create a %d-sided polygon, you must supply 0, 1, or %d colors = 0, 4, or %d floats\n",
			sides, sides, sides*4)
		panic("Invalid number of colors.")

	}

	offset := rand.Float64() * tau

	for i := 0; i < sides; i++ {
		angle := anglePerSide*float64(i) + offset

		y, x := math.Sincos(angle)
		fx, fy := float32(x), float32(y)

		texCoords = append(texCoords, 0.5-0.5*fx, 0.5-0.5*fy)
		verts = append(verts, fx, fy, 0.0)

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

	return &Polygon{
		verts:     verts,
		indeces:   indeces,
		texCoords: texCoords,
		colors:    actualColors,
		texture:   texture,
	}
}

func (Pol *Polygon) Draw(window *glfw.Window) {
	Pol.Initialize()
	Pol.shader.Use()
	gl.BindVertexArray(Pol.vao)
	gl.DrawElements(gl.TRIANGLES, int32(len(Pol.indeces)), gl.UNSIGNED_INT, nil)
	AssertGLOK("Draw", "D")
}

func (Pol *Polygon) Initialize() {
	if Pol.ready {
		return
	}
	Pol.initShader()
	Pol.shader.Use()
	Pol.initBuffers()
	Pol.texture.Initialize()

	loc, _ := Pol.shader.getUniformLocation("uniTexture")
	Pol.texture.Setuniform(loc)

	Pol.ready = true

	AssertGLOK("Polygon.Initialize")
	log.Println("Shaders and Buffers are Ready")
}

func (Pol *Polygon) initShader() {
	if Pol.ready {
		return
	}
	if Pol.shader != nil {
		log.Printf("Shader was initialized before the buffers. %v\n", Pol)
		log.Printf("Possibly trying to init a previously destroyed polygon\n")
		panic("Trying to initialize a dirty object. Shader was already initialized")
	}

	Pol.shader = CreateProgramFromFiles("polygon.vert", "polygon.frag")

	// tell GL to use the shader
	log.Println("Shaders compiled and loaded")
}

// Make a Vertex Array Object and return its ID
func (Pol *Polygon) initBuffers() {
	if Pol.ready {
		return
	}
	if len(Pol.verts) < 3 {
		panic(fmt.Errorf("expected at least 3 elements in the verts array, only got %d", len(Pol.verts)))
	}

	// Vertex Buffer Object
	gl.GenBuffers(1, &Pol.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, Pol.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(Pol.verts), gl.Ptr(Pol.verts), gl.STATIC_DRAW)

	// Vertex Array Object
	gl.GenVertexArrays(1, &Pol.vao)
	gl.BindVertexArray(Pol.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, Pol.vao)

	Pol.shader.EnableVertexAttribArray("iVert")
	Pol.shader.VertexAttribPointer("iVert", 3, gl.FLOAT, false, 0, nil)

	// Element/Index Buffer Object
	gl.GenBuffers(1, &Pol.indexBuffer)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, Pol.indexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(Pol.indeces)*4, gl.Ptr(Pol.indeces), gl.STATIC_DRAW)

	//
	// Color Buffer Object
	gl.GenBuffers(1, &Pol.cbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, Pol.cbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(Pol.colors)*4, gl.Ptr(Pol.colors), gl.STATIC_DRAW)

	Pol.shader.EnableVertexAttribArray("iColor")
	Pol.shader.VertexAttribPointer("iColor", 4, gl.FLOAT, false, 0, nil)

	//
	// Texture Buffers

	gl.GenBuffers(1, &Pol.tbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, Pol.tbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(Pol.texCoords)*4, gl.Ptr(Pol.texCoords), gl.STATIC_DRAW)

	Pol.shader.EnableVertexAttribArray("iTexCoord")
	Pol.shader.VertexAttribPointer("iTexCoord", 2, gl.FLOAT, false, 0, nil)

	gl.BindVertexArray(0)
	Pol.shader.DisableVertexAttribArray("iVert")
	Pol.shader.DisableVertexAttribArray("iTexCoord")
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

func (Pol *Polygon) Destroy() {

	Pol.ready = false

	// TODO: keep all buffers in an array so we can delete them all with one call.
	gl.DeleteBuffers(1, &Pol.indexBuffer)
	gl.DeleteBuffers(1, &Pol.vao)
	gl.DeleteBuffers(1, &Pol.vbo)
	gl.DeleteBuffers(1, &Pol.cbo)
	gl.DeleteBuffers(1, &Pol.tbo)
	AssertGLOK("Polygon.Destroy")

	Pol.shader.Destroy()
	Pol.shader = nil

	Pol.texture.Unbind()
	Pol.texture = nil

	AssertGLOK("Polygon.Destroy")
}
