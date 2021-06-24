package main

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func setDarkPreferences(c echo.Context) (err error) {
	set := c.Request().URL.Query().Get("set")
	back := c.Request().URL.Query().Get("back")

	cookie, err := c.Request().Cookie("dark")
	if err != nil && err != http.ErrNoCookie {
		return err
	}

	if (set == "true" || set == "false" || set == "reset") && cookie == nil {
		cookie = &http.Cookie{
			Name: "dark",
		}
	}

	if set != "" {
		switch set {
		case "true":
			{
				cookie.Value = "true"
				break
			}
		case "false":
			{
				cookie.Value = "false"
				break
			}
		case "reset":
			{
				cookie.Value = ""
				cookie.Expires = time.Now()
			}
		}
	}

	if cookie != nil {
		log.Println("writing cookie: " + cookie.Value)
		c.SetCookie(cookie)
	}

	if back == "" {
		back = "/"
	}

	return c.Redirect(302, back)
}
