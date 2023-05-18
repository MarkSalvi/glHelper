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
	vertexModTime, err := getModifiedTime(vertexPath)
	if err != nil {
		panic(err)
	}
	fragmentModTime, err := getModifiedTime(fragmentPath)
	if err != nil {
		panic(err)
	}
	result := &Shader{id, vertexPath, fragmentPath, vertexModTime, fragmentModTime}

	return result, nil
}

func (shader *Shader) Use() {
	UseProgram(shader.id)
}

func getModifiedTime(filepath string) (time.Time, error) {
	file, err := os.Stat(filepath)
	if err != nil {
		return time.Time{}, err
	}
	return file.ModTime(), nil
}

func (shader *Shader) CheckShaderForChanges() error {

	vertexModTime, err := getModifiedTime(shader.vertPath)
	if err != nil {
		return err
	}
	fragmentModTime, err := getModifiedTime(shader.fragPath)
	if err != nil {
		return err
	}
	if !vertexModTime.Equal(shader.vertexModified) || !fragmentModTime.Equal(shader.fragmentModified) {
		id, err := CreateProgram(shader.vertPath, shader.fragPath)
		if err != nil {
			fmt.Println(err)
		} else {
			gl.DeleteProgram(uint32(shader.id))
			shader.id = id
		}
	}
	return nil
}
