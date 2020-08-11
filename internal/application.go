package internal

import (
	"fmt"
	"github.com/Rhymen/go-whatsapp"
	"github.com/gorilla/mux"
	"github.com/pbusenius/WhatsAppChatExport/internal/export"
	"github.com/pbusenius/WhatsAppChatExport/internal/message"
	"github.com/pbusenius/WhatsAppChatExport/internal/sessions"
	"github.com/pbusenius/WhatsAppChatExport/internal/whatsappHandler"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

type ApplicationServer struct {
	port           		string
	whatsappConnection  *whatsapp.Conn
	whatsappHandler     whatsappHandler.Handler
	router         		*mux.Router
	loggedIn       		bool
	initialChats        map[string]struct{}
}

func (s *ApplicationServer) Run() {
	srv := &http.Server{
		Handler:      s.router,
		Addr:         fmt.Sprintf("127.0.0.1:%s", s.port),

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("http server crashed: %e", err)
	}
}

func (s *ApplicationServer) initRoutes() {
	s.router.HandleFunc("/", s.indexHandler)
	s.router.HandleFunc("/login", s.loginHandler)
	s.router.HandleFunc("/chats", s.chatsFormHandler)
	s.router.HandleFunc("/chats/export", s.exportChats)

	s.router.PathPrefix("/export/").Handler(http.StripPrefix("/export/", http.FileServer(http.Dir("."))))
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web"))))
}

func (s *ApplicationServer) exportChats(w http.ResponseWriter, r *http.Request) {
	if s.loggedIn {
		var chatId string
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
		mobileNumber := r.Form["mobile"][0]
		fileTypesToExport := r.Form["files"]

		for jid := range s.initialChats {
			if strings.Contains(jid, mobileNumber) {
				chatId = jid
				break
			}
		}

		handler := &message.Handler{Conn: s.whatsappConnection, ContentType: fileTypesToExport, ChatID: chatId}
		s.whatsappConnection.LoadFullChatHistory(chatId, 300, time.Millisecond*300, handler)

		sort.SliceStable(handler.Messages, func(i, j int) bool {
			return handler.Messages[i].MessageTimestamp < handler.Messages[j].MessageTimestamp
		})

		err = serveTemplate(w, "web/templates/chats_out_online.html", handler.Messages)
		if err != nil {
			fmt.Println(err)
		}

		err = export.ToHtmlFile("chats_out.html", handler.Messages)
		if err != nil {
			fmt.Println(err)
		}

		err = export.ToJson("chat_out.json", handler.Messages)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		http.Redirect(w, r, "/login", 302)
	}
}

func (s *ApplicationServer) loginHandler(w http.ResponseWriter, r *http.Request) {
	wac, err := whatsapp.NewConn(1 * time.Second)
	if err != nil {
		log.Fatal(err)
	}
	s.whatsappConnection = wac

	err = s.whatsappConnection.SetClientName("test", "test", "1.0")
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		err = sessions.Login(s.whatsappConnection)
		if err != nil {
			fmt.Println(err)
		}
	}()

	err = serveTemplate(w, "web/templates/login.html", nil)
	if err != nil {
		fmt.Println(err)
	}

	h := whatsappHandler.Handler{Conn: s.whatsappConnection, Chats: make(map[string]struct{})}

	s.whatsappConnection.AddHandler(&h)

	<-time.After(5 * time.Second)

	s.initialChats = h.Chats

	s.loggedIn = true
}

func (s *ApplicationServer) indexHandler(w http.ResponseWriter, r *http.Request)  {
	if s.loggedIn {
		http.Redirect(w, r, "/chats", 302)
	} else {
		http.Redirect(w, r, "/login", 302)
	}
}

func (s *ApplicationServer) chatsFormHandler(w http.ResponseWriter, r *http.Request) {
	if s.loggedIn {
		var mobileNumbers []string

		for chatID := range s.initialChats {
			number := strings.Split(strings.Split(chatID, "@")[0], "-")[0]
			mobileNumbers = append(mobileNumbers, number)
		}

		err := serveTemplate(w, "web/templates/chat_form.html", struct {
			MobileNumber []string
		}{mobileNumbers})
		if err != nil {
			fmt.Println(err)
		}
	} else {
		http.Redirect(w, r, "/login", 302)
	}
}

func (s *ApplicationServer) downloadProfilePictures(messages []message.Message) error {
	chatIds := make(map[string]bool)

	for _, m := range messages {
		if !chatIds[m.MessageSender] {
			chatIds[m.MessageSender] = true
		}
	}

	for chatId := range chatIds {
		profileImage, err :=  s.whatsappConnection.GetProfilePicThumb(chatId)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Download Profile Picture!!!")
		fmt.Println(profileImage)
	}

	return nil
}

func (s *ApplicationServer) Close() error {
	_, err := s.whatsappConnection.Disconnect()
	if err != nil {
		return fmt.Errorf("could not disconnect: %v", err)
	}

	return nil
}


func serveTemplate(w http.ResponseWriter, file string, data interface{}) error {
	// TODO: convert to packr.Box -> all static files(html, css) are in the final binary
	t, err := template.ParseFiles(file)
	if err != nil {
		return err
	}

	err = t.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}


func NewApplicationServer(httpPort string) (*ApplicationServer, error) {
	var s ApplicationServer

	s.port = httpPort
	s.router = mux.NewRouter()
	s.initRoutes()

	return &s, nil
}