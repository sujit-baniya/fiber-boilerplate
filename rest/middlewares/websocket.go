package middlewares

import (
	"errors"
	"github.com/sujit-baniya/xid"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Source @url:https://github.com/gorilla/websocket/blob/master/conn.go#L61
// The message types are defined in RFC 6455, section 11.8.
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1
	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2
	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8
	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9
	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

// Supported event list
const (
	// Fired when a Text/Binary message is received
	EventMessage = "message"
	// More details here:
	// @url https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API/Writing_WebSocket_servers#Pings_and_Pongs_The_Heartbeat_of_WebSockets
	EventPing = "ping"
	EventPong = "pong"
	// Fired on disconnection
	// The error provided in disconnection event
	// is defined in RFC 6455, section 11.7.
	// @url https://github.com/gofiber/websocket/blob/cd4720c435de415b864d975a9ca23a47eaf081ef/websocket.go#L192
	EventDisconnect = "disconnect"
	// Fired on first connection
	EventConnect = "connect"
	// Fired when the connection is actively closed from the server
	EventClose = "close"
	// Fired when some error appears useful also for debugging websockets
	EventError = "error"
)

var (
	// The addressed ws connection is not available anymore
	// error data is the uuid of that connection
	ErrorInvalidConnection = errors.New("message cannot be delivered invalid/gone connection")
)

// Raw form of websocket message
type message struct {
	// Message type
	mType int
	// Message data
	data []byte
}

// Event Payload is the object that
// stores all the information about the event and
// the connection
type EventPayload struct {
	// The connection object
	Kws *Websocket
	// The name of the event
	Name string
	// Unique connection UUID
	SocketUUID string
	// Optional websocket attributes
	SocketAttributes map[string]string
	// Optional error when are fired events like
	// - Disconnect
	// - Error
	Error error
	// Data is used on Message and on Error event
	Data []byte
}

type Websocket struct {
	// The Fiber.Websocket connection
	ws *websocket.Conn
	// Define if the connection is alive or not
	isAlive bool
	// Queue of messages sent from the socket
	queue map[string]message
	// Attributes map collection for the connection
	attributes map[string]string
	// Unique id of the connection
	UUID string
	// Wrap Fiber Locals function
	Locals func(key string) interface{}
	// Wrap Fiber Params function
	Params func(key string, defaultValue ...string) string
	// Wrap Fiber Query function
	Query func(key string, defaultValue ...string) string
	// Wrap Fiber Cookies function
	Cookies func(key string, defaultValue ...string) string
	// Deprecated: Old in-callback function
	OnConnect func()
	// Deprecated: Old in-callback function
	OnMessage func(data []byte)
	// Deprecated: Old in-callback function
	OnDisconnect func()
}

// Pool with the active connections
var pool = make(map[string]*Websocket)

// List of the listeners for the events
var listeners = make(map[string][]func(payload *EventPayload))

func New(callback func(kws *Websocket)) func(*fiber.Ctx) error {
	return websocket.New(func(c *websocket.Conn) {

		kws := &Websocket{
			ws: c,
			Locals: func(key string) interface{} {
				return c.Locals(key)
			},
			Params: func(key string, defaultValue ...string) string {
				return c.Params(key, defaultValue...)
			},
			Query: func(key string, defaultValue ...string) string {
				return c.Query(key, defaultValue...)
			},
			Cookies: func(key string, defaultValue ...string) string {
				return c.Cookies(key, defaultValue...)
			},
			queue:      make(map[string]message),
			attributes: make(map[string]string),
			isAlive:    true,
		}

		//Generate uuid
		kws.UUID = kws.createUUID()

		// register the connection into the pool
		pool[kws.UUID] = kws

		// execute the callback of the socket initialization
		callback(kws)

		kws.fireEvent(EventConnect, nil, nil)

		// fire event also via function
		if kws.OnConnect != nil {
			kws.OnConnect()
		}
		// Run the loop for the given connection
		kws.run()
	})
}

// Set a specific attribute for the specific socket connection
func (kws *Websocket) SetAttribute(key string, attribute string) {
	kws.attributes[key] = attribute
}

// Get a specific attribute from the socket attributes
func (kws *Websocket) GetAttribute(key string) string {
	return kws.attributes[key]
}

// Emit the message to a specific socket uuids list
func (kws *Websocket) EmitToList(uuids []string, message []byte) {
	for _, uuid := range uuids {
		err := kws.EmitTo(uuid, message)
		if err != nil {
			kws.fireEvent(EventError, message, err)
		}
	}
}

// Emit to a specific socket connection
func (kws *Websocket) EmitTo(uuid string, message []byte) error {

	if !connectionExists(uuid) {
		kws.fireEvent(EventError, []byte(uuid), ErrorInvalidConnection)
		return ErrorInvalidConnection
	}

	if !pool[uuid].isAlive {
		kws.fireEvent(EventError, []byte(uuid), ErrorInvalidConnection)
		return ErrorInvalidConnection
	}
	pool[uuid].Emit(message)
	return nil
}

// Broadcast to all the active connections
// except avoid broadcasting the message to itself
func (kws *Websocket) Broadcast(message []byte, except bool) {

	for uuid, _ := range pool {
		if except && kws.UUID == uuid {
			continue
		}
		err := kws.EmitTo(uuid, message)
		if err != nil {
			kws.fireEvent(EventError, message, err)
		}
	}
}

// Emit/Write the message into the given connection
func (kws *Websocket) Emit(message []byte) {
	kws.write(TextMessage, message)
}

// Actively close the connection from the server
func (kws *Websocket) Close() {
	kws.write(CloseMessage, []byte("Connection closed"))
	kws.fireEvent(EventClose, nil, nil)
	kws.isAlive = false
}

// pong writes a control message to the client
func (kws *Websocket) pong() {
	for range time.Tick(5 * time.Second) {
		kws.write(PongMessage, []byte{})
	}
}

// Add in message queue
func (kws *Websocket) write(messageType int, messageBytes []byte) {
	kws.queue[kws.randomUUID()] = message{
		mType: messageType,
		data:  messageBytes,
	}
}

// Start Pong/Read/Write functions
func (kws *Websocket) run() {

	go kws.pong()
	go kws.read()

	// every millisecond send messages from the queue
	for range time.Tick(1 * time.Millisecond) {
		if !kws.isAlive {
			break
		}
		if len(kws.queue) == 0 {
			continue
		}
		for uuid, message := range kws.queue {

			err := kws.ws.WriteMessage(message.mType, message.data)
			if err != nil {
				kws.disconnected(err)
			}
			delete(kws.queue, uuid)
		}
	}
}

// Listen for incoming messages
// and filter by message type
func (kws *Websocket) read() {
	for range time.Tick(10 * time.Millisecond) {
		if !kws.isAlive {
			break
		}
		mtype, msg, err := kws.ws.ReadMessage()

		if mtype == PingMessage {
			kws.fireEvent(EventPing, nil, nil)
			continue
		}

		if mtype == PongMessage {
			kws.fireEvent(EventPong, nil, nil)
			continue
		}

		if mtype == CloseMessage {
			kws.disconnected(nil)
			break
		}

		if err != nil {
			kws.disconnected(err)
			break
		}

		// We have a message and we fire the message event
		kws.fireEvent(EventMessage, msg, nil)
		if kws.OnMessage != nil {
			kws.OnMessage(msg)
		}
		continue
	}
}

// When the connection closes, disconnected method
// handle also the OnDisconnect() event
func (kws *Websocket) disconnected(err error) {
	kws.fireEvent(EventDisconnect, nil, err)
	kws.isAlive = false

	// Fire error event if the connection is
	// disconnected by an error
	if err != nil {
		kws.fireEvent(EventError, nil, err)
	}

	// Close the connection from the server side
	if kws.ws.Conn != nil {
		err = kws.ws.Close()

		if err != nil {
			kws.fireEvent(EventError, nil, err)
		}
	}

	// fire the specific on disconnect event defined in the initial New callback
	if kws.OnDisconnect != nil {
		kws.OnDisconnect()
	}

	// Remove the socket from the pool
	delete(pool, kws.UUID)
}

// Create random UUID for each connection
func (kws *Websocket) createUUID() string {
	uuid := kws.randomUUID()

	//make sure about the uniqueness of the uuid
	if connectionExists(uuid) {
		return kws.createUUID()
	}
	return uuid
}

func (kws *Websocket) randomUUID() string {
	return xid.New().String()
}

// Checks if there is at least a listener for a given event
// and loop over the callbacks registered
func (kws *Websocket) fireEvent(event string, data []byte, error error) {
	callbacks, ok := listeners[event]

	if ok {
		for _, callback := range callbacks {
			callback(&EventPayload{
				Kws:              kws,
				Name:             event,
				SocketUUID:       kws.UUID,
				SocketAttributes: kws.attributes,
				Data:             data,
				Error:            error,
			})
		}
	}
}

// Add listener callback for an event into the listeners list
func On(event string, callback func(payload *EventPayload)) {
	listeners[event] = append(listeners[event], callback)
}

func connectionExists(uuid string) bool {
	_, ok := pool[uuid]

	return ok
}
