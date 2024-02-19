package main

import (
    "io"
    "net"
    "log"
)

func main() {
    server, err := net.Listen("tcp", ":1234")
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
