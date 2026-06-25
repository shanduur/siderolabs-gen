// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package xyaml_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/siderolabs/gen/xyaml"
)

type A struct {
	Field string            `yaml:"field"`
	Map   map[string]string `yaml:"map"`
	Slice []A               `yaml:"slice"`
}

// argValue is a custom unmarshaler accepting either a scalar or a sequence.
type argValue struct {
	str  string
	list []string
}

func (a *argValue) UnmarshalYAML(unmarshal func(any) error) error {
	var s string
	if err := unmarshal(&s); err == nil {
		a.str = s

		return nil
	}

	return unmarshal(&a.list)
}

type withArgs struct {
	Args map[string]argValue `yaml:"args"`
}

//go:embed testdata/valid.yaml
var valid []byte

//go:embed testdata/invalid.yaml
var invalid []byte

//go:embed testdata/invalid-nested.yaml
var invalidNested []byte

func TestUnmarshalStrict(t *testing.T) {
	for _, tt := range []struct {
		name string
		err  string
		data []byte
	}{
		{
			name: "valid",
			data: valid,
		},
		{
			name: "invalid",
			data: invalid,
			err:  "unknown keys",
		},
		{
			name: "invalid nested",
			data: invalidNested,
			err:  "this",
		},
		{
			name: "empty",
			data: []byte{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var a A

			err := xyaml.UnmarshalStrict(tt.data, &a)

			if tt.err != "" {
				require.ErrorContains(t, err, tt.err)

				return
			}

			if len(tt.data) != 0 {
				require.NotEmpty(t, a)
			}

			require.NoError(t, err)
		})
	}
}

// TestUnmarshalStrictCustomUnmarshaler checks that a node shape accepted by a custom unmarshaler is not rejected.
func TestUnmarshalStrictCustomUnmarshaler(t *testing.T) {
	for _, tt := range []struct {
		name string
		data string
	}{
		{
			name: "scalar arg value",
			data: "args:\n  issuer: https://one\n",
		},
		{
			name: "sequence arg value",
			data: "args:\n  issuer:\n    - https://one\n    - https://two\n",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var w withArgs

			require.NoError(t, xyaml.UnmarshalStrict([]byte(tt.data), &w))
		})
	}
}
