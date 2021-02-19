package models

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Set up
	fmt.Println("About to run models tests...")
	// Run tests in package:
	exitCode := m.Run()
	// Actions post test
	fmt.Println("Model tests have completed")
	// Exit:
	os.Exit(exitCode)
}
