package system

import (
	"testing"
)

type MockExiter struct {
	Code int
}

func (e *MockExiter) Exit(code int) {
	e.Code = code
}

func TestExit(t *testing.T) {
	mockExiter := MockExiter{}
	DefaultExitor = &mockExiter

	Exit(1)

	exp := 1
	res := mockExiter.Code

	if res != exp {
		t.Errorf("expected exit code %d, got %d", exp, res)
	}
	
}