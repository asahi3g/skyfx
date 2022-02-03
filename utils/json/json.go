package json

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"skyfx/utils"
)

const (
	JSON_TOKEN_INVALID = -1
	JSON_TOKEN_NULL    = iota
	JSON_TOKEN_DELIM
	JSON_TOKEN_BOOL
	JSON_TOKEN_F64
	JSON_TOKEN_NUMBER
	JSON_TOKEN_STR
	JSON_DELIM_SQUARE_LEFT  = 91
	JSON_DELIM_SQUARE_RIGHT = 93
	JSON_DELIM_CURLY_LEFT   = 123
	JSON_DELIM_CURLY_RIGHT  = 125
)

type JSONFile struct {
	file        *os.File
	reader      *bufio.Reader
	decoder     *json.Decoder
	token       interface{}
	tokenType   int32
	tokenDelim  json.Delim
	tokenBool   bool
	tokenF64    float64
	tokenNumber json.Number
	tokenStr    string
}

var jsons []JSONFile
var freeJsons []int32

// Open the named json file for reading, returns an int32 identifying the json cxgo.
func Open(path string) (handle int32) {
	handle = int32(-1)

	file, err := utils.CXOpenFile(path)
	if err == nil {
		freeCount := len(freeJsons)
		if freeCount > 0 {
			freeCount--
			handle = int32(freeJsons[freeCount])
			freeJsons = freeJsons[:freeCount]
		} else {
			handle = int32(len(jsons))
			jsons = append(jsons, JSONFile{})
		}

		if handle < 0 || handle >= int32(len(jsons)) {
			panic("internal error")
		}

		var jsonFile JSONFile
		jsonFile.file = file
		jsonFile.reader = bufio.NewReader(file)
		jsonFile.decoder = json.NewDecoder(jsonFile.reader)
		jsonFile.decoder.UseNumber()

		jsons[handle] = jsonFile
	}
	return
}

// Close json cxgo (and all underlying resources) idendified by it's int32 handle.
func Close(handle int32) (success bool) {
	if jsonFile := validJsonFile(handle); jsonFile != nil {
		if err := jsonFile.file.Close(); err != nil {
			panic(err)
		}

		jsons[handle] = JSONFile{}
		freeJsons = append(freeJsons, handle)
		success = true
	}

	return
}

// More return true if there is another element in the current array or object being parsed.
func TokenMore(handle int32) (more bool, success bool) {
	if jsonFile := validJsonFile(handle); jsonFile != nil {
		more = jsonFile.decoder.More()
		success = true
	}
	return
}

// Token parses the next token.
func TokenNext(handle int32) (tokenType int32, success bool) {
	tokenType = int32(JSON_TOKEN_INVALID)

	if jsonFile := validJsonFile(handle); jsonFile != nil {
		token, err := jsonFile.decoder.Token()
		if err == io.EOF {
			tokenType = JSON_TOKEN_NULL
			success = true
		} else if err == nil {
			jsonFile.token = token
			switch value := token.(type) {
			case json.Delim:
				tokenType = JSON_TOKEN_DELIM
				jsonFile.tokenDelim = value
				success = true
			case bool:
				tokenType = JSON_TOKEN_BOOL
				jsonFile.tokenBool = value
				success = true
			case float64:
				tokenType = JSON_TOKEN_F64
				jsonFile.tokenF64 = value
				success = true
			case json.Number:
				tokenType = JSON_TOKEN_NUMBER
				jsonFile.tokenNumber = value
				success = true
			case string:
				tokenType = JSON_TOKEN_STR
				jsonFile.tokenStr = value
				success = true
			default:
				if value == nil {
					tokenType = JSON_TOKEN_NULL
					success = true
				}
			}
		}
		jsonFile.tokenType = tokenType
	}
	return
}

// Type returns the type of the current token.
func TokenType(handle int32) (tokenType int32, success bool) {
	tokenType = int32(JSON_TOKEN_INVALID)

	if jsonFile := validJsonFile(handle); jsonFile != nil {
		tokenType = jsonFile.tokenType
		success = true
	}
	return
}

// Delim returns current token as an int32 delimiter.
func TokenDelim(handle int32) (tokenDelim int32, success bool) {
	tokenDelim = int32(JSON_TOKEN_INVALID)

	if jsonFile := validJsonFile(handle); jsonFile != nil {
		if jsonFile.tokenType == JSON_TOKEN_DELIM {
			tokenDelim = int32(jsonFile.tokenDelim)
			success = true
		}
	}
	return
}

// Bool returns current token as a bool value.
func TokenBool(handle int32) (tokenBool bool, success bool) {
	if jsonFile := validJsonFile(handle); jsonFile != nil {
		if jsonFile.tokenType == JSON_TOKEN_BOOL {
			tokenBool = jsonFile.tokenBool
			success = true
		}
	}
	return
}

