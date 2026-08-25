package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gorm.io/gorm"
)

type App struct {
	db     *gorm.DB
	source string
	dest   string
}

func (a *App) run() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if a.command(line) {
			return
		}
	}
	fmt.Println()
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "input error:", err)
	}
}

func (a *App) command(line string) bool {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return false
	}
	switch parts[0] {
	case "/source":
		a.handleSource(strings.TrimSpace(strings.TrimPrefix(line, "/source")))
	case "/dest":
		a.handleDest(strings.TrimSpace(strings.TrimPrefix(line, "/dest")))
	case "/list":
		a.handleList(parts[1:])
	case "/move":
		a.handleMove(strings.TrimSpace(strings.TrimPrefix(line, "/move")))
	case "/check":
		a.check()
	case "/delete":
		a.handleDelete(parts[1:])
	case "/exit":
		fmt.Println("Bye.")
		return true
	default:
		fmt.Println("Unknown command. Use /source, /dest, /list, /move, /check, /delete, or /exit.")
	}
	return false
}
