package pkg

import (
	"testing"
)

func TestParseCommandLine(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		expected    []string
		expectError bool
	}{
		{
			name:     "simple command 1",
			command:  `ssh -o StrictHostKeyChecking=no ${SERVER_OWNER}@${SERVER_IP} 'cd ${PROJECT_PATH} && pwd && sudo -u www -H bash -c "git reset --hard && git pull" && php think swoole reload'`,
			expected: []string{"ssh", "-o", "StrictHostKeyChecking=no", "${SERVER_OWNER}@${SERVER_IP}", `cd ${PROJECT_PATH} && pwd && sudo -u www -H bash -c "git reset --hard && git pull" && php think swoole reload`},
		},
		{
			name:    "simple command 2",
			command: `bash -c 'rsync -rtv --exclude .git --delete  -e "ssh -o StrictHostKeyChecking=no -p ${SERVER_PORT} -i ~/.ssh/id_rsa" --rsync-path="mkdir -p ${PROJECT_SYMLINK_PATH} && rsync" repository/project_${PROJECT_ID}/dist/build/h5/   ${SERVER_OWNER}@${SERVER_IP}:${PROJECT_SYMLINK_PATH} && ssh -o StrictHostKeyChecking=no -p ${SERVER_PORT} -i ~/.ssh/id_rsa ${SERVER_OWNER}@${SERVER_IP} "ln -sfn ${PROJECT_SYMLINK_PATH} ${PROJECT_PATH}"'`,
			expected: []string{
				"bash",
				"-c",
				`rsync -rtv --exclude .git --delete  -e "ssh -o StrictHostKeyChecking=no -p ${SERVER_PORT} -i ~/.ssh/id_rsa" --rsync-path="mkdir -p ${PROJECT_SYMLINK_PATH} && rsync" repository/project_${PROJECT_ID}/dist/build/h5/   ${SERVER_OWNER}@${SERVER_IP}:${PROJECT_SYMLINK_PATH} && ssh -o StrictHostKeyChecking=no -p ${SERVER_PORT} -i ~/.ssh/id_rsa ${SERVER_OWNER}@${SERVER_IP} "ln -sfn ${PROJECT_SYMLINK_PATH} ${PROJECT_PATH}"`,
			},
		},
		{
			name:     "command with double quotes",
			command:  `echo "hello world"`,
			expected: []string{"echo", "hello world"},
		},
		{
			name:     "command with single quotes",
			command:  `echo 'hello world'`,
			expected: []string{"echo", "hello world"},
		},
		{
			name:     "command with escaped space",
			command:  `echo hello\ world`,
			expected: []string{"echo", "hello world"},
		},
		{
			name:     "command with escaped quote in single quotes",
			command:  `echo 'hello'"'"'world'"'"`,
			expected: []string{"echo", `hello'world'`},
		},
		{
			name:     "empty command",
			command:  "",
			expected: []string{},
		},
		{
			name:     "command with only spaces",
			command:  "   ",
			expected: []string{},
		},
		{
			name:        "unclosed double quote",
			command:     `echo "hello world`,
			expectError: true,
		},
		{
			name:        "unclosed single quote",
			command:     `echo 'hello world`,
			expectError: true,
		},
		{
			name:        "dangling escape",
			command:     `echo hello\`,
			expectError: true,
		},
		{
			name:     "command with equals sign",
			command:  `key=value command`,
			expected: []string{"key=value", "command"},
		},
		{
			name:     "complex command with multiple escapes and quotes",
			command:  `cp "file with spaces.txt" 'another\ file.txt'`,
			expected: []string{"cp", "file with spaces.txt", `another\ file.txt`},
		},
		{
			name:     "command starting with escape character",
			command:  `\echo hello`,
			expected: []string{"echo", "hello"},
		},

		// ----------------------------
		// New cases
		// ----------------------------
		{
			name:     "windows path inside quotes",
			command:  `cmd "C:\\Program Files\\My App\\run.exe" /silent`,
			expected: []string{"cmd", `C:\\Program Files\\My App\\run.exe`, "/silent"},
		},
		{
			name:     "windows path outside quotes",
			command:  `cmd C:\\Windows\\System32\\calc.exe`,
			expected: []string{"cmd", `C:\Windows\System32\calc.exe`},
		},
		{
			name:     "unix path with spaces",
			command:  `cp /path/to/"my file.txt" ./`,
			expected: []string{"cp", "/path/to/my file.txt", "./"},
		},
		{
			name:     "nested quotes mixed",
			command:  `echo "a'b'c"`,
			expected: []string{"echo", `a'b'c`},
		},
		{
			name:     "nested double inside single",
			command:  `echo 'a"b"c'`,
			expected: []string{"echo", `a"b"c`},
		},
		{
			name:     "interleaving quoted and unquoted",
			command:  `echo abc"def"ghi`,
			expected: []string{"echo", "abcdefghi"},
		},
		{
			name:     "escaped backslash",
			command:  `echo hello\\world`,
			expected: []string{"echo", `hello\world`},
		},
		{
			name:     "escaped backslash inside quotes",
			command:  `echo "hello\\world"`,
			expected: []string{"echo", `hello\\world`},
		},
		{
			name:     "multiple consecutive spaces",
			command:  `echo  a   b    c`,
			expected: []string{"echo", "a", "b", "c"},
		},
		{
			name:     "key value with spaces in value",
			command:  `--env="production mode" --debug=false`,
			expected: []string{"--env=production mode", "--debug=false"},
		},
		{
			name:     "key value separated",
			command:  `--env production`,
			expected: []string{"--env", "production"},
		},
		{
			name:     "emoji inside quotes",
			command:  `echo "hello 🌏"`,
			expected: []string{"echo", "hello 🌏"},
		},
		{
			name:     "unicode",
			command:  `echo "你好 世界"`,
			expected: []string{"echo", "你好 世界"},
		},
		{
			name:        "complex nested with escapes",
			command:     `echo "a\"b'c" 'd\'e"f'`,
			expectError: true,
		},
		{
			name:     "multiple commands via semicolon",
			command:  `sh -c "echo 123; echo 456"`,
			expected: []string{"sh", "-c", "echo 123; echo 456"},
		},
		{
			name:     "escaped semicolon",
			command:  `echo hello\;world`,
			expected: []string{"echo", "hello;world"},
		},

		{
			name:     "single dash arg",
			command:  `cmd - -a`,
			expected: []string{"cmd", "-", "-a"},
		},
		{
			name:     "mixed complex shell command",
			command:  `bash -c "echo "test123" && echo 'ok'"`,
			expected: []string{"bash", "-c", `echo test123 && echo 'ok'`},
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseCommandLine(tt.command)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("result len %d != expected %d", len(result), len(tt.expected))
				return
			}

			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("arg[%d] = %q, expected %q", i, result[i], tt.expected[i])
				}
			}
		})
	}
}
