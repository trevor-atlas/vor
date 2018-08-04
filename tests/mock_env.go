package tests

// sad path
type sadMockGetter struct{}

func (m sadMockGetter) String(string) string {
	return ""
}

func (m sadMockGetter) Bool(string) bool {
	return false
}

// happy path
type happyMockGetter struct{}

func (m happyMockGetter) String(string) string {
	return "/usr/local/bin/git"
}

func (m happyMockGetter) Bool(string) bool {
	panic("implement me")
}
