package glHelper

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

func GetVersion() string {
	return gl.GoStr(gl.GetString(gl.VERSION))
}
