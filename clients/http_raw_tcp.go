package clients

import (
    "log"
    "net"
    "fmt"
    "bufio"
    "io"
    "net/url"
    "time"

)


// var request_keepalive = []byte("GET / HTTP/1.1\r\n\r\n")
// var request_close = []byte("GET / HTTP/1.1\r\nConnection: close\r\n\r\n")

func req_tcp(conn net.Conn, buffer []byte, request_str []byte)( text []byte, err error ){

// fmt.Println("flagvar has value ", *keepAlive)
    reader := bufio.NewReader(conn)
    //reader := conn
    // if keepAlive {
    //     // log.Println("Normal")
    //     //fmt.Fprintf(conn, "GET / HTTP/1.1\r\n\r\n")
    //     (conn).Write(request_keepalive)
    // } else {
    //     // log.Println("Conn, close")
    //     //fmt.Fprintf(conn, "GET / HTTP/1.1\r\nConnection: close\r\n\r\n")
    //     (conn).Write(request_close)
    // }
    // 
    
    (conn).Write(request_str)


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

type SimpleChunkedReader struct {
    conn net.Conn
    url url.URL
    buffer []byte
    KeepaliveConn bool
    request_keepalive []byte
    request_close     []byte


}

func  ( oc *SimpleChunkedReader )   Call(addr string) bool {
    var err error
    url, err := url.Parse(addr)
    if err != nil {
        log.Fatal("Couldn't parse URL: ", err)
    }
    oc.url = *url 
    oc.buffer = make([]byte, 128, 128)
    oc.KeepaliveConn = true
    fmt.Printf("Dial  %s, requesting %s\n", oc.url.Host, oc.url.Path)

    oc.request_keepalive = []byte(fmt.Sprintf("GET %s HTTP/1.1\r\n\r\n", oc.url.Path))
    oc.request_close = []byte(fmt.Sprintf("GET %s HTTP/1.1\r\nConnection: close\r\n\r\n", oc.url.Path) )


    oc.conn, err = net.DialTimeout("tcp", oc.url.Host, time.Duration(300* 1000000000))
     if err != nil {
        log.Fatal("Couldn't Dial ", oc.url.Host, err)
    }
    if err != nil {
        return false
    }

    return true

}

 func  (oc *SimpleChunkedReader )  Request(buffer []byte) ( text []byte, err error ){
    if(oc.KeepaliveConn){
        return req_tcp( oc.conn, buffer, oc.request_keepalive)
    } else {
        oc.conn, err = net.DialTimeout("tcp", oc.url.Host, time.Duration(300* 1000000000))
        return req_tcp( oc.conn, buffer, oc.request_close)

     }
 }
  func  (oc *SimpleChunkedReader )  Close() {
    req_tcp( oc.conn, oc.buffer, oc.request_close)
    oc.conn.Close()
  }


