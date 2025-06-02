// SPDX-FileCopyrightText: 2025 caixw
//
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"io"
)

var frameworks = map[string]framework{
	"solid": &solid{},
}

// 每一种前端框架需要实现的接口
type framework interface {
	// 框架名称
	name() string

	// 输出符合组件属性参数的方法
	//
	// 返回可导出的属性列表，这此属性会在 index.ts 中导出以及在图标集中导入。
	writeProps(io.Writer) ([]string, error)

	// 输出单个图标的组件表示形式
	writeIcon(io.Writer, *Set, string, *Icon) error
}

//--------------------------- solid ----------------------------------

type solid struct{}

func (*solid) name() string { return "solid" }

func (*solid) writeProps(w io.Writer) ([]string, error) {
	props := "Props"
	_, err := fmt.Fprintf(w, `export type %s = VoidProps<{
	height?: string;
	width?: string;
	colors?: Map<string, string>;
}>;

`, props)

	return []string{props}, err
}

func (*solid) writeIcon(w io.Writer, s *Set, name string, icon *Icon) error {
	_, err := io.WriteString(w, "export function "+toCamel(name)+"(props: Props): JSX.Element {\n	return ")
	if err != nil {
		return err
	}

	if err := s.write(w, icon); err != nil {
		return err
	}

	_, err = io.WriteString(w, "}\n\n")
	return err
}
