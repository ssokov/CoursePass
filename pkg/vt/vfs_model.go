//nolint:dupl
package vt

import (
	"time"

	"courses/pkg/db"
)

type VfsFile struct {
	ID         int       `json:"id"`
	FolderID   int       `json:"folderId" validate:"required"`
	Title      string    `json:"title" validate:"required,max=255"`
	Path       string    `json:"path" validate:"required,max=255"`
	Params     *string   `json:"params"`
	IsFavorite *bool     `json:"isFavorite"`
	MimeType   string    `json:"mimeType" validate:"required,max=255"`
	FileSize   *int      `json:"fileSize"`
	FileExists bool      `json:"fileExists" validate:"required"`
	CreatedAt  time.Time `json:"createdAt"`
	StatusID   int       `json:"statusId" validate:"required,status"`

	Folder *VfsFolderSummary `json:"folder"`
	Status *Status           `json:"status"`
}

func (vf *VfsFile) ToDB() *db.VfsFile {
	if vf == nil {
		return nil
	}

	vfsFile := &db.VfsFile{
		ID:         vf.ID,
		FolderID:   vf.FolderID,
		Title:      vf.Title,
		Path:       vf.Path,
		Params:     vf.Params,
		IsFavorite: vf.IsFavorite,
		MimeType:   vf.MimeType,
		FileSize:   vf.FileSize,
		FileExists: vf.FileExists,
		CreatedAt:  vf.CreatedAt,
		StatusID:   vf.StatusID,
	}

	return vfsFile
}

type VfsFileSearch struct {
	ID         *int       `json:"id"`
	FolderID   *int       `json:"folderId"`
	Title      *string    `json:"title"`
	Path       *string    `json:"path"`
	Params     *string    `json:"params"`
	IsFavorite *bool      `json:"isFavorite"`
	MimeType   *string    `json:"mimeType"`
	FileSize   *int       `json:"fileSize"`
	FileExists *bool      `json:"fileExists"`
	CreatedAt  *time.Time `json:"createdAt"`
	StatusID   *int       `json:"statusId"`
	IDs        []int      `json:"ids"`
}

func (vfs *VfsFileSearch) ToDB() *db.VfsFileSearch {
	if vfs == nil {
		return nil
	}

	return &db.VfsFileSearch{
		ID:            vfs.ID,
		FolderID:      vfs.FolderID,
		TitleILike:    vfs.Title,
		PathILike:     vfs.Path,
		ParamsILike:   vfs.Params,
		IsFavorite:    vfs.IsFavorite,
		MimeTypeILike: vfs.MimeType,
		FileSize:      vfs.FileSize,
		FileExists:    vfs.FileExists,
		CreatedAt:     vfs.CreatedAt,
		StatusID:      vfs.StatusID,
		IDs:           vfs.IDs,
	}
}

type VfsFileSummary struct {
	ID         int       `json:"id"`
	FolderID   int       `json:"folderId"`
	Title      string    `json:"title"`
	Path       string    `json:"path"`
	Params     *string   `json:"params"`
	IsFavorite *bool     `json:"isFavorite"`
	MimeType   string    `json:"mimeType"`
	FileSize   *int      `json:"fileSize"`
	FileExists bool      `json:"fileExists"`
	CreatedAt  time.Time `json:"createdAt"`

	Folder *VfsFolderSummary `json:"folder"`
	Status *Status           `json:"status"`
}

type VfsFolder struct {
	ID             int       `json:"id"`
	ParentFolderID *int      `json:"parentFolderId"`
	Title          string    `json:"title" validate:"required,max=255"`
	IsFavorite     *bool     `json:"isFavorite"`
	CreatedAt      time.Time `json:"createdAt"`
	StatusID       int       `json:"statusId" validate:"required,status"`

	ParentFolder *VfsFolderSummary `json:"parentFolder"`
	Status       *Status           `json:"status"`
}

func (vf *VfsFolder) ToDB() *db.VfsFolder {
	if vf == nil {
		return nil
	}

	vfsFolder := &db.VfsFolder{
		ID:             vf.ID,
		ParentFolderID: vf.ParentFolderID,
		Title:          vf.Title,
		IsFavorite:     vf.IsFavorite,
		CreatedAt:      vf.CreatedAt,
		StatusID:       vf.StatusID,
	}

	return vfsFolder
}

type VfsFolderSearch struct {
	ID             *int       `json:"id"`
	ParentFolderID *int       `json:"parentFolderId"`
	Title          *string    `json:"title"`
	IsFavorite     *bool      `json:"isFavorite"`
	CreatedAt      *time.Time `json:"createdAt"`
	StatusID       *int       `json:"statusId"`
	IDs            []int      `json:"ids"`
}

func (vfs *VfsFolderSearch) ToDB() *db.VfsFolderSearch {
	if vfs == nil {
		return nil
	}

	return &db.VfsFolderSearch{
		ID:             vfs.ID,
		ParentFolderID: vfs.ParentFolderID,
		TitleILike:     vfs.Title,
		IsFavorite:     vfs.IsFavorite,
		CreatedAt:      vfs.CreatedAt,
		StatusID:       vfs.StatusID,
		IDs:            vfs.IDs,
	}
}

type VfsFolderSummary struct {
	ID             int       `json:"id"`
	ParentFolderID *int      `json:"parentFolderId"`
	Title          string    `json:"title"`
	IsFavorite     *bool     `json:"isFavorite"`
	CreatedAt      time.Time `json:"createdAt"`

	ParentFolder *VfsFolderSummary `json:"parentFolder"`
	Status       *Status           `json:"status"`
}
