package repository

type TestRepo struct {
}

func (t TestRepo) SayHi() string {
	return "Hello from repo!"
}
