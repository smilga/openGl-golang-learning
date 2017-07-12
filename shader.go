package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Shader struct {
	Program uint32
}

func (s *Shader) Use() {
	gl.UseProgram(s.Program)
}

func NewProgram(vPath, fPath string) (Shader, error) {
	vertexShader, err := compileShader(readShader(vPath), gl.VERTEX_SHADER)
	if err != nil {
		return Shader{0}, err
	}

	fragmentShader, err := compileShader(readShader(fPath), gl.FRAGMENT_SHADER)
	if err != nil {
		return Shader{0}, err
	}

	p := gl.CreateProgram()

	gl.AttachShader(p, vertexShader)
	gl.AttachShader(p, fragmentShader)
	gl.LinkProgram(p)

	var status int32
	gl.GetProgramiv(p, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(p, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(p, logLength, nil, gl.Str(log))

		return Shader{0}, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return Shader{p}, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func readShader(path string) string {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Unable to read shader file: %v", err)
	}

	return string(bs)
}
