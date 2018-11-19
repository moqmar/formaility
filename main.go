package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"teahub.io/momar/config"
)

var cfg = config.Open("formaility.yaml")

var dialer = gomail.NewDialer(cfg.Get("smtp.host").String(), cfg.Get("smtp.port").Int(), cfg.Get("smtp.username").String(), cfg.Get("smtp.password").String())

func form(c *gin.Context) {
	form := cfg.Get("forms").Child(c.Param("id"))
	if !form.Exists {
		c.String(400, "There is no form with this ID.\n")
		return
	}

	// TODO: File support
	c.Request.ParseForm()

	m := gomail.NewMessage()
	m.SetHeader("From", cfg.Get("smtp.from").String())
	m.SetHeader("To", form.Get("to").StringList()...)
	m.SetHeader("Subject", "New Formaility!")
	m.SetHeader("Content-Type", "text/plain; charset=UTF-8")

	msg := ""
	for f, v := range c.Request.Form {
		if f == "redirect" {
			continue
		}
		msg += fmt.Sprintf("[[[ %s ]]]\r\n%s\r\n\r\n", f, strings.Join(v, "\r\n"))
	}
	msg += fmt.Sprintf("[[[ Sender Information ]]]\r\nIP: %s\r\nUA: %s", c.Request.RemoteAddr, c.Request.UserAgent())
	m.SetBody("text/plain", msg)

	err := dialer.DialAndSend(m)
	if err != nil {
		log.Printf("%s\n", err)
		c.String(500, "An internal error occured.\n")
		return
	}

	c.Redirect(302, c.Request.Form["redirect"][0])
}

func ui(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.String(200, `<title>`+c.Param("f")+`</title><meta charset="utf-8"><form method="post">
		<input type="hidden" name="redirect" value="https://duckduckgo.com/">
		<input type="text" name="test" placeholder="test">
		<input type="submit">
	</form>`)
}

func main() {
	r := gin.Default()
	r.POST("/f/:id", form)
	r.GET("/f/:id", ui)
	err := r.Run()
	if err != nil {
		panic(err)
	}
}
