package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r, domain)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users []User

func getUsers(r io.Reader, domain string) (result users, err error) {
	result = make(users, 0)

	scanner := bufio.NewScanner(r)
	var user *User
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "."+domain) {
			if err = json.Unmarshal(scanner.Bytes(), &user); err != nil {
				return
			}
			result = append(result, *user)
			*user = User{}
		}
	}

	err = scanner.Err()
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		if strings.Contains(user.Email, domain) {
			d := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[d]++
		}
	}
	return result, nil
}