// Float64 returns current token as float64 value.
func TokenF64(handle int32) (tokenF64 float64, success bool) {
	if jsonFile := validJsonFile(handle); jsonFile != nil {
		if jsonFile.tokenType == JSON_TOKEN_F64 {
			tokenF64 = jsonFile.tokenF64
			success = true
		} else if jsonFile.tokenType == JSON_TOKEN_NUMBER {
			var err error
			if tokenF64, err = jsonFile.tokenNumber.Float64(); err == nil {
				success = true
			}
		}
	}
	return
}

// Int64 returns current token as int64 value.
func TokenI64(handle int32) (tokenI64 int64, success bool) {
	if jsonFile := validJsonFile(handle); jsonFile != nil {
		if jsonFile.tokenType == JSON_TOKEN_NUMBER {
			var err error
			if tokenI64, err = jsonFile.tokenNumber.Int64(); err == nil {
				success = true
			}
		}
	}
	return
}

// Str returns current token as string value.
func TokenStr(handle int32) (tokenStr string, success bool) {
	if jsonFile := validJsonFile(handle); jsonFile != nil {
		if jsonFile.tokenType == JSON_TOKEN_STR {
			tokenStr = jsonFile.tokenStr
			success = true
		}
	}
	return
}

// helper function used to validate json handle from int32
func validJsonFile(handle int32) *JSONFile {
	if handle >= 0 && handle < int32(len(jsons)) && jsons[handle].file != nil {
		return &jsons[handle]
	}
	return nil
}

// Read next token in json file and checks if it's a delimiter matching the delim value, returns true on success
func MatchDelim(file int32, delim int32) (success bool) {
	var value int32
	if ReadDelim(file, &value) {
		success = value == delim
	}
	return
}

// Read and interpret next token in json file as delim value, returns true on success
func ReadDelim(file int32, value *int32) (success bool) {
	//var tokenType int32
	_ /*tokenType*/, success = TokenNext(file)
	if success {
		*value, success = TokenDelim(file)
		if success {
			return
		}
	}

	//debugToken(file, tokenType)
	return
}

// Read and interpret next token in json file as string value, returns true on success
func ReadStr(file int32, value *string) (success bool) {
	//var tokenType int32
	_ /*tokenType*/, success = TokenNext(file)
	if success {
		*value, success = TokenStr(file)
		if success {
			return
		}
	}

	//debugToken(file, tokenType)
	return
}

// Read and interpret next token in json file as int32 value, returns true on success
func ReadBool(file int32, value *bool) (success bool) {
	//var tokenType int32
	_ /*tokenType*/, success = TokenNext(file)
	if success {
		*value, success = TokenBool(file)
		if success {
			return
		}
	}

	//debugToken(file, tokenType)
	return
}

// Read and interpret next token in json file as i64 value, returns true on success
func ReadI64(file int32, value *int64) (success bool) {
	//var tokenType int32
	_ /*tokenType*/, success = TokenNext(file)
	if success {
		*value, success = TokenI64(file)
		if success {
			return
		}
	}

	//debugToken(file, tokenType)
	return
}

// Read and interpret next token in json file as float64 value, returns true on success
func ReadF64(file int32, value *float64) (success bool) {
	//var tokenType int32
	_ /*tokenType*/, success = TokenNext(file)
	if success {
		*value, success = TokenF64(file)
		if success {
			return
		}
	}

	//debugToken(file, tokenType)
	return
}

// Read and interpret next token in json file as int32 value, returns true on success
func ReadI32(file int32, value *int32) (success bool) {
	//var tokenType int32
	_ /*tokenType*/, success = TokenNext(file)
	if success {
		var valueI64 int64
		valueI64, success = TokenI64(file)
		if success {
			*value = int32(valueI64)
			return
		}
	}

	//debugToken(file, tokenType)
	return
}

// Read and interpret next token in json file as float32 value, returns true on success
func ReadF32(file int32, value *float32) (success bool) {
	//var tokenType int32
	_ /*tokenType*/, success = TokenNext(file)
	if success {
		var valueF64 float64
		valueF64, success = TokenF64(file)
		if success {
			*value = float32(valueF64)
			return
		}
	}

	//debugToken(file, tokenType)
	return
}

// Read and interpret next array in json file as [3]float32, returns true on success
// This function move the current token
func ReadF32Vec3(file int32, array *[3]float32) (success bool) {
	if MatchDelim(file, JSON_DELIM_SQUARE_LEFT) == false {
		return
	}

	var i int32
	for i < 3 {
		var more bool
		more, success = TokenMore(file)
		if more == false || success == false {
			success = false
			return
		}
		var value float32
		if ReadF32(file, &value) == false {
			success = false
			return
		}
		(*array)[i] = value
		i++
	}

	success = MatchDelim(file, JSON_DELIM_SQUARE_RIGHT)
	return
}

