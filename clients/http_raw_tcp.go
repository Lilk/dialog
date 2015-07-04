package clients

import (
    "log"
    "net"
    "fmt"
    "bufio"
    "io"
)


var request_keepalive = []byte("GET / HTTP/1.1\r\n\r\n")
var request_close = []byte("GET / HTTP/1.1\r\nConnection: close\r\n\r\n")

func Req_tcp(conn net.Conn, buffer []byte, keepAlive bool )( text []byte, err error ){

// fmt.Println("flagvar has value ", *keepAlive)
    reader := bufio.NewReader(conn)
    //reader := conn
    if keepAlive {
        // log.Println("Normal")
        //fmt.Fprintf(conn, "GET / HTTP/1.1\r\n\r\n")
        (conn).Write(request_keepalive)
    } else {
        // log.Println("Conn, close")
        //fmt.Fprintf(conn, "GET / HTTP/1.1\r\nConnection: close\r\n\r\n")
        (conn).Write(request_close)
    }


    line, err := reader.ReadSlice('\n') 
    for len(line) > 2 {
        // fmt.Printf("Reading line (%d)(%v): %s", len(line), err, line)
        line, err = reader.ReadSlice('\n') 
        // for err == io.EOF {
        //  reader.Reset(*conn)
        //  reader.ReadSlice('\n')
        // }
    }
    chars := 0
    offset := 0

    fmt.Fscanf(reader, "%x\r\n", &chars)
    // fmt.Printf("Reading %d chars.\n", chars)
    for chars > 0 {
        // _, _ := 
        io.ReadFull(reader, buffer[offset:offset+chars+2])
        // reader.Read(buffer[offset:offset+chars+2])
        // fmt.Printf("read chunk (%d):[%s]\n", offset + chars, buffer[:offset+chars])
        offset += chars + 2
        // fmt.Printf("read chunk (%d):[%s]\n", n, buffer[:n])
        _, err = fmt.Fscanf(reader, "%x\r\n", &chars)
        if(err != nil) {
            log.Println("Error in scanf from chunk size", err)
            log.Printf("Read: %s", buffer[:offset])
            return buffer[:offset], err
        }
        // fmt.Printf("Reading %d chars.\n", chars)
    }

    text = buffer[:offset]
    // fmt.Println(string(text))
    return
}
