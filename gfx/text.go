package gfx

import (
	"fmt"
	"os"
	"skyfx/math"
	v4 "skyfx/math/v4"
	"skyfx/utils"
	"unicode/utf8"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/skycoin/gltext"
)

// Constants ...
const (
	TEXT_LEFT_TO_RIGHT = gltext.LeftToRight
)

var fonts map[string]*gltext.Font = make(map[string]*gltext.Font, 0)

func loadTrueType(name string, file *os.File, scale int32, low int32, high int32, dir int32, fixedPipeline bool) {
	if theFont, err := gltext.LoadTruetype(file,
		scale, rune(low), rune(high),
		gltext.Direction(dir), fixedPipeline); err == nil {
		fonts[name] = theFont
	}
}

// func opGltextLoadTrueType(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
// 	loadTrueType(prgrm, inputs, outputs, true)
// }

func opGltextLoadTrueTypeCore() {
}

// func opGltextPrintf(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
// 	if err := fonts[inputs[0].Get_str(prgrm)].Printf(inputs[1].Get_f32(prgrm), inputs[2].Get_f32(prgrm), inputs[3].Get_str(prgrm)); err != nil {
// 		panic(err)
// 	}
// }

// func opGltextMetrics(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
// 	width, height := fonts[inputs[0].Get_str(prgrm)].Metrics(inputs[1].Get_str(prgrm))

// 	outputs[0].Set_i32(prgrm, int32(width))
// 	outputs[1].Set_i32(prgrm, int32(height))
// }

// func opGltextTexture(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
// 	outputs[0].Set_i32(prgrm, int32(fonts[inputs[0].Get_str(prgrm)].Texture()))
// }

// func opGltextNextGlyph(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) { // refactor
// 	font := fonts[inputs[0].Get_str(prgrm)]
// 	str := inputs[1].Get_str(prgrm)
// 	var index int = int(inputs[2].Get_i32(prgrm))
// 	var runeValue rune = -1
// 	var width int = -1
// 	var x int = 0
// 	var y int = 0
// 	var w int = 0
// 	var h int = 0
// 	var advance int = 0
// 	if index < len(str) {
// 		runeValue, width = utf8.DecodeRuneInString(str[index:])
// 		g := font.Glyphs()[runeValue-font.Low()]
// 		x = g.X
// 		y = g.Y
// 		w = g.Width
// 		h = g.Height
// 		advance = g.Advance
// 	}

// 	outputs[0].Set_i32(prgrm, int32(runeValue-font.Low()))
// 	outputs[1].Set_i32(prgrm, int32(width))
// 	outputs[2].Set_i32(prgrm, int32(x))
// 	outputs[3].Set_i32(prgrm, int32(y))
// 	outputs[4].Set_i32(prgrm, int32(w))
// 	outputs[5].Set_i32(prgrm, int32(h))
// 	outputs[6].Set_i32(prgrm, int32(advance))
// }

// func opGltextGlyphBounds(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
// 	font := fonts[inputs[0].Get_str(prgrm)]
// 	var maxGlyphWidth, maxGlyphHeight int = font.GlyphBounds()
// 	outputs[0].Set_i32(prgrm, int32(maxGlyphWidth))
// 	outputs[1].Set_i32(prgrm, int32(maxGlyphHeight))
// }

// func opGltextGlyphMetrics(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) { // refactor
// 	width, height := fonts[inputs[0].Get_str(prgrm)].GlyphMetrics(uint32(inputs[1].Get_i32(prgrm)))

// 	outputs[0].Set_i32(prgrm, int32(width))
// 	outputs[1].Set_i32(prgrm, int32(height))
// }

// func opGltextGlyphInfo(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) { // refactor
// 	font := fonts[inputs[0].Get_str(prgrm)]
// 	glyph := inputs[1].Get_i32(prgrm)
// 	var x int = 0
// 	var y int = 0
// 	var w int = 0
// 	var h int = 0
// 	var advance int = 0
// 	g := font.Glyphs()[glyph]
// 	x = g.X
// 	y = g.Y
// 	w = g.Width
// 	h = g.Height
// 	advance = g.Advance

// 	outputs[0].Set_i32(prgrm, int32(x))
// 	outputs[1].Set_i32(prgrm, int32(y))
// 	outputs[2].Set_i32(prgrm, int32(w))
// 	outputs[3].Set_i32(prgrm, int32(h))
// 	outputs[4].Set_i32(prgrm, int32(advance))
// }

// MeasureGlyph ...
func MeasureGlyph(name string, glyph int32) (w int32, h int32) {
	width, height := fonts[name].GlyphMetrics(uint32(glyph))
	w = int32(width)
	h = int32(height)
	return
}

