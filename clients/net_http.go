package clients
import (
    "net"
    "net/http"
    "io/ioutil"
    "time"
)
func req_client(client *http.Client, addr string) ( buffer []byte, ts time.Time, err error ){
    resp, err := client.Get(addr)
    ts = time.Now()
    // checkerr(err)
    buffer, err = ioutil.ReadAll(resp.Body)
    resp.Body.Close()
    return
}
    

    // secs := 30
    // tr := &http.Transport{
    //  DisableCompression: true,
    //  Dial: (&net.Dialer{
    //            Timeout:   time.Duration(secs*1000000000),
    //            KeepAlive: time.Duration(30*1000000000),
    //          }).Dial,
    // }
    // client := &http.Client{Transport: tr}
    // 
    


    // reader := bufio.NewReader(conn)


    // resp, err := client.Get(addr)
    // if  err != nil {
    //              log.Fatal(err)
    //              wg.Done()
    //              return

    // }
    // ioutil.ReadAll(resp.Body)
    // req_client(client, addr)
    // 
    // 
type NetHttp struct {
   tr *http.Transport
   client *http.Client
   // url net.URL
   addr string
}

func  ( nh *NetHttp )   Call(addr string) bool {
     secs := 30
     nh.tr = &http.Transport{
     DisableCompression: true,
     Dial: (&net.Dialer{
               Timeout:   time.Duration(secs*1000000000),
               KeepAlive: time.Duration(30*1000000000),
             }).Dial,
    }
    nh.client = &http.Client{Transport: nh.tr}
    // url, _ = url.Parse(addr)
    // nh.url = *url
    nh.addr = addr
    return true

}

 func  (nh *NetHttp )     Request(buffer []byte) ( text []byte, ts time.Time, err error ){
    text, ts, err = req_client(nh.client, nh.addr)
    return
 }
  func  (nh *NetHttp )  Close() {
    nh.tr.CloseIdleConnections()
  }