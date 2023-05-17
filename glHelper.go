package glHelper

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"strings"
)

type ShaderID uint32
type ProgramID uint32
type VAOID uint32
type VBOID uint32

func GetVersion() string {
	return gl.GoStr(gl.GetString(gl.VERSION))
}

func CreateShader(shaderSource string, shaderType uint32) ShaderID {

	shaderId := gl.CreateShader(shaderType)
	shaderSource += "\x00"
	csource, free := gl.Strs(shaderSource)
	gl.ShaderSource(shaderId, 1, csource, nil)
	free()
	gl.CompileShader(shaderId)
	var status int32
	gl.GetShaderiv(shaderId, gl.COMPILE_STATUS, &status)

	if status == gl.FALSE {
		var logLenght int32
		gl.GetShaderiv(shaderId, gl.INFO_LOG_LENGTH, &logLenght)
		log := strings.Repeat("\x00", int(logLenght+1))
		gl.GetShaderInfoLog(shaderId, logLenght, nil, gl.Str(log))
		panic("Failed to compile Shader: \n" + log)
	}
	return ShaderID(shaderId)
}

func CreateProgram(vert ShaderID, frag ShaderID) ProgramID {
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, uint32(vert))
	gl.AttachShader(shaderProgram, uint32(frag))
	gl.LinkProgram(shaderProgram)

	var success int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLenght int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLenght)
		log := strings.Repeat("\x00", int(logLenght+1))
		gl.GetProgramInfoLog(shaderProgram, logLenght, nil, gl.Str(log))
		panic("Failed to link Shader Program: \n" + log)
	}
	gl.DeleteShader(uint32(vert))
	gl.DeleteShader(uint32(frag))

	return ProgramID(shaderProgram)
}

func GenBindBuffer(target uint32) VBOID {
	var VBO uint32
	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(target, VBO)

	return VBOID(VBO)
}

func GenBindVertexArray() VAOID {
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)
	return VAOID(VAO)
}

func BindVertexArray(vao VAOID) {
	gl.BindVertexArray(uint32(vao))
}

// todo add generic instead of float
func BufferData(target uint32, data []float32, usage uint32) {
	gl.BufferData(target, len(data)*4, gl.Ptr(data), usage)
}

func UseProgram(prog ProgramID) {
	gl.UseProgram(uint32(prog))
}
