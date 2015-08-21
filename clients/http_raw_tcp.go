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

func req_tcp(conn net.Conn, buffer []byte, request_str []byte) ( text []byte, ts time.Time, err error ){

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
        ts = time.Now()


    line, err := reader.ReadSlice('\n') 
    // fmt.Printf("///Reading line (%d)(%v): %s", len(line), err, line)

    lines_read := 1
    for len(line) > 2 || lines_read < 2{
        // fmt.Printf("Reading line (%d)(%v): %s", len(line), err, line)
        line, err = reader.ReadSlice('\n') 
        lines_read++
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
         // n, errr := 
         io.ReadFull(reader, buffer[offset:offset+chars+2])
         // fmt.Printf("Io Readfull returned: %v %v \n ", n, errr)
        // reader.Read(buffer[offset:offset+chars+2])
         // fmt.Printf("read chunk (%d):[%s]\n", offset + chars, buffer[:offset+chars])
        // fmt.Printf("read chunk (%d):[%s]\n", chars, buffer[offset:offset+chars])

        offset += chars //+ 2
        // fmt.Printf("read chunk (%d):[%s]\n", n, buffer[:n])
        _, err = fmt.Fscanf(reader, "%x\r\n", &chars)
        // fmt.Printf("Reading %d chars.\n", chars)

        if(err != nil) {
            log.Println("Error in scanf from chunk size", err)
            log.Printf("Read: %s\n", buffer[:offset])
            // return buffer[:offset], ts, err
            text = buffer[:offset]
            return
        }
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
    had_error bool 


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
    oc.had_error = false
    // fmt.Printf("Dial  %s, requesting %s\n", oc.url.Host, oc.url.Path)

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

 func  (oc *SimpleChunkedReader )  Request(buffer []byte) ( text []byte, ts time.Time,  err error ){
    if(oc.KeepaliveConn){
        // fmt.Printf("Requesting keepalive\n")
        text, ts, err = req_tcp( oc.conn, buffer, oc.request_keepalive)
        if err != nil {
            oc.had_error = true
        }
        return
    } else {
        // fmt.Printf("Requesting new conn'n")

        ts = time.Now()
        oc.conn, err = net.DialTimeout("tcp", oc.url.Host, time.Duration(300* 1000000000))
        text, _, err = req_tcp( oc.conn, buffer, oc.request_close)
        return

     }
 }
  func  (oc *SimpleChunkedReader )  Close() {
    if !oc.had_error {
        req_tcp( oc.conn, oc.buffer, oc.request_close)
    } 
    oc.conn.Close()
  }


