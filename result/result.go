package result

import(
"sort"
"log"
"fmt"
"time"
)

type Result struct {
    TotLatency time.Duration
    N_latencySamples int
    LatencySamples []Sample //[]time.Duration
    Duration time.Duration
    N_errors int

    sorted bool
}
type Sample struct {
    Latency time.Duration
    IssueTime time.Time
    Response string
}


type ResultSummary struct {
    Clients int
    AverageThroughput float64
    N_errors int
}

// func (res *Result) AddLatencySample(latency time.Duration){
//     sample := Sample{latency}
//     res.TotLatency += sample.Latency
//     res.N_latencySamples++
//     res.LatencySamples = append(res.LatencySamples, sample)
//     res.sorted = false
// }

func (res *Result) AddSample(sample Sample){
    res.TotLatency += sample.Latency
    res.N_latencySamples++
    res.LatencySamples = append(res.LatencySamples, sample)
    res.sorted = false
}

func (res *Result) NumberOfSamples() int{
    return res.N_latencySamples
}
func (res *Result) AverageThroughput() float64 {
    return float64(res.NumberOfSamples())/res.Duration.Seconds()
}
func (res *Result) CombineWith(other Result){
    res.TotLatency += other.TotLatency
    res.N_latencySamples += other.N_latencySamples
    res.LatencySamples = append(res.LatencySamples, other.LatencySamples...)
    res.sorted = false

    res.N_errors += other.N_errors
}
type durationSlice []time.Duration

func (slice durationSlice) Len() int {return len(slice)}
func (slice durationSlice) Less(i, j int) bool { return slice[i] < slice[j] }
func (slice durationSlice) Swap(i, j int)  { slice[i], slice[j] = slice[j], slice[i] }

type sampleSlice []Sample
func (slice sampleSlice) Len() int {return len(slice)}
func (slice sampleSlice) Less(i, j int) bool { return slice[i].Latency < slice[j].Latency }
func (slice sampleSlice) Swap(i, j int)  { slice[i], slice[j] = slice[j], slice[i] }

func (res *Result) AverageLatency() time.Duration {
    return time.Duration(res.TotLatency.Nanoseconds()/int64(res.N_latencySamples))
}
func (res* Result) Sort() {
    sort.Sort(sampleSlice(res.LatencySamples))
    res.sorted = true
}
func (res* Result) Percentile(percentile float64) time.Duration {
    if !res.sorted {
        res.Sort()
    }
    position := int(percentile * float64(res.N_latencySamples))
    return res.LatencySamples[position].Latency
}

func NewResult(duration time.Duration, rate float64) (result Result) {
    expectedSamples := int(rate * duration.Seconds())
    result.Duration = duration
    result.TotLatency = 0
    result.N_latencySamples = 0
    result.LatencySamples = make([]Sample, 0, 2*expectedSamples) //time.Duration

    result.N_errors = 0
    return
}

func  Microseconds(d time.Duration) int64{
    return (d.Nanoseconds()+500)/1000
}


func PrintResult(result Result, clients int){
    duration := result.Duration
    log.Printf("TOTAL %d responses in %v, Rate: %v, avgLatency: %v\n", result.NumberOfSamples(), duration, result.AverageThroughput(), result.AverageLatency())
    result.Sort()
    log.Printf("Latency distribution:\n")
    tails := []float64{0.5, 0.9, 0.99}
    for _, v := range tails {
        log.Printf("\t %d: \t%v\n", int(v*100), result.Percentile(v))
    }
    fmt.Print(clients, "\t")
    fmt.Printf("%.1f\t%d", result.AverageThroughput(), Microseconds(result.AverageLatency()))
    for _, v := range tails {
        fmt.Printf("\t%d",Microseconds(result.Percentile(v)))
    }
    fmt.Printf("\t%d\n", result.N_errors)
}