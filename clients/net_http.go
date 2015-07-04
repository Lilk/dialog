package clients
import (
    "net/http"
    "io/ioutil"
)
func req_client(client *http.Client, addr string) ( buffer []byte, err error ){
    resp, err := client.Get(addr)
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