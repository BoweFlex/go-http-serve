package main

import (
    "fmt"
    "io" // Used for echo server, maybe not needed longterm?
    "flag" // Allows handling named CLI parameters, easier to work with than os.Args
    "net" // Needed for creating a listener
    "log"
    "os"
    "os/signal"
)

func shutdownServer() {
    // This works, but may be nicer to manually log closing and close the server
    // Would need to receive *net.Listener
    // Worth noting this isn't really needed for os signal handling, but I want to use
    // the same process to handle a /shutdown page (impractical but good practice)
    panic("Shutdown request received, goodbye")
}

func main() {
    var port int
    var host string
    var network string

    flag.IntVar(&port, "port", 1234, "Provide a port number")
    flag.StringVar(&host, "host", "", "Provide a directory")
    flag.StringVar(&network, "network", "tcp", "Provide a directory")

    flag.Parse()

    fullHost := fmt.Sprintf("%v:%v", host, port)
    server, err := net.Listen(network, fullHost)

    if err != nil {
        log.Fatal(err)
    }
    defer server.Close()

    log.Printf("Listening on %v, type 'Ctrl+c' to exit", fullHost)

    // This allows a more graceful shutdown than just killing the process, but I'm not sure if I love it.
    // If I'm understanding routines right, this "server" either eats 3 threads for one request or blocks and does nothing.
    shutdownRequests := make(chan os.Signal, 1)
    signal.Notify(shutdownRequests, os.Interrupt, os.Kill)
    go func() {
        <-shutdownRequests
        shutdownServer()
    }()

    for {
        connection, err := server.Accept()
        if err != nil {
            log.Fatal(err)
        }

        // Launching a new routine allows us to handle multiple connections
        // But still only as many as num cores - 1
        go func(c net.Conn) {
            io.Copy(c, c)

            c.Close()
        }(connection)
    }

}
