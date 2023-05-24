package http

import "log"

func InitHttpService() {
	router := NewRouter()
	go func() {
		if err := router.Run(":9999"); err != nil {
			log.Fatal("failed to start http service:", err)
		}
	}()
}
