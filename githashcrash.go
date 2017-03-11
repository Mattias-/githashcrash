package main

import (
    "bytes"
	"crypto/sha1"
    "encoding/hex"
	"fmt"
    "log"
    "os"
    "io"
    "math/rand"
    "os/exec"
    "regexp"
    "strings"
    "time"
)

func main() {
    args := os.Args
    var targetHash = regexp.MustCompile(args[1])

    var obj []byte
    if len(args) == 3 {
        obj = []byte(args[2])
    } else {
        obj, _ = exec.Command("git", "cat-file", "-p", "HEAD").Output()
    }
    if !bytes.HasSuffix(obj, []byte("\n")) {
       obj = append(obj, "\n"...)
    }

    lines := strings.Split(string(obj), "\n")
    author_line := strings.Split(lines[1], " ")
    author_date := strings.Join(author_line[len(author_line)-2:], " ")
    author_email := strings.Trim(author_line[len(author_line)-3:len(author_line)-2][0], "<>")
    author_name := strings.Join(author_line[1:len(author_line)-3], " ")
    committer_line := strings.Split(lines[2], " ")
    committer_date := strings.Join(committer_line[len(committer_line)-2:], " ")
    committer_email := strings.Trim(committer_line[len(committer_line)-3:len(committer_line)-2][0], "<>")
    committer_name := strings.Join(committer_line[1:len(committer_line)-3], " ")

    seed := fmt.Sprintf("%x", rand.Intn(99999))
    val, ok := os.LookupEnv("GITHASHCRASH_SEED")
    if ok {
        seed = val
    }

    start := time.Now()
    matched := false
    extra := ""
    i := 1
    for !matched {
        i++
        extra = fmt.Sprintf("%s-%d", seed, i)

        h := sha1.New()

        io.WriteString(h, fmt.Sprintf("commit %d\x00", len(obj) + len(extra) + 1))
        h.Write(obj)
        io.WriteString(h, extra)
        io.WriteString(h, "\n")

        encodedStr := hex.EncodeToString(h.Sum(nil))
        matched = targetHash.MatchString(encodedStr)
        if i % 100000 == 0{
            log.Println(i)
        }
    }
    elapsed := time.Since(start)

    log.Println("Time:", elapsed)
    log.Println("Commits tested:", i)
    log.Println("Tests per second:", float64(i)/elapsed.Seconds())

    fmt.Println("Recreate with:")
    envString := strings.Join([]string{
        "export",
        fmt.Sprintf("GIT_AUTHOR_DATE='%s'", author_date),
        fmt.Sprintf("GIT_AUTHOR_NAME='%s'", author_name),
        fmt.Sprintf("GIT_AUTHOR_EMAIL='%s'", author_email),
        fmt.Sprintf("GIT_COMMITTER_DATE='%s'", committer_date),
        fmt.Sprintf("GIT_COMMITTER_NAME='%s'", committer_name),
        fmt.Sprintf("GIT_COMMITTER_EMAIL='%s'", committer_email),
    }, " ")
    fmt.Printf("(%s; printf '%%s\\n%s' \"$(git show -s --format=%%B)\" | git commit --amend -F -)\n", envString, extra)

}
