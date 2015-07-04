package distr

import (
    "log"
    "net"
    "bufio"
    "encoding/json"
    "github.com/Lilk/dialog/result"
    "github.com/Lilk/dialog/core"
)


func handle_client(conn net.Conn) {
    reader := bufio.NewReader(conn)
    dec := json.NewDecoder(reader)
    enc := json.NewEncoder(conn)

    var p core.TestParameters
     if err := dec.Decode(&p); err != nil {
        log.Println(err)
        return
     }
        log.Println("Received test parameters from commander.")
        // log.Printf("Hitting %s at %f reqs/s by %d clients during %v.", p.Addr, p.Rate, p.Clients, p.Duration)


        globalResult, sync := core.Spawn_workers(p)
        
        sync.WaitReady() // ready.Wait();
        conn.Write([]byte("DONE\n"))
        log.Println("Sent DONE notification to commander.")
        line, err := reader.ReadString('\n');
        if err != nil || line != "GO\n"{
            log.Print("Did not receive GO command:", line, err, "\n")
            return
        }
        log.Println("Received go command from commander, starting loading")
        sync.Go() // start.Done()
        
        sync.WaitDone() //wg.Wait()

        summary := result.ResultSummary{p.Clients, globalResult.AverageThroughput(), globalResult.N_errors}
        enc.Encode(summary)
        result.PrintResult(*globalResult, p.Clients)
        log.Printf("Send result: %v\n", summary)
        conn.Close()
        log.Println("Closed connection to master.\n")
}

func Start_server(){
    service := ":9988"
    tcpAddr, error := net.ResolveTCPAddr("tcp", service)
    if error != nil {
        log.Println("Error: Could not resolve address")
    } else {
        netListen, error := net.Listen(tcpAddr.Network(), tcpAddr.String())
        if error != nil {
            log.Println(error)
        } else {
            defer netListen.Close()
 
            for {
                log.Println("Waiting for a client.")
                conn, error := netListen.Accept()
                if error != nil {
                    log.Println("Client error: ", error)
                } else {
                    handle_client(conn)                    
                }
            }
        }
    }
}