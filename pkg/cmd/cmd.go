package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// Exec executes the given command with the given arguments returning the output and optionally prints output to console.
func Exec(cmdName string, cmdArgs []string, printOutput bool) (string, error) {
	return ExecWithStringInput(cmdName, cmdArgs, "", printOutput)
}

// ExecWithStringInput executes the given command with the given arguments and optionally prints output to console.
func ExecWithStringInput(cmdName string, cmdArgs []string, input string, printOutput bool) (string, error) {
	cmd := exec.Command(cmdName, cmdArgs...)
	if input != "" {
		cmd.Stdin = strings.NewReader(input)
	}
	cmdStr := strings.Trim(fmt.Sprintf("%s", cmdArgs), "[]")
	if printOutput {
		if input != "" {
			fmt.Printf("**** Exec '%s %v' with input: '%s' \n", cmdName, cmdStr, input)

		} else {
			fmt.Printf("**** Exec '%s %v'\n", cmdName, cmdStr)
		}
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrapf(err, "Failed '%s %s", cmdName, cmdStr)
	}
	if printOutput {
		fmt.Print(string(out))
	}
	return string(out), nil

}
