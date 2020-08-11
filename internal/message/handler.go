package message

import (
	"fmt"
	"github.com/Rhymen/go-whatsapp"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)


type Handler struct {
	Conn        *whatsapp.Conn
	ChatID      string
	Messages    []Message
	ContentType []string
}

type Message struct {
	MessageType      string `json:"messageType"`
	MessageSender    string `json:"messageSender"`
	MessageTimestamp int64 `json:"messageTimestamp"`
	MessageTime      string `json:"messageTime"`
	MessageContent   string `json:"messageContent"`
	FromMe           bool `json:"fromMe"`
}


func saveFile(filename string, data []byte) error {
	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func getFilename(messageType string, messageInfoId string, contentType string) string {
	fileType := strings.Split(strings.Split(messageType, "/")[1], ";")[0]

	filename := fmt.Sprintf("%v/%v/%v.%v", "files", contentType, messageInfoId, fileType)

	return filename
}

func setMessageFields(messageInfo whatsapp.MessageInfo, m *Message, filename string, contentType string) {
	if messageInfo.FromMe {
		m.FromMe = true
	} else{
		m.FromMe = false
	}

	m.MessageSender = messageInfo.RemoteJid
	m.MessageTimestamp = int64(messageInfo.Timestamp)
	m.MessageTime = time.Unix(m.MessageTimestamp, 0).Format("02.01.2006 15:04:05")
	m.MessageType = contentType
	m.MessageContent = filename
}


func (h *Handler) contentTypeIsWanted(fileType string) bool {
	result := false

	for _, fileWanted := range h.ContentType {
		if fileType == fileWanted {
			result = true
		}
	}

	return result
}

func (h *Handler) HandleTextMessage(message whatsapp.TextMessage) {
	if h.contentTypeIsWanted("text") {
		var m Message

		if message.Info.FromMe {
			m.MessageSender = h.Conn.Info.Wid
		} else {
			if message.Info.Source.Participant != nil {
				m.MessageSender = *message.Info.Source.Participant
			} else {
				m.MessageSender = message.Info.RemoteJid
			}
		}

		if message.Info.FromMe {
			m.FromMe = true
		} else{
			m.FromMe = false
		}

		setMessageFields(message.Info, &m, "", "text")

		m.MessageContent = message.Text

		h.Messages = append(h.Messages, m)
	} else {
		fmt.Printf("unwanted content type found: text")
	}
}

func (h *Handler) HandleImageMessage(message whatsapp.ImageMessage) {
	if h.contentTypeIsWanted("image") {
		var m Message

		data, err := message.Download()
		if err != nil {
			if err != whatsapp.ErrMediaDownloadFailedWith410 && err != whatsapp.ErrMediaDownloadFailedWith404 {
				return
			}
			if _, err = h.Conn.LoadMediaInfo(message.Info.RemoteJid, message.Info.Id, strconv.FormatBool(message.Info.FromMe)); err == nil {
				data, err = message.Download()
				if err != nil {
					return
				}
			}
		}

		filename := getFilename(message.Type, message.Info.Id, "image")

		err = saveFile(filename, data)
		if err != nil {
			fmt.Println(err)
		}

		setMessageFields(message.Info, &m, filename, "image")

		m.MessageContent = filename

		h.Messages = append(h.Messages, m)
	}
}

func (h *Handler) HandleAudioMessage(message whatsapp.AudioMessage) {
	if h.contentTypeIsWanted("audio") {
		var m Message

		data, err := message.Download()
		if err != nil {
			if err != whatsapp.ErrMediaDownloadFailedWith410 && err != whatsapp.ErrMediaDownloadFailedWith404 {
				return
			}
			if _, err = h.Conn.LoadMediaInfo(message.Info.RemoteJid, message.Info.Id, strconv.FormatBool(message.Info.FromMe)); err == nil {
				data, err = message.Download()
				if err != nil {
					return
				}
			}
		}

		filename := getFilename(message.Type, message.Info.Id, "audio")

		err = saveFile(filename, data)
		if err != nil {
			fmt.Println(err)
		}

		if message.Info.FromMe {
			m.FromMe = true
		} else{
			m.FromMe = false
		}

		setMessageFields(message.Info, &m, filename, "audio")

		h.Messages = append(h.Messages, m)

	}
}

func (h *Handler) HandleVideoMessage(message whatsapp.VideoMessage) {
	if h.contentTypeIsWanted("video") {
		var m Message

		data, err := message.Download()
		if err != nil {
			if err != whatsapp.ErrMediaDownloadFailedWith410 && err != whatsapp.ErrMediaDownloadFailedWith404 {
				return
			}
			if _, err = h.Conn.LoadMediaInfo(message.Info.RemoteJid, message.Info.Id, strconv.FormatBool(message.Info.FromMe)); err == nil {
				data, err = message.Download()
				if err != nil {
					return
				}
			}
		}

		filename := getFilename(message.Type, message.Info.Id, "video")

		err = saveFile(filename, data)
		if err != nil {
			fmt.Println(err)
			return
		}

		if message.Info.FromMe {
			m.FromMe = true
		} else{
			m.FromMe = false
		}

		setMessageFields(message.Info, &m, filename, "video")

		h.Messages = append(h.Messages, m)
	}
}

func (h *Handler) HandleDocumentMessage(message whatsapp.DocumentMessage) {
	if h.contentTypeIsWanted("document") {
		var m Message

		data, err := message.Download()
		if err != nil {
			if err != whatsapp.ErrMediaDownloadFailedWith410 && err != whatsapp.ErrMediaDownloadFailedWith404 {
				return
			}
			if _, err = h.Conn.LoadMediaInfo(message.Info.RemoteJid, message.Info.Id, strconv.FormatBool(message.Info.FromMe)); err == nil {
				data, err = message.Download()
				if err != nil {
					return
				}
			}
		}

		filename := getFilename(message.Type, message.Info.Id, "document")

		err = saveFile(filename, data)
		if err != nil {
			fmt.Println(data)
		}

		if message.Info.FromMe {
			m.FromMe = true
		} else{
			m.FromMe = false
		}

		setMessageFields(message.Info, &m, filename, "document")

		h.Messages = append(h.Messages, m)
	}
}

func (h *Handler) HandleError(err error) {
	return
}