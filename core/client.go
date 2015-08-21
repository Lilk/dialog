package core


import (
    "github.com/Lilk/dialog/clients"
    "time"
)



type Client interface {
    Call(addr string) bool
    Request(buffer []byte) ( text []byte, ts time.Time, err error )
    Close() 
}


type ClientConstructor func() Client


func SimpleChunkedReader() Client {
    return new(clients.SimpleChunkedReader)
}

func NetHttp() Client {
    return new(clients.NetHttp)
}