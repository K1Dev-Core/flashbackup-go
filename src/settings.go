package main

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func (a *App) loadSettings() {
	var setting Setting
	if err := a.db.Order("idx DESC").First(&setting).Error; err == nil {
		if isDir(setting.Source) {
			a.source = setting.Source
		}
		if isDir(setting.Dest) {
			a.dest = setting.Dest
		}
	}
}

func (a *App) saveSettings() error {
	if a.source == "" || a.dest == "" {
		return nil
	}
	var setting Setting
	if err := a.db.Order("idx DESC").First(&setting).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return a.db.Create(&Setting{Source: a.source, Dest: a.dest}).Error
	} else if err == nil {
		return a.db.Model(&setting).Updates(map[string]any{"source": a.source, "dest": a.dest}).Error
	} else {
		return err
	}
}

func (a *App) handleSource(raw string) {
	path, err := directory(raw)
	if err != nil {
		fmt.Println("Warning:", err)
		return
	}
	a.source = path
	if err := a.saveSettings(); err != nil {
		fmt.Println("Warning:", err)
		return
	}
	fmt.Println("Source:", path)
}

func (a *App) handleDest(raw string) {
	path, err := directory(raw)
	if err != nil {
		fmt.Println("Warning:", err)
		return
	}
	a.dest = path
	if err := a.saveSettings(); err != nil {
		fmt.Println("Warning:", err)
		return
	}
	fmt.Println("Destination:", path)
	a.integrity(path)
}
