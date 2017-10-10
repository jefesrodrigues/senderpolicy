package main

import (
	"fmt"
	"io"
	"net"
	"net/textproto"
	"strings"
	"log/syslog"
)

func main() {

	var request string
	var action string
	var sender string
	var sasl_username string
	var queue_id string
	var recipient_count string

	logger, e := syslog.New(syslog.LOG_MAIL, "SenderPolicy")
	if e != nil {
		logger.Err(e.Error())
	}

	server := "127.0.0.1:7778"
	listener, err := net.Listen("tcp", server)

	if err != nil {
		logger.Err(err.Error())
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			logger.Err(err.Error())
			break
		}

		go func() {
			textproto := textproto.NewConn(conn)
			lines := []string{}
			for {
				line, err := textproto.ReadLine()
				if err == io.EOF {
					break
				}
				if err != nil {
					logger.Err(err.Error())
				}
				if line == "" {
					break
				}
				lines = append(lines, line)
			}

			for _, line := range lines {
				headers := strings.Split(line, "=")
				switch headers[0] {
				case "request":
					request = strings.Trim(headers[1], " \r\t\n")
				case "sender":
					sender = strings.Trim(headers[1], " \r\t\n")
				case "sasl_username":
					sasl_username = strings.Trim(headers[1], " \r\t\n")
				case "recipient_count":
					recipient_count = strings.Trim(headers[1], " \r\t\n")
				case "queue_id":
					queue_id = strings.Trim(headers[1], " \r\t\n")
				}

				action = "dunno"

				if request == "smtpd_access_policy" {
					if sender == "" || sasl_username == "" {
						action = "reject"
					}

					if sasl_username == sender {
						action = "ok"
					} else {
						action = "reject"
					}
				}
			}
			logger.Info("sendercheck: " + action + ", Sender: " + sender + ", sasl_username: " + sasl_username + ", Recipient: " + recipient_count + ", Queue-ID: " + queue_id)
			fmt.Fprintf(conn, "action=%s\n\n",action)
			conn.Close()
		}()
	}
}