package manifestspec

import (
	"errors"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

var (
	ErrInvalidPath      = errors.New("manifest path is invalid")
	ErrDuplicatePath    = errors.New("manifest contains duplicate canonical path")
	ErrSymbolicLink     = errors.New("symbolic links are not supported")
	windowsDrivePattern = regexp.MustCompile(`^[A-Za-z]:`)
)

func CanonicalPath(value string) (string, error) {
	if value == "" || strings.ContainsRune(value, '\x00') || !utf8.ValidString(value) {
		return "", ErrInvalidPath
	}
	value = strings.ReplaceAll(value, `\`, "/")
	if strings.HasPrefix(value, "/") || strings.HasPrefix(value, "//") || windowsDrivePattern.MatchString(value) {
		return "", ErrInvalidPath
	}
	segments := strings.Split(value, "/")
	for _, segment := range segments {
		if segment == "" || segment == "." || segment == ".." {
			return "", ErrInvalidPath
		}
	}
	value = norm.NFC.String(value)
	cleaned := path.Clean(value)
	if cleaned == "." || cleaned != value || strings.HasPrefix(cleaned, "../") {
		return "", ErrInvalidPath
	}
	return cleaned, nil
}

func RelativeCanonicalPath(root, filename string) (string, error) {
	relative, err := filepath.Rel(root, filename)
	if err != nil {
		return "", ErrInvalidPath
	}
	return CanonicalPath(relative)
}
