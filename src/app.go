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
	fmt.Println("=== Flash Drive Backup CLI (GORM Engine) ===")
	fmt.Println("ระบบสำรองไฟล์จาก Source ไปยัง Destination")
	fmt.Println("พิมพ์ /help เพื่อดูคำสั่ง หรือ /exit เพื่อออก")
	if a.source != "" || a.dest != "" {
		fmt.Println("สถานะปัจจุบัน:")
		if a.source != "" {
			fmt.Println("  Source:", a.source)
		}
		if a.dest != "" {
			fmt.Println("  Destination:", a.dest)
		}
	}
	fmt.Println()
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
	case "/help":
		a.help()
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
		fmt.Println("Program terminated. Goodbye!")
		return true
	default:
		fmt.Println("Unknown command. Type /help to see available commands.")
	}
	return false
}

func (a *App) help() {
	fmt.Println("คำสั่งที่ใช้ได้:")
	fmt.Println("  /source <path>              ตั้งค่าโฟลเดอร์ต้นทาง")
	fmt.Println("  /dest <path>                ตั้งค่าโฟลเดอร์ปลายทาง")
	fmt.Println("  /list source                แสดงไฟล์ใน source")
	fmt.Println("  /list dest                  แสดงไฟล์ใน destination")
	fmt.Println("  /list db [path]             แสดงประวัติจากฐานข้อมูล")
	fmt.Println("  /move <file1>, <file2>      ย้ายไฟล์ที่เลือก")
	fmt.Println("  /move all                   ย้ายไฟล์ทั้งหมด")
	fmt.Println("  /check                      ตรวจสอบไฟล์กับประวัติใน DB")
	fmt.Println("  /delete dest <file>, ...    ลบไฟล์และประวัติ")
	fmt.Println("  /delete dest all            ลบไฟล์ทั้งหมดใน destination")
	fmt.Println("  /exit                       ออกจากโปรแกรม")
}
