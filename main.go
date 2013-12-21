package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	user := os.Getenv("SENDLOG_USER")
	pass := os.Getenv("SENDLOG_PASSWD")
	addr := os.Getenv("SENDLOG_ADDR")
	server := strings.Split(addr, ":")[0]

	auth := smtp.PlainAuth("", user, pass, server)
	err := smtp.SendMail(addr, auth, user, []string{"rwcarlsen@gmail.com"}, buildMsg())
	if err != nil {
		log.Fatal(err)
	}
}

var lognames = []string{
	"JOURNAL_ERRS.txt",
	"CRON_RUNS.txt",
	"NET_SERVICES.txt",
	"ORPHANS.txt",
	"backup-summary.log",
	"nasbackup.log",
	"localbackup.log",
}

func buildMsg() (msg []byte) {
	subj := fmt.Sprintf("Subject: rwc nightly system logs (%v)\r\n\r\n", time.Now())
	buf := bytes.NewBufferString(subj)
	for _, fname := range lognames {
		fmt.Fprintf(buf, "------ %v ------\n", fname)
		fpath := filepath.Join(os.Getenv("HOME"), "logs", fname)
		data, err := ioutil.ReadFile(fpath)
		if err != nil {
			buf.WriteString(err.Error())
		} else {
			lines := bytes.Split(data, []byte("\n"))
			for i, line := range lines {
				buf.Write(line)
				buf.WriteString("\n")
				if i == 100 {
					fmt.Fprintf(buf, "... %v lines of output truncated\n", len(lines)-100)
					break
				}
			}
		}
		buf.WriteString("\r\n\r\n")
	}
	return buf.Bytes()
}
