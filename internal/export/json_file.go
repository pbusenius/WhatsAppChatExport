package export

import (
	"encoding/json"
	"fmt"
	"github.com/pbusenius/WhatsAppChatExport/internal/message"
	"os"
)

func ToJson(outFile string, messages []message.Message) error {
	file, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	data, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
