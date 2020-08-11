package export

import (
	"fmt"
	"github.com/pbusenius/WhatsAppChatExport/internal/message"
	"html/template"
	"os"
)


func ToHtmlFile(htmlOut string, messages []message.Message) error {
	t, err := template.ParseFiles("web/templates/chats_out.html")
	if err != nil {
		return err
	}

    file, err := os.Create(htmlOut)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	err = t.Execute(file, messages)
	if err != nil {
		return err
	}

	return nil
}
