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
			"split":   strings.Split,
			"join":    strings.Join,
			"trim":    strings.Trim,
			"replace": strings.Replace,
		}).ParseFiles(src)

	if err != nil {
		return err
	}

	parentDir := filepath.Dir(dest)
	tmpFile, err := ioutil.TempFile(parentDir, fmt.Sprint(".%s.tmp-", srcName))
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
