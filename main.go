package main

import (
	_ "image/jpeg"
	_ "image/png"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl" // OR: github.com/go-gl/gl/v2.1/gl
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width  = 800
	height = 600
)

var (
	cameraPos   = mgl32.Vec3{0.0, 0.0, 3.0}
	cameraFront = mgl32.Vec3{0.0, 0.0, -1.0}
	cameraUp    = mgl32.Vec3{0.0, 1.0, 0.0}

	deltaTime   = float32(1)
	lastFrame   = float32(0)
	cameraSpeed = float32(0.5) * deltaTime

	firstMouse = true
	mouseSens  = float64(0.1)
	yaw        = float64(-90)
	pitch      = float64(0)
	lastX      = float64(width / 2)
	lastY      = float64(height / 2)
	fov        = float64(45)
)

var (
	triangle1 = []float32{
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,

		-0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
	}
)

func init() {
	runtime.LockOSThread()
}

func main() {

	window := NewWindow(width, height, "New window")

	window.SetMouseCallback(NewMouseHandler())

	defer window.Terminate()

	if err := gl.Init(); err != nil {
		log.Println(err)
	}

	shader, err := NewProgram("shaders/shader.vert", "shaders/shader.frag")
	if err != nil {
		log.Println(err)
	}

	var vao, vbo uint32

	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	texture, err := LoadTexture("cherry.png")
	if err != nil {
		log.Println(err)
	}

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(triangle1), gl.Ptr(triangle1), gl.STATIC_DRAW)

	//have to reconfigure attribute pointers
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	uformLoc := gl.GetUniformLocation(shader.Program, gl.Str("ourTexture\x00"))

	modelLoc := gl.GetUniformLocation(shader.Program, gl.Str("model\x00"))
	viewLoc := gl.GetUniformLocation(shader.Program, gl.Str("view\x00"))
	projectionLoc := gl.GetUniformLocation(shader.Program, gl.Str("projection\x00"))

	cubePos := []mgl32.Vec3{
		//mgl32.Vec3{0.0, 0.0, 0.0},
		mgl32.Vec3{2.0, 5.0, -15.0},
		mgl32.Vec3{-1.5, -2.2, -2.5},
		mgl32.Vec3{-3.8, -2.0, -12.3},
		mgl32.Vec3{2.4, -0.4, -3.5},
		mgl32.Vec3{-1.7, 3.0, -7.5},
		mgl32.Vec3{1.3, -2.0, -2.5},
		mgl32.Vec3{1.5, 2.0, -2.5},
		mgl32.Vec3{1.5, 0.2, -1.5},
		mgl32.Vec3{-1.3, 1.0, -1.5},
	}

	for !window.ShouldClose() {
		deltaTime = float32(window.GetTime()) - lastFrame
		lastFrame = float32(window.GetTime())

		processInput(window)

		projectionTransf := mgl32.Perspective(mgl32.DegToRad(45), width/height, 0.1, 100)

		dir := cameraPos.Add(cameraFront)

		view := mgl32.LookAt(cameraPos.X(), cameraPos.Y(), cameraPos.Z(),
			dir.X(), dir.Y(), dir.Z(),
			cameraUp.X(), cameraUp.Y(), cameraUp.Z(),
		)

		gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])
		gl.UniformMatrix4fv(projectionLoc, 1, false, &projectionTransf[0])

		//	gl.UniformMatrix4fv(transfomLoc, 1, false, &modelTransf[0])
		gl.Uniform1i(uformLoc, 0)

		gl.ClearColor(0.3, 0.4, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.Enable(gl.DEPTH_TEST)

		shader.Use()

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)

		gl.BindVertexArray(vao)

		for i := 0; i < len(cubePos); i++ {
			model := mgl32.Translate3D(cubePos[i].X(), cubePos[i].Y(), cubePos[i].Z())
			angle := float32(20 * i)

			if i%3 == 0 {
				angle = float32(glfw.GetTime()) * 25
			}

			rotate := mgl32.HomogRotate3D(mgl32.DegToRad(angle), mgl32.Vec3{1.0, 0.3, 0.5})

			modelAll := model.Mul4(rotate)
			gl.UniformMatrix4fv(modelLoc, 1, false, &modelAll[0])
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		window.SwapBuffers()
		window.PollEvents()
	}

}

func processInput(window Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
	if window.GetKey(glfw.KeyW) == glfw.Press {
		cameraPos = cameraPos.Add(cameraFront.Mul(cameraSpeed))
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		cameraPos = cameraPos.Sub(cameraFront.Mul(cameraSpeed))
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		cameraPos = cameraPos.Sub(cameraFront.Cross(cameraUp).Normalize().Mul(cameraSpeed))
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		cameraPos = cameraPos.Add(cameraFront.Cross(cameraUp).Normalize().Mul(cameraSpeed))
	}

}
