package main

import (
	"fmt"
	"strings"
)

type gitUser struct {
	name  string
	email string
	date  string
}

func parseUserLine(line string) gitUser {
	author := gitUser{}
	authorLine := strings.Split(line, " ")
	author.date = strings.Join(authorLine[len(authorLine)-2:], " ")
	author.email = strings.Trim(authorLine[len(authorLine)-3 : len(authorLine)-2][0], "<>")
	author.name = strings.Join(authorLine[1:len(authorLine)-3], " ")
	return author
}

func parseObj(obj []byte) (gitUser, gitUser) {
	lines := strings.Split(string(obj), "\n")
	var author gitUser
	var committer gitUser
	for _, line := range lines {
		if strings.HasPrefix(line, "author ") {
			author = parseUserLine(line)
		}
		if strings.HasPrefix(line, "committer ") {
			committer = parseUserLine(line)
		}
	}
	return author, committer
}

func printRecreate(author gitUser, committer gitUser, extra string) {
	fmt.Println("Recreate with:")
	envString := strings.Join([]string{
		"export",
		fmt.Sprintf("GIT_AUTHOR_DATE='%s'", author.date),
		fmt.Sprintf("GIT_AUTHOR_NAME='%s'", author.name),
		fmt.Sprintf("GIT_AUTHOR_EMAIL='%s'", author.email),
		fmt.Sprintf("GIT_COMMITTER_DATE='%s'", committer.date),
		fmt.Sprintf("GIT_COMMITTER_NAME='%s'", committer.name),
		fmt.Sprintf("GIT_COMMITTER_EMAIL='%s'", committer.email),
	}, " ")
	fmt.Printf("(%s; printf '%%s\\n%s' \"$(git show -s --format=%%B)\" | git commit --amend -F -)\n", envString, extra)
}
