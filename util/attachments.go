package util

import "log"

var (
	imageTypes = []string{"image/png", "image/jpeg", "image/bmp", "image/webp", "image/tiff"}
	audioTypes = []string{"audio/mpeg", "audio/wav", "audio/ogg", "audio/flac", "audio/aac", "audio/opus"}
	videoTypes = []string{"video/mp4", "video/webm", "video/mov"}
)

/*
Check to make sure content type is valid
*/
func CheckContentType(c string, l []string) bool {
	for _, item := range l {
		if item == c {
			log.Println("[INFO] Found valid ContentType")
			return true
		}
	}
	log.Println("[ERROR] Invalid ContentType")
	return false
}
