// SPDX-FileCopyrightText: 2025 caixw
//
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"io"
	"strings"
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
	presetProps := "presetProps"
	_, err := fmt.Fprintf(w, `export type %s = VoidProps<{
	/**
	 * 图标的高度，默认为 1rem
	 */
	height?: string;

	/**
	 * 图标的宽度，默认为 1rem
	 */
	width?: string;
}>;

export const %s: Props = {
	height: '1rem',
	width: '1rem'
} as const;

`, props, presetProps)

	return []string{props, presetProps}, err
}

const solidComponentString = `export function %s(props: Props): JSX.Element {
	props = mergeProps(presetProps, props);
	return <svg xmlns="http://www.w3.org/2000/svg" width={props.width} height={props.height} viewBox="%g %g %g %g">
		%s
	</svg>;
}

`

func (*solid) writeIcon(w io.Writer, s *Set, name string, icon *Icon) error {
	transforms, err := icon.transform()
	if err != nil {
		return err
	}

	body := icon.Body
	if len(transforms) > 0 {
		body = `<g transform="` + strings.Join(transforms, " ") + `">` + body + "</g>"
	}

	io.WriteString(w, "// "+name+"\n")
	width, height := icon.size(s)
	_, err = fmt.Fprintf(w, solidComponentString, toCamel(s.Prefix)+toCamel(name), icon.Left, icon.Top, width, height, body)
	return err
}
