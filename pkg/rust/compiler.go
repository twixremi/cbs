package rust

import (
	"fmt"
	"os/exec"
)

type RustHelper struct{}

func (r *RustHelper) Build(path string, release bool) error {
	args := []string{"build"}
	if release {
		args = append(args, "--release")
	}

	fmt.Printf("[RUST] Running cargo %v in %s\n", args, path)
	cmd := exec.Command("cargo", args...)
	cmd.Dir = path
	
	// We want to see cargo output
	cmd.Stdout = nil // or os.Stdout
	cmd.Stderr = nil // or os.Stderr

	return cmd.Run()
}

func (r *RustHelper) RunTests(path string) error {
	fmt.Printf("[RUST] Running cargo test in %s\n", path)
	cmd := exec.Command("cargo", "test")
	cmd.Dir = path
	return cmd.Run()
}
