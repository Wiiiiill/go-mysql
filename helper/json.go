package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ReadJSON read json to data
func ReadJSON(v interface{}, filename string) error {

	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	return nil
}

// WriteJSON write data to json
func WriteJSON(v interface{}, filename string, force bool) error {

	if force || !FileExist(filename) {

		data, err := json.MarshalIndent(v, "", "  ")

		if err != nil {
			return err
		}

		dir := filepath.Dir(filename)

		if !DirExist(dir) {
			err = os.MkdirAll(dir, 0700)

			if err != nil {
				return err
			}
		}

		return ioutil.WriteFile(filename, data, 0600)
	}

	return os.ErrExist
}

func SliceToString(sl []string) string {
	var result string

	if len(sl) > 0 {
		for i := 0; i < len(sl)-1; i++ {
			result += sl[i]
			result += ","
		}

		result += sl[len(sl)-1]

	}

	return result
}

func SliceUnit64ToString(sl []uint64) string {
	var result string

	if len(sl) > 0 {
		for i := 0; i < len(sl)-1; i++ {

			result += strconv.FormatUint(sl[i], 10)
			result += ","
		}

		result += strconv.FormatUint(sl[len(sl)-1], 10)

	}

	return result
}

func StringToUnit64Arr(s string) []uint64 {

	if s == "" {
		return []uint64{}
	}

	uintArr := []uint64{}
	sArr := strings.Split(s, ",")
	for i := 0; i < len(sArr); i++ {
		u, err := strconv.ParseUint(sArr[i], 10, 64)
		if err != nil {
			fmt.Printf("字符串转uint64有误：%v \n", err)
			return nil
		}
		uintArr = append(uintArr, u)
	}

	return uintArr
}
