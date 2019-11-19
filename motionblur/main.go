// Render a square with and without motion blur.
package main

import ( // {{{
	"encoding/binary"
	"fmt"
	"log"
	"math"
	//"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/goxjs/gl/glutil"
	"github.com/goxjs/glfw"
	//"github.com/paulbellamy/ratecounter"
	"golang.org/x/mobile/exp/f32"
) // }}}

type glInfo struct { // {{{
	window          *glfw.Window
	windowSize      [2]int
	cursorPos       [2]float32
	lastCursorPos   [2]float32
	pMatrixUniform  gl.Uniform
	mvMatrixUniform gl.Uniform
	drawTriangle    func(triangle [9]float32, velocity mgl32.Vec3)
} // }}}

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func initWindow(g *glInfo) error { // {{{
	if err := glfw.Init(gl.ContextWatcher); err != nil {
		return err
	}
	defer glfw.Terminate()

	g.windowSize = [2]int{1024, 768}

	// Antialias
	//glfw.WindowHint(glfw.Samples, 8)

	if w, err := glfw.CreateWindow(g.windowSize[0], g.windowSize[1], "", nil, nil); err != nil {
		return err
	} else {
		g.window = w
	}
	g.window.MakeContextCurrent()

	fmt.Printf("OpenGL: %s %s %s; %v samples.\n", gl.GetString(gl.VENDOR), gl.GetString(gl.RENDERER), gl.GetString(gl.VERSION), gl.GetInteger(gl.SAMPLES))
	fmt.Printf("GLSL: %s.\n", gl.GetString(gl.SHADING_LANGUAGE_VERSION))

	// Set callbacks.
	g.cursorPos = [2]float32{float32(g.windowSize[0]) / 2, float32(g.windowSize[1]) / 2}
	g.lastCursorPos = g.cursorPos
	cursorPosCallback := func(_ *glfw.Window, x, y float64) {
		g.cursorPos[0], g.cursorPos[1] = float32(x), float32(y)
	}
	g.window.SetCursorPosCallback(cursorPosCallback)

	framebufferSizeCallback := func(w *glfw.Window, framebufferSize0, framebufferSize1 int) {
		gl.Viewport(0, 0, framebufferSize0, framebufferSize1)

		g.windowSize[0], g.windowSize[1] = w.GetSize()
	}
	g.window.SetFramebufferSizeCallback(framebufferSizeCallback)
	{
		var framebufferSize [2]int
		framebufferSize[0], framebufferSize[1] = g.window.GetFramebufferSize()
		framebufferSizeCallback(g.window, framebufferSize[0], framebufferSize[1])
	}

	// Set OpenGL options.
	gl.ClearColor(0, 0, 0, 1)
	gl.Enable(gl.CULL_FACE)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	gl.Enable(gl.BLEND)
	return nil
} // }}}

func run() error {

	// 初期化
	var g glInfo
	if err := initWindow(&g); err != nil {
		return err
	}

	// 何やってるか理解してない処理をまとめた関数f
	// 理解するごとに処理が減る
	if err := f(&g); err != nil {
		return err
	}

	// N角形のポリゴンを作る
	triangle := makeTriangle(100, 50)

	// 描画ループ
	for !g.window.ShouldClose() {

		gl.Clear(gl.COLOR_BUFFER_BIT)

		// 初期化 単位行列にする
		pMatrix := mgl32.Ortho2D(0, float32(g.windowSize[0]), float32(g.windowSize[1]), 0)
		// Square with motion blur on the left.
		{
			mvMatrix := mgl32.Translate3D(g.cursorPos[0], g.cursorPos[1], 0)

			gl.UniformMatrix4fv(g.pMatrixUniform, pMatrix[:])
			gl.UniformMatrix4fv(g.mvMatrixUniform, mvMatrix[:])

			velocity := mgl32.Vec3{
				g.cursorPos[0] - g.lastCursorPos[0],
				g.cursorPos[1] - g.lastCursorPos[1],
				0,
			}

			for _, t := range triangle {
				g.drawTriangle(t, velocity)
			}
		}
		g.lastCursorPos = g.cursorPos
		g.window.SwapBuffers()
		glfw.PollEvents()
	}

	return nil
}

func makeTriangle(n int, r float64) [][9]float32 {
	t := make([][9]float32, n)
	pre := [2]float32{float32(r), 0}
	for i := 1; i <= n; i++ {
		// 座標を計算
		rate := float64(i) / float64(n)
		x := r * math.Cos(2.0*math.Pi*rate)
		y := r * math.Sin(2.0*math.Pi*rate)
		t[i-1] = [9]float32{
			0, 0, 0,
			float32(x), float32(y), 0,
			pre[0], pre[1], 0,
		}
		pre[0], pre[1] = float32(x), float32(y)
	}
	return t
}

