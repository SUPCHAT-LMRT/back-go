package mail

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"mime"
	"mime/multipart"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Mailer struct {
	Host string
	Tls  bool
	Port int
	From string
	smtp.Auth
}

func NewMailer(host string, tls bool, port int, username, password string) func() *Mailer {
	return func() *Mailer {
		return &Mailer{
			Host: host,
			Tls:  tls,
			Port: port,
			From: username,
			Auth: smtp.PlainAuth("", username, password, host),
		}
	}
}

func (m *Mailer) Send(payload *Message) (err error) {
	// Set GMAIL Message-ID header
	payload.AddHeader("Message-ID", "<"+uuid.New().String()+"@"+m.Host+">")

	// Connect to the SMTP Server
	conn, err := m.createConnection()
	if err != nil {
		return
	}

	c, err := smtp.NewClient(conn, m.Host)
	if err != nil {
		return
	}

	// Auth
	if err = c.Auth(m.Auth); err != nil {
		return
	}

	// To && From
	if err = c.Mail(payload.From); err != nil {
		return
	}

	if err = c.Rcpt(strings.Join(payload.To, ",")); err != nil {
		return
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return
	}

	_, err = w.Write(payload.Bytes())
	if err != nil {
		return
	}

	err = w.Close()
	if err != nil {
		return
	}

	return c.Quit()
}

func (m *Mailer) SendAsync(payload *Message) chan<- error {
	errCh := make(chan error)

	go func() {
		errCh <- m.Send(payload)
	}()

	return errCh
}

func (m *Mailer) createConnection() (net.Conn, error) {
	if m.Tls {
		tlsconfig := tls.Config{
			InsecureSkipVerify: true,
			ServerName:         m.Host,
		}
		return tls.Dial("tcp", fmt.Sprintf("%s:%d", m.Host, m.Port), &tlsconfig)
	}

	return net.Dial("tcp", fmt.Sprintf("%s:%d", m.Host, m.Port))
}

// Attachment represents an email attachment.
type Attachment struct {
	Filename string
	Data     []byte
	Inline   bool
}

// Header represents an additional email header.
type Header struct {
	Key   string
	Value string
}

// Message represents a smtp message.
type Message struct {
	From            string
	To              []string
	Cc              []string
	Bcc             []string
	ReplyTo         string
	Subject         string
	Body            string
	BodyContentType string
	Headers         []Header
	Attachments     map[string]*Attachment
}

func (m *Message) attach(filePath string, inline bool) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	_, filename := filepath.Split(filePath)

	m.Attachments[filename] = &Attachment{
		Filename: filename,
		Data:     data,
		Inline:   inline,
	}

	return nil
}

func (m *Message) SetFrom(address string) string {
	m.From = address
	return m.From
}

func (m *Message) AddTo(address string) []string {
	m.To = append(m.To, address)
	return m.To
}

func (m *Message) AddCc(address string) []string {
	m.Cc = append(m.Cc, address)
	return m.Cc
}

func (m *Message) AddBcc(address string) []string {
	m.Bcc = append(m.Bcc, address)
	return m.Bcc
}

// AttachBuffer attaches a binary attachment.
func (m *Message) AttachBuffer(filename string, buf []byte, inline bool) error {
	m.Attachments[filename] = &Attachment{
		Filename: filename,
		Data:     buf,
		Inline:   inline,
	}
	return nil
}

// Attach attaches a file.
func (m *Message) Attach(filePath string) error {
	return m.attach(filePath, false)
}

// Inline includes a file as an inline attachment.
func (m *Message) Inline(filePath string) error {
	return m.attach(filePath, true)
}

// Ads a Header to message
func (m *Message) AddHeader(key string, value string) Header {
	newHeader := Header{Key: key, Value: value}
	m.Headers = append(m.Headers, newHeader)
	return newHeader
}

func newMessage(subject string, body string, bodyContentType string) *Message {
	m := &Message{Subject: subject, Body: body, BodyContentType: bodyContentType}

	m.Attachments = make(map[string]*Attachment)

	return m
}

// NewMessage returns a new Message that can compose an email with attachments
func NewMessage(subject string, body string) *Message {
	return newMessage(subject, body, "text/plain")
}

