package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func readSMTPResponse(reader *bufio.Reader) (string, error) {
	resp, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp), nil
}

func logMessage(message string) {
	fmt.Println(message)
}

func sendEmail(email, mxServer, port, fromdomain, from, subject, body string) error {
	var conn net.Conn
	var err error

	dialDirectly := func() (net.Conn, error) {
		conn, err := net.DialTimeout("tcp", mxServer+":"+port, 10*time.Second)
		if err != nil {
			return nil, fmt.Errorf("could not connect to SMTP server %s on port %s: %v", mxServer, port, err)
		}
		return conn, nil
	}

	conn, err = dialDirectly()
	if err != nil {
		return err
	}
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Read initial response
	resp, err := readSMTPResponse(reader)
	if err != nil {
		return fmt.Errorf("could not read initial response: %v", err)
	}
	logMessage(fmt.Sprintf("Initial response: %s", resp))

	// Send EHLO command
	fmt.Fprintf(conn, "EHLO "+fromdomain+"\r\n")
	resp, err = readSMTPResponse(reader)
	if err != nil {
		return fmt.Errorf("could not read EHLO response: %v", err)
	}
	logMessage(fmt.Sprintf("EHLO response: %s", resp))

	// Send MAIL FROM command
	fmt.Fprintf(conn, "MAIL FROM:<"+from+">\r\n")
	resp, err = readSMTPResponse(reader)
	if err != nil {
		return fmt.Errorf("could not read MAIL FROM response: %v", err)
	}
	logMessage(fmt.Sprintf("MAIL FROM response: %s", resp))

	// Send RCPT TO command
	fmt.Fprintf(conn, "RCPT TO:<%s>\r\n", email)
	resp, err = readSMTPResponse(reader)
	if err != nil {
		return fmt.Errorf("could not read RCPT TO response: %v", err)
	}
	logMessage(fmt.Sprintf("RCPT TO response: %s", resp))

	// Send DATA command
	fmt.Fprintf(conn, "DATA\r\n")
	resp, err = readSMTPResponse(reader)
	if err != nil {
		return fmt.Errorf("could not read DATA response: %v", err)
	}
	logMessage(fmt.Sprintf("DATA response: %s", resp))

	// Send email content
	emailContent := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n.\r\n", from, email, subject, body)
	fmt.Fprintf(conn, emailContent)
	resp, err = readSMTPResponse(reader)
	if err != nil {
		return fmt.Errorf("could not read final response: %v", err)
	}
	logMessage(fmt.Sprintf("Email send response: %s", resp))

	if !strings.HasPrefix(resp, "250") {
		return fmt.Errorf("failed to send email, server response: %s", resp)
	}

	return nil
}

func main() {
	// Exemplo de uso da função sendEmail
	err := sendEmail("ju@julianabenfatti.com.br", "mx.b.locaweb.com.br", "25", "boxarmazenagens.com.br", "remetente@boxarmazenagens.com.br", "Assunto do Email", "Corpo do email.")
	if err != nil {
		fmt.Printf("Erro ao enviar email: %v\n", err)
	} else {
		fmt.Println("Email enviado com sucesso!")
	}
}
