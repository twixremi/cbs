package runner

import (
	"crypto/sha256"
	"crow-build-system/config"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

type Engine struct {
	File *config.CrowFile
}

func NewEngine(cf *config.CrowFile) *Engine {
	return &Engine{File: cf}
}

func (e *Engine) Run(taskName string) error {
	task, ok := e.File.Config.Tasks[taskName]
	if !ok {
		return fmt.Errorf("task '%s' not found", taskName)
	}

	// Calculate input hash
	inputHash, err := e.calculateInputsHash(task.Inputs)
	if err != nil {
		return err
	}

	// Check if UP-TO-DATE
	lock, exists := e.File.Locks.Tasks[taskName]
	if exists && lock.Success && lock.InputHash == inputHash && inputHash != "" {
		fmt.Printf("[CROW] Task '%s' is UP-TO-DATE (skipped)\n", taskName)
		return nil
	}

	fmt.Printf("[RUNNER] Executing task: %s\n", taskName)
	
	cmd := exec.Command("bash", "-c", task.Command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		e.File.Locks.Tasks[taskName] = config.TaskLock{InputHash: inputHash, Success: false}
		return fmt.Errorf("task failed: %v", err)
	}

	// Save success state
	e.File.Locks.Tasks[taskName] = config.TaskLock{InputHash: inputHash, Success: true}
	e.File.Save()

	fmt.Printf("[RUNNER] Task '%s' completed successfully.\n", taskName)
	return nil
}

func (e *Engine) calculateInputsHash(inputs []string) (string, error) {
	if len(inputs) == 0 {
		return "", nil
	}

	var allFiles []string
	for _, pattern := range inputs {
		matches, err := filepath.Glob(pattern)
		if err != nil || len(matches) == 0 {
			// If no matches, treat as a single file path
			allFiles = append(allFiles, pattern)
			continue
		}
		allFiles = append(allFiles, matches...)
	}

	sort.Strings(allFiles)
	
	h := sha256.New()
	for _, path := range allFiles {
		f, err := os.Open(path)
		if err != nil {
			h.Write([]byte(path))
			continue
		}
		
		if _, err := io.Copy(h, f); err != nil {
			f.Close()
			return "", err
		}
		f.Close()
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
