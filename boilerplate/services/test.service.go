package services

import "env-test-app/repository"

type TestService struct {
	Repo repository.TestRepo
}

func (t TestService) CallService() string {
	fromRepo := t.Repo.SayHi()
	return fromRepo
}
