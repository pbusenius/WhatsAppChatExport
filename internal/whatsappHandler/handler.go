package whatsappHandler

import (
	"github.com/Rhymen/go-whatsapp"
	"github.com/Rhymen/go-whatsapp/binary/proto"
	"github.com/pbusenius/WhatsAppChatExport/internal/sessions"
	"github.com/skip2/go-qrcode"
	"log"
	"time"
)

type Handler struct {
	Conn     *whatsapp.Conn
	Chats 	 map[string]struct{}
	QrCode   string
}

func (h *Handler) Login() error {
	session, err := sessions.ReadSession()
	if err == nil {
		session, err = h.Conn.RestoreWithSession(session)
		if err != nil {
			err = sessions.DeleteSession()
			if err != nil {
				return err
			}
			return h.Login()

		}
	} else {
		qr := make(chan string)
		go func() {
			qrData := <-qr
			err := qrcode.WriteFile(qrData, qrcode.Medium, 256, "qr-code.png")
			log.Printf("could not create qr-file: %v", err)
		}()
		session, err = h.Conn.Login(qr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) ShouldCallSynchronously() bool {
	return true
}

func (h *Handler) HandleRawMessage(message *proto.WebMessageInfo) {
	if message != nil && message.Key.RemoteJid != nil {
		h.Chats[*message.Key.RemoteJid] = struct{}{}
	}
}

func (h *Handler) HandleError(err error) {

	if e, ok := err.(*whatsapp.ErrConnectionFailed); ok {
		log.Printf("Connection failed, underlying error: %v", e.Err)
		log.Println("Waiting 30sec...")
		<-time.After(30 * time.Second)
		log.Println("Reconnecting...")
		err := h.Conn.Restore()
		if err != nil {
			log.Fatalf("Restore failed: %v", err)
		}
	}
}

func (h *Handler) Close() (whatsapp.Session, error) {
	session, err := h.Conn.Disconnect()
	return session, err
}


func NewWhatsappClient() (*Handler, error) {
	var wc Handler

	wac, err := whatsapp.NewConn(5 * time.Second)
	if err != nil {
		return nil, err
	}

	wc.Conn = wac

	err = wc.Conn.SetClientName("WhatsAppChatExport by SchnuckSolutions", "WhatsAppChatExport", "1.0")
	if err != nil {
		return nil, err
	}

	return &wc, nil
}
