package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"gorm.io/gorm"
)

func (a *App) handleMove(raw string) {
	if a.source == "" {
		fmt.Println("Warning: set a source with /source first.")
		return
	}
	if a.dest == "" {
		fmt.Println("Warning: set a destination with /dest first.")
		return
	}
	names, err := requestedNames(raw)
	if err != nil {
		fmt.Println("Warning:", err)
		return
	}
	if len(names) == 1 && strings.EqualFold(names[0], "all") {
		entries, readErr := os.ReadDir(a.source)
		if readErr != nil {
			fmt.Println("Warning:", readErr)
			return
		}
		names = names[:0]
		for _, entry := range entries {
			if entry.Type().IsRegular() {
				names = append(names, entry.Name())
			}
		}
	}
	start := time.Now()
	moved := 0
	for _, name := range names {
		if err := a.moveOne(name); err != nil {
			fmt.Printf("%s: %v\n", name, err)
			continue
		}
		moved++
		fmt.Println("Moved:", name)
	}
	fmt.Printf("Moved %d file(s) in %.2f ms.\n", moved, float64(time.Since(start).Nanoseconds())/1e6)
}

func (a *App) moveOne(name string) error {
	if err := validFilename(name); err != nil {
		return err
	}
	src := filepath.Join(a.source, name)
	dst := filepath.Join(a.dest, name)
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("source file: %w", err)
	}
	if !info.Mode().IsRegular() {
		return errors.New("source is not a regular file")
	}
	if _, err := os.Stat(dst); err == nil {
		return errors.New("destination file already exists")
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("check destination: %w", err)
	}
	var existing File
	if err := a.db.Where("dest = ? AND filename = ?", a.dest, name).First(&existing).Error; err == nil {
		return errors.New("file already exists in database history")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("check database: %w", err)
	}
	if err := transfer(src, dst); err != nil {
		return err
	}
	if err := a.db.Create(&File{Dest: a.dest, Filename: name}).Error; err != nil {
		if rollbackErr := transfer(dst, src); rollbackErr != nil {
			return fmt.Errorf("save history: %v; rollback: %v", err, rollbackErr)
		}
		return fmt.Errorf("save history: %w", err)
	}
	return nil
}

func transfer(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	} else if !errors.Is(err, syscall.EXDEV) {
		return fmt.Errorf("move: %w", err)
	}
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("copy source: %w", err)
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("copy destination: %w", err)
	}
	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	if copyErr != nil {
		_ = os.Remove(dst)
		return fmt.Errorf("copy: %w", copyErr)
	}
	if closeErr != nil {
		_ = os.Remove(dst)
		return fmt.Errorf("close destination: %w", closeErr)
	}
	if err := os.Remove(src); err != nil {
		_ = os.Remove(dst)
		return fmt.Errorf("remove source: %w", err)
	}
	return nil
}
