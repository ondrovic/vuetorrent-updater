// system_test.go
package system

import (
	"errors"
	"os"
	"os/exec"
	"testing"
)

// MockExiter is a mock implementation of Exitor that does not call os.Exit
type MockExiter struct {
	Code   int
	Called bool
}

// Exit sets the exit code but does not exit the program
func (e *MockExiter) Exit(code int) {
	e.Code = code
	e.Called = true
}

func TestExit(t *testing.T) {
	// Replace the DefaultExitor with a MockExiter
	mockExiter := &MockExiter{}
	DefaultExitor = mockExiter

	// Call the Exit function
	Exit(1)

	// Check that the MockExiter recorded the correct exit code
	if mockExiter.Code != 1 {
		t.Errorf("expected exit code 1, got %d", mockExiter.Code)
	}
}

func TestCheckError(t *testing.T) {
	// Replace the DefaultExitor with a MockExiter
	mockExiter := &MockExiter{}
	DefaultExitor = mockExiter

	// Define an error to pass to CheckError
	err := errors.New("test error")

	// Call the CheckError function with the error
	CheckError(err, 1)

	// Check that the MockExiter recorded the correct exit code
	if mockExiter.Code != 1 {
		t.Errorf("expected exit code 1, got %d", mockExiter.Code)
	}
	if !mockExiter.Called {
		t.Errorf("expected exit to be called")
	}

	// Call the CheckError function without an error
	mockExiter.Code = 0
	mockExiter.Called = false
	CheckError(nil, 1)

	// Check that the MockExiter did not call exit
	if mockExiter.Code != 0 {
		t.Errorf("expected exit code to remain 0, got %d", mockExiter.Code)
	}
	if mockExiter.Called {
		t.Errorf("expected exit not to be called")
	}
}

func TestOSExiter(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		DefaultExitor = OSExiter{}
		Exit(1)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestOSExiter")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return // we expect a non-zero exit code
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
