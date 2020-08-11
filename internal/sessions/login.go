package sessions

import (
	"fmt"
	"github.com/Rhymen/go-whatsapp"
	"github.com/skip2/go-qrcode"
)

func Login(wac *whatsapp.Conn) error {
	qr := make(chan string)
	go func() {
		qrData := <-qr
		err := qrcode.WriteFile(qrData, qrcode.Medium, 256, "web/qr-code.png")
		if err != nil {
			fmt.Println(err)
		}
	}()
	_, err := wac.Login(qr)
	if err != nil {
		return fmt.Errorf("login nicht möglich: %v\n", err)
	}

	return nil
}
