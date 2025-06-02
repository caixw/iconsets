// SPDX-FileCopyrightText: 2025 caixw
//
// SPDX-License-Identifier: MIT

package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const url = "https://github.com/iconify/icon-sets/archive/refs/tags/%s.zip"

func main() {
	id := "bytesize"
	ver := "2.2.342"

	z, err := download(ver)
	if err != nil {
		panic(err)
	}

	pkg, err := newPkg(&solid{}, z, ver)
	if err != nil {
		panic(err)
	}
	defer pkg.close()

	if err = pkg.createIconSet(id); err != nil {
		panic(err)
	}
}

func download(ver string) (*zip.ReadCloser, error) {
	file := "./download/icon-sets-" + ver + ".zip"

	if _, err := os.Stat(file); err != nil && !errors.Is(err, os.ErrExist) {
		u := fmt.Sprintf(url, ver)
		fmt.Printf("下载文件: %s\n", u)

		resp, err := http.Get(u)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if err = os.WriteFile(file, buf, os.ModePerm); err != nil {
			return nil, err
		}

		fmt.Println("下载完成")
	} else {
		fmt.Println("文件已经下载")
	}

	return zip.OpenReader(file)
}
