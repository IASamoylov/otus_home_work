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
	return countDomains(u[:len(u)])
}

type users []User

func getUsers(r io.Reader, domain string) (result users, err error) {
	result = make(users, 0)

	scanner := bufio.NewScanner(r)
	var user *User
	b := []byte("." + domain + "\",\"Phone\"")
	for scanner.Scan() {
		if isSubArray(scanner.Bytes(), b, len(scanner.Bytes()), len(b)) {
			if err = json.Unmarshal(scanner.Bytes(), &user); err != nil {
				return
			}

			result = append(result, *user)
		}
	}

	err = scanner.Err()
	return
}

func countDomains(u []User) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		d := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
		result[d]++
	}
	return result, nil
}

func isSubArray(a, b []byte, n, m int) bool {
	i := 0
	j := 0

	for i < n && j < m {
		if a[i] == b[j] {
			i++
			j++

			if j == m {
				return true
			}
		} else {
			i = i - j + 1
			j = 0
		}
	}

	return false
}
