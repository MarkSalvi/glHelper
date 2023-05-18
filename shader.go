package glHelper

import (
	"errors"
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
		if err == errors.New("no such file or directory") {
			fmt.Println(time.Now(), " ", filepath, " : ", err)
		} else {
			panic(err)
		}
	}
	return file.ModTime()
}

func (shader *Shader) CheckShaderForChanges() {

	vertexModTime := getModifiedTime(shader.vertPath)
	fragmentModTime := getModifiedTime(shader.fragPath)
	if !vertexModTime.Equal(shader.vertexModified) || !fragmentModTime.Equal(shader.fragmentModified) {
		id, err := CreateProgram(shader.vertPath, shader.fragPath)
		if err != nil {
			fmt.Println(err)
		} else {
			gl.DeleteProgram(uint32(shader.id))
			shader.id = id
		}
	}

}
