package crawl_vp

import(
    "net/smtp"
    "text/template"
    "bytes"
    "log"
    "fmt"
    "encoding/json"
    "strconv"
    "time"
    "sort"
)

type EmailUser struct {
  Username    string
  Password    string
  EmailServer string
  Port        int
}

type SmtpTemplateData struct {
  From    string
  To      string
  Subject string
  Body    string
}


const emailTemplate = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}

{{.Body}}

Sincerely,

{{.From}}
`

func SendChangedAndNew(changedWines ListOfWines,newWines ListOfWines){
    var doc bytes.Buffer 
    var err error   
    fmt.Println("Jeg er kuul")
    emailUser:=&EmailUser{Username:"crawlbeerandwine@gmail.com",Password:"CrawlBonanza", // change password to make it work
        EmailServer:"smtp.gmail.com",Port:587}
    auth := smtp.PlainAuth("",
      emailUser.Username,
      emailUser.Password,
      emailUser.EmailServer)


    t := template.New("emailTemplate")
    t, err = t.Parse(emailTemplate)
    if err != nil {
      log.Print("error trying to parse mail template")
    }

    sort.Sort(changedWines)

    b1,_:=json.MarshalIndent(changedWines,"","  ")
    b2,_:=json.MarshalIndent(newWines,"","  ")
    chW:=string(b1)
    nW:=string(b2)
    bdy:="Viner med prisforandring:\n\n"+chW+"\n\n\n\n\n\n\n\n\n\nNye viner:\n\n"+nW

    
    
    
    context := &SmtpTemplateData{
      "SmtpEmailSender",
      "erlend.aune.1983@gmail.com",
      "Viner med prisforandring og nye viner: "+time.Now().String(),
      bdy}
    err = t.Execute(&doc, context)
    if err != nil {
      log.Print("error trying to execute mail template")
    }
    //fmt.Println(string(doc.Bytes()))

    err = smtp.SendMail(emailUser.EmailServer+":"+strconv.Itoa(emailUser.Port), // in our case, "smtp.google.com:587"
        auth,
        emailUser.Username,
        []string{"erlend.aune.1983+beer@gmail.com"},
        doc.Bytes())

    if err != nil {
      log.Print("ERROR: attempting to send a mail ", err)
    }
}









