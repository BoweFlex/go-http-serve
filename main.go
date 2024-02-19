package main

import (
    "fmt"
    "io"
    "flag" // Allows handling named CLI parameters, easier to work with than os.Args
    "net"
    "log"
)

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
