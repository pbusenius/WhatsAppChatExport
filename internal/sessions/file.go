package sessions

import (
	"encoding/gob"
	"fmt"
	"github.com/Rhymen/go-whatsapp"
	"os"
)


const (
	SessionFile = "whatsappSession.gob"
)

func SessionFileExists() bool {
	_, err := os.Stat(SessionFile)
	if err != nil {
		return false
	}

	return true
}

func DeleteSession() error {
	return os.Remove(SessionFile)
}

func WriteSession(session whatsapp.Session) error {
	file, err := os.Create(SessionFile)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}
	return nil
}

func ReadSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open(SessionFile)
	if err != nil {
		return session, err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}
	return session, nil
}