// MeasureText ...
func MeasureText(name string, text string) (w int32, h int32) {
	width, height := fonts[name].Metrics(text)
	w = int32(width)
	h = int32(height)
	return
}

// LoadTrueType ...
func LoadTrueType(name string, path string, scale int32, min int32, max int32, dir int32) (texture uint32) {
	file, err := utils.CXOpenFile(path)
	defer file.Close()
	if err != nil {
		panic(fmt.Sprintf("texture %q not found on disk: %v\n", path, err))
	}

	loadTrueType(name, file, scale, min, max, dir, false)
	texture = fonts[name].Texture()
	return
}

// GlyphBounds ...
func GlyphBounds(name string) (width int32, height int32) {
	var maxGlyphWidth, maxGlyphHeight int = fonts[name].GlyphBounds()

	width = int32(maxGlyphWidth)
	height = int32(maxGlyphHeight)
	return
}

// GlyphInfo ...
func GlyphInfo(name string, glyph int32) (x int32, y int32, w int32, h int32, a int32) {
	font := fonts[name]
	g := font.Glyphs()[glyph]
	x = int32(g.X)
	y = int32(g.Y)
	w = int32(g.Width)
	h = int32(g.Height)
	a = int32(g.Advance)
	return
}

// NextGlyph ...
func NextGlyph(name string, text string, index int32) (r int32, s int32, x int32, y int32, w int32, h int32, a int32) {
	font := fonts[name]
	s = -1
	var runeValue rune = -1
	var width int = -1
	if int(index) < len(text) {
		runeValue, width = utf8.DecodeRuneInString(text[index:])
		g := font.Glyphs()[runeValue-font.Low()]
		x = int32(g.X)
		y = int32(g.Y)
		w = int32(g.Width)
		h = int32(g.Height)
		a = int32(g.Advance)
	}
	r = runeValue - font.Low()
	s = int32(width)

	return
}

// TextureCreateFont ...
func TextureCreateFont(name string, path string, scale int32, min int32, max int32, dir int32, mipmap int32) (out TextureId) {
	var glName uint32 = LoadTrueType(name, path, scale, min, max, dir)
	out = textureCreate(glName, gl.TEXTURE_2D, FORMAT_R8_G8_B8_A8, 0, 0, mipmap, false)
	return
}

// MeshAppendText ...
func MeshAppendText(id MeshId, texture TextureId, name string, position math.V2, scale math.V2, color math.V4, text string, debug bool, color0 math.V4, color1 math.V4, clip math.V4, depth float32) {
	//	fmt.Printf("NAME %s\n", name)
	//	if (name == "awesomeBold_25") {
	utils.PanicIfNot(IsValidMesh(id), "invalid id")
	var size int32 = 0

	// var maxGlyphW int32 = 0
	var maxGlyphH int32 = 0
	_ /*maxGlyphW*/, maxGlyphH = GlyphBounds(name)
	// var maxGlyphWidth float32 = scale.X * float32(maxGlyphW)
	var maxGlyphHeight float32 = scale.Y * float32(maxGlyphH)

	// var textW int32
	var textH int32
	_ /*textW*/, textH = MeasureText(name, text)

	// var textWidth float32 = scale.X * float32(textW)
	var textHeight float32 = scale.Y * float32(textH)

	var quadX float32 = position.X
	var quadY float32 = position.Y - (maxGlyphHeight - textHeight) // - ratio * float32(textH) / gfx_height
	//if (center == true) {
	//	quadX = -1.0 + (2.0 - textWidth / gfx_width) / 2.0 // / gfx_width
	//	quadY =  0.0 - textHeight / (gfx_height * 2.0) // / gfx_height
	//}

	var tw float32 = TextureWidthF32(texture)
	var th float32 = TextureHeightF32(texture)
	var dummyIndex int32 = 0
	// var debugColor math.V4

	//	fmt.Printf("TEXT : %s\n", text)

	var index int32 = 0
	for size >= 0 {

		var rune int32
		var glyphX int32
		var glyphY int32
		// var glyphWidth int32
		var glyphHeight int32
		var glyphAdvance int32

		rune, size, glyphX, glyphY, _ /*glyphWidth*/, glyphHeight, glyphAdvance = NextGlyph(name, text, index)
		if size >= 0 {

			//			fmt.Printf("RUNE\n")
			var quadAdvance float32 = scale.X * float32(glyphAdvance) // / gfx_width
			//var quadWidth float32 = scale.X * float32(glyphWidth)     // / gfx_width
			var quadHeight float32 = textHeight //scale.Y * float32(glyphHeight)// / gfx_height
			var qu0 float32 = float32(glyphX) / tw
			var qv0 float32 = float32(glyphY) / th

			//		  fmt.Printf("U %f, V %f, GH %f, MGH %f, TH %f\n", qu0, qv0, float32(glyphHeight), maxGlyphHeight, textHeight)
			var qu1 float32 = qu0 + float32(glyphAdvance)/tw
			var qv1 float32 = qv0 + float32(glyphHeight)/th

			utils.PanicIfNot(rune >= 0, "invalid rune")
			if debug == true {
				color = color0
				// debugColor = color1
				if (dummyIndex % 2) == 0 {
					color = color1
					// debugColor = color0
				}
			}

			//MeshAppendOrthoQuad(g_debugMesh, v4.Make(quadX, quadY, quadAdvance, quadHeight), v4.Make(qu0, qv0, qu1, qv1), red, clip, depth)
			MeshAppendOrthoQuad(id, v4.Make(quadX, quadY, quadAdvance, quadHeight), v4.Make(qu0, qv0, qu1, qv1), color, clip, depth)
			index = index + size
			quadX = quadX + quadAdvance
			dummyIndex = dummyIndex + 1
		} // else {
		//  fmt.Printf("no rune")
		//}
	}
	//  }
}

