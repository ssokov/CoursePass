package vt

import (
	"courses/pkg/db"
)

func NewVfsFile(in *db.VfsFile) *VfsFile {
	if in == nil {
		return nil
	}

	vfsFile := &VfsFile{
		ID:         in.ID,
		FolderID:   in.FolderID,
		Title:      in.Title,
		Path:       in.Path,
		Params:     in.Params,
		IsFavorite: in.IsFavorite,
		MimeType:   in.MimeType,
		FileSize:   in.FileSize,
		FileExists: in.FileExists,
		CreatedAt:  in.CreatedAt,
		StatusID:   in.StatusID,

		Folder: NewVfsFolderSummary(in.Folder),
		Status: NewStatus(in.StatusID),
	}

	return vfsFile
}

func NewVfsFileSummary(in *db.VfsFile) *VfsFileSummary {
	if in == nil {
		return nil
	}

	return &VfsFileSummary{
		ID:         in.ID,
		FolderID:   in.FolderID,
		Title:      in.Title,
		Path:       in.Path,
		Params:     in.Params,
		IsFavorite: in.IsFavorite,
		MimeType:   in.MimeType,
		FileSize:   in.FileSize,
		FileExists: in.FileExists,
		CreatedAt:  in.CreatedAt,

		Folder: NewVfsFolderSummary(in.Folder),
		Status: NewStatus(in.StatusID),
	}
}

func NewVfsFolder(in *db.VfsFolder) *VfsFolder {
	if in == nil {
		return nil
	}

	vfsFolder := &VfsFolder{
		ID:             in.ID,
		ParentFolderID: in.ParentFolderID,
		Title:          in.Title,
		IsFavorite:     in.IsFavorite,
		CreatedAt:      in.CreatedAt,
		StatusID:       in.StatusID,

		ParentFolder: NewVfsFolderSummary(in.ParentFolder),
		Status:       NewStatus(in.StatusID),
	}

	return vfsFolder
}

func NewVfsFolderSummary(in *db.VfsFolder) *VfsFolderSummary {
	if in == nil {
		return nil
	}

	return &VfsFolderSummary{
		ID:             in.ID,
		ParentFolderID: in.ParentFolderID,
		Title:          in.Title,
		IsFavorite:     in.IsFavorite,
		CreatedAt:      in.CreatedAt,

		ParentFolder: NewVfsFolderSummary(in.ParentFolder),
		Status:       NewStatus(in.StatusID),
	}
}
