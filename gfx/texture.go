package gfx

import (
	// "fmt"
	"bufio"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"skyfx/math"
	"skyfx/types"
	"skyfx/utils"
	"strings"
	"unsafe"

	gomath "math"

	"github.com/go-gl/gl/v3.3-core/gl"
)

var gifs map[string]*gif.GIF = make(map[string]*gif.GIF, 0)
var textures map[string]texture = make(map[string]texture, 0)

func decodeImg(file *os.File, cpuCopy bool) (data []byte, width uint32, height uint32, pixels []float32) {
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	data = rgba.Pix
	width = uint32(rgba.Rect.Size().X)
	height = uint32(rgba.Rect.Size().Y)
	if cpuCopy {
		pixels = make([]float32, width*height*4)
		var x uint32
		var y uint32
		for y = 0; y < height; y++ {
			yoffset := y * width * 4
			for x = 0; x < width; x++ {
				var xoffset = yoffset + x*4
				color := rgba.At(int(x), int(y))
				r, g, b, a := color.RGBA()
				pixels[xoffset] = float32(r) / 65535.0
				pixels[xoffset+1] = float32(g) / 65535.0
				pixels[xoffset+2] = float32(b) / 65535.0
				pixels[xoffset+3] = float32(a) / 65535.0
			}
		}
	}
	return
}

const (
	HDR_NONE = iota
	HDR_32_RLE_RGBE
	MINLEN = 8
	MAXLEN = 0x7fff
	R      = 0
	G      = 1
	B      = 2
	E      = 3
)

func unpack(file *os.File, width types.Pointer, line []byte) bool {
	if width < MINLEN || width > MAXLEN {
		return unpack_(file, width, line)
	}

	file.Read(line[:4])
	if line[R] != 2 {
		file.Seek(-4, io.SeekCurrent)
		return unpack_(file, width, line)
	}

	if line[G] != 2 || (line[B]&128) != 0 {
		return unpack_(file, width-1, line[4:])
	}

	var b [1]byte
	for i := types.Pointer(0); i < 4; i++ {
		for j := types.Pointer(0); j < width; {
			file.Read(b[:])
			count := types.Pointer(b[0])
			if count > 128 {
				count &= 127
				file.Read(b[:])
				var value int = int(b[0])
				for c := types.Pointer(0); c < count; c++ {
					line[j+c+i] = byte(value)
				}
			} else {
				for c := types.Pointer(0); c < count; c++ {
					offset := j + c + i
					file.Read(line[offset : offset+1])
				}
			}
		}
	}
	return true
}

func unpack_(file *os.File, width types.Pointer, line []byte) bool {
	var rshift uint
	var repeat [4]byte
	for width > 0 {
		file.Read(line[0:4])
		if line[R] == 1 && line[G] == 1 && line[B] == 1 {
			for i := line[E] << rshift; i > 0; i-- {
				copy(line[0:4], repeat[:])
				line = line[4:]
				width--
			}
			rshift += 8
		} else {
			copy(repeat[:], line[0:4])
			line = line[4:]
			width--
			rshift = 0
		}
	}
	return true
}

