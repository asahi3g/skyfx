package utils

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	// "bytes"
	// "encoding/binary"
	// "math"
	// "os"
	// "os/exec"
	// "strings"
	// "syscall"
	// "time"
	// "github.com/skycoin/cx/cx/ast"
	// "github.com/skycoin/cx/cx/constants"
	// "github.com/skycoin/cx/cx/globals"
	// "github.com/skycoin/cx/cx/types"
	// "github.com/skycoin/cx/cx/util"
)
var workingDir string
var logFile bool = true

// CXSetWorkingDir ...
func CXSetWorkingDir(dir string) {
	workingDir = dir
}

// CXLogFile ...
func CXLogFile(enable bool) {
	logFile = enable
}

// CXOpenFile ...
func CXOpenFile(filename string) (*os.File, error) {
	filename = filepath.Join(workingDir, filename)

	if logFile {
		fmt.Printf("CXOpenFile: Opening '%s'\n", filename)
	}

	file, err := os.Open(filename)
	if logFile && err != nil {
		fmt.Printf("CXOpenFile: Failed to open '%s': %v\n", filename, err)
	}
	return file, err
}

// CXCreateFile ...
func CXCreateFile(filename string) (*os.File, error) {
	filename = filepath.Join(workingDir, filename)

	if logFile {
		fmt.Printf("Creating file : '%s', '%s'\n", workingDir, filename)
	}

	file, err := os.Create(filepath.Join(workingDir, filename))
	if logFile && err != nil {
		fmt.Printf("Failed to create file : '%s', '%s', err '%v'\n", workingDir, filename, err)
	}

	return file, err
}

// CXRemoveFile ...
func CXRemoveFile(path string) error {
	if logFile {
		fmt.Printf("Removing file : '%s', '%s'\n", workingDir, path)
	}

	err := os.Remove(fmt.Sprintf("%s%s", workingDir, path))

	if logFile && err != nil {
		fmt.Printf("Failed to remove file : '%s', '%s', err '%v'\n", workingDir, path, err)
	}

	return err
}

// CXReadFile ...
func CXReadFile(path string) ([]byte, error) {
	if logFile {
		fmt.Printf("Reading file : '%s', '%s'\n", workingDir, path)
	}

	bytes, err := ioutil.ReadFile(fmt.Sprintf("%s%s", workingDir, path))

	if logFile && err != nil {
		fmt.Printf("Failed to read file : '%s', '%s', err '%v'\n", workingDir, path, err)
	}

	return bytes, err
}

// CXStatFile ...
func CXStatFile(path string) (os.FileInfo, error) {
	if logFile {
		fmt.Printf("Stating file : '%s', '%s'\n", workingDir, path)
	}

	fileInfo, err := os.Stat(fmt.Sprintf("%s%s", workingDir, path))

	if logFile && err != nil {
		fmt.Printf("Failed to stat file : '%s', '%s', err '%v'\n", workingDir, path, err)
	}

	return fileInfo, err
}

// CXMkdirAll ...
func CXMkdirAll(path string, perm os.FileMode) error {
	if logFile {
		fmt.Printf("Creating dir : '%s'\n", path)
	}

	err := os.MkdirAll(fmt.Sprintf("%s%s", workingDir, path), perm)

	if logFile && err != nil {
		fmt.Printf("Failed to create dir : '%s', '%s', err '%v'\n", workingDir, path, err)
	}

	return err
}


const (
	OS_SEEK_SET = iota
	OS_SEEK_CUR
	OS_SEEK_END
)

var openFiles []*os.File
var freeFiles []int32

// helper function used to validate file handle from int32
func ValidFile(handle int32) *os.File {
	if handle >= 0 && handle < int32(len(openFiles)) && openFiles[handle] != nil {
		return openFiles[handle]
	}
	return nil
}

// func LogFile(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
// 	util.CXLogFile(inputs[0].Get_bool(prgrm))
// }

// func ReadAllText(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
// 	success := false

// 	if byts, err := util.CXReadFile(inputs[0].Get_str(prgrm)); err == nil {
// 		outputs[0].Set_str(prgrm, string(byts))
// 		success = true
// 	}

// 	outputs[1].Set_bool(prgrm, success)
// }

