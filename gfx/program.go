package gfx

import (
	"fmt"
	"skyfx/utils"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// import "gl"

// Globals ...
var g_programs []uint32

// Program ...
type Program struct {
	name uint32
}

// NullProgram ...
func NullProgram() (out Program) {
	out.name = 0
	return
}

func compileShader(shaderType uint32, source string) (out uint32) {
	out = gl.CreateShader(shaderType)
	utils.PanicIf(GlError(), "gl.CreateShader")

	glSrcs, freeFn := gl.Strs(source + "\x00")
	defer freeFn()

	gl.ShaderSource(out, 1, glSrcs, nil)
	utils.PanicIf(GlError(), "gl.ShaderSource")

	gl.CompileShader(out)
	utils.PanicIf(GlError(), "gl.CompileShader")

	var status int32
	gl.GetShaderiv(out, gl.COMPILE_STATUS, &status)
	utils.PanicIf(GlError(), "gl.GetShaderiv")
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(out, gl.INFO_LOG_LENGTH, &logLength)
		utils.PanicIf(GlError(), "gl.GetShaderiv")

		var glLog []uint8 = make([]uint8, logLength+1)
		var glLogLength int32
		gl.GetShaderInfoLog(out, logLength+1, &glLogLength, &glLog[0])
		utils.PanicIf(GlError(), "gl.GetShaderInfoLog")

		var log string = string(glLog)
		fmt.Printf("Failed to compile: %v\n", log)
	}
	return
}

// ProgramCreate ...
func ProgramCreate(attributes []Location, samplers []Location, uniforms []Location, vertex string, pixel string) (out Program) {
	var vertexShader uint32 = compileShader(gl.VERTEX_SHADER, vertex)
	var pixelShader uint32 = compileShader(gl.FRAGMENT_SHADER, pixel)

	// program
	out.name = gl.CreateProgram()
	utils.PanicIf(GlError(), "gl.CreateProgram")
	g_programs = append(g_programs, out.name)

	var attributeCount int = len(attributes)
	for i := 0; i < attributeCount; i++ {
		if attributes[i].IsValid() {
			// TODO REMOVE
			// glAttribute, freeFn := gl.Strs(attributes[i].name + "\x00")
			// defer freeFn()

			var location uint32 = attributes[i].location
			ProgramBindAttribute(out, location, attributes[i].name)
		}
	}

	gl.AttachShader(out.name, vertexShader)
	utils.PanicIf(GlError(), "gl.AttachShader")

	gl.AttachShader(out.name, pixelShader)
	utils.PanicIf(GlError(), "gl.AttachShader")

	gl.LinkProgram(out.name)
	utils.PanicIf(GlError(), "gl.LinkProgram")

	var status int32
	gl.GetProgramiv(out.name, gl.LINK_STATUS, &status)
	utils.PanicIf(GlError(), "gl.GetProgramiv")
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(out.name, gl.INFO_LOG_LENGTH, &logLength)
		utils.PanicIf(GlError(), "gl.GetProgramiv")

		var glLog []uint8 = make([]uint8, logLength+1)
		var glLogLength int32
		gl.GetProgramInfoLog(out.name, logLength+1, &glLogLength, &glLog[0])
		utils.PanicIf(GlError(), "gl.GetProgramInfoLog")

		var log string = string(glLog)
		fmt.Printf("Failed to link: %v\n", log)

		fmt.Printf("//----------------------------------------------\n")
		fmt.Printf("VERTEX_SHADER:\n%s\n", vertex)
		fmt.Printf("//----------------------------------------------\n")
		fmt.Printf("PIXEL_SHADER:\n%s\n", pixel)
	}

	gl.DetachShader(out.name, vertexShader)
	utils.PanicIf(GlError(), "gl.LinkProgram")

	gl.DetachShader(out.name, pixelShader)
	utils.PanicIf(GlError(), "gl.DetachShader")

	gl.DeleteShader(vertexShader)
	utils.PanicIf(GlError(), "gl.DeleteShader")

	gl.DeleteShader(pixelShader)
	utils.PanicIf(GlError(), "gl.DeleteShader")

	ProgramUse(out)

	var samplerCount int = len(samplers)
	for i := 0; i < samplerCount; i++ {
		var name string = samplers[i].name
		if name != "" {
			samplers[i].location = ProgramTryBindSampler(out, name)
		}
	}

	var uniformCount int = len(uniforms)
	for i := 0; i < uniformCount; i++ {
		var name string = uniforms[i].name
		if name != "" {
			uniforms[i].location = ProgramTryBindUniform(out, name)
		}
	}
	return
}

// ProgramUse ...
func ProgramUse(program Program) {
	gl.UseProgram(program.name)
	utils.PanicIf(GlError(), "gl.UseProgram")
}

// ProgramBindAttribute ...
func ProgramBindAttribute(program Program, location uint32, name string) {
	glName, freeFn := gl.Strs(name + "\x00")
	defer freeFn()
	gl.BindAttribLocation(program.name, location, *glName)
	utils.PanicIf(GlError(), "gl.BindAttributeLocation")
}

// ProgramBindUniform ...
func ProgramBindUniform(program Program, name string) (uniform uint32) {
	uniform = ProgramTryBindUniform(program, name)
	utils.PanicIfNot(uniform >= 0, "uniform>= 0")
	return
}

// ProgramTryBindUniform ...
func ProgramTryBindUniform(program Program, name string) (uniform uint32) {
	glName, freeFn := gl.Strs(name + "\x00")
	defer freeFn()
	uniform = uint32(gl.GetUniformLocation(program.name, *glName)) // ##2 crash find a way to test if string is present in gl.Strs
	utils.PanicIf(GlError(), "gl.GetUniformLocation")
	return
}

// ProgramBindSampler ...
func ProgramBindSampler(program Program, name string) (sampler uint32) {
	sampler = ProgramBindUniform(program, name)
	return
}

// ProgramTryBindSampler ...
func ProgramTryBindSampler(program Program, name string) (sampler uint32) {
	sampler = ProgramTryBindUniform(program, name)
	return
}