func decodeHdr(file *os.File) (data []byte, i32Width uint32, i32Height uint32) {
	data = nil
	i32Width = 0
	i32Height = 0

	var format int
	scanner := bufio.NewScanner(file)

	var pos int64
	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		pos += int64(advance)
		return
	}

	scanner.Split(scanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "#?RADIANCE" {
		} else if strings.HasPrefix(line, "#") {
		} else if strings.HasPrefix(line, "FORMAT=") {
			var sformat string
			if n, err := fmt.Sscanf(line, "FORMAT=%s\n", &sformat); n != 1 && err != nil {
				fmt.Printf("Failed to scan format : err '%s'\n", err)
				return
			}
			if sformat == "32-bit_rle_rgbe" {
				format = HDR_32_RLE_RGBE
			}
		} else if len(line) == 0 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Failed to scan : err %v\n", scanner.Err())
	}

	if format != HDR_32_RLE_RGBE {
		fmt.Printf("Invalid format %d\n", format)
		return
	}

	file.Seek(pos, 0)

	var iwidth int
	var iheight int
	if n, err := fmt.Fscanf(file, "-Y %d +X %d\n", &iwidth, &iheight); n != 2 || err != nil {
		fmt.Printf("Failed to scan width and height : err '%s'\n", err)
		return
	}

	i32Width = uint32(iwidth)
	i32Height = uint32(iheight)

	width := types.Cast_int_to_ptr(iwidth)
	height := types.Cast_int_to_ptr(iheight)

	//var colors []float32 = make([]float32, width*height*3)
	var line []byte = make([]byte, width*4)
	data = make([]byte, width*height*3*4)

	for y := types.Pointer(0); y < height; y++ {
		if unpack(file, width, line) == false {
			fmt.Printf("Failed to unpack line %d\n", y)
			return
		}

		yoffset := y /*(height - y - 1)*/ * width * 3 * 4
		for x := types.Pointer(0); x < width; x++ {
			loffset := x * 4
			exponent := gomath.Pow(2.0, float64(int(line[loffset+3])-128))
			xoffset := yoffset + x*3*4
			r := float32(exponent * float64(line[loffset]) / 256.0)
			g := float32(exponent * float64(line[loffset+1]) / 256.0)
			b := float32(exponent * float64(line[loffset+2]) / 256.0)

			types.Write_f32(data, xoffset, r)
			types.Write_f32(data, xoffset+4, g)
			types.Write_f32(data, xoffset+8, b)
		}
	}
	return
}

func uploadTexture(path string, target uint32, level uint32, cpuCopy bool) {
	file, err := utils.CXOpenFile(path)
	defer file.Close()
	if err != nil {
		panic(fmt.Sprintf("texture %q not found on disk: %v\n", path, err))
	}

	ext := filepath.Ext(path)
	var data []byte
	var internalFormat int32
	var inputFormat uint32
	var inputType uint32
	var width uint32
	var height uint32
	var pixels []float32
	if ext == ".png" || ext == ".jpeg" || ext == ".jpg" {
		internalFormat = gl.RGBA8
		inputFormat = gl.RGBA
		inputType = gl.UNSIGNED_BYTE
		data, width, height, pixels = decodeImg(file, cpuCopy)
		if cpuCopy {
		}
		if len(pixels) > 0 {
			var texture texture
			texture.pixels = pixels
			texture.width = width
			texture.height = height
			texture.path = path
			texture.level = level
			textures[path] = texture
		}
	} else if ext == ".hdr" {
		internalFormat = gl.RGB16F
		inputFormat = gl.RGB
		inputType = gl.FLOAT
		data, width, height = decodeHdr(file)
	}

	if len(data) > 0 {
		gl.TexImage2D(
			target,
			int32(level),
			internalFormat,
			int32(width),
			int32(height),
			0,
			inputFormat,
			inputType,
			unsafe.Pointer(&data[0]))
	}
}

func opGlNewTexture(path string) uint32 {
	var texture uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	uploadTexture(path, gl.TEXTURE_2D, 0, false)

	return texture
}

func opGlNewTextureCube(pattern string, extension string) uint32 {
	var texture uint32
	gl.Enable(gl.TEXTURE_CUBE_MAP)
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, texture)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	var faces []string = []string{"posx", "negx", "posy", "negy", "posz", "negz"}
	for i := 0; i < 6; i++ {
		uploadTexture(fmt.Sprintf("%s%s%s", pattern, faces[i], extension), uint32(gl.TEXTURE_CUBE_MAP_POSITIVE_X+i), 0, false)
	}
	return texture
}

func opCxReleaseTexture(path string) {
	textures[path] = texture{}
}

func opCxTextureGetPixel(path string, x uint32, y uint32) (r float32, g float32, b float32, a float32) {

	if texture, ok := textures[path]; ok {
		var yoffset = y * texture.width * 4
		var xoffset = yoffset + x*4
		pixels := texture.pixels
		r = pixels[xoffset]
		g = pixels[xoffset+1]
		b = pixels[xoffset+2]
		a = pixels[xoffset+3]
	}
	return
}

func opGlUploadImageToTexture(path string, target uint32, level uint32, cpuCopy bool) {
	uploadTexture(path, target, level, cpuCopy)
}

