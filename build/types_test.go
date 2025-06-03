// SPDX-FileCopyrightText: 2025 caixw
//
// SPDX-License-Identifier: MIT

package main

import (
	"testing"

	"github.com/issue9/assert/v4"
)

func TestToCamel(t *testing.T) {
	a := assert.New(t, false)
	a.Equal(toCamel("abc"), "Abc").
		Equal(toCamel("abc-def"), "AbcDef")
}

func TestIcon_transform(t *testing.T) {
	a := assert.New(t, false)

	i := &Icon{}
	trans, err := i.transform()
	a.NotError(err).Length(trans, 0)

	i = &Icon{HFlip: true}
	trans, err = i.transform()
	a.NotError(err).Equal(trans, []string{"scale(-1,1)"})

	i = &Icon{HFlip: true, VFlip: true}
	trans, err = i.transform()
	a.NotError(err).Equal(trans, []string{"scale(-1,-1)"})

	i = &Icon{VFlip: true}
	trans, err = i.transform()
	a.NotError(err).Equal(trans, []string{"scale(1,-1)"})

	i = &Icon{VFlip: true, Rotate: 1}
	trans, err = i.transform()
	a.NotError(err).Equal(trans, []string{"scale(1,-1)", "rotate(90)"})

	i = &Icon{VFlip: true, Rotate: 2}
	trans, err = i.transform()
	a.NotError(err).Equal(trans, []string{"scale(1,-1)", "rotate(180)"})

	i = &Icon{VFlip: true, Rotate: 3}
	trans, err = i.transform()
	a.NotError(err).Equal(trans, []string{"scale(1,-1)", "rotate(270)"})

	i = &Icon{VFlip: true, Rotate: 4}
	trans, err = i.transform()
	a.ErrorString(err, "rotate 值 4 无效").Length(trans, 0)
}

func TestIcon_size(t *testing.T) {
	a := assert.New(t, false)
	s := &Set{}

	icon := &Icon{}
	w, h := icon.size(s)
	a.Equal(w, 16).Equal(h, 16)

	icon = &Icon{Height: 48}
	w, h = icon.size(s)
	a.Equal(w, 16).Equal(h, 48)

	// 带默认值

	s = &Set{Width: 32, Height: 48}
	icon = &Icon{}
	w, h = icon.size(s)
	a.Equal(w, 32).Equal(h, 48)

	icon = &Icon{Height: 52}
	w, h = icon.size(s)
	a.Equal(w, 32).Equal(h, 52)
}
