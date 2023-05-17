package glHelper

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"os"
	"time"
)

type Shader struct {
	id               ProgramID
	vertPath         string
	fragPath         string
	vertexModified   time.Time
	fragmentModified time.Time
}

func NewShader(vertexPath, fragmentPath string) (*Shader, error) {
	id, err := CreateProgram(vertexPath, fragmentPath)
	if err != nil {
		return nil, err
	}
	result := &Shader{id, vertexPath, fragmentPath, getModifiedTime(vertexPath), getModifiedTime(fragmentPath)}

	return result, nil
}

func (shader *Shader) Use() {
	UseProgram(shader.id)
}

func getModifiedTime(filepath string) time.Time {
	file, err := os.Stat(filepath)
	if err != nil {
		panic(err)
	}
	return file.ModTime()
}

func (shaderInfo *Shader) CheckShaderForChanges() {

	vertexModTime := getModifiedTime(shaderInfo.vertPath)
	fragmentModTime := getModifiedTime(shaderInfo.fragPath)
	if !vertexModTime.Equal(shaderInfo.vertexModified) || fragmentModTime.Equal(shaderInfo.fragmentModified) {
		id, err := CreateProgram(shaderInfo.vertPath, shaderInfo.fragPath)
		if err != nil {
			fmt.Println(err)
		} else {
			gl.DeleteProgram(uint32(shaderInfo.id))
			shaderInfo.id = id
		}
	}

}