func opGlNewGIF(path string) (count int32, loop int32, width int32, height int32) {

	file, err := utils.CXOpenFile(path)
	defer file.Close()
	if err != nil {
		panic(fmt.Sprintf("file not found %q, %v", path, err))
	}

	reader := bufio.NewReader(file)
	gif, err := gif.DecodeAll(reader)
	if err != nil {
		panic(fmt.Sprintf("failed to decode file %q, %v", path, err))
	}

	gifs[path] = gif

	count = int32(len(gif.Image))
	loop = int32(gif.LoopCount)
	width = int32(gif.Config.Width)
	height = int32(gif.Config.Height)
	return
}

func opGlFreeGIF(path string) {
	gifs[path] = nil
}

func opGlGIFFrameToTexture(path string, frame int32, texture int32) (delay int32, disposal int32) {
	gif := gifs[path]
	img := gif.Image[frame]
	delay = int32(gif.Delay[frame])
	disposal = int32(gif.Disposal[frame])

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	gl.BindTexture(gl.TEXTURE_2D, uint32(texture))
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		unsafe.Pointer(&rgba.Pix[0]))

	return
}

// // TODO : texture array
// // TODO : texture 3d
// // TODO : texture compression
// // TODO : convert to pot
// // TODO : generate mipmap
// // TODO : compression

// // Constants ...
const (
	FORMAT_NONE             int32 = 0
	FORMAT_DEPTH24          int32 = 1
	FORMAT_STENCIL8         int32 = 2
	FORMAT_DEPTH24_STENCIL8 int32 = 3
	FORMAT_R8_G8_B8_A8      int32 = 4
	FORMAT_RGB_16F          int32 = 5
	FORMAT_RGBA_16F         int32 = 6
	FORMAT_R8               int32 = 7
)

// Globals ...
var g_textures []Texture
var g_nil []uint32

// TextureId ...
type TextureId struct {
	texture int32
}

type texture struct {
	path   string
	width  uint32
	height uint32
	level  uint32
	pixels []float32
}

// Texture ...
type Texture struct {
	id             TextureId
	path           string
	extension      string
	name           uint32
	width          uint32
	height         uint32
	widthF32       float32
	heightF32      float32
	resizable      bool
	mipmap         int32
	mipmapCount    uint32
	format         int32
	internalFormat uint32
	pixelFormat    uint32
	pixelType      uint32
	target         uint32
	face           uint32
}

// GetString ...
func GfxString(value int32) (out string) {
	if value == FORMAT_NONE {
		out = "FORMAT_NONE"
	} else if value == FORMAT_DEPTH24 {
		out = "FORMAT_DEPTH24"
	} else if value == FORMAT_R8_G8_B8_A8 {
		out = "FORMAT_R8_G8_B8_A8"
	} else if value == FORMAT_RGB_16F {
		out = "FORMAT_RGB_16F"
	} else if value == FORMAT_RGBA_16F {
		out = "FORMAT_RGBA_16F"
	} else if value == FORMAT_R8 {
		out = "FORMAT_R8"
	} else if value == FORMAT_STENCIL8 {
		out = "FORMAT_STENCIL8"
	} else if value == FORMAT_DEPTH24_STENCIL8 {
		out = "FORMAT_DEPTH24_STENCIL8"
	} else {
		out = "unknown gfx enum"
	}
	return
}

// InvalidTexture ...
func InvalidTexture() (out TextureId) {
	out.texture = -1
	return
}

// IsValidTexture ...
func IsValidTexture(id TextureId) (out bool) {
	out = int(id.texture) >= 0 && int(id.texture) < len(g_textures)
	return
}

// TexturePrint ...
func TexturePrint(message string, id TextureId) {
	utils.PanicIfNot(IsValidTexture(id), "invalid id")
	TextureBind(id)
	var face uint32 = g_textures[id.texture].target
	var width int32
	gl.GetTexLevelParameteriv(face, 0, gl.TEXTURE_WIDTH, &width)
	utils.PanicIf(GlError(), "gl.GetTexLevelParameteriv")

	var height int32
	gl.GetTexLevelParameteriv(face, 0, gl.TEXTURE_HEIGHT, &height)
	utils.PanicIf(GlError(), "gl.GetTexLevelParameteriv")
	fmt.Printf("%s : path '%s', name %d, width %d, height %d, mipmap %d, target %s, format %s, internal %d, %s, pixel %s, type %s, glWidth %d, glHeight %d\n",
		message,
		g_textures[id.texture].path,
		g_textures[id.texture].name,
		g_textures[id.texture].width,
		g_textures[id.texture].height,
		g_textures[id.texture].mipmap,
		GlString(g_textures[id.texture].target),
		GfxString(g_textures[id.texture].format),
		g_textures[id.texture].internalFormat,
		GlString(g_textures[id.texture].internalFormat),
		GlString(g_textures[id.texture].pixelFormat),
		GlString(g_textures[id.texture].pixelType),
		width, height)
}

