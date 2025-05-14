package utils

import "time"

type ShortListFileDropbox struct {
	Name       string
	Path       string
	Revision   string
	Size       int
	Modified   string
	SortByThis time.Time
}
type ShortListFileLocal struct {
	Size     int64
	FileName string
	Dir      string
	Modified time.Time
}
