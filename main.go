package main

import (
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
	"localbackup.txt",
	"nasbackup.txt",
	"NET_SERVICES.txt",
	"ORPHANS.txt",
}

func buildMsg() (msg []byte) {
	subj := fmt.Sprintf("Subject: rwc nightly system logs (%v)\r\n\r\n", time.Now())
	msg = []byte(subj)
	for _, fname := range lognames {
		msg = append(msg, []byte("------ "+fname+" ------\n")...)
		fpath := filepath.Join(os.Getenv("HOME"), "logs", fname)
		data, err := ioutil.ReadFile(fpath)
		if err != nil {
			msg = append(msg, []byte(err.Error())...)
		} else {
			msg = append(msg, data...)
		}
		msg = append(msg, []byte("\r\n\r\n")...)
	}
	return msg
}
