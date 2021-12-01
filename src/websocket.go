package main

import (
  "fmt"
  "github.com/kataras/iris"
  "github.com/kataras/iris/websocket"
  "sync"
)

var (
 Conn = make(map[websocket.Connection]bool)
 myChatRoom = "room1"
 mutex = new(sync.Mutex)
 ConnProperties = make(map[websocket.Connection]string)
 updateStream = make(chan string)
)

func setupWebsocket(app *iris.Application) {
  // create our echo websocket server
  ws := websocket.New(websocket.Config{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
  })
  ws.OnConnection(handleConnection)

  // register the server on an endpoint.
  // see the inline javascript code in the websockets.html,
  // this endpoint is used to connect to the server.
  app.Get("/echo", ws.Handler())
  // serve the javascript built'n client-side library,
  // see websockets.html script tags, this path is used.
  app.Any("/iris-ws.js", websocket.ClientHandler())
}

func handleConnection(c websocket.Connection) {
  mutex.Lock()
  Conn[c] = true
  mutex.Unlock()
  // Read events from browser
  c.On("chat", func(msg string) {
    // Print the message to the console, c.Context() is the iris's http context.
    fmt.Printf("%s sent: %s\n", c.Context().RemoteAddr(), msg)
    // Write message back to the client message owner with:
    // c.Emit("chat", msg)
    // Write message to all except this client with:
    c.To(websocket.Broadcast).Emit("chat", msg)
  })
  c.On("properties", func(msg string) {
    fmt.Printf("%s sent: %s\n", c.Context().RemoteAddr(), msg)
    mutex.Lock()
    fmt.Printf("************************************************************** %s\n", msg)
    ConnProperties[c] = msg
    c.Join(msg)
    mutex.Unlock()
  })
  /*
  go func() { 
    for {
      select {
      case <-updateStream:
        fmt.Printf("-------------------- Sendint message\n")
        c.To(websocket.Broadcast).Emit("chat", "<img src='https://folduka.imim.science/files/logo-imim.png'>")
      }
    }
  }()
  */
}
