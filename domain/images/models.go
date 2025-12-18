package images

import "time"

type ImageData struct {
	fileName    string
	filePath    string
	cameraModel string
	timestamp   time.Time
	DestPath    string
}

var imageFileTypes = []string{"jpg", "jpeg", "raw", "cr3", "cr2"}
