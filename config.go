package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/samuelotter/i3ipc"
)

type Rule struct {
	Event   i3ipc.EventType
	Change  string
	Action  Action
	Args    []string
}

type Config struct {
	Rules []Rule
}

type Action int
const (
	ActionIgnore Action = iota
	ActionExec
)

var actionNames = map[string]Action {
	"ignore": ActionIgnore,
	"exec": ActionExec,
}

func ReadConfiguration(configPath string) (*Config, error) {
	if configPath == "" {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}
		configPath = path.Join(usr.HomeDir, ".i3event")
	}

	file, err := os.Open(configPath)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	var rules []Rule
	line := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line += 1
		str := strings.TrimSpace(scanner.Text())
		switch {
		case strings.HasPrefix(str, "#"):
			// Ignore comment
			continue
		case str == "":
			// Ignore empty lines
			continue
		case strings.HasPrefix(str, "bindevent"):
			tokens := strings.Fields(str)
			if len(tokens) < 4 {
				return nil, errors.New("Incomplete command bindevent, expected bindevent <event> <change> <action> [args..]")
			}
			action, ok := actionNames[tokens[3]]
			if !ok {
				return nil, fmt.Errorf("Invalid action: %s", tokens[3])
			}
			rule := Rule{
				Event:  eventTypes[tokens[1]],
				Change: tokens[2],
				Action: action,
				Args:   tokens[4:],
			}
			rules = append(rules, rule)
		default:
			return nil,
			fmt.Errorf("Invalid token at line %d: %s", line, str)
		}
	}
	return &Config{
		Rules: rules,
	}, nil
}

