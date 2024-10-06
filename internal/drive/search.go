package drive

import (
	"regexp"
	"strings"
)

var SearchFunction = struct {
	FormatSearchKeyword func(keyword string) string
}{
	FormatSearchKeyword: formatSearchKeyword,
}

func formatSearchKeyword(keyword string) string {
	if keyword == "" {
		return ""
	}

	re := regexp.MustCompile(`(!=)|['\"=<>/\\:]`)
	keyword = re.ReplaceAllString(keyword, "")

	re = regexp.MustCompile(`[,，|(){}]`)
	keyword = re.ReplaceAllString(keyword, " ")

	return strings.TrimSpace(keyword)
}

var DriveFixedTerms = struct {
	DefaultFileFields string
	GDRootType        struct {
		UserDrive  int
		ShareDrive int
	}
	FolderMimeType string
}{
	DefaultFileFields: "parents,id,name,mimeType,modifiedTime,createdTime,fileExtension,size",
	GDRootType: struct {
		UserDrive  int
		ShareDrive int
	}{
		UserDrive:  0,
		ShareDrive: 1,
	},
	FolderMimeType: "application/vnd.google-apps.folder",
}
