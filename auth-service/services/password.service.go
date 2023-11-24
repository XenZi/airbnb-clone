package services

import (
	"bufio"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct {
	blacklistedPasswords map[string]string
}

func NewPasswordService() *PasswordService {
	passwordService := PasswordService{}
	passwordService.blacklistedPasswords = make(map[string]string)
	passwordService.readPasswordsFromFile()
	return &passwordService
}

func (p PasswordService) HashPassword(password string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func (p PasswordService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (p PasswordService) readPasswordsFromFile() {
	file, err := os.Open("data/blacklist-passwords.txt")
	if err != nil {
		log.Println("Error while reading blacklist passwords, here is why: ", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		password := strings.TrimSpace(line)
		p.blacklistedPasswords[password] = password
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error while reading blacklist passwords, here is why: ", err)
		return
	}
}

func (p PasswordService) CheckPasswordExistanceInBlacklist(password string) bool {
	return p.blacklistedPasswords[password] != ""
}