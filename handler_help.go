package main

import (
	"fmt"
	"sort"
)

func handlerHelp(s *state, cmd command, cmds *commands) error {
	fmt.Println("Available commands:")

	names := make([]string, 0, len(cmds.descriptions))
	for name := range cmds.descriptions {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		fmt.Printf("  %-12s %s\n", name, cmds.descriptions[name])
	}
	return nil
}
