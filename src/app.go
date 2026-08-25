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

func (a *App) printHeader() {
	fmt.Println("Interactive Console — พิมพ์ /help เพื่อดูคำสั่ง หรือ /exit เพื่อออก")
	if a.dest == "" || !isDir(a.dest) {
		fmt.Println("dest เป็นค่าว่าง กรุณาตั้งค่าด้วย /dest <path>")
	}
}

func (a *App) run() {
	if interactiveTerminal() {
		a.runInteractive()
		return
	}
	a.runScanner()
}

func (a *App) runScanner() {
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
	case "/help":
		a.help()
	case "/source":
		a.handleSource(strings.TrimSpace(strings.TrimPrefix(line, "/source")))
	case "/dest":
		a.handleDest(strings.TrimSpace(strings.TrimPrefix(line, "/dest")))
	case "/set":
		a.handleSet(parts[1:])
	case "/add":
		a.handleAdd(parts[1:])
	case "/settings":
		a.handleSettings()
	case "/list":
		a.handleList(parts[1:])
	case "/move":
		a.handleMove(strings.TrimSpace(strings.TrimPrefix(line, "/move")))
	case "/check":
		a.check()
	case "/clean":
		a.handleClean()
	case "/delete":
		a.handleDelete(parts[1:])
	case "/exit":
		fmt.Println("ออกจากโปรแกรม")
		return true
	default:
		fmt.Println("Unknown command. Type /help to see available commands.")
	}
	return false
}

func (a *App) help() {
	fmt.Println("คำสั่งที่ใช้ได้: /add, /check, /clean, /delete, /dest, /exit, /help, /list, /move, /set, /settings, /source")
}
