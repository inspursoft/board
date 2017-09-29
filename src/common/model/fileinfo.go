package model

type FileInfo struct {
	Path     string `json:"path"`
	FileName string `json:"file_name"`
	Size     int64  `json:"size"`
}