// TextureName ...
func TextureName(id TextureId) (out uint32) {
	utils.PanicIfNot(IsValidTexture(id), "")
	out = g_textures[id.texture].name
	return
}

// TextureTarget ...
func TextureTarget(id TextureId) (out uint32) {
	utils.PanicIfNot(IsValidTexture(id), "")
	out = g_textures[id.texture].target
	return
}

// TextureGetMipmapCount ...
func TextureGetMipmapCount(id TextureId) (out uint32) {
	utils.PanicIfNot(IsValidTexture(id), "")
	out = g_textures[id.texture].mipmapCount
	return
}

// TextureWidth ...
func TextureWidth(id TextureId) (out uint32) {
	utils.PanicIfNot(IsValidTexture(id), "")
	out = g_textures[id.texture].width
	return
}

// TextureHeight ...
func TextureHeight(id TextureId) (out uint32) {
	utils.PanicIfNot(IsValidTexture(id), "")
	out = g_textures[id.texture].height
	return
}

// TextureWidthF32 ...
func TextureWidthF32(id TextureId) (out float32) {
	utils.PanicIfNot(IsValidTexture(id), "")
	out = g_textures[id.texture].widthF32
	return
}

// TextureHeightF32 ...
func TextureHeightF32(id TextureId) (out float32) {
	utils.PanicIfNot(IsValidTexture(id), "")
	out = g_textures[id.texture].heightF32
	return
}

/*// TextureSetFormat ...
func TextureSetFormat(id TextureId, internalFormat int32, pixelFormat int32, pixelType int32) {
	utils.PanicIfNot(IsValidTexture(id), "")
	g_textures[id.texture].internalFormat = internalFormat
	g_textures[id.texture].pixelFormat = pixelFormat
	g_textures[id.texture].pixelType = pixelType
}

// TextureSetMipmap ...
func TextureSetMipmap(id TextureId, mipmap int32) {
	utils.PanicIfNot(IsValidTexture(id), "")
	g_textures[id.texture].mipmap = mipmap
}

// TextureSetResizable ...
func TextureSetResizable(id TextureId, resizable bool) {
	utils.PanicIfNot(IsValidTexture(id), "")
	g_textures[id.texture].resizable = resizable
}*/

