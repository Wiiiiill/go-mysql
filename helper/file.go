package helper

import (
	"os"
	"path/filepath"
	"strings"
)

// FileExist check file
func FileExist(filename string) bool {
	info, err := os.Stat(filename)

	if err != nil && os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

// DirExist check dir
func DirExist(dir string) bool {
	info, err := os.Stat(dir)

	if err != nil && os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

// HomeDir join path with app home dir
func HomeDir(path string) (string, error) {
	if !filepath.IsAbs(path) {
		dir, err := os.Getwd()

		if err != nil {
			return "", err
		}

		path = filepath.Join(dir, path)
	}

	return path, nil
}

const (
	_uploadFolderName = "upload"
	_folderName       = "images"
	_pathSpaceMark    = "/"
)

// GetImageFolderPath 保存图片的文件夹路径
func GetImageFolderPath() string {
	return AppendPath(GetUploadFolderPath(), _folderName)
}

// GetUploadFolderPath 保存上传文件的文件夹路径
func GetUploadFolderPath() string {
	dir, _ := os.Getwd()
	return AppendPath(dir, _uploadFolderName)
}

// pathExists 判断路径是否存在
func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// appendPath 拼接路径
func AppendPath(rootPath string, paths ...string) string {

	builder := strings.Builder{}
	builder.WriteString(rootPath)

	for _, path := range paths {
		builder.WriteString(_pathSpaceMark)
		builder.WriteString(path)
	}

	fullPath := builder.String()
	fullPath = strings.Replace(fullPath, "//", "/", -1)
	return fullPath
}
