package main

import (
	"fmt"
	"github.com/pbusenius/WhatsAppChatExport/internal"
	"log"
)

func main() {
	application, err := internal.NewApplicationServer("8080")
	if err != nil {
		log.Fatalf("could not create application: %v", err)
	}
	defer func() {
		err := application.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("WhatsAppChatExport is running on http://localhost:8080/")

	application.Run()
}



