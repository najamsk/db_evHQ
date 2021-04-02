package utils

import (
	"gopkg.in/gomail.v2"
	_"strings"
	"fmt"
	_"bytes"
	"github.com/revel/revel"
	"html/template"
	"bytes"
	"crypto/tls"
	
	
)

type Mailer struct {
}

//BasicAuth ead these from config. and setup a func in utils to export this func
func (s *Mailer)SendEmail(emailto []string, emailcc []string, subject string, emailbody string) (error){
	fmt.Println("emailto:", emailto)
	fmt.Println("emailcc:", emailcc)
	m := gomail.NewMessage()

	//format to email address list
	addresses_to := make([]string, len(emailto))
	for i, recipient := range emailto {
		addresses_to[i] = m.FormatAddress(recipient, "")
	}
	
	m.SetHeader("To", addresses_to...)

	//format cc email address list
	fmt.Println("len(emailcc):", len(emailcc))
	if(len(emailcc)>0){
	addresses_cc := make([]string, len(emailcc))
	for i, recipient := range emailcc {
		addresses_cc[i] = m.FormatAddress(recipient, "")
	}
	
	m.SetHeader("Cc", addresses_cc...)
}
	
	m.SetHeader("From",revel.Config.StringDefault("hq.smtp.emailFrom",""))
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", emailbody)
		fmt.Println( "host:",revel.Config.StringDefault("hq.smtp.host",""))
		fmt.Println( "port:",revel.Config.IntDefault("hq.smtp.port",465))
		fmt.Println( "user:", revel.Config.StringDefault("hq.smtp.user",""))
		fmt.Println( "password:",revel.Config.StringDefault("hq.smtp.password",""))
		fmt.Println("mail from:",revel.Config.StringDefault("hq.smtp.emailFrom",""))
	d := gomail.NewDialer(revel.Config.StringDefault("hq.smtp.host",""), revel.Config.IntDefault("hq.smtp.port",465), revel.Config.StringDefault("hq.smtp.user",""), revel.Config.StringDefault("hq.smtp.password",""))
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("autherr",err)
		//panic(err)
		return err
	}
	return nil
}

func (s *Mailer)ParseEmailTemplate(data interface{}, templateName string) (string, error){

	t := template.New(templateName)
	var tBasePath = revel.Config.StringDefault("hq.emailTemplates.folderbasepath","");
	fmt.Println("hye path",tBasePath+"/"+templateName)

	var err error
	t, err = t.ParseFiles(tBasePath +"/"+ templateName)
	if err != nil {
		fmt.Println("eroor At  t.ParseFiles",err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		fmt.Println(err)
	}

	result := tpl.String()
	return result, err;
}
