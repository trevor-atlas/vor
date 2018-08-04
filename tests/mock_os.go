package tests

import "time"

type happyMockOSHandler struct{}

func (m happyMockOSHandler) Confirm(string) bool {
	return true
}

func (m happyMockOSHandler) Exit(string) {}

func (m happyMockOSHandler) Exists(str string) (bool, error) {
	return true, nil
}

func (m happyMockOSHandler) Exec(string) (string, error) {
	return "", nil
}

func (m happyMockOSHandler) ExecutionTimer(time.Time, string) {}

type sadMockOSHandler struct{}

func (m sadMockOSHandler) Confirm(string) bool {
	return false
}

func (m sadMockOSHandler) Exit(string) {}

func (m sadMockOSHandler) Exists(string) (bool, error) {
	return false, nil
}

func (m sadMockOSHandler) Exec(string) (string, error) {
	return "", *new(error)
}

func (m sadMockOSHandler) ExecutionTimer(time.Time, string) {}
