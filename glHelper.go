package glHelper

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"strings"
)

func GetVersion() string {
	return gl.GoStr(gl.GetString(gl.VERSION))
}

func MakeShader(shaderSource string, shaderType uint32) uint32 {

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
	return shaderId
}
