package router

import (
	"archive/zip"
)

// ArchiveRouter can serve content from an archive
type ArchiveRouter struct {
	archive   *zip.ReadCloser
	filepaths map[string]bool
}

// NewArchiveRouter creates a router that can serve static content
func NewArchiveRouter(archivePath string) (*ArchiveRouter, error) {
	closer, err := zip.OpenReader(archivePath)
	if err != nil {
		return nil, err
	}
	router := ArchiveRouter{archive: closer}
	return &router, nil
}
