package models

import "time"

func Now() int64 {
	return time.Now().UTC().Unix()
}
