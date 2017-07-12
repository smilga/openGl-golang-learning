package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"log"
)

type Window struct {
	Width  int
	Height int
	*glfw.Window
}

func (w *Window) Terminate() {
	glfw.Terminate()
}

func (w *Window) GetTime() float64 {
	return glfw.GetTime()
}

func (w *Window) PollEvents() {
	glfw.PollEvents()
}

func (w *Window) SetMouseCallback(c glfw.CursorPosCallback) {
	w.SetCursorPosCallback(c)
}

func NewWindow(width int, height int, title string) Window {
	if err := glfw.Init(); err != nil {
		log.Fatalln("Failed to init glfw", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	w, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}

	w.MakeContextCurrent()

	w.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	window := Window{
		width,
		height,
		w,
	}

	return window
}
