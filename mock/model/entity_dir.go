package model

type DirResponse struct {
	Path string `json:"path"`

	Dirs []DirParam `json:"dirs"`
}

type DirParam struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
	Size  int64  `json:"size"`
}

type DirInfoResponse struct {
	Path      string `json:"path"`
	DirCount  int64  `json:"dirCount"`
	FileCount int64  `json:"fileCount"`
	TotalSize int64  `json:"totalSize"`
}
