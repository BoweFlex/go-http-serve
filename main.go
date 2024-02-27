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

func handleRequest(conn net.Conn, endpoint string, success, shutdown chan<- bool) {
    if endpoint == "/shutdown" {
        conn.Write([]byte("Shutdown request received, goodbye!"))
        shutdown <- true
    } else {
        io.Copy(conn, conn)
        conn.Close()
        success <- true
    }
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
    osShutdown := make(chan os.Signal, 1)
    signal.Notify(osShutdown, os.Interrupt, os.Kill)

    success, shutdown := make(chan bool), make(chan bool)
    for {
        go func() {
            select {
            case <-osShutdown:
                break
            case <-shutdown:
                break
            case <-success:
            }
        }()
        connection, err := server.Accept()
        if err != nil {
            log.Fatal(err)
        }

        // Launching a new routine allows us to handle multiple connections
        // But still only as many as num cores - 1
        /* This currently doesn't work because endpoint is not defined.
        I'm having trouble finding out how to handle this without using the http package,
        so I will do some more reading and either transition to using that 
        or figure out a way to make this work without it. */
        go handleRequest(connection, endpoint, success, shutdown)
    }

}
