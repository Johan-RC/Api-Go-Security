// Package mailer envía correos transaccionales vía SMTP (librería estándar).
// Soporta:
//   - Gmail (smtp.gmail.com): puerto 587 con STARTTLS + autenticación,
//     o puerto 465 con TLS implícito.
//   - Servidores locales sin TLS (MailHog, dev): envío en claro sin auth.
package mailer

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
)

// Config concentra la configuración SMTP.
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

// Enabled indica si el envío de correo está configurado (host y remitente).
func (c Config) Enabled() bool {
	return c.Host != "" && c.From != ""
}

// Send envía un correo de texto plano.
func (c Config) Send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", c.Host, c.Port)
	fromAddr := fromEmail(c.From)
	msg := []byte(buildMessage(c.From, to, subject, body))

	if c.Port == "465" {
		return sendSSL(addr, c.Host, c.Username, c.Password, fromAddr, to, msg)
	}

	return sendSTARTTLS(addr, c.Host, c.Username, c.Password, fromAddr, to, msg)
}

// sendSTARTTLS conecta por el puerto 587 (o 25): usa STARTTLS si el servidor
// lo ofrece y autentica cuando hay credenciales. Sin TLS y sin credenciales
// es compatible con servidores locales de desarrollo.
func sendSTARTTLS(addr, host, username, password, from, to string, msg []byte) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Quit()

	if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(&tls.Config{ServerName: host}); err != nil {
			return err
		}
	}

	if username != "" {
		auth := smtp.PlainAuth("", username, password, host)
		if err := client.Auth(auth); err != nil {
			return err
		}
	}

	return deliver(client, from, to, msg)
}

// sendSSL conecta por el puerto 465 (SMTP implícito TLS).
func sendSSL(addr, host, username, password, from, to string, msg []byte) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	tlsConn := tls.Client(conn, &tls.Config{ServerName: host})
	if err := tlsConn.Handshake(); err != nil {
		return err
	}
	defer tlsConn.Close()

	client, err := smtp.NewClient(tlsConn, host)
	if err != nil {
		return err
	}
	defer client.Quit()

	if username != "" {
		auth := smtp.PlainAuth("", username, password, host)
		if err := client.Auth(auth); err != nil {
			return err
		}
	}

	return deliver(client, from, to, msg)
}

func deliver(client *smtp.Client, from, to string, msg []byte) error {
	if err := client.Mail(from); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	return w.Close()
}

// fromEmail extrae únicamente la dirección de correo de un remitente que
// puede declararse con nombre: `Nombre <correo@dominio.com>`.
func fromEmail(from string) string {
	start := strings.LastIndex(from, "<")
	end := strings.LastIndex(from, ">")
	if start >= 0 && end > start {
		return from[start+1 : end]
	}
	return strings.TrimSpace(from)
}

// buildMessage construye el mensaje RFC 5322 con Content-Type utf-8.
func buildMessage(from, to, subject, body string) string {
	var b strings.Builder
	b.WriteString("From: " + from + "\r\n")
	b.WriteString("To: " + to + "\r\n")
	b.WriteString("Subject: " + subject + "\r\n")
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	b.WriteString("\r\n")
	b.WriteString(body)
	return b.String()
}