package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sandromai/go-http-server/routes"
	"github.com/sandromai/go-http-server/utils"
)

func main() {
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

	http.HandleFunc("/routes/admins/login", routes.AdminLogin)
	http.HandleFunc("/routes/admins/register", routes.AdminRegister)
	http.HandleFunc("/routes/admins/update", routes.AdminUpdate)

	http.HandleFunc("/routes/emailSettings/update", routes.EmailSettingsUpdate)
	http.HandleFunc("/routes/emailSettings/list", routes.EmailSettingsList)

	http.HandleFunc("/routes/loginTokens/create", routes.LoginTokenCreate)

	err := http.ListenAndServe(":3333", nil)

	if err != nil {
		currentPath, err := os.Executable()

		if err != nil {
			panic(err)
		}

		appErr := (&utils.Logger{
			FolderPath:     filepath.Join(filepath.Dir(currentPath), "logs"),
			FileName:       "main.log",
			MessagesPrefix: "[MAIN]:",
		}).Save(fmt.Sprintf("Server closed: %v", err.Error()))

		if appErr != nil {
			panic(appErr.Message)
		}
	}
}
