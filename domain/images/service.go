package images

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/exif2"
	"go.uber.org/zap"
)

func GetPhoto(_ *zap.Logger, path string) (ImageData, error) {
	var i ImageData
	f, err := os.Open(path)
	if err != nil {
		return i, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	e, err := imagemeta.Decode(f)
	if err != nil {
		return i, fmt.Errorf("failed to decode image: %w", err)
	}

	sepPath := strings.Split(path, "/")
	return toImageData(e, sepPath[len(sepPath)-1], path), nil
}

func toImageData(e exif2.Exif, name, path string) ImageData {
	return ImageData{
		fileName:    name,
		filePath:    path,
		cameraModel: e.Model,
		timestamp:   e.DateTimeOriginal(),
	}
}

func GetImageTypes() []string {
	return imageFileTypes
}

func (i ImageData) GetFileName() string {
	return i.fileName
}

func (i ImageData) GetFilePath() string {
	return i.filePath
}

func (i ImageData) GetCameraModel() string {
	return strings.ToLower(i.cameraModel)
}

func (i ImageData) GetTimestamp() time.Time {
	return i.timestamp
}
