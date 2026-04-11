package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"example/TermCord/shared"

	"github.com/coder/websocket"
	"golang.org/x/time/rate"
)

var users []string

// type usersOnline struct {
// 	Users users `json:"online_users"`
// }

func main() {
	termcord := newChatServer()
	http.ListenAndServe(":8080", termcord)
	fmt.Print("Brewing some coffe...and a great terminal-based chat!")
}

type chatServer struct {
	// subscriberMessageBuffer controls the max number
	// of messages that can be queued for a subscriber
	// before it is kicked.
	//
	// Defaults to 16.
	subscriberMessageBuffer int

	// publishLimiter controls the rate limit applied to the publish endpoint.
	//
	// Defaults to one publish every 100ms with a burst of 8.
	publishLimiter *rate.Limiter

	// logf controls where logs are sent.
	// Defaults to log.Printf.
	logf func(f string, v ...any)

	// serveMux routes the various endpoints to the appropriate handler.
	serveMux http.ServeMux

	subscribersMu sync.Mutex
	subscribers   map[*subscriber]struct{} // todo need to understand maps better
}

func newChatServer() *chatServer {
	cs := &chatServer{
		subscriberMessageBuffer: 16,
		logf:                    log.Printf,
		subscribers:             make(map[*subscriber]struct{}),
		publishLimiter:          rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
	}
	cs.serveMux.Handle("/", http.FileServer(http.Dir(".")))
	cs.serveMux.HandleFunc("/subscribe", cs.subscribeHandler)
	cs.serveMux.HandleFunc("/publish", cs.publishHandler)

	return cs
}

// subscriber represents a subscriber.
// Messages are sent on the msgs channel and if the client
// cannot keep up with the messages, closeSlow is called.
type message struct {
	data      []byte
	publisher string
}

type subscriber struct {
	msgs      chan message
	heartbeat chan []byte
	userID    string
	closeSlow func()
}

func (cs *chatServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cs.serveMux.ServeHTTP(w, r)
}

// subscribeHandler accepts the WebSocket connection and then subscribes
// it to all future messages.
func (cs *chatServer) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	err := cs.subscribe(w, r)
	if errors.Is(err, context.Canceled) {
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
	if err != nil {
		cs.logf("%v", err)
		return
	}
}

// publishHandler reads the request body with a limit of 8192 bytes and then publishes
// the received message.
func (cs *chatServer) publishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	body := http.MaxBytesReader(w, r.Body, 8192)

	msg, err := io.ReadAll(body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
		return
	}

	cs.logf("message: %s", msg)
	cs.publish(msg, "")

	w.WriteHeader(http.StatusAccepted)
}

// subscribe subscribes the given WebSocket to all broadcast messages.
// It creates a subscriber with a buffered msgs chan to give some room to slower
// connections and then registers the subscriber. It then listens for all messages
// and writes them to the WebSocket. If the context is cancelled or
// an error occurs, it returns and deletes the subscription.
//
// It uses CloseRead to keep reading from the connection to process control
// messages and cancel the context if the connection drops.
func (cs *chatServer) subscribe(w http.ResponseWriter, r *http.Request) error {
	var mu sync.Mutex
	var c *websocket.Conn
	var closed bool
	var userID string

	userID = r.URL.Query().Get("username")

	s := &subscriber{
		msgs:   make(chan message, cs.subscriberMessageBuffer),
		userID: userID,
		closeSlow: func() {
			mu.Lock()
			defer mu.Unlock()
			closed = true
			if c != nil {
				c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
			}
		},
	}
	cs.addSubscriber(s)
	defer cs.deleteSubscriber(s)
	defer func() {
		userLeft := shared.UserLeft{UserID: s.userID}
		bytes, err := json.Marshal(userLeft)
		if err != nil {
			fmt.Printf("error marshaling userLeft")
			return
		}
		packet := shared.Packet{Type: "UserLeft", Data: bytes}
		msg, err := json.Marshal(packet)
		if err != nil {
			fmt.Printf("error marshaling packet")
			return
		}
		cs.publish(msg, s.userID)
	}()

	c2, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
	}
	mu.Lock()
	if closed {
		mu.Unlock()
		return net.ErrClosed
	}
	c = c2
	mu.Unlock()
	defer c.CloseNow()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//need to move to a onJoin() function
	userJoin := shared.UserJoined{UserID: s.userID}
	bytes, err := json.Marshal(userJoin)
	if err != nil {
		fmt.Printf("error marshaling userJoin")
	}
	packet := shared.Packet{Type: "UserJoined", Data: bytes}
	var msg []byte
	msg, err = json.Marshal(packet)
	if err != nil {
		fmt.Printf("error marshaling packet")
	}
	cs.publish(msg, s.userID)

	// Read loop in a separate goroutine
	go func() {
		defer cancel()
		for {
			_, msg, err := c.Read(ctx)
			if err != nil {
				return
			}

			var packet shared.Packet
			json.Unmarshal(msg, &packet)

			switch packet.Type {
			case "PostMsg":
				cs.publish(msg, s.userID)
				var chat shared.PostMsg
				json.Unmarshal(packet.Data, &chat)
				log.Printf("msg: %s", chat)

			case "heartbeat":
				var hb shared.HeartBeat
				json.Unmarshal(packet.Data, &hb)
				log.Printf("heartbeat: %s", hb)

			}
		}
	}()

	// Write loop stays the same
	for {
		select {
		case msg := <-s.msgs:
			if msg.publisher == s.userID {
				continue
			}
			err := writeTimeout(ctx, time.Second*5, c, msg.data)
			if err != nil {
				return err
			}
		case heartbeat := <-s.heartbeat:
			err := writeTimeout(ctx, time.Second*5, c, heartbeat)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}

}

// publish publishes the msg to all subscribers.
// It never blocks and so messages to slow subscribers
// are dropped.
func (cs *chatServer) publish(msg []byte, publisher string) {
	cs.subscribersMu.Lock()
	defer cs.subscribersMu.Unlock()

	cs.publishLimiter.Wait(context.Background())

	for s := range cs.subscribers {
		select {
		case s.msgs <- message{data: msg, publisher: publisher}:
		default:
			go s.closeSlow()
		}
	}
}

// addSubscriber registers a subscriber.
func (cs *chatServer) addSubscriber(s *subscriber) {
	cs.subscribersMu.Lock()
	cs.subscribers[s] = struct{}{}
	cs.subscribersMu.Unlock()
}

// deleteSubscriber deletes the given subscriber.
func (cs *chatServer) deleteSubscriber(s *subscriber) {
	cs.subscribersMu.Lock()
	delete(cs.subscribers, s)
	cs.subscribersMu.Unlock()
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)
}
