package main

import (
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
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
	"PKGLIST.txt",
	"ORPHANS.txt",
	"NET_SERVICES.txt",
	"localbackup.txt",
	"nasbackup.txt",
	"JOURNAL_ERRS.txt",
	"CRON_RUNS.txt",
}

func buildMsg() (msg []byte) {
	for _, fname := range lognames {
		msg = append(msg, []byte("------ "+fname+" ------\n")...)
		fpath := filepath.Join(os.Getenv("HOME"), "logs", fname)
		data, err := ioutil.ReadFile(fpath)
		if err != nil {
			msg = append(msg, []byte(err.Error())...)
		} else {
			msg = append(msg, data...)
		}
		msg = append(msg, []byte("\n\n")...)
	}
	return msg
}
