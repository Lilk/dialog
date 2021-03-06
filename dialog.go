package main

import (
    "flag"
    "log"
    "runtime"
    "math/rand"
    "time"
    "os"
    "runtime/pprof"
    "github.com/Lilk/dialog/result"
    "github.com/Lilk/dialog/core"
    "github.com/Lilk/dialog/distr"
    "io/ioutil"
    "strings"

)


var threads = flag.Int("t", 1, "Number of goroutines to spawn")
var server = flag.Bool("s", false, "Start in Server mode")
var commander = flag.Bool("c", false, "Start in Commander mode")
var serverfile = flag.String("serverfile", "", "List of servers to command.")

var p_addr = flag.String("address", "http://localhost", "address")
var p_rate = flag.Int("rate", 1000, "rate")
var p_duration = flag.String("time", "10s", "duration time")
var keepAlive = flag.Bool("k", true, "To keep an open connection")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var dataOutput = flag.String("outputfile", "", "write samples to file")




var clientType  = core.SimpleChunkedReader
//var clientType  = core.NetHttp


func main() {
    flag.Parse()
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }

    rand.Seed(time.Now().Unix())
    // numprocs, numprocsdef := 0, 0
    
    hw_threads := runtime.NumCPU()
    // if *threads < hw_threads {
    //  hw_threads = *threads
    // }
    // numprocsdef =
     runtime.GOMAXPROCS(hw_threads)
//  numprocsdef = runtime.GOMAXPROCS(4)
    numprocs := runtime.GOMAXPROCS(0)

    if  *server {
        distr.StartServer(clientType)
        return
    }




    log.Printf("====================================================\n")
    hostname, _ := os.Hostname()
    log.Printf("Starting go load generator with %d HW threads at %s \n", numprocs, hostname)
    duration, err := time.ParseDuration(*p_duration)
    if(err != nil){
            panic(err)
    }
    // log.Printf("Hitting %s at %d reqs/s by %d clients during %v.", *p_addr, *p_rate, *threads, duration)
    parameters := core.TestParameters{ Addr: *p_addr, Rate: float64(*p_rate),Duration: duration,Clients: *threads, SaveSamples: (*dataOutput != ""), SampleFile: *dataOutput}

    var globalResult result.Result
    if *commander {
        serverfilecontents, err := ioutil.ReadFile(*serverfile)
        if err != nil {
            panic(err)
        }
        lines := strings.Split(string(serverfilecontents), "\n")
        var servers []string
        for _, line := range lines {
            if strings.TrimSpace(line) != "" {
                servers = append(servers, line)
            }
        }

        log.Printf("Using serverfile %s with contents (%d) %v\n", *serverfile, len(servers), servers)
        globalResult = distr.StartCommander( parameters, servers, clientType)

    } else {
        globalResult = core.StartTest( parameters, clientType)
        result.PrintResult(globalResult, *threads)
        if(parameters.SaveSamples){
            result.SaveToFile(&globalResult, parameters.SampleFile)
        }

    }
    
    // if totalResponses > 0
    // log.Printf("TOTAL %d responses in %v, Rate: %v, avgLatency: %v\n", totalResponses, duration, float64(totalResponses)/duration.Seconds(), time.Duration(avgLatency))
    
}
