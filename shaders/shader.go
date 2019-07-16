package shaders

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

//ShaderProgram is a generic shader program
type ShaderProgram struct {
	ID               uint32
	vertexShaderID   uint32
	fragmentShaderID uint32
}

//StaticShader is a shader program that contains a static shader
type StaticShader struct {
	ShaderProgram *ShaderProgram
}

//NewStaticShader creates a static shader
func NewStaticShader(vertexFile string, fragmentFile string) (*StaticShader, error) {
	ss := new(StaticShader)
	var err error
	ss.ShaderProgram, err = NewShaderProgram(vertexFile, fragmentFile)
	if err != nil {
		return nil, err
	}

	ss.ShaderProgram.BindAttribute(0, "position")
	ss.ShaderProgram.BindAttribute(1, "textureCoords")
	return ss, nil
}

//NewShaderProgram creates a new shader program given the vertex file and fragment file
func NewShaderProgram(vertexFile string, fragmentFile string) (*ShaderProgram, error) {
	s := new(ShaderProgram)
	var err error

	s.vertexShaderID, err = loadShaderFromFile(vertexFile, gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}
	s.fragmentShaderID, err = loadShaderFromFile(fragmentFile, gl.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}

	s.ID = gl.CreateProgram()
	gl.AttachShader(s.ID, s.vertexShaderID)
	gl.AttachShader(s.ID, s.fragmentShaderID)
	gl.LinkProgram(s.ID)
	gl.ValidateProgram(s.ID)

	return s, nil
}

//Start starts the shader program
func (sp *ShaderProgram) Start() {
	gl.UseProgram(sp.ID)
}

//Stop stops the shader program
func (sp *ShaderProgram) Stop() {
	gl.UseProgram(0)
}

//CleanUp stops the program, detaches shaders and finally deletes them as well as the program
func (sp *ShaderProgram) CleanUp() {
	sp.Stop()
	gl.DetachShader(sp.ID, sp.vertexShaderID)
	gl.DetachShader(sp.ID, sp.fragmentShaderID)
	gl.DeleteShader(sp.vertexShaderID)
	gl.DeleteShader(sp.fragmentShaderID)
	gl.DeleteProgram(sp.ID)
}

//BindAttribute binds an attribute to the shader program
func (sp *ShaderProgram) BindAttribute(attrib uint32, name string) {
	gl.BindAttribLocation(sp.ID, attrib, gl.Str(name+"\x00"))
}

//LoadShader loads a shader file from system and compiles it
func loadShaderFromFile(fileName string, shaderType uint32) (uint32, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	source, err := ioutil.ReadAll(file)
	if err != nil {
		return 0, err
	}

	id, err := compileShader(string(source)+"\x00", shaderType)
	return id, err
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
