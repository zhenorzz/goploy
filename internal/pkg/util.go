// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package pkg

import (
	"fmt"
	"regexp"
	"strings"
)

// GetScriptExt return script extension default bash
func GetScriptExt(scriptMode string) string {
	switch scriptMode {
	case "sh", "zsh", "bash":
		return "sh"
	case "php":
		return "php"
	case "python":
		return "py"
	case "cmd":
		return "bat"
	default:
		return "sh"
	}
}

// ParseCommandLine parse cmd arg
func ParseCommandLine(command string) ([]string, error) {
	var args []string
	var current strings.Builder

	inQuotes := false
	quoteChar := byte(0)
	escapeNext := false

	for i := 0; i < len(command); i++ {
		c := command[i]

		if escapeNext {
			current.WriteByte(c)
			escapeNext = false
			continue
		}

		if c == '\\' && !inQuotes {
			escapeNext = true
			continue
		}

		if inQuotes {
			if c == quoteChar {
				inQuotes = false
				quoteChar = 0
			} else {
				current.WriteByte(c)
			}
			continue
		}

		if c == '"' || c == '\'' {
			inQuotes = true
			quoteChar = c
			continue
		}

		if c == ' ' || c == '\t' {
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
			continue
		}

		current.WriteByte(c)
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	if inQuotes {
		return nil, fmt.Errorf("unclosed quote in command line: %s", command)
	}

	// 检查未处理的转义字符
	if escapeNext {
		return nil, fmt.Errorf("dangling escape character at end of command line: %s", command)
	}

	return args, nil
}

func ClearNewline(str string) string {
	return strings.TrimRight(strings.Replace(str, "\r\n", "\n", -1), "\n")
}

func IsFilePath(path string) bool {
	pathPattern := `^\/(?:[^\/]+\/)*[^\/]+(?:\.[^\/]+)?$`
	regex, _ := regexp.Compile(pathPattern)

	if !regex.MatchString(path) {
		return false
	}

	return true
}
