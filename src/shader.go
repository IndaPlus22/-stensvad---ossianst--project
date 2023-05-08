package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Shader struct {
	id uint32
}

func NewShader(filePath string) Shader {
	filePath = "../res/shaders/" + filePath
	vertexSource, fragmentSource, err := parseShader(filePath)
	if err != nil {
		panic(err)
	}

	id, err := newProgram(vertexSource, fragmentSource)
	if err != nil {
		panic(err)
	}

	return Shader{id}
}

func parseShader(filePath string) (vertexShader, fragmentShader string, err error) {
	// Read content of file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", "", err
	}

	// Two builders for storing vertex and fragment shaders
	var sb [2]strings.Builder
	var currentShader *strings.Builder

	// Iterate over each line in content
	for _, line := range strings.Split(string(content), "\n") {
		// Check if line is first line of a shader
		if strings.HasPrefix(line, "#shader") {
			// Determine shader type based on line content
			if strings.Contains(line, "vertex") {
				currentShader = &sb[0] // Set current shader builder to vertex
			} else if strings.Contains(line, "fragment") {
				currentShader = &sb[1] // Set current shader builder to fragment
			}
		} else if currentShader != nil {
			// If we are inside a shader block, append line to current shader builder
			currentShader.WriteString(line + "\n")
		}
	}

	// Return vertex and fragment shaders with null terminators
	return sb[0].String() + "\x00", sb[1].String() + "\x00", nil
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	// Compile vertex shader
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	// Compile fragment shader
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	// Attach shaders to program
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)

	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)

	// Check for errors
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	// Delete shaders as they are no longer needed
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	// Convert source string to C-style string
	csources, free := gl.Strs(source)

	// Set shader source code
	gl.ShaderSource(shader, 1, csources, nil)

	// Free C-style string resources
	free()

	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)

	// Check for errors
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func (s *Shader) setUniform1i(name string, value int32) {
	location := gl.GetUniformLocation(s.id, gl.Str(name+"\x00"))
	gl.Uniform1i(location, value)
}

func (s *Shader) setUniform1f(name string, value float32) {
	location := gl.GetUniformLocation(s.id, gl.Str(name+"\x00"))
	gl.Uniform1f(location, value)
}

func (s *Shader) setUniform2f(name string, v0, v1 float32) {
	location := gl.GetUniformLocation(s.id, gl.Str(name+"\x00"))
	gl.Uniform2f(location, v0, v1)
}

func (s *Shader) setUniform3f(name string, v0, v1, v2 float32) {
	location := gl.GetUniformLocation(s.id, gl.Str(name+"\x00"))
	gl.Uniform3f(location, v0, v1, v2)
}

func (s *Shader) setUniform4f(name string, v0, v1, v2, v3 float32) {
	location := gl.GetUniformLocation(s.id, gl.Str(name+"\x00"))
	gl.Uniform4f(location, v0, v1, v2, v3)
}

func (s *Shader) setUniformMat4fv(name string, matrix mgl32.Mat4) {
	location := gl.GetUniformLocation(s.id, gl.Str(name+"\x00"))
	gl.UniformMatrix4fv(location, 1, false, &matrix[0])
}

func (s *Shader) bind() {
	gl.UseProgram(s.id)
}

func (s *Shader) unbind() {
	gl.UseProgram(0)
}
