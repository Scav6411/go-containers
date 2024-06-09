package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	// Parse command-line arguments to decide whether to run as parent or child
	switch os.Args[1] {
	case "run":
		parent() // If "run" argument is provided, run as parent process
	case "child":
		child() // If "child" argument is provided, run as child process (inside container)
	default:
		panic("wat should I do") // If neither "run" nor "child" is provided, panic
	}
}

// Parent function sets up container environment and runs child process inside it
func parent() {
	// Set up command to run same program with "child" argument plus additional arguments
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	// Configure command to run with new namespaces (UTS, PID, and mount namespace)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	// Set standard input, output, and error streams to use same streams as parent process
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}

// Child function sets up container filesystem and runs specified command inside it
func child() {
	// Mount new filesystem from "rootfs" directory to "/rootfs"
	must(syscall.Mount("rootfs", "rootfs", "", syscall.MS_BIND, ""))
	// Create new directory "/rootfs/oldrootfs" to store original root filesystem
	must(os.MkdirAll("rootfs/oldrootfs", 0700))
	// Pivot root filesystem to newly mounted filesystem
	must(syscall.PivotRoot("rootfs", "rootfs/oldrootfs"))
	// Change current working directory to "/"
	must(os.Chdir("/"))

	// Run specified command (usually the application to run inside container)
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}

// Must function is a simple error handling function that panics if an error is encountered
func must(err error) {
	if err != nil {
		panic(err)
	}
}
