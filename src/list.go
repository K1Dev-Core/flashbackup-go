package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (a *App) handleList(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: /list source|dest|db [path]")
		return
	}
	kind := strings.ToLower(args[0])
	path := ""
	if len(args) > 1 {
		path = strings.Join(args[1:], " ")
	}
	switch kind {
	case "source":
		if path == "" {
			path = a.source
		}
		a.listDirectory(path, "source")
	case "dest":
		if path == "" {
			path = a.dest
		}
		a.listDirectory(path, "destination")
	case "db":
		if path == "" {
			path = a.dest
		}
		path, err := directory(path)
		if err != nil {
			fmt.Println("Warning: destination is required for /list db.")
			return
		}
		var files []File
		if err := a.db.Where("dest = ?", path).Order("idx").Find(&files).Error; err != nil {
			fmt.Println("Warning:", err)
			return
		}
		if len(files) == 0 {
			fmt.Println("No backup history for", path)
			return
		}
		for _, file := range files {
			fmt.Println(file.Filename)
		}
	default:
		fmt.Println("Usage: /list source|dest|db [path]")
	}
}

func (a *App) listDirectory(raw, label string) {
	path, err := directory(raw)
	if err != nil {
		fmt.Println("Warning:", err)
		return
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Warning:", err)
		return
	}
	count := 0
	for _, entry := range entries {
		if entry.Type().IsRegular() {
			fmt.Println(entry.Name())
			count++
		}
	}
	if count == 0 {
		fmt.Println("No files in", label, path)
	}
}

func (a *App) check() {
	if a.dest == "" {
		fmt.Println("Warning: set a destination with /dest first.")
		return
	}
	a.integrity(a.dest)
}

func (a *App) integrity(dest string) {
	var files []File
	if err := a.db.Where("dest = ?", dest).Order("idx").Find(&files).Error; err != nil {
		fmt.Println("Warning:", err)
		return
	}
	missing := 0
	for _, file := range files {
		if _, err := os.Stat(filepath.Join(dest, file.Filename)); errors.Is(err, os.ErrNotExist) {
			fmt.Println("Missing:", file.Filename)
			missing++
		}
	}
	if missing == 0 {
		fmt.Println("Integrity check passed.")
	} else {
		fmt.Printf("Integrity check found %d missing file(s).\n", missing)
	}
}