func f(g *glInfo) error { // {{{

	// Init shaders.
	program, err := glutil.CreateProgram(vertexSource, fragmentSource)
	if err != nil {
		return err
	}

	gl.ValidateProgram(program)
	if gl.GetProgrami(program, gl.VALIDATE_STATUS) != gl.TRUE {
		return fmt.Errorf("gl validate status: %s", gl.GetProgramInfoLog(program))
	}

	gl.UseProgram(program)

	g.pMatrixUniform = gl.GetUniformLocation(program, "uPMatrix")
	g.mvMatrixUniform = gl.GetUniformLocation(program, "uMVMatrix")

	tri0v0 := gl.GetUniformLocation(program, "tri0v0")
	tri0v1 := gl.GetUniformLocation(program, "tri0v1")
	tri0v2 := gl.GetUniformLocation(program, "tri0v2")
	tri1v0 := gl.GetUniformLocation(program, "tri1v0")
	tri1v1 := gl.GetUniformLocation(program, "tri1v1")
	tri1v2 := gl.GetUniformLocation(program, "tri1v2")

	vertexPositionAttrib := gl.GetAttribLocation(program, "aVertexPosition")
	gl.EnableVertexAttribArray(vertexPositionAttrib)

	triangleVertexPositionBuffer := gl.CreateBuffer()

	// drawTriangle draws a triangle, consisting of 3 vertices, with motion blur corresponding
	// to the provided velocity. The triangle vertices specify its final position (at t = 1.0,
	// the end of frame), and its velocity is used to compute where the triangle is coming from
	// (at t = 0.0, the start of frame).
	g.drawTriangle = func(triangle [9]float32, velocity mgl32.Vec3) {
		triangle0 := triangle
		for i := 0; i < 3*3; i++ {
			triangle0[i] -= velocity[i%3]
		}
		triangle1 := triangle

		gl.Uniform3f(tri0v0, triangle0[0], triangle0[1], triangle0[2])
		gl.Uniform3f(tri0v1, triangle0[3], triangle0[4], triangle0[5])
		gl.Uniform3f(tri0v2, triangle0[6], triangle0[7], triangle0[8])
		gl.Uniform3f(tri1v0, triangle1[0], triangle1[1], triangle1[2])
		gl.Uniform3f(tri1v1, triangle1[3], triangle1[4], triangle1[5])
		gl.Uniform3f(tri1v2, triangle1[6], triangle1[7], triangle1[8])

		{
			gl.BindBuffer(gl.ARRAY_BUFFER, triangleVertexPositionBuffer)
			vertices := f32.Bytes(binary.LittleEndian,
				triangle0[0], triangle0[1], triangle0[2],
				triangle0[3], triangle0[4], triangle0[5],
				triangle0[6], triangle0[7], triangle0[8],
				triangle1[0], triangle1[1], triangle1[2],
				triangle1[6], triangle1[7], triangle1[8],
				triangle1[3], triangle1[4], triangle1[5],
			)
			gl.BufferData(gl.ARRAY_BUFFER, vertices, gl.DYNAMIC_DRAW)
			itemSize := 3
			itemCount := 6

			gl.VertexAttribPointer(vertexPositionAttrib, itemSize, gl.FLOAT, false, 0, 0)
			gl.DrawArrays(gl.TRIANGLES, 0, itemCount)
		}

		{
			gl.BindBuffer(gl.ARRAY_BUFFER, triangleVertexPositionBuffer)
			vertices := f32.Bytes(binary.LittleEndian,
				triangle0[0], triangle0[1], triangle0[2],
				triangle1[0], triangle1[1], triangle1[2],
				triangle0[3], triangle0[4], triangle0[5],
				triangle1[3], triangle1[4], triangle1[5],
				triangle0[6], triangle0[7], triangle0[8],
				triangle1[6], triangle1[7], triangle1[8],
				triangle0[0], triangle0[1], triangle0[2],
				triangle1[0], triangle1[1], triangle1[2],
			)
			gl.BufferData(gl.ARRAY_BUFFER, vertices, gl.DYNAMIC_DRAW)
			itemSize := 3
			itemCount := 8

			gl.VertexAttribPointer(vertexPositionAttrib, itemSize, gl.FLOAT, false, 0, 0)
			gl.DrawArrays(gl.TRIANGLE_STRIP, 0, itemCount)
		}
	}

	if err := gl.GetError(); err != 0 {
		return fmt.Errorf("gl error: %v", err)
	}
	return nil
} // }}}