// NewHTMLMessage returns a new Message that can compose an HTML email with attachments
func NewHTMLMessage(subject string, body string) *Message {
	return newMessage(subject, body, "text/html")
}

// Tolist returns all the recipients of the email
func (m *Message) Tolist() []string {
	rcptList := []string{}

	toList, _ := mail.ParseAddressList(strings.Join(m.To, ","))
	for _, to := range toList {
		rcptList = append(rcptList, to.Address)
	}

	ccList, _ := mail.ParseAddressList(strings.Join(m.Cc, ","))
	for _, cc := range ccList {
		rcptList = append(rcptList, cc.Address)
	}

	bccList, _ := mail.ParseAddressList(strings.Join(m.Bcc, ","))
	for _, bcc := range bccList {
		rcptList = append(rcptList, bcc.Address)
	}

	return rcptList
}

// Bytes returns the mail data
func (m *Message) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString("From: " + m.From + "\r\n")

	t := time.Now()
	buf.WriteString("Date: " + t.Format(time.RFC1123Z) + "\r\n")

	buf.WriteString("To: " + strings.Join(m.To, ",") + "\r\n")
	if len(m.Cc) > 0 {
		buf.WriteString("Cc: " + strings.Join(m.Cc, ",") + "\r\n")
	}

	// fix Encode
	coder := base64.StdEncoding
	subject := "=?UTF-8?B?" + coder.EncodeToString([]byte(m.Subject)) + "?="
	buf.WriteString("Subject: " + subject + "\r\n")

	if len(m.ReplyTo) > 0 {
		buf.WriteString("Reply-To: " + m.ReplyTo + "\r\n")
	}

	buf.WriteString("MIME-Version: 1.0\r\n")

	// Add custom headers
	if len(m.Headers) > 0 {
		for _, header := range m.Headers {
			fmt.Fprintf(buf, "%s: %s\r\n", header.Key, header.Value)
		}
	}

	writer := multipart.NewWriter(bytes.NewBuffer(nil))
	boundary := writer.Boundary()

	if len(m.Attachments) > 0 {
		buf.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + "\r\n")
		buf.WriteString("\r\n--" + boundary + "\r\n")
	}

	fmt.Fprintf(buf, "Content-Type: %s; charset=utf-8\r\n\r\n", m.BodyContentType)
	buf.WriteString(m.Body)
	buf.WriteString("\r\n")

	if len(m.Attachments) > 0 {
		for _, attachment := range m.Attachments {
			buf.WriteString("\r\n\r\n--" + boundary + "\r\n")

			if attachment.Inline {
				buf.WriteString("Content-Type: message/rfc822\r\n")
				buf.WriteString(
					"Content-Disposition: inline; filename=\"" + attachment.Filename + "\"\r\n\r\n",
				)

				buf.Write(attachment.Data)
			} else {
				ext := filepath.Ext(attachment.Filename)
				mimetype := mime.TypeByExtension(ext)
				if mimetype != "" {
					mime := fmt.Sprintf("Content-Type: %s\r\n", mimetype)
					buf.WriteString(mime)
				} else {
					buf.WriteString("Content-Type: application/octet-stream\r\n")
				}
				buf.WriteString("Content-Transfer-Encoding: base64\r\n")

				buf.WriteString("Content-Disposition: attachment; filename=\"=?UTF-8?B?")
				buf.WriteString(coder.EncodeToString([]byte(attachment.Filename)))
				buf.WriteString("?=\"\r\n\r\n")

				b := make([]byte, base64.StdEncoding.EncodedLen(len(attachment.Data)))
				base64.StdEncoding.Encode(b, attachment.Data)

				// write base64 content in lines of up to 76 chars (RFC 2045)
				for i, l := 0, len(b); i < l; i++ {
					buf.WriteByte(b[i])
					if (i+1)%76 == 0 {
						buf.WriteString("\r\n")
					}
				}
			}

			buf.WriteString("\r\n--" + boundary)
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}

// ref: https://stackoverflow.com/questions/11065913/send-email-through-unencrypted-connection
type unencryptedAuth struct {
	smtp.Auth
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}
