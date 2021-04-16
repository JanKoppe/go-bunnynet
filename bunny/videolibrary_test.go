// Copyright (c) 2021 Jan Koppe
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"testing"
)

func TestReadVideoLibraries(t *testing.T) {
	c, err := NewClient("")
	if err != nil {
		t.Errorf(err.Error())
	}

	lib, err := c.AddVideoLibrary("go-bunnynet-testread", []string{"SYD"})
	if err != nil {
		t.Errorf(err.Error())
	}

	videoLibraries, err := c.ListVideoLibraries()
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = c.GetVideoLibrary((*videoLibraries)[0].ID)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = c.DeleteVideoLibrary(lib.ID)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestCrudVideoLibrary(t *testing.T) {
	c, err := NewClient("")
	if err != nil {
		t.Errorf(err.Error())
	}

	lib, err := c.AddVideoLibrary("go-bunnynet-test", []string{"SYD"})
	if err != nil {
		t.Errorf(err.Error())
	}

	lib.KeepOriginalFiles = false

	lib, err = c.UpdateVideoLibrary((*lib))
	if err != nil {
		t.Errorf(err.Error())
	}

	err = c.AddVideoLibraryAllowedReferrer(lib.ID, "bunny.net")
	if err != nil {
		t.Errorf(err.Error())
	}

	err = c.RemoveVideoLibraryAllowedReferrer(lib.ID, "bunny.net")
	if err != nil {
		t.Errorf(err.Error())
	}

	err = c.AddVideoLibraryBlockedReferrer(lib.ID, "bunny.net")
	if err != nil {
		t.Errorf(err.Error())
	}

	err = c.RemoveVideoLibraryBlockedReferrer(lib.ID, "bunny.net")
	if err != nil {
		t.Errorf(err.Error())
	}

	err = c.DeleteVideoLibrary(lib.ID)
	if err != nil {
		t.Errorf(err.Error())
	}
}
