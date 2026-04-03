package util

import (
	"file/internal/common"
	"file/internal/conf"
	"slices"
	"strings"
)

const (
	_mimeTypeOctetStream    = "application/octet-stream"
	_mimeTypeApplicationZip = "application/zip"
	_mimeTypeVideoMp4       = "video/mp4"
	_mimeTypeMpeg           = "audio/mpeg"
)

var ignoredMimeTypes = []string{
	_mimeTypeOctetStream,
	_mimeTypeApplicationZip,
	_mimeTypeVideoMp4,
	_mimeTypeMpeg,
}

// AutoDetectFileTypeByExt detects the file type based on its extension.
func AutoDetectFileTypeByExt(cfg *conf.Data_Extension, ext string) string {
	switch {
	case slices.Contains(cfg.Photo.Exts, ext):
		return common.FileGroupPhoto
	case slices.Contains(cfg.Video.Exts, ext):
		return common.FileGroupVideo
	case slices.Contains(cfg.Audio.Exts, ext):
		return common.FileGroupAudio
	case slices.Contains(cfg.Compress.Exts, ext):
		return common.FileGroupCompress
	case slices.Contains(cfg.Document.Exts, ext):
		return common.FileGroupDocument
	default:
		return common.FileGroupOther
	}
}

func AutoDetectFileType(cfg *conf.Data_Extension, mimeType, ext string) string {
	mimeType = normalizeMimeType(mimeType)
	// 1. Check mime type
	if mimeType != "" && !slices.Contains(ignoredMimeTypes, mimeType) {
		switch {
		case matchMime(mimeType, cfg.Photo.MimeTypes):
			return common.FileGroupPhoto
		case matchMime(mimeType, cfg.Video.MimeTypes):
			return common.FileGroupVideo
		case matchMime(mimeType, cfg.Audio.MimeTypes):
			return common.FileGroupAudio
		case matchMime(mimeType, cfg.Compress.MimeTypes):
			return common.FileGroupCompress
		case matchMime(mimeType, cfg.Document.MimeTypes):
			return common.FileGroupDocument
		}
	}
	// fallback EXT
	switch {
	case matchExt(ext, cfg.Photo.Exts):
		return common.FileGroupPhoto
	case matchExt(ext, cfg.Video.Exts):
		return common.FileGroupVideo
	case matchExt(ext, cfg.Audio.Exts):
		return common.FileGroupAudio
	case matchExt(ext, cfg.Compress.Exts):
		return common.FileGroupCompress
	case matchExt(ext, cfg.Document.Exts):
		return common.FileGroupDocument
	}
	return common.FileGroupOther
}

func matchMime(mimeType string, patterns []string) bool {
	for _, p := range patterns {
		if strings.HasSuffix(p, "/*") {
			prefix := strings.TrimSuffix(p, "*")
			if strings.HasPrefix(mimeType, prefix) {
				return true
			}
		} else if mimeType == p {
			return true
		}
	}
	return false
}

func matchExt(ext string, exts []string) bool {
	ext = strings.ToLower(ext)
	return slices.Contains(exts, ext)
}

func normalizeMimeType(mimeType string) string {
	mimeType = strings.ToLower(strings.TrimSpace(mimeType))
	// ignore ; charset=utf-8 if exists
	if idx := strings.Index(mimeType, ";"); idx >= 0 {
		mimeType = mimeType[:idx]
	}
	return mimeType
}

// func GenerateUniqueNameV2(
// 	fileRepo file.Repository,
// 	folderID *uuid.UUID,
// 	ownerID string,
// 	originalName string,
// ) (string, error) {
// 	ext := path.Ext(originalName)
// 	base := strings.TrimSuffix(originalName, ext)

// 	names, err := fileRepo.ListCandidateNames(folderID, ownerID, base, ext)
// 	if err != nil {
// 		return "", err
// 	}

// 	used := map[int]bool{}
// 	hasOriginal := false

// 	for _, name := range names {
// 		if name == originalName {
// 			hasOriginal = true
// 			continue
// 		}

// 		nExt := path.Ext(name)
// 		nBase := strings.TrimSuffix(name, nExt)
// 		if nExt != ext {
// 			continue
// 		}

// 		m := suffixRegex.FindStringSubmatch(nBase)
// 		if len(m) != 3 {
// 			continue
// 		}
// 		if m[1] != base {
// 			continue
// 		}

// 		num, err := strconv.Atoi(m[2])
// 		if err != nil {
// 			continue
// 		}
// 		used[num] = true
// 	}

// 	if !hasOriginal {
// 		return originalName, nil
// 	}

// 	for i := 1; ; i++ {
// 		if !used[i] {
// 			return fmt.Sprintf("%s(%d)%s", base, i, ext), nil
// 		}
// 	}
// }
