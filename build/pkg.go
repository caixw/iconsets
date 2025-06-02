// SPDX-FileCopyrightText: 2025 caixw
//
// SPDX-License-Identifier: MIT

package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 前端项目的根目录
const root = "../packages"

const fileHeader = "// 当前文件由工具自动生成，如无必要请勿手动修改！\n\n"

const propsFile = "_props" // 尽量避免与其它图标集重名。

type pkg struct {
	fx     framework
	outDir string // 输出的根目录
	zip    *zip.ReadCloser
	ver    string
	index  *os.File
	props  string
}

func newPkg(fx framework, z *zip.ReadCloser, ver string) (*pkg, error) {
	outDir := filepath.Join(root, fx.name(), "src")
	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		return nil, err
	}

	fmt.Printf("创建 %s\n", filepath.Join(fx.name(), "src", "index.ts"))
	index, err := os.Create(filepath.Join(outDir, "index.ts"))
	if err != nil {
		return nil, err
	}
	if _, err = index.WriteString(fileHeader); err != nil {
		return nil, err
	}

	props, err := createProps(outDir, fx)
	if err != nil {
		return nil, err
	}

	var prop string
	if len(props) > 0 {
		prop := strings.Join(props, ", ")
		if _, err = fmt.Fprintf(index, "export type { %s } from './%s';\n\n", prop, propsFile); err != nil {
			return nil, err
		}
	}

	return &pkg{
		outDir: outDir,
		fx:     fx,
		zip:    z,
		ver:    ver,
		index:  index,
		props:  prop,
	}, nil
}

func createProps(outDir string, fx framework) ([]string, error) {
	fmt.Printf("创建 %s\n", filepath.Join(fx.name(), "src", propsFile+".ts"))

	f, err := os.Create(filepath.Join(outDir, propsFile+".ts"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err = f.WriteString(fileHeader); err != nil {
		return nil, err
	}

	if _, err = f.WriteString("import { VoidProps } from 'solid-js';\n\n"); err != nil {
		return nil, err
	}

	props, err := fx.writeProps(f)
	if err != nil {
		return nil, err
	}

	return props, nil
}

// 为当前框架创建一图标集
func (p *pkg) createIconSet(iconset string) error {
	fmt.Printf("准备创建图标集 %s\n", iconset)
	defer fmt.Printf("完成创建图标集 %s\n", iconset)

	f, err := p.zip.Open("icon-sets-" + p.ver + "/json/" + iconset + ".json")
	if err != nil {
		return err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	set := &Set{}
	if err := json.Unmarshal(data, set); err != nil {
		return err
	}

	if err := p.genComponents(set, filepath.Join(p.outDir, iconset+".tsx")); err != nil {
		return err
	}

	_, err = fmt.Fprintf(p.index, "export * as %s from './%s'\n", iconset, iconset)
	return err
}

func (p *pkg) genComponents(set *Set, root string) error {
	f, err := os.Create(root)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.WriteString(fileHeader); err != nil {
		return err
	}

	_, err = fmt.Fprintf(f, "// 当前文件是根据一组图标集为框架 %s 生成的组件列表\n// 作者为：%s <%s>\n// 许可证为：%s <%s>\n\n", p.fx.name(), set.Info.Author.Name, set.Info.Author.URL, set.Info.License.SPDX, set.Info.License.URL)
	if err != nil {
		return err
	}

	if _, err = f.WriteString("import { JSX } from 'solid-js';\n\n"); err != nil {
		return err
	}

	if p.props != "" {
		if _, err = fmt.Fprintf(f, "import { %s } from './%s';\n\n", p.props, propsFile); err != nil {
			return err
		}
	}

	for name, i := range set.Icons {
		if err := p.fx.writeIcon(f, name, i); err != nil {
			return err
		}
	}

	return nil
}

func (p *pkg) close() error {
	if _, err := p.index.WriteString("\n"); err != nil {
		return err
	}

	if err := p.index.Close(); err != nil {
		return err
	}

	return nil
}
