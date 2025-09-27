package fileserverhandler

import (
	"fmt"
	"math"
)

func getFileSizeInString(size int) string {
	unit := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	index := 0

	fileSize := float64(size)

	for fileSize > 1024 && index < len(unit) {
		fileSize = fileSize / 1024
		index += 1
	}

	if index == len(unit) {
		return "Too large to show"
	} else {
		return fmt.Sprintf("%.2f", math.Round(fileSize*100) / 100) + unit[index]
	}

}