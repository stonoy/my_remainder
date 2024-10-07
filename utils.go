package main

import (
	"fmt"
	"time"
)

// func parseID(idStr string) (uuid.UUID, error) {
// 	id, err := uuid.Parse(idStr)
// 	return id, err
// }

func strToTime(timeStr string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-01T15:04",
	}

	for _, layout := range layouts {
		theTime, err := time.Parse(layout, timeStr)

		if err == nil {
			return theTime, nil
		}
	}

	return time.Time{}, fmt.Errorf("can not parse the given time")
}