// MeshAppendGlyph ...
func MeshAppendGlyph(id MeshId, texture TextureId, name string, position math.V2, scale math.V2, color math.V4, glyph int32, debug bool, color0 math.V4, color1 math.V4, clip math.V4, depth float32) {
	//	fmt.Printf("NAME %s\n", name)
	utils.PanicIfNot(IsValidMesh(id), "invalid id")
	// var size int32 = 0

	var maxGlyphW int32 = 0
	var maxGlyphH int32 = 0
	maxGlyphW, maxGlyphH = GlyphBounds(name)
	var maxGlyphWidth float32 = float32(maxGlyphW)
	var maxGlyphHeight float32 = float32(maxGlyphH)
	maxGlyphWidth = scale.X * maxGlyphWidth
	maxGlyphHeight = scale.Y * maxGlyphHeight

	// var textW int32
	var textH int32
	_ /*textW*/, textH = MeasureGlyph(name, glyph)

	// var textWidth float32 = scale.X * float32(textW)
	var textHeight float32 = scale.Y * float32(textH)

	var quadX float32 = position.X
	var quadY float32 = position.Y - (maxGlyphHeight - textHeight) // - ratio * float32(textH) / gfx_height
	//if (center == true) {
	//	quadX = -1.0 + (2.0 - textWidth / gfx_width) / 2.0 // / gfx_width
	//	quadY =  0.0 - textHeight / (gfx_height * 2.0) // / gfx_height
	//}

	var tw float32 = TextureWidthF32(texture)
	var th float32 = TextureHeightF32(texture)
	var dummyIndex int32 = 0
	// var debugColor math.V4

	// var rune int32
	var glyphX int32
	var glyphY int32
	// var glyphWidth int32
	var glyphHeight int32
	var glyphAdvance int32

	glyphX, glyphY, _ /*glyphWidth*/, glyphHeight, glyphAdvance = GlyphInfo(name, glyph)
	//	fmt.Printf("RUNE\n")
	var quadAdvance float32 = scale.X * float32(glyphAdvance) // / gfx_width
	// var quadWidth float32 = scale.X * float32(glyphWidth)     // / gfx_width
	var quadHeight float32 = textHeight //scale.Y * float32(glyphHeight)// / gfx_height
	var qu0 float32 = float32(glyphX) / tw
	var qv0 float32 = float32(glyphY) / th

	//   fmt.Printf("U %f, V %f, GH %f, MGH %f, TH %f\n", qu0, qv0, float32(glyphHeight), maxGlyphHeight, textHeight)
	var qu1 float32 = qu0 + float32(glyphAdvance)/tw
	var qv1 float32 = qv0 + float32(glyphHeight)/th

	if debug == true {
		color = color0
		// debugColor = color1
		if (dummyIndex % 2) == 0 {
			color = color1
			// debugColor = color0
		}
	}

	//MeshAppendOrthoQuad(g_debugMesh, v4.Make(quadX, quadY, quadAdvance, quadHeight), v4.Make(qu0, qv0, qu1, qv1), red, clip, depth)
	MeshAppendOrthoQuad(id, v4.Make(quadX, quadY, quadAdvance, quadHeight), v4.Make(qu0, qv0, qu1, qv1), color, clip, depth)
	quadX = quadX + quadAdvance
	dummyIndex = dummyIndex + 1
}
