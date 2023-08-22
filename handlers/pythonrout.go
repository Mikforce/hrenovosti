package handlers

import (
	"fmt"
	"os"
	"os/exec"
)

func Gopython() {
	// Specify the Python script file to execute
	pythonScript := "/home/mik/Документы/GitHub/hrenovosti/parse_python/parsria.py"

	// Prepare the command to run the Python script
	cmd := exec.Command("/usr/bin/python3", pythonScript)

	// Set up pipes for standard output and error
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the Python script
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
func Gopythontwo() {
	// Specify the Python script file to execute
	pythonScript := "/home/mik/Документы/GitHub/hrenovosti/parse_python/parspanorama.py"

	// Prepare the command to run the Python script
	cmd := exec.Command("/usr/bin/python3", pythonScript)

	// Set up pipes for standard output and error
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the Python script
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