func getFileHandle(file *os.File) int32 {
	handle := int32(-1)
	freeCount := len(freeFiles)
	if freeCount > 0 {
		freeCount--
		handle = int32(freeFiles[freeCount])
		freeFiles = freeFiles[:freeCount]
	} else {
		handle = int32(len(openFiles))
		openFiles = append(openFiles, nil)
	}

	if handle < 0 || handle >= int32(len(openFiles)) {
		panic("internal error")
	}

	openFiles[handle] = file
	return handle
}

func Open(path string) (handle int32) {
	handle = int32(-1)
	if file, err := CXOpenFile(path); err == nil {
		handle = getFileHandle(file)
	}
	return
}

func Create(path string) (handle int32) {
	handle = int32(-1)
	if file, err := CXCreateFile(path); err == nil {
		handle = getFileHandle(file)
	}
	return
}

func Close(handle int32) (success bool) {
	if file := ValidFile(handle); file != nil {
		if err := file.Close(); err == nil {
			success = true
		}

		openFiles[handle] = nil
		freeFiles = append(freeFiles, handle)
	}

	return
}

func Seek(handle int32, ioffset int64, whence int32) (offset int64) {
	offset = int64(-1)
	if file := ValidFile(handle); file != nil {
		var err error
		if offset, err = file.Seek(ioffset, int(whence)); err != nil {
			offset = -1
		}
	}
	return
}

// func ReadStr(handle int32) (value string, success bool) {
// 	var len uint64
// 	if file := ValidFile(handle); file != nil {
// 		if err := binary.Read(file, binary.LittleEndian, &len); err == nil {
// 			bytes := make([]byte, len)
// 			if err := binary.Read(file, binary.LittleEndian, &bytes); err == nil {
// 				value = string(bytes)
// 				success = true
// 			}
// 		}
// 	}
// 	return
// }

func ReadF64(handle int32) (value float64, success bool) {
	if file := ValidFile(handle); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}
	return
}

func ReadF32(handle int32) (value float32, success bool) {
	if file := ValidFile(handle); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}
	return
}

func ReadUI64(handle int32) (value uint64, success bool) {
	if file := ValidFile(handle); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}
	return
}

func ReadUI32(handle int32) (value uint32, success bool) {
	if file := ValidFile(handle); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}
	return
}

func ReadUI16(handle int32) (value uint16, success bool) {
	if file := ValidFile(handle); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}
	return
}

func ReadUI8(handle int32) (value int8, success bool) {
	if file := ValidFile(handle); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}
	return
}

func ReadI64(handle int32) (value int64, success bool) {
	if file := ValidFile(handle); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}
	return
}

func ReadI32(handle int32) (value int32, success bool) {
	if file := ValidFile(handle); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}
	return
}

func ReadI16(handle int32) (value int16, success bool) {
	if file := ValidFile(handle); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}
	return
}

func ReadI8(handle int32) (value int8, success bool) {
	if file := ValidFile(handle); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}
	return
}

func ReadBOOL(handle int32) (value bool, success bool) {
	if file := ValidFile(handle); file != nil {
		if err := binary.Read(file, binary.LittleEndian, &value); err == nil {
			success = true
		}
	}
	return
}

// func WriteStr(handle int32, value string) (success bool) {
// 	if file := ValidFile(handle); file != nil {
// 		len := len(value)
// 		if err := binary.Write(file, binary.LittleEndian, uint64(len)); err == nil {
// 			if err := binary.Write(file, binary.LittleEndian, []byte(value)); err == nil {
// 				success = true
// 			}
// 		}
// 	}

// 	return
// }

func WriteF64(handle int32, value float64) (success bool) {
	if file := ValidFile(handle); file != nil {
		if err := binary.Write(file, binary.LittleEndian, value); err == nil {
			success = true
		}
	}

	return
}

func WriteF32(handle int32, value float32) (success bool) {
	if file := ValidFile(handle); file != nil {
		if err := binary.Write(file, binary.LittleEndian, value); err == nil {
			success = true
		}
	}
	
	return
}

func WriteUI64(handle int32, value uint64) (success bool) {
	if file := ValidFile(handle); file != nil {
		if err := binary.Write(file, binary.LittleEndian, value); err == nil {
			success = true
		}
	}

	return
}

func WriteUI32(handle int32, value uint32) (success bool) {
	if file := ValidFile(handle); file != nil {
		if err := binary.Write(file, binary.LittleEndian, value); err == nil {
			success = true
		}
	}

	return
}