func textureCreate0(path string, target uint32, format int32, width uint32, height uint32, mipmap int32, resizable bool, cpuCopy bool) (out TextureId) {
	var name uint32
	gl.GenTextures(1, &name) // TODO : change gl.GenTextures signature
	utils.PanicIf(GlError(), "gl.GenTextures")
	var extension string
	var pattern string = path
	if path != "" {
		var extensionIndex int = strings.LastIndex(path, ".")
		var pathLen int = len(path)
		if extensionIndex >= 0 && extensionIndex < pathLen {
			extension = path[extensionIndex:pathLen]
			pattern = path[:extensionIndex]
		}

		var textureCount int = len(g_textures)
		for t := 0; t < textureCount; t++ {
			if g_textures[t].path == path &&
				g_textures[t].extension == extension &&
				g_textures[t].format == format {
				out.texture = int32(t)
				return
			}
		}

		gl.BindTexture(target, name)
		utils.PanicIf(GlError(), "gl.BindTexture")

		gl.TexParameteri(target, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		utils.PanicIf(GlError(), "gl.TexParameteri")

		gl.TexParameteri(target, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
		utils.PanicIf(GlError(), "gl.TexParameteri")

		gl.TexParameteri(target, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		utils.PanicIf(GlError(), "gl.TexParameteri")

		gl.TexParameteri(target, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		utils.PanicIf(GlError(), "gl.TexParameteri")

		if target == gl.TEXTURE_2D {
			uploadTexture(path, target, 0, cpuCopy)
		} else if target == gl.TEXTURE_CUBE_MAP {
			gl.TexParameteri(target, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
			uploadTexture(fmt.Sprintf("%sright%s", pattern, extension), gl.TEXTURE_CUBE_MAP_POSITIVE_X, 0, cpuCopy)
			uploadTexture(fmt.Sprintf("%sleft%s", pattern, extension), gl.TEXTURE_CUBE_MAP_NEGATIVE_X, 0, cpuCopy)
			uploadTexture(fmt.Sprintf("%stop%s", pattern, extension), gl.TEXTURE_CUBE_MAP_POSITIVE_Y, 0, cpuCopy)
			uploadTexture(fmt.Sprintf("%sbottom%s", pattern, extension), gl.TEXTURE_CUBE_MAP_NEGATIVE_Y, 0, cpuCopy)
			uploadTexture(fmt.Sprintf("%sfront%s", pattern, extension), gl.TEXTURE_CUBE_MAP_POSITIVE_Z, 0, cpuCopy)
			uploadTexture(fmt.Sprintf("%sback%s", pattern, extension), gl.TEXTURE_CUBE_MAP_NEGATIVE_Z, 0, cpuCopy)
		}
		utils.PanicIf(GlError(), "gl.NewTexture")
	}

	out = textureCreate(name, target, format, width, height, mipmap, resizable)

	if path != "" {
		if mipmap > 0 {
			var mipmapCount uint32 = TextureGetMipmapCount(out)
			if target == gl.TEXTURE_2D {
				for i := uint32(1); i < mipmapCount; i++ {
					uploadTexture(fmt.Sprintf("%s_%d%s", pattern, i, extension), target, i, cpuCopy)
				}
			} else if target == gl.TEXTURE_CUBE_MAP {
				for i := uint32(1); i < mipmapCount; i++ {
					uploadTexture(fmt.Sprintf("%sright_%d%s", pattern, i, extension), gl.TEXTURE_CUBE_MAP_POSITIVE_X, i, cpuCopy)
					uploadTexture(fmt.Sprintf("%sleft_%d%s", pattern, i, extension), gl.TEXTURE_CUBE_MAP_NEGATIVE_X, i, cpuCopy)
					uploadTexture(fmt.Sprintf("%stop_%d%s", pattern, i, extension), gl.TEXTURE_CUBE_MAP_POSITIVE_Y, i, cpuCopy)
					uploadTexture(fmt.Sprintf("%sbottom_%d%s", pattern, i, extension), gl.TEXTURE_CUBE_MAP_NEGATIVE_Y, i, cpuCopy)
					uploadTexture(fmt.Sprintf("%sfront_%d%s", pattern, i, extension), gl.TEXTURE_CUBE_MAP_POSITIVE_Z, i, cpuCopy)
					uploadTexture(fmt.Sprintf("%sfront_%d%s", pattern, i, extension), gl.TEXTURE_CUBE_MAP_NEGATIVE_Z, i, cpuCopy)
				}
			}
		}
	}

	g_textures[out.texture].path = path
	g_textures[out.texture].extension = extension
	return
}

// TextureCreate ...
func TextureCreate(path string, format int32, width uint32, height uint32, mipmap int32, resizable bool, cpuCopy bool) (out TextureId) {
	out = textureCreate0(path, gl.TEXTURE_2D, format, width, height, mipmap, resizable, cpuCopy)
	return
}

// TextureCreateCube ...
func TextureCreateCube(path string, format int32, width uint32, height uint32, mipmap int32, resizable bool) (out TextureId) {
	out = textureCreate0(path, gl.TEXTURE_CUBE_MAP, format, width, height, mipmap, resizable, false)
	return
}

// TextureGetPath ...
func TextureGetPath(id TextureId) (out string) {
	utils.PanicIfNot(IsValidTexture(id), "invalid id")
	out = g_textures[id.texture].path
	return
}

// TextureGetPixel ...
func TextureGetPixel(id TextureId, x uint32, y uint32) (r float32, g float32, b float32, a float32) {
	utils.PanicIfNot(IsValidTexture(id), "invalid id")

	r, g, b, a = opCxTextureGetPixel(g_textures[id.texture].path, x, y)
	return
}

// TextureResize ...
func TextureResize(id TextureId, width uint32, height uint32) {
	utils.PanicIfNot(IsValidTexture(id), "invalid id")
	//gl.ActiveTexture(gl.TEXTURE0)
	TextureBind(id)
	//utils.PanicIfNot(g_textures[id.texture].resizable, "texture is not resizable")
	//utils.PanicIf(g_textures[id.texture].mipmap != 0, "not implemented")

	var target uint32 = g_textures[id.texture].target
	if g_textures[id.texture].width != width || g_textures[id.texture].height != height {
		if target == gl.TEXTURE_2D {
			gl.TexImage2D(gl.TEXTURE_2D, 0,
				int32(g_textures[id.texture].internalFormat),
				int32(width),
				int32(height),
				0,
				g_textures[id.texture].pixelFormat,
				g_textures[id.texture].pixelType,
				nil)
			utils.PanicIf(GlError(), "gl.TexImage2D")
			textureResize(id, width, height)
		} else if target == gl.TEXTURE_CUBE_MAP {
			for i := 0; i < 6; i++ {
				gl.TexImage2D(uint32(gl.TEXTURE_CUBE_MAP_POSITIVE_X+i), 0,
					int32(g_textures[id.texture].internalFormat),
					int32(width),
					int32(height),
					0,
					g_textures[id.texture].pixelFormat,
					g_textures[id.texture].pixelType,
					nil)
				utils.PanicIf(GlError(), "gl.TexImage2D")
				textureResize(id, width, height)
			}
		}
	}
}

// TextureBind ...
func TextureBind(id TextureId) {
	utils.PanicIfNot(IsValidTexture(id), "")
	bindTexture(g_textures[id.texture].target, g_textures[id.texture].name)
}

// TextureFilter ...
func TextureFilter(id TextureId, min int32, mag int32) {
	utils.PanicIfNot(IsValidTexture(id), "")
	TextureBind(id)

	var target uint32 = g_textures[id.texture].target
	gl.TexParameteri(target, gl.TEXTURE_MIN_FILTER, min)
	utils.PanicIf(GlError(), "gl.TexParameteri")

	gl.TexParameteri(target, gl.TEXTURE_MAG_FILTER, mag)
	utils.PanicIf(GlError(), "gl.TexParameteri")

}

// TextureWrap ...
func TextureWrap(id TextureId, wrapS int32, wrapT int32, wrapR int32) {
	utils.PanicIfNot(IsValidTexture(id), "")
	TextureBind(id)

	var target uint32 = g_textures[id.texture].target
	gl.TexParameteri(target, gl.TEXTURE_WRAP_S, wrapS)
	utils.PanicIf(GlError(), "gl.TexParameteri")

	gl.TexParameteri(target, gl.TEXTURE_WRAP_T, wrapT)
	utils.PanicIf(GlError(), "gl.TexParameteri")

	gl.TexParameteri(target, gl.TEXTURE_WRAP_R, wrapR)
	utils.PanicIf(GlError(), "gl.TexParameteri")
}

// TextureSamplerState ...
func TextureSamplerState(id TextureId, sampler SamplerState) {
	TextureFilter(id, sampler.min, sampler.mag)
	TextureWrap(id, sampler.s, sampler.t, sampler.r)
}

// TextureGenerateMipmap ...
func TextureGenerateMipmap(id TextureId) {
	utils.PanicIfNot(IsValidTexture(id), "invalid texture")

	TextureBind(id)

	var target uint32 = g_textures[id.texture].target
	gl.GenerateMipmap(target)
	utils.PanicIf(GlError(), "gl.GenerateMipmap")

	g_textures[id.texture].mipmap = 0
}

// DestroyTextures ...
func DestroyTextures() {
	for i := 0; i < len(g_textures); i++ {
		gl.DeleteTextures(1, &g_textures[i].name)
		utils.PanicIf(GlError(), "g.DeleteTextures")
	}
}

// ResizeTextures ...
/*func ResizeTextures(width int32, height int32)() {
	fmt.Printf("COUNT %d\n", len(g_textures))
	for i := 0 ; i < len(g_textures); i++ {
		TextureResize(g_textures[i].id, width, height)
	}
}*/

func textureResize(id TextureId, width uint32, height uint32) {
	utils.PanicIfNot(IsValidTexture(id), "")
	g_textures[id.texture].width = width
	g_textures[id.texture].height = height
	g_textures[id.texture].widthF32 = float32(width)
	g_textures[id.texture].heightF32 = float32(height)
}

func getGlFormat(format int32) (internalFormat uint32, pixelFormat uint32, pixelType uint32) { // TODO : remove pixelType, only used when uploading data.
	if format == FORMAT_DEPTH24 {
		internalFormat = gl.DEPTH_COMPONENT24
		pixelFormat = gl.DEPTH_COMPONENT
		pixelType = gl.UNSIGNED_INT
	} else if format == FORMAT_DEPTH24_STENCIL8 {
		internalFormat = gl.DEPTH24_STENCIL8
		pixelFormat = gl.DEPTH_STENCIL
		pixelType = gl.UNSIGNED_INT_24_8
	} else if format == FORMAT_STENCIL8 {
		internalFormat = gl.STENCIL_INDEX8
		pixelFormat = gl.STENCIL_INDEX
		pixelType = gl.UNSIGNED_BYTE
	} else if format == FORMAT_R8_G8_B8_A8 {
		internalFormat = gl.RGBA8
		pixelFormat = gl.RGBA
		pixelType = gl.UNSIGNED_BYTE
	} else if format == FORMAT_RGB_16F {
		internalFormat = gl.RGB16F
		pixelFormat = gl.RGB
		pixelType = gl.FLOAT
	} else if format == FORMAT_RGBA_16F {
		internalFormat = gl.RGBA16F
		pixelFormat = gl.RGBA
		pixelType = gl.HALF_FLOAT
	} else if format == FORMAT_R8 {
		internalFormat = gl.R8
		pixelFormat = gl.RED
		pixelType = gl.UNSIGNED_BYTE
	}
	return
}

func textureCreate(name uint32, target uint32, format int32, width uint32, height uint32, mipmap int32, resizable bool) (out TextureId) {
	out.texture = int32(len(g_textures))

	var texture Texture
	texture.id = out
	texture.resizable = resizable
	texture.mipmap = mipmap
	texture.target = target
	texture.format = format
	texture.internalFormat, texture.pixelFormat, texture.pixelType = getGlFormat(format)
	texture.name = name
	texture.target = target

	var face uint32 = 0
	if target == gl.TEXTURE_2D {
		face = gl.TEXTURE_2D
	} else if target == gl.TEXTURE_CUBE_MAP {
		face = gl.TEXTURE_CUBE_MAP_POSITIVE_X
	}

	texture.face = face
	g_textures = append(g_textures, texture)

	gl.BindTexture(target, texture.name)
	utils.PanicIf(GlError(), "gl.BindTexture")

	if width > 0 && height > 0 {
		TextureResize(out, width, height)
	} else {
		var iwidth int32
		var iheight int32
		gl.GetTexLevelParameteriv(face, 0, gl.TEXTURE_WIDTH, &iwidth)
		utils.PanicIf(GlError(), "gl.GetTexLevelParameteriv")
		gl.GetTexLevelParameteriv(face, 0, gl.TEXTURE_HEIGHT, &iheight)
		utils.PanicIf(GlError(), "gl.GetTexLevelParameteriv")
		width = uint32(iwidth)
		height = uint32(iheight)
		textureResize(out, width, height)
	}

	if texture.mipmap != 0 {
		TextureGenerateMipmap(out)
	}

	var textureWidth uint32 = TextureWidth(out)
	var textureHeight uint32 = TextureHeight(out)

	if texture.mipmap != 0 {
		var mipmapCount uint32 = 0
		for (textureWidth > 0) && (textureHeight > 0) {
			textureWidth = textureWidth / 2
			textureHeight = textureHeight / 2
			mipmapCount = mipmapCount + 1
		}
		mipmapCount = uint32(math.Min_i32(int32(mipmapCount), math.Abs_i32(mipmap)))
		g_textures[out.texture].mipmapCount = mipmapCount
	}

	utils.PanicIfNot(IsValidTexture(out), "")
	return
}
