package main

import (
	_ "embed"
	"net/http"
	"time"

	"github.com/sandromai/go-http-server/routes"
)

//go:embed templates/emails/loginToken.min.html
var loginTokenTemplate string

func main() {
	timezone, err := time.LoadLocation("America/Sao_Paulo")

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/" {
			writer.WriteHeader(404)

			return
		}

		if request.Method != "GET" {
			writer.WriteHeader(405)

			return
		}

		writer.Write([]byte("<h1>Hello world!</h1>"))
	})

	adminRoutes := &routes.Admin{}

	http.HandleFunc("/routes/admins/login", adminRoutes.Login)
	http.HandleFunc("/routes/admins/register", adminRoutes.Register)
	http.HandleFunc("/routes/admins/update", adminRoutes.Update)

	emailSettingRoutes := &routes.EmailSetting{}

	http.HandleFunc("/routes/emailSettings/list", emailSettingRoutes.List)
	http.HandleFunc("/routes/emailSettings/update", emailSettingRoutes.Update)

	loginTokenRoutes := &routes.LoginToken{
		Template: loginTokenTemplate,
		Timezone: timezone,
	}

	http.HandleFunc("/routes/loginTokens/create", loginTokenRoutes.Create)
	http.HandleFunc("/routes/loginTokens/check", loginTokenRoutes.Check)
	http.HandleFunc("/routes/loginTokens/deny", loginTokenRoutes.Deny)
	http.HandleFunc("/routes/loginTokens/authorize", loginTokenRoutes.Authorize)

	userRoutes := &routes.User{
		Timezone: timezone,
	}

	http.HandleFunc("/routes/users/authenticate", userRoutes.Authenticate)
	http.HandleFunc("/routes/users/ban/", userRoutes.Ban)
	http.HandleFunc("/routes/users/unban/", userRoutes.Unban)

	userTokenRoutes := &routes.UserToken{
		Timezone: timezone,
	}

	http.HandleFunc("/routes/userTokens/disconnect/", userTokenRoutes.Disconnect)

	err = http.ListenAndServe(":3333", nil)

	if err != nil {
		panic(err)
	}
}
