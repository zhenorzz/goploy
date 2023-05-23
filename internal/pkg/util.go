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
	state := "start"
	current := ""
	quote := "\""
	escapeNext := true
	for i := 0; i < len(command); i++ {
		c := command[i]

		if state == "quotes" {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = "start"
			}
			continue
		}

		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if c == '"' {
			state = "quotes"
			quote = string(c)
			continue
		}

		if state == "arg" {
			if c == ' ' || c == '=' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '=' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}

	if state == "quotes" {
		return []string{}, fmt.Errorf("unclosed quote in command line: %s", command)
	}

	if current != "" {
		args = append(args, current)
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
