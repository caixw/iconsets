// SPDX-FileCopyrightText: 2025 caixw
//
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Set struct {
	Prefix string `json:"prefix"`
	Info   *struct {
		Author *struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"author"`
		License *struct {
			SPDX string `json:"spdx"`
			URL  string `json:"url"`
		} `json:"license"`
	} `json:"info"`

	Palette bool `json:"palette"` // 是否不可自定义颜色，true 表示不可自定义。

	Icons   map[string]*Icon  `json:"icons"`
	Aliases map[string]*Alias `json:"aliases"`
}

type Icon struct {
	Body string `json:"body"`

	// viewBox
	Left   int `json:"left,omitempty"`   // 0
	Top    int `json:"top,omitempty"`    // 0
	Width  int `json:"width,omitempty"`  // 16
	Height int `json:"height,omitempty"` // 16

	// transform
	Rotate int  `json:"rotate"` // 0, [0,90]
	HFlip  bool `json:"hFlip"`  // false
	VFlip  bool `json:"vFlip"`  // false
}

type Alias struct {
	Icon
	Parent string `json:"parent"`
}

func toCamel(name string) string {
	words := strings.Split(name, "-")
	for i, w := range words {
		words[i] = string(unicode.ToUpper(rune(w[0]))) + w[1:]
	}
	return strings.Join(words, "")
}

const svgFormat = `<svg xmlns="http://www.w3.org/2000/svg" width={props.width} height={props.height} viewBox="%d %d %d %d">
	%s
</svg>`

// 转换为 svg 图片
func (i *Icon) write(w io.Writer) (err error) {
	if i.Width == 0 {
		i.Width = 16
	}
	if i.Height == 0 {
		i.Height = 16
	}

	transforms := []string{}
	switch { // 翻转属性
	case i.HFlip && i.VFlip:
		transforms = append(transforms, "scale(-1,-1)")
	case i.HFlip:
		transforms = append(transforms, "scale(-1,1)")
	case i.VFlip:
		transforms = append(transforms, "scale(1,-1)")
	}

	switch i.Rotate {
	case 1:
		transforms = append(transforms, "rotate(90)")
	case 2:
		transforms = append(transforms, "rotate(180)")
	case 3:
		transforms = append(transforms, "rotate(270)")
	default:
		return fmt.Errorf("rotate 值 %d 无效", i.Rotate)
	}

	_, err = fmt.Fprintf(w, svgFormat, i.Left, i.Top, i.Width, i.Height, i.Body)
	return err
}