// Read next array in json file as [4]float32, returns true on success
// This function move the current token
func ReadF32Vec4(file int32, array *[4]float32) (success bool) {
	if MatchDelim(file, JSON_DELIM_SQUARE_LEFT) == false {
		return
	}

	var i int32
	for i < 4 {
		var more bool
		more, success = TokenMore(file)
		if more == false || success == false {
			success = false
			return
		}
		var value float32
		if ReadF32(file, &value) == false {
			success = false
			return
		}
		(*array)[i] = value
		i++
	}

	success = MatchDelim(file, JSON_DELIM_SQUARE_RIGHT)
	return
}

// Read next array in json file as []int32, returns true on success
// This function move the current token
func ReadI32Slice(file int32, array *[]int32) (success bool) {
	if MatchDelim(file, JSON_DELIM_SQUARE_LEFT) == false {
		return
	}

	var more bool = true
	for more == true {
		more, success = TokenMore(file)
		if success == false {
			return
		}
		if more {
			var value int32
			if ReadI32(file, &value) == false {
				success = false
				return
			}
			*array = append(*array, value)
		}
	}

	success = MatchDelim(file, JSON_DELIM_SQUARE_RIGHT)
	return
}

// Read next array in json file as []float32, returns true on success
// This function move the current token
func ReadF32Slice(file int32, array *[]float32) (success bool) {
	if MatchDelim(file, JSON_DELIM_SQUARE_LEFT) == false {
		return
	}

	var more bool = true
	for more == true {
		more, success = TokenMore(file)
		if success == false {
			return
		}
		if more {
			var value float32
			if ReadF32(file, &value) == false {
				success = false
				return
			}
			*array = append(*array, value)
		}
	}

	success = MatchDelim(file, JSON_DELIM_SQUARE_RIGHT)
	return
}

// Read next array in json file as []string, returns true on success
// This function move the current token
func ReadStrSlice(file int32, array *[]string) (success bool) {
	if MatchDelim(file, JSON_DELIM_SQUARE_LEFT) == false {
		return
	}

	var more bool = true
	for more == true {
		more, success = TokenMore(file)
		if success == false {
			return
		}
		if more {
			var value string
			if ReadStr(file, &value) == false {
				success = false
				return
			}
			*array = append(*array, value)
		}
	}

	success = MatchDelim(file, JSON_DELIM_SQUARE_RIGHT)
	return
}

// debug helper
func DebugToken(file int32, t int32) {
	fmt.Printf("DEBUG_JSON_TYPE %d\n", t)
	if t == JSON_TOKEN_DELIM {
		fmt.Printf("DEBUG_JSON_DELIM\n")
		var value int32
		var success bool
		value, success = TokenDelim(file)
		if success {
			if value == JSON_DELIM_CURLY_LEFT {
				fmt.Printf("{\n")
			} else if value == JSON_DELIM_CURLY_RIGHT {
				fmt.Printf("}\n")
			} else if value == JSON_DELIM_SQUARE_LEFT {
				fmt.Printf("[\n")
			} else if value == JSON_DELIM_SQUARE_RIGHT {
				fmt.Printf("]\n")
			} else {
				fmt.Printf("invalid delimiter\n")
			}
		} else {
			utils.PanicIf(false, "failed to parse delimiter value")
		}
	} else if t == JSON_TOKEN_BOOL {
		fmt.Printf("DEBUG_JSON_BOOL\n")
		var value bool
		var success bool
		value, success = TokenBool(file)
		if success {
			if value {
				fmt.Printf("true\n")
			} else {
				fmt.Printf("false\n")
			}
		} else {
			utils.PanicIf(false, "failed to parse bool value")
		}
	} else if t == JSON_TOKEN_F64 {
		fmt.Printf("DEBUG_JSON_F64\n")
		var value float64
		var success bool
		value, success = TokenF64(file)
		if success {
			fmt.Printf("%f\n", value)
		} else {
			utils.PanicIf(false, "failed to parse float64 value")
		}
	} else if t == JSON_TOKEN_NUMBER {
		fmt.Printf("DEBUG_JSON_NUMBER\n")
		var value float64
		var success bool
		value, success = TokenF64(file)
		if success {
			fmt.Printf("%f\n", value)
		} else {
			utils.PanicIf(false, "failed to parse number value")
		}
	} else if t == JSON_TOKEN_STR {
		fmt.Printf("DEBUG_JSON_STR\n")
		var value string
		var success bool
		value, success = TokenStr(file)
		if success {
			fmt.Printf("%s\n", value)
		} else {
			utils.PanicIf(false, "failed to parse string value")
		}
	} else {
		utils.PanicIf(false, fmt.Sprintf("invalid token type : %d", t))
	}
}
