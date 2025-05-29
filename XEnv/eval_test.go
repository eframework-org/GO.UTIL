// Copyright (c) 2025 EFramework Organization. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XEnv

import (
	"os"
	"runtime"
	"sync"
	"testing"

	"github.com/eframework-org/GO.UTIL/XString"
	"github.com/stretchr/testify/assert"
)

func TestVarsEval(t *testing.T) {
	originalArgs := os.Args
	originalEnv := os.Getenv("TEST_VAR")

	defer func() {
		os.Args = originalArgs
		os.Setenv("TEST_VAR", originalEnv)
	}()

	tests := []struct {
		name     string
		setup    func()
		input    string
		expected string
	}{
		{
			name: "Command Line Arg",
			setup: func() {
				os.Args = []string{"prog", "--test=value"}
			},
			input:    "prefix ${Env.test} suffix",
			expected: "prefix value suffix",
		},
		{
			name: "Environment Variable",
			setup: func() {
				os.Setenv("TEST_VAR", "env_value")
			},
			input:    "prefix ${Env.TEST_VAR} suffix",
			expected: "prefix env_value suffix",
		},
		{
			name: "Arg Precedence Over Env",
			setup: func() {
				os.Args = []string{"prog", "--TEST_VAR=arg_value"}
				os.Setenv("TEST_VAR", "env_value")
			},
			input:    "${Env.TEST_VAR}",
			expected: "arg_value",
		},
		{
			name: "Multiple Mixed Sources",
			setup: func() {
				os.Args = []string{"prog", "--arg1=val1"}
				os.Setenv("ENV1", "val2")
			},
			input:    "${Env.arg1} and ${Env.ENV1}",
			expected: "val1 and val2",
		},
		{
			name: "Missing Variable",
			setup: func() {
				os.Args = []string{"prog"}
			},
			input:    "hello ${Env.missing}",
			expected: "hello ${Env.missing}(Unknown)",
		},
		{
			name: "Recursive Variables",
			setup: func() {
				os.Args = []string{"prog", "--var1=${Env.var2}", "--var2=${Env.var1}"}
			},
			input:    "${Env.var1}",
			expected: "${Env.var1}(Recursive)",
		},
		{
			name: "Nested Variables",
			setup: func() {
				os.Args = []string{"prog", "--outer=value"}
			},
			input:    "nested ${Env.outer${Env.inner}}",
			expected: "nested ${Env.outer${Env.inner}(Nested)}",
		},
		{
			name: "Empty Value",
			setup: func() {
				os.Args = []string{"prog", "--empty="}
			},
			input:    "test${Env.empty}end",
			expected: "test${Env.empty}(Unknown)end",
		},
		{
			input:    "${Env.LocalPath}",
			expected: LocalPath(),
		},
		{
			input:    "${Env.AssetPath}",
			expected: AssetPath(),
		},
		{
			input:    "${Env.UserName}",
			expected: PrefsAuthorDefault,
		},
		{
			input:    "${Env.Platform}",
			expected: Platform(),
		},
		{
			input:    "${Env.App}",
			expected: App(),
		},
		{
			input:    "${Env.Mode}",
			expected: Mode(),
		},
		{
			input:    "${Env.Solution}",
			expected: Solution(),
		},
		{
			input:    "${Env.Project}",
			expected: Project(),
		},
		{
			input:    "${Env.Product}",
			expected: Product(),
		},
		{
			input:    "${Env.Channel}",
			expected: Channel(),
		},
		{
			input:    "${Env.Version}",
			expected: Version(),
		},
		{
			input:    "${Env.Author}",
			expected: Author(),
		},
		{
			input:    "${Env.Secret}",
			expected: Secret(),
		},
		{
			input:    "${Env.NumCPU}",
			expected: XString.ToString(runtime.NumCPU()),
		},
	}

	for _, tt := range tests {
		name := tt.name
		if name == "" {
			name = tt.input
		}
		t.Run(name, func(t *testing.T) {
			argsCache = nil
			argsCacheInit = sync.Once{}
			if tt.setup != nil {
				tt.setup()
			}

			result := Vars().Eval(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
