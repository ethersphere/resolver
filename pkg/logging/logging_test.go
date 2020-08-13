// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package logging_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/ethersphere/resolver/pkg/logging"
)

var testLevelStrings = []string{
	"silent",
	"error",
	"warn",
	"info",
	"debug",
	"trace",
}

func TestWithOutput(t *testing.T) {
	l := logging.New(
		logging.WithOutput(ioutil.Discard),
	)

	impl := logging.GetImpl(l)
	want := ioutil.Discard
	got := impl.Out
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestWithLevel(t *testing.T) {
	t.Run("info", func(t *testing.T) {
		l := logging.New(
			logging.WithLevel(logging.InfoLevel),
		)

		want := logging.Level(logging.InfoLevel)
		got := l.Level
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("silent", func(t *testing.T) {
		l := logging.New(
			logging.WithLevel(logging.SilentLevel),
		)

		wantLvl := logging.Level(logging.SilentLevel)
		gotLvl := l.Level
		if gotLvl != wantLvl {
			t.Errorf("got %v, want %v", gotLvl, wantLvl)
		}

		impl := logging.GetImpl(l)
		wantOut := ioutil.Discard
		gotOut := impl.Out
		if gotOut != wantOut {
			t.Errorf("got %v, want %v", gotOut, wantOut)

		}
	})
}

func TestLevelToString(t *testing.T) {
	for idx, lvl := range logging.AllLevels {
		got := lvl.String()
		want := testLevelStrings[idx]
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	}

	got := (logging.Level)(99).String()
	want := "unknown"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}

}

func TestLevelUnmarshalText(t *testing.T) {
	var u logging.Level
	for _, level := range logging.AllLevels {
		if err := u.UnmarshalText([]byte(level.String())); err != nil {
			t.Fatal(err)
		}
		if u != level {
			t.Errorf("want %d, got %d", u, level)
		}
	}

	badVal := "bogus"
	wantErr := fmt.Errorf("not a valid Level: %q", badVal)
	gotErr := u.UnmarshalText([]byte(badVal))
	if wantErr.Error() != gotErr.Error() {
		t.Errorf("invalid error: got %v, want %v", gotErr, wantErr)
	}
}