func WriteUI16(handle int32, value uint16) (success bool) {
	if file := ValidFile(handle); file != nil {
		if err := binary.Write(file, binary.LittleEndian, value); err == nil {
			success = true
		}
	}

	return
}

func WriteUI8(handle int32, value uint8) (success bool) {
	if file := ValidFile(handle); file != nil {
		if err := binary.Write(file, binary.LittleEndian, value); err == nil {
			success = true
		}
	}

	return
}

func WriteI64(handle int32, value int64) (success bool) {
	if file := ValidFile(handle); file != nil {
		if err := binary.Write(file, binary.LittleEndian, value); err == nil {
			success = true
		}
	}

	return
}

func WriteI32(handle int32, value int32) (success bool) {
	if file := ValidFile(handle); file != nil {
		if err := binary.Write(file, binary.LittleEndian, value); err == nil {
			success = true
		}
	}

	return
}

func WriteI16(handle int32, value int16) (success bool) {
	if file := ValidFile(handle); file != nil {
		if err := binary.Write(file, binary.LittleEndian, value); err == nil {
			success = true
		}
	}

	return
}

func WriteI8(handle int32, value int8) (success bool) {
	if file := ValidFile(handle); file != nil {
		if err := binary.Write(file, binary.LittleEndian, value); err == nil {
			success = true
		}
	}

	return
}

func WriteBool(handle int32, value bool) (success bool) {
	if file := ValidFile(handle); file != nil {
		if err := binary.Write(file, binary.LittleEndian, value); err == nil {
			success = true
		}
	}

	return
}

func ReadF64Slice(handle int32, dest []float64, count uint64) (success bool) {
	if count > 0 {
		if file := ValidFile(handle); file != nil {
			for i := uint64(0); i < count; i++ {
				if err := binary.Read(file, binary.LittleEndian, &dest[i]); err == nil {
					success = true
				}
			}
		}
	}

	return
}

func ReadF32Slice(handle int32, dest []float32, count uint64) (success bool) {
	if count > 0 {
		if file := ValidFile(handle); file != nil {
			for i := uint64(0); i < count; i++ {
				if err := binary.Read(file, binary.LittleEndian, &dest[i]); err == nil {
					success = true
				}
			}
		}
	}

	return
}

func ReadUI64Slice(handle int32, dest []uint64, count uint64) (success bool) {
	if count > 0 {
		if file := ValidFile(handle); file != nil {
			for i := uint64(0); i < count; i++ {
				if err := binary.Read(file, binary.LittleEndian, &dest[i]); err == nil {
					success = true
				}
			}
		}
	}

	return
}

func ReadUI32Slice(handle int32, dest []uint32, count uint64) (success bool) {
	if count > 0 {
		if file := ValidFile(handle); file != nil {
			for i := uint64(0); i < count; i++ {
				if err := binary.Read(file, binary.LittleEndian, &dest[i]); err == nil {
					success = true
				}
			}
		}
	}

	return
}

func ReadUI16Slice(handle int32, dest []uint16, count uint64) (success bool) {
	if count > 0 {
		if file := ValidFile(handle); file != nil {
			for i := uint64(0); i < count; i++ {
				if err := binary.Read(file, binary.LittleEndian, &dest[i]); err == nil {
					success = true
				}
			}
		}
	}

	return
}

func ReadUI8Slice(handle int32, dest []uint8, count uint64) (success bool) {
	if count > 0 {
		if file := ValidFile(handle); file != nil {
			for i := uint64(0); i < count; i++ {
				if err := binary.Read(file, binary.LittleEndian, &dest[i]); err == nil {
					success = true
				}
			}
		}
	}

	return
}

func ReadI64Slice(handle int32, dest []int64, count uint64) (success bool) {
	if count > 0 {
		if file := ValidFile(handle); file != nil {
			for i := uint64(0); i < count; i++ {
				if err := binary.Read(file, binary.LittleEndian, &dest[i]); err == nil {
					success = true
				}
			}
		}
	}

	return
}

func ReadI32Slice(handle int32, dest []int32, count uint64) (success bool) {
	if count > 0 {
		if file := ValidFile(handle); file != nil {
			for i := uint64(0); i < count; i++ {
				if err := binary.Read(file, binary.LittleEndian, &dest[i]); err == nil {
					success = true
				}
			}
		}
	}

	return
}

