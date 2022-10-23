package utility

import (
	"log"

	"github.com/Kucoin/kucoin-go-sdk"
)

func WsTargetSymbol(s *kucoin.ApiService, targetSymbol string, <- string, ) {
	rsp, err := s.WebSocketPublicToken()
	if err != nil {
		// Handle error
		return
	}

	tk := &kucoin.WebSocketTokenModel{}
	if err := rsp.ReadData(tk); err != nil {
		// Handle error
		return
	}

	c := s.NewWebSocketClient(tk)

	mc, errorChan, err := c.Connect()
	if err != nil {
		// Handle error
		return
	}
	allTicker := kucoin.NewSubscribeMessage("/market/ticker:$targetSymbol", false)
	closeCh1 := kucoin.NewUnsubscribeMessage("/market/ticker:$targetSymbol", false)

	if err := c.Subscribe(allTicker); err != nil {
		// Handle error
		return
	}

	var i = 0
	for {
		select {
		case err := <-errorChan:
			c.Stop() // Stop subscribing the WebSocket feed
			log.Printf("Error: %s", err.Error())
			// Handle error
			return
		case msg := <-mc:
			// log.Printf("Received: %s", kucoin.ToJsonString(m))
			t := &kucoin.TickerLevel1Model{}
			if err := msg.ReadData(t); err != nil {
				log.Printf("Failure to read: %s", err.Error())
				return
			}
			// form data and provide to channel
		
		case closeMsg := <-closeChan:
				log.Println("Unsubscribe from target cpin")
				if err = c.Unsubscribe(closeCh1); err != nil {
					log.Printf("Error: %s", err.Error())
					// Handle error
					return
				}
			}
		}
	}
}
