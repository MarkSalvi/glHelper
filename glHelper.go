package glHelper

import (
	"errors"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"image/png"
	"math"
	"os"
	"strings"
)

type ShaderID uint32
type ProgramID uint32
type BufferID uint32
type TextureID uint32

type Number interface {
	[]uint32 | []float32
}

func GetVersion() string {
	return gl.GoStr(gl.GetString(gl.VERSION))
}

func LoadShader(path string, shaderType uint32) (ShaderID, error) {
	shaderFile, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	shaderFileStr := string(shaderFile)
	shaderID, err := CreateShader(shaderFileStr, shaderType)
	if err != nil {
		return 0, err
	}

	return shaderID, nil
}

func CreateShader(shaderSource string, shaderType uint32) (ShaderID, error) {

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
		fmt.Println("Failed to compile Shader: \n" + log)
		return 0, errors.New("Failed to compile Shader")
	}

	return ShaderID(shaderId), nil
}

func CreateProgram(vertPath string, fragPath string) (ProgramID, error) {

	vert, err := LoadShader(vertPath, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	frag, err := LoadShader(fragPath, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
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
		return 0, errors.New("FAiled to link Program" + log)

	}
	gl.DeleteShader(uint32(vert))
	gl.DeleteShader(uint32(frag))

	return ProgramID(shaderProgram), nil
}

func GenBindBuffer(target uint32) BufferID {
	var buffer uint32
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(target, buffer)

	return BufferID(buffer)
}

func GenBindVertexArray() BufferID {
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)
	return BufferID(VAO)
}

func BindVertexArray(vao BufferID) {
	gl.BindVertexArray(uint32(vao))
}

// todo add generic instead of float
func BufferData[N Number](target uint32, data N, usage uint32) {
	gl.BufferData(target, len(data)*4, gl.Ptr(data), usage)
}

func UseProgram(prog ProgramID) {
	gl.UseProgram(uint32(prog))
}

// todo SDL2_image or std_image
func LoadTexture(filename string) TextureID {
	infile, err := os.Open(filename)
	defer infile.Close()
	if err != nil {
		panic(err)
	}

	img, err := png.Decode(infile)
	if err != nil {
		panic(err)
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y
	pixels := make([]byte, w*h*4)
	bIndex := 0

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[bIndex] = byte(r / 256)
			bIndex++
			pixels[bIndex] = byte(g / 256)
			bIndex++
			pixels[bIndex] = byte(b / 256)
			bIndex++
			pixels[bIndex] = byte(a / 256)
			bIndex++
		}
	}

	texture := GenBindTexture()
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	return texture

}

func GenBindTexture() TextureID {
	var texID uint32
	gl.GenTextures(1, &texID)
	gl.BindTexture(gl.TEXTURE_2D, texID)
	return TextureID(texID)
}

func BindTexture(id TextureID) {
	gl.BindTexture(gl.TEXTURE_2D, uint32(id))
}

// Rotate3DMat4 returns a 4x4 (non-homogeneous) Matrix that rotates by angle about the axis provided in the vector and an error
//
// vec3{1.0,0.0,0.0} rotates about axis x
//
// vec3{0.0,1.0,0.0} rotates about axis y
//
// vec3{0.0,0.0,1.0} rotates about axis z
func Rotate3DMat4(angle float32, axis mgl32.Vec3) (mgl32.Mat4, error) {
	sin, cos := float32(math.Sin(float64(angle))), float32(math.Cos(float64(angle)))
	switch axis {
	case mgl32.Vec3{1.0, 0.0, 0.0}:
		return mgl32.Mat4{1, 0, 0, 0, 0, cos, sin, 0, 0, -sin, cos, 0, 0, 0, 0, 1}, nil
	case mgl32.Vec3{0.0, 1.0, 0.0}:
		return mgl32.Mat4{cos, 0, -sin, 0, 0, 1, 0, 0, sin, 0, cos, 0, 0, 0, 0, 1}, nil
	case mgl32.Vec3{0.0, 0.0, 1.0}:
		return mgl32.Mat4{cos, sin, 0, 0, -sin, cos, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}, nil
	default:
		return mgl32.Ident4(), errors.New("Matrix rotation error")
	}
}

// MulVec3 performs vector multiplication between two vec3 vectors, returning the product of them
func MulVec3(a, b mgl32.Vec3) mgl32.Vec3 {
	return mgl32.Vec3{a.X() * b.X(), a.Y() * b.Y(), a.Z() * b.Z()}
}