func ReadI16Slice(handle int32, dest []int16, count uint64) (success bool) {
	if count > 0 {
		if file := ValidFile(handle); file != nil {
			for i := uint64(0); i < count; i++ {
				if err := binary.Read(file, binary.LittleEndian, &dest[i]); err == nil {
					success = true
				}
			}
		}
	}

	return
}

func ReadI8Slice(handle int32, dest []int8, count uint64) (success bool) {
	if count > 0 {
		if file := ValidFile(handle); file != nil {
			for i := uint64(0); i < count; i++ {
				if err := binary.Read(file, binary.LittleEndian, &dest[i]); err == nil {
					success = true
				}
			}
		}
	}

	return
}

func WriteF64Slice(handle int32, value []float64) (success bool) {
	if file := ValidFile(handle); file != nil {
		if value != nil {
			count := len(value)
			for i := 0; i < count; i++ {
				if err := binary.Write(file, binary.LittleEndian, value[i]); err == nil {
					success = true
				}
			}
		}
	}
	return
}

func WriteF32Slice(handle int32, value []float32) (success bool) {
	if file := ValidFile(handle); file != nil {
		if value != nil {
			count := len(value)
			for i := 0; i < count; i++ {
				if err := binary.Write(file, binary.LittleEndian, value[i]); err == nil {
					success = true
				}
			}
		}
	}
	return
}

func WriteUI64Slice(handle int32, value []uint64) (success bool) {
	if file := ValidFile(handle); file != nil {
		if value != nil {
			count := len(value)
			for i := 0; i < count; i++ {
				if err := binary.Write(file, binary.LittleEndian, value[i]); err == nil {
					success = true
				}
			}
		}
	}
	return
}

func WriteUI32Slice(handle int32, value []uint32) (success bool) {
	if file := ValidFile(handle); file != nil {
		if value != nil {
			count := len(value)
			for i := 0; i < count; i++ {
				if err := binary.Write(file, binary.LittleEndian, value[i]); err == nil {
					success = true
				}
			}
		}
	}
	return
}

func WriteUI16Slice(handle int32, value []uint16) (success bool) {
	if file := ValidFile(handle); file != nil {
		if value != nil {
			count := len(value)
			for i := 0; i < count; i++ {
				if err := binary.Write(file, binary.LittleEndian, value[i]); err == nil {
					success = true
				}
			}
		}
	}
	return
}

func WriteUI8Slice(handle int32, value []uint8) (success bool) {
	if file := ValidFile(handle); file != nil {
		if value != nil {
			count := len(value)
			for i := 0; i < count; i++ {
				if err := binary.Write(file, binary.LittleEndian, value[i]); err == nil {
					success = true
				}
			}
		}
	}
	return
}


func WriteI64Slice(handle int32, value []int64) (success bool) {
	if file := ValidFile(handle); file != nil {
		if value != nil {
			count := len(value)
			for i := 0; i < count; i++ {
				if err := binary.Write(file, binary.LittleEndian, value[i]); err == nil {
					success = true
				}
			}
		}
	}
	return
}

func WriteI32Slice(handle int32, value []int32) (success bool) {
	if file := ValidFile(handle); file != nil {
		if value != nil {
			count := len(value)
			for i := 0; i < count; i++ {
				if err := binary.Write(file, binary.LittleEndian, value[i]); err == nil {
					success = true
				}
			}
		}
	}
	return
}

func WriteI16Slice(handle int32, value []int16) (success bool) {
	if file := ValidFile(handle); file != nil {
		if value != nil {
			count := len(value)
			for i := 0; i < count; i++ {
				if err := binary.Write(file, binary.LittleEndian, value[i]); err == nil {
					success = true
				}
			}
		}
	}
	return
}

func WriteI8Slice(handle int32, value []int8) (success bool) {
	if file := ValidFile(handle); file != nil {
		if value != nil {
			count := len(value)
			for i := 0; i < count; i++ {
				if err := binary.Write(file, binary.LittleEndian, value[i]); err == nil {
					success = true
				}
			}
		}
	}
	return
}


// func GetWorkingDirectory(prgrm *ast.CXProgram, inputs []ast.CXValue, outputs []ast.CXValue) {
// 	//outputs[0].Set_str(prgrm,prgrm,cxcore.PROGRAM.Path)
// 	outputs[0].Set_str(prgrm, globals.CxProgramPath)
// }

