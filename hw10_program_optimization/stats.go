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
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	scanner := bufio.NewScanner(r)
	i := 0
	var user *User
	for scanner.Scan() {
		if err = json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return
		}
		result[i] = *user
		i++
	}

	err = scanner.Err()
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		if strings.Contains(user.Email, domain) {
			d := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			num := result[d]
			num++
			result[d] = num
		}
	}
	return result, nil
}
