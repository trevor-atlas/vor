package tests

type sadMockGetter struct {
	Countstring int
	Countbool   int
}

func (m sadMockGetter) String(string) string {
	m.Countstring += 1
	return ""
}

func (m sadMockGetter) Bool(string) bool {
	m.Countbool += 1
	return false
}

type happyMockGetter struct {
	Countstring int
	Countbool   int
}

func (m happyMockGetter) String(string) string {
	m.Countstring += 1
	return "/usr/local/bin/git"
}

func (m happyMockGetter) Bool(string) bool {
	m.Countbool += 1
	panic("implement me")
}
