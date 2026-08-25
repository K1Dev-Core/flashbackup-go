package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func (a *App) handleDelete(args []string) {
	if a.dest == "" {
		fmt.Println("Warning: set a destination with /dest first.")
		return
	}
	if len(args) < 2 || strings.ToLower(args[0]) != "dest" {
		fmt.Println("Usage: /delete dest <file1>, ...|all")
		return
	}
	names, err := requestedNames(strings.Join(args[1:], " "))
	if err != nil {
		fmt.Println("Warning:", err)
		return
	}
	if len(names) == 1 && strings.EqualFold(names[0], "all") {
		names = a.destinationNames()
	}
	for _, name := range names {
		if err := a.deleteOne(name); err != nil {
			fmt.Printf("%s: %v\n", name, err)
			continue
		}
		fmt.Println("Deleted:", name)
	}
}

func (a *App) destinationNames() []string {
	seen := map[string]bool{}
	var files []File
	_ = a.db.Where("dest = ?", a.dest).Find(&files).Error
	for _, file := range files {
		seen[file.Filename] = true
	}
	if entries, err := os.ReadDir(a.dest); err == nil {
		for _, entry := range entries {
			if entry.Type().IsRegular() {
				seen[entry.Name()] = true
			}
		}
	}
	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (a *App) deleteOne(name string) error {
	if err := validFilename(name); err != nil {
		return err
	}
	path := filepath.Join(a.dest, name)
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove file: %w", err)
	}
	if err := a.db.Where("dest = ? AND filename = ?", a.dest, name).Delete(&File{}).Error; err != nil {
		return fmt.Errorf("remove history: %w", err)
	}
	return nil
}
