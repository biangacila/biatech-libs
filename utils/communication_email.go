package utils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/biangacila/luvungula-go/global"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type EmailRequest struct {
	Org        string
	From       string
	To         string
	Subject    string
	Body       string
	Files      []string
	ReplyTo    string
	SenderName string
}
type Email struct {
	Username    string
	Sender      string
	Receiver    string
	Subject     string
	Body        string
	ReplyTo     string
	FromEmail   string
	FromCompany string

	HasAttached  string
	Attaches     string
	AttachedType string

	SmtpInfo SmtpInfo
}

const (
	DirTempAttachedDownload = "download-email-files"
)

func (obj *Email) Send() error {
	smtpInfo := obj.SmtpInfo
	/*
		smtpInfo, err := GetSmtpInfo()
		global.DisplayObject("smtpInfo", smtpInfo)
		if err != nil {
			return err
		}*/
	port, _ := strconv.Atoi(smtpInfo.Port)
	d := gomail.NewDialer(smtpInfo.Host, port, smtpInfo.Username, smtpInfo.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	m := gomail.NewMessage(gomail.SetCharset("ISO-8859-1"), gomail.SetEncoding(gomail.Base64))

	fromName := fmt.Sprintf("%s<%s>", obj.FromCompany, obj.FromEmail)

	fmt.Println("fromName :)--> ", fromName)

	m.SetHeaders(map[string][]string{
		"From":    {fromName},
		"To":      {obj.Receiver},
		"Subject": {obj.Subject},
	})
	m.SetBody("text/html", obj.Body)
	/*
		LET MAKE ATTACHMENT REQUE
	*/

	global.CreateFolderUploadIfNotExist(DirTempAttachedDownload)

	var lsitdelete []string
	attachedList := strings.Split(obj.Attaches, ";")

	if obj.HasAttached == "yes" {
		if obj.AttachedType == "local" {
			for _, filename := range attachedList {
				m.Attach(filename)
			}
		}

		fmt.Println("obj.AttachedType>>>> ", obj.AttachedType)

		if obj.AttachedType == "external" {
			for _, fileUrl := range attachedList {
				mybite, filename := global.GetHttpFileContent2(fileUrl)
				attFile := fmt.Sprintf("./%s/%s", DirTempAttachedDownload, filename)
				ioutil.WriteFile(attFile, mybite, 0644)
				m.Attach(attFile)
				timer := time.NewTimer(1 * time.Second)
				lsitdelete = append(lsitdelete, attFile)
				<-timer.C
			}
		}
	}

	err := d.DialAndSend(m)
	if err != nil {
		return err
	}
	fmt.Println("@-=> EMAIL SEND REPORT : ", err)

	timer2 := time.NewTimer(time.Second * 60)
	go func() {
		<-timer2.C
		for _, fname := range lsitdelete {
			var err = os.Remove(fname)
			global.CheckError(err)
		}
	}()

	return nil
}
func SendSampleEmail(in EmailRequest) error {
	hub := Email{}
	hub.FromEmail = "easidoc@easipath.com"
	if in.From != "" {
		hub.FromEmail = in.From
	}
	hub.Receiver = in.To
	hub.Subject = in.Subject
	hub.Body = in.Body
	hub.ReplyTo = "easidoc@easipath.com"
	if in.ReplyTo != "" {
		hub.ReplyTo = in.ReplyTo
	}
	hub.Sender = in.SenderName
	hub.FromCompany = strings.ToUpper(in.SenderName)

	str, _ := json.Marshal(hub)
	fmt.Println(string(str))
	return hub.Send()
}
