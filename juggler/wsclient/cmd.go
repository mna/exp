package main

import "sort"

type cmd struct {
	Help string
	Run  func(...string)
}

var helpCmd = &cmd{
	Help: "print this message",

	Run: func(_ ...string) {
		keys := make([]string, 0, len(commands))
		for k := range commands {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			print("%s\t%s", k, commands[k].Help)
		}
	},
}
