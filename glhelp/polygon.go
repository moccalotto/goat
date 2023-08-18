package glhelp

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type Polygon struct {
	shader    *ShaderProgram
	texture   *Texture
	indeces   []uint32  // array of indices that point into the verts array
	verts     []float32 // array of floats that constitute our verts
	texCoords []float32 // texture coordinates
	colors    []float32 // color at each vert position
	ready     bool      // have the buffers been initialized
	vao       uint32    // Vector Attribute Object
	buffers   [4]uint32 // our 4 buffers
}

type PolygonChanel chan *Polygon

func GoCreatePolygon(result PolygonChanel, sides int, colors []float32, texFilePath string) {
	go func() {
		result <- CreatePolygon(sides, colors, texFilePath)
	}()
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

func (Pol *Polygon) Draw() {
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

	Pol.shader = CreateProgramFromFiles("shaders/polygon.vert", "shaders/polygon.frag")

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

	// Create out buffers
	// VBO: 0	vertex
	// EBO: 1	element
	// CBO: 2	color
	// TBO: 3	texture
	gl.GenBuffers(int32(len(Pol.buffers)), &Pol.buffers[0])

	// Vertex Array Object
	gl.GenVertexArrays(1, &Pol.vao)
	gl.BindVertexArray(Pol.vao)

	// Vertex Buffer Object
	gl.BindBuffer(gl.ARRAY_BUFFER, Pol.buffers[0])
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(Pol.verts), gl.Ptr(Pol.verts), gl.STATIC_DRAW)

	Pol.shader.EnableVertexAttribArray("iVert")
	Pol.shader.VertexAttribPointer("iVert", 3, gl.FLOAT, false, 0, nil)

	// Element/Index Buffer Object
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, Pol.buffers[1])
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(Pol.indeces)*4, gl.Ptr(Pol.indeces), gl.STATIC_DRAW)

	//
	// Color Buffer Object
	gl.BindBuffer(gl.ARRAY_BUFFER, Pol.buffers[2])
	gl.BufferData(gl.ARRAY_BUFFER, len(Pol.colors)*4, gl.Ptr(Pol.colors), gl.STATIC_DRAW)

	Pol.shader.EnableVertexAttribArray("iColor")
	Pol.shader.VertexAttribPointer("iColor", 4, gl.FLOAT, false, 0, nil)

	//
	// Texture Buffers

	gl.BindBuffer(gl.ARRAY_BUFFER, Pol.buffers[3])
	gl.BufferData(gl.ARRAY_BUFFER, len(Pol.texCoords)*4, gl.Ptr(Pol.texCoords), gl.STATIC_DRAW)

	Pol.shader.EnableVertexAttribArray("iTexCoord")
	Pol.shader.VertexAttribPointer("iTexCoord", 2, gl.FLOAT, false, 0, nil)

	// cleanup

	Pol.shader.DisableVertexAttribArray("iVert")
	Pol.shader.DisableVertexAttribArray("iTexCoord")

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

func (Pol *Polygon) Destroy() {

	Pol.ready = false

	gl.DeleteBuffers(int32(len(Pol.buffers)), &Pol.buffers[0])
	AssertGLOK("Polygon.Destroy")

	Pol.shader.Destroy()
	Pol.shader = nil

	Pol.texture.Unbind()
	Pol.texture = nil

	AssertGLOK("Polygon.Destroy")
}
