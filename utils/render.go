package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func RenderTemplate(src, dest string, data interface{}) error {
	srcName := filepath.Base(src)
	tpl, err := template.New(srcName).
		Funcs(template.FuncMap{
			"split":        split,
			"join":         join,
			"trim":         trim,
			"keyOrDefault": keyOrDefault,
		}).ParseFiles(src)

	if err != nil {
		return err
	}

	parentDir := filepath.Dir(dest)
	tmpFile, err := ioutil.TempFile(parentDir, fmt.Sprint(".%s.tmp-", srcName))
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	err = tpl.Execute(tmpFile, data)
	if err != nil {
		return err
	}

	err = os.Rename(tmpFile.Name(), dest)
	if err != nil {
		return err
	}

	return nil
}

func keyOrDefault(needle string, fallback string, haystack map[string]string) string {
	if value, ok := haystack[needle]; ok {
		return value
	}

	return fallback
}

// Revert arguments so this could be used in pipeline ex. "foo,bar" | split "," | join " "
func split(sep, s string) []string {
	return strings.Split(s, sep)
}

func join(sep string, a []string) string {
	return strings.Join(a, sep)
}

func trim(cutset string, s string) string {
	return strings.Trim(s, cutset)
}
