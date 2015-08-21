package main

import (
	"flag"
	"fmt"
	"github.com/Lilk/dialog/result"
	"log"
	"os"
	// "github.com/Lilk/dialog/core"
	// "github.com/Lilk/dialog/distr"
	"bufio"
	// "path"
	// "encoding"
	"encoding/gob"

	// "encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"
)

var dataInput = flag.String("file", "", "read samples from file")
var responseSequence = flag.Bool("t", false, "response is a sequence ")

var printData = flag.Bool("pipe", false, "pipe requests in issueorder on stdout")

type Sample struct {
	// Latency time.Duration
	IssueTime      time.Time
	SequenceNumber int
	ProcessingTime time.Duration
}

type sampleByIssue []Sample

func (slice sampleByIssue) Len() int           { return len(slice) }
func (slice sampleByIssue) Less(i, j int) bool { return slice[i].IssueTime.Before(slice[j].IssueTime) }
func (slice sampleByIssue) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }

type resultSampleByIssue []result.Sample

func (slice resultSampleByIssue) Len() int { return len(slice) }
func (slice resultSampleByIssue) Less(i, j int) bool {
	return slice[i].IssueTime.Before(slice[j].IssueTime)
}
func (slice resultSampleByIssue) Swap(i, j int) { slice[i], slice[j] = slice[j], slice[i] }

type sampleByProcessing []Sample

func (slice sampleByProcessing) Len() int { return len(slice) }
func (slice sampleByProcessing) Less(i, j int) bool {
	if *responseSequence {
		return slice[i].SequenceNumber < slice[j].SequenceNumber
	} else { //response is a timestamp
		return slice[i].ProcessingTime < slice[j].ProcessingTime
	}

}
func (slice sampleByProcessing) Swap(i, j int) { slice[i], slice[j] = slice[j], slice[i] }

func readFromFile(fileName string) (res result.Result) {
	inputFile, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	bufR := bufio.NewReader(inputFile)

	// var decoder encoding.Decoder
	// Since this is a binary format large parts of it will be unreadable

	// if path.Ext(fileName) == "json" {
	// log.Println("Reading input as JSON")
	// decoder := json.NewDecoder(bufR)
	// } else if path.Ext(fileName) == "gob" {
	// log.Println("Reading input as GOB")
	decoder := gob.NewDecoder(bufR)
	// }

	// Write to the file
	if err := decoder.Decode(&res); err != nil {
		panic(err)
	}
	inputFile.Close()
	return
}

func reorderings(res *result.Result) (int, int, int) {
	byIssue, byProcessing := make([]Sample, len(res.LatencySamples)), make([]Sample, len(res.LatencySamples))
	// copy(res.LatencySamples, byIssue)
	// copy(res.LatencySamples, byProcessing)
	for i, sample := range res.LatencySamples {
		byIssue[i].IssueTime, byProcessing[i].IssueTime = sample.IssueTime, sample.IssueTime
		if *responseSequence {
			seq, err := strconv.ParseInt(sample.Response[0:len(sample.Response)-1], 10, 32)
			if err != nil {
				panic(err)
			}
			byIssue[i].SequenceNumber, byProcessing[i].SequenceNumber = int(seq), int(seq)
		} else { //response is a timestamp
			split := strings.Split(sample.Response, ",")
			seconds, err := strconv.ParseInt(split[0], 10, 64)
			if err != nil {
				panic(err)
			}

			nanoseconds, err := strconv.ParseInt(split[1][0:len(split[1])-1], 10, 64)
			if err != nil {
				panic(err)
			}
			duration := time.Duration(int64(time.Second)*seconds + nanoseconds)
			byIssue[i].ProcessingTime, byProcessing[i].ProcessingTime = duration, duration
		}

	}

	sort.Sort(sampleByIssue(byIssue))
	sort.Sort(sampleByProcessing(byProcessing))

	totalReordering, reorderings, maxV := 0, 0, 0

	// lastIssue := byIssue[0].IssueTime
	// lastProcessing := byProcessing[0].ProcessingTime
	// lastSequence := byProcessing[0].SequenceNumber

	all_reord := make([]int, len(res.LatencySamples))

	// for i, _ := range byIssue {
	//     if byIssue[i].IssueTime.Before(lastIssue) {
	//         log.Fatal("Wrong issue sort\n")
	//     } else {
	//         if byIssue[i].IssueTime == lastIssue  && i > 0{
	//                 log.Printf("Double IssueTime: %d,  %v,  %v\n", i, byIssue[i].IssueTime, lastIssue)
	//             }
	//         lastIssue = byIssue[i].IssueTime
	//     }
	//     if ! *responseSequence {
	//         if byProcessing[i].ProcessingTime < lastProcessing {
	//             log.Fatal("Wrong processing (time) sort\n")
	//         } else {
	//             if byProcessing[i].ProcessingTime == lastProcessing  && i > 0{
	//                 log.Printf("%v %v\n", byProcessing[i].ProcessingTime, lastProcessing)
	//                 if i > 0{
	//                  log.Printf("Double ProcessingTime: %d, %v, %v\n", i, byProcessing[i], byProcessing[i-1])
	//                 }
	//             }
	//             lastProcessing = byProcessing[i].ProcessingTime
	//         }
	//     } else { //response is a timestamp
	//         if byProcessing[i].SequenceNumber < lastSequence {
	//             log.Fatal("Wrong processing (sequence) sort\n")
	//         } else {
	//             if byProcessing[i].SequenceNumber == lastSequence  && i > 0{
	//                 log.Printf("Double SequenceNumber: %d, %v, %v\n", i, byProcessing[i].SequenceNumber, lastSequence)
	//             }
	//             lastSequence = byProcessing[i].SequenceNumber
	//         }
	//     }

	// }
	// log.Println("Passed sorting tests.")

	for issueIndex, sample := range byIssue {
		processingIndex := sort.Search(len(byProcessing), func(i int) bool {
			if *responseSequence {
				return byProcessing[i].SequenceNumber >= sample.SequenceNumber
			} else { //response is a timestamp
				return byProcessing[i].ProcessingTime >= sample.ProcessingTime
			}
			// return data[i] >= 23
		})
		diff := processingIndex - issueIndex
		if diff > 0 {
			totalReordering += diff
			reorderings++
			if diff > maxV {
				maxV = diff
			}
			all_reord[issueIndex] = diff
		}

	}

	sort.Ints(all_reord)
	reord_pp := func(pp float64) int {
		return all_reord[int(pp*float64(len(all_reord)))]
	}
	log.Printf("Percentiles of reorderings, 50th: %d, 99th: %d (%d)\n", reord_pp(0.5), reord_pp(0.99), all_reord[len(all_reord)-1])

	return totalReordering, reorderings, maxV
}

func main() {
	flag.Parse()
	res := readFromFile(*dataInput)
	if *printData {
		sort.Sort(resultSampleByIssue(res.LatencySamples))
		fmt.Println(len(res.LatencySamples))
		for _, sample := range res.LatencySamples {
			fmt.Print(sample.Response)
		}
	} else {
		log.Printf("%d total responses\n", len(res.LatencySamples))
		// result.PrintResult(res, -1)
		// log.Println("First response: %v", res.LatencySamples[0].Response)

		totalReordering, reorderings, maxV := reorderings(&res)
		fmt.Printf("Reordering: %d, %d, (%d)\n", totalReordering, reorderings, maxV)
	}

}
