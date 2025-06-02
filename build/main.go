// SPDX-FileCopyrightText: 2025 caixw
//
// SPDX-License-Identifier: MIT

package main

import (
	"archive/zip"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

const url = "https://github.com/iconify/icon-sets/archive/refs/tags/%s.zip"

func main() {
	id := flag.String("iconset", "", "指定图标集，为空表示所有。")
	ver := flag.String("ver", "", "指定图标集的版本，不能为空。")
	fx := flag.String("fx", "", "指定适用的框架，不能为空。")
	flag.Parse()

	f, found := frameworks[*fx]
	if !found {
		panic(fmt.Sprintf("不支持参数 fx 指定的框架：%s\n", *fx))
	}

	z, err := download(*ver)
	if err != nil {
		panic(err)
	}

	pkg, err := newPkg(f, z, *ver)
	if err != nil {
		panic(err)
	}
	defer pkg.close()

	if err = pkg.createIconSets(*id); err != nil {
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
