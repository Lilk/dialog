package distr


import (
    "log"
    "time"
    "net"
    "sync"
    "fmt"
    "bufio"
    "encoding/json"
    "github.com/Lilk/dialog/result"
    "github.com/Lilk/dialog/core"
)


func serverHandler(serverAddr string, p core.TestParameters, serverReady, serverStart *sync.WaitGroup, ret chan result.ResultSummary ){
    conn, err := net.Dial("tcp", serverAddr)
    defer conn.Close()

    if err != nil {
        log.Fatal("Couldn't call server",  err)
    }

    reader := bufio.NewReader(conn)
    dec := json.NewDecoder(reader)
    enc := json.NewEncoder(conn)
    enc.Encode(p)


    line, err := reader.ReadString('\n');
    if err != nil || line != "DONE\n" {
        log.Fatal("Did not receive DONE ACK: ", line, err, "\n")
    } 
    serverReady.Done()
    log.Printf("%s done spawning connections, waiting for rest of servers\n", serverAddr)
    serverStart.Wait()
    log.Println("Sending GO")

    conn.Write([]byte("GO\n"))
    var summary result.ResultSummary

    if err := dec.Decode(&summary); err != nil {
        log.Println("Ans from %s:", serverAddr, err)
    }

    log.Printf("From %s: %v\n", serverAddr, summary)

    ret <- summary

    

}

func StartCommander(p core.TestParameters, cc core.ClientConstructor ) result.Result  {
    servers := []string{"192.168.7.1:9988","192.168.8.1:9988","192.168.9.1:9988", "192.168.10.1:9988"}
    if p.Clients - 4 <= 0 || p.Rate - 1000 <= 0 {
        res := core.StartTest(p, cc)
        result.PrintResult(res, p.Clients)
        return res
    }
    localParams := p
    localParams.Clients, localParams.Rate = 4, 1000;

    

    p.Clients, p.Rate = (p.Clients - localParams.Clients)/len(servers), (p.Rate - localParams.Rate)/float64(len(servers)) ;
    log.Printf("Params to send: %v\n", p)
    //TODO : divide P properlu if #servers > 1
    var serverReady, serverStart sync.WaitGroup
    

    serverReady.Add(len(servers))
    serverStart.Add(1)

    returnChannel := make(chan result.ResultSummary)
    for _, address := range servers{
        time.Sleep(time.Duration(30* 1000000))
        log.Println("spawning", address)
        go serverHandler(address, p, &serverReady, &serverStart, returnChannel )
    }

    // wg.Add(localParams.Clients)
    // ready.Add(localParams.Clients)
    // start.Add(1)
    fmt.Printf("localParams %v\n", localParams)
    globalResult, localSync := core.SpawnWorkers(cc, localParams)
    
    localSync.WaitReady()// ready.Wait()
    serverReady.Wait()
    serverStart.Done()

    time.Sleep(time.Duration(100* 1000)) //Sleep 50 µs ~ 1/2 RTT
    //start.Done()
    // wg.Wait()
    localSync.Go()
    localSync.WaitDone()


    rate := globalResult.AverageThroughput()
    clients := localParams.Clients
    errors := globalResult.N_errors

    log.Printf("LOCAL: %f (%d) \n", rate, globalResult.N_latencySamples)

    for range servers{
        summary := <-returnChannel
        rate  += summary.AverageThroughput
        clients += summary.Clients
        errors += summary.N_errors
    }

    // printResult(*globalResult, localParams.Clients)
    // fmt.Println()


    fmt.Print(clients, "\t")
    globalResult.Sort()
    tails := []float64{0.5, 0.9, 0.99}
    fmt.Printf("%.1f\t%d",  rate, result.Microseconds(globalResult.AverageLatency()))
    for _, v := range tails {
        fmt.Printf("\t%d",result.Microseconds(globalResult.Percentile(v)))
    }
    fmt.Printf("\t%d\n", errors)
    return *globalResult

}
