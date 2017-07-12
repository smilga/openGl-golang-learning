package main

import (
	"math"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func NewMouseHandler() func(w *glfw.Window, xpos float64, ypos float64) {

	f := func(w *glfw.Window, xpos float64, ypos float64) {

		if firstMouse {
			lastX = xpos
			lastY = ypos
			firstMouse = false
		}

		xoffset := xpos - lastX
		yoffset := lastY - ypos
		lastX = xpos
		lastY = ypos

		xoffset *= mouseSens
		yoffset *= mouseSens

		yaw += xoffset
		pitch += yoffset

		if pitch > float64(89) {
			pitch = float64(89)
		}
		if pitch < float64(-89) {
			pitch = float64(89)
		}

		var front = mgl32.Vec3{}
		front[0] = float32(math.Cos(float64(mgl32.DegToRad(float32(yaw)))) * math.Cos(float64(mgl32.DegToRad(float32(pitch)))))
		front[1] = float32(math.Sin(float64(mgl32.DegToRad(float32(pitch)))))
		front[2] = float32(math.Sin(float64(mgl32.DegToRad(float32(yaw)))) * math.Cos(float64(mgl32.DegToRad(float32(pitch)))))

		cameraFront = front.Normalize()
	}
	return f
}
