// SPDX-FileCopyrightText: 2025 caixw
//
// SPDX-License-Identifier: MIT

package main

import (
	"bytes"
	"testing"

	"github.com/issue9/assert/v4"
)

func TestToCamel(t *testing.T) {
	a := assert.New(t, false)
	a.Equal(toCamel("abc"), "Abc").
		Equal(toCamel("abc-def"), "AbcDef")
}

func TestIcon_write(t *testing.T) {
	a := assert.New(t, false)
	w := &bytes.Buffer{}

	icon := &Icon{
		Body: "<g />",
	}
	a.NotError(icon.write(w))
	a.Equal(w.String(), `<svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" viewBox="0 0 16 16">
	<g />
</svg>`)

	w.Reset()
	icon = &Icon{
		Body:   "<g />",
		Height: 32,
	}
	a.NotError(icon.write(w))
	a.Equal(w.String(), `<svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" viewBox="0 0 16 32">
	<g />
</svg>`)
}
