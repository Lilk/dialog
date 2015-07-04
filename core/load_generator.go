
package core

import (
    "log"
    "math/rand"
    "time"
    "sync"
    "net"
    "github.com/Lilk/dialog/result"
    "github.com/Lilk/dialog/clients"
)


// var wg sync.WaitGroup   
// var ready sync.WaitGroup
// var start sync.WaitGroup
// 
var globalSync Sync

var mt sync.Mutex

type TestParameters struct {
    Addr string
    Rate float64
    Clients int
    Duration time.Duration
}


func printStarted(){
    // start.Wait()
    globalSync.start.Wait()
    log.Printf("Load simulation started\n")
}



func Start_test( p TestParameters) result.Result {
   

    globalResult, _ := Spawn_workers(p)
    
    globalSync.WaitReady()
    globalSync.Go()
    globalSync.WaitDone()

    // ready.Wait();
    // start.Done()



    // wg.Wait()
    return *globalResult
}


func time_to_deadline(r *rand.Rand, desiredRate float64) (next_deadline time.Duration) {
    //ExpFloat64 returns an exponentially distributed float64 in the range (0, +math.MaxFloat64]
    // with an exponential distribution whose rate parameter (lambda)
    // is 1 and whose mean is 1/lambda (1) from the default Source. 
    //To produce a distribution with a different rate parameter, callers can adjust the output using:
    // sample := rand.ExpFloat64() / desiredRateParameter

    // desiredRateNanos := desiredRate/10e8
    sample := r.ExpFloat64() / desiredRate
    return  time.Duration (int64( sample * 10e8))
}

func checkerr(err error) {
    if err != nil {
        //panic(err)
        log.Fatal(err)
        panic(err)
    }

}


func client(addr string, rate float64, duration time.Duration, globalResult *result.Result) {
    localResult := result.NewResult(duration, rate)

    buffer := make([]byte, 65536, 65536)
    conn, err := net.DialTimeout("tcp", addr[7:len(addr)-1], time.Duration(300* 1000000000))
    checkerr(err)
    
    _, err = clients.Req_tcp( conn, buffer, true) //*keepAlive
    checkerr(err)


    // log.Printf("Client target rate %v\n", rate)
    loop_time  := time.Duration(50000)

    random_generator := rand.New(rand.NewSource(time.Now().UnixNano()))

    next_deadline := time_to_deadline(random_generator, rate) - loop_time

    globalSync.signalReady()//ready.Done()
    globalSync.awaitGo() // start.Wait()


    start_time := time.Now()
    t_read := start_time
    t_last_read := t_read
    had_error := false

    
    for( time.Since(start_time) < duration ){
        if(start_time.Add(duration).Before(t_read.Add(next_deadline))){
            time.Sleep(start_time.Add(duration).Sub(t_read))
            break
        }
        time.Sleep(next_deadline);

        t_start := time.Now()
        if( true  || had_error){ //!*keepAlive
            // log.Println("Opening new connection")
            conn.Close()
            conn, err = net.DialTimeout("tcp", addr[7:len(addr)-1], time.Duration(30* 1000000000))
            checkerr(err)
            had_error = false
            // reader = bufio.NewReader(conn)
        }
        _, err = clients.Req_tcp( conn, buffer, true)//*keepAlive
        if(err != nil){
            log.Println("Read ERROR:", err)
            localResult.N_errors++
            had_error = true
            // wg.Done()
            // return
        }
        

        t_last_read = t_read
        t_read = time.Now()

        time_spent := t_read.Sub(t_start)
        // req_time += time_spent
        localResult.AddSample(time_spent)


        this_loop := t_read.Sub(t_last_read) - time_spent - next_deadline
        loop_time = time.Duration( int64(  0.95 * float64( loop_time.Nanoseconds() ) + 0.05* float64(this_loop.Nanoseconds())))

        next_deadline = time_to_deadline(random_generator, rate) - time_spent - loop_time
        
        //log.Printf("Got response in %v, %v \n", t_read.Sub(t_start), t_int.Sub(t_start))
        // log.Printf("Got response in %v, %v \n%s", t_read.Sub(t_start), t_int.Sub(t_start), body)
    }

    clients.Req_tcp( conn, buffer, false)
    conn.Close()

    // fmt.Printf("Loop time: %v\n", loop_time)
    mt.Lock()
    globalResult.CombineWith(localResult)
    mt.Unlock()
    // log.Printf("Got %d responses in %v (%v), Rate: %v\n", responses, time.Since(start_time), req_time, float64(responses)/duration.Seconds())
    
    globalSync.signalDone() // wg.Done()

}

func Spawn_workers( p TestParameters ) (*result.Result, *Sync) {
    

    var addr, rate, duration , clients =  p.Addr, p.Rate, p.Duration, p.Clients
    log.Printf("Hitting %s at %f reqs/s by %d clients during %v.", addr, rate, clients, duration)

    
    globalSync = NewSync(clients)

    batches := clients / 1024
    globalResult := result.NewResult(duration, rate)

    thread_rate := float64(rate) / float64(clients)
    go printStarted()
    start_threads := func (n_threads int, rate float64){
        for i:=0; i < n_threads; i++ {
            go client(addr, rate , duration, &globalResult)
        }
        return
    }
    for i := 0; i < batches; i++{
        start_threads(1024, thread_rate)
        time.Sleep( time.Duration( 100* 1000000  ) ) //100 millisecs
    }
    start_threads( clients % 1024, thread_rate )
    return &globalResult, &globalSync
}