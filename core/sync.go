package core
import (
    "sync"
)

type Sync struct {
    done sync.WaitGroup   
    ready sync.WaitGroup
    start sync.WaitGroup
}


func (sync *Sync) signalReady() {
    sync.ready.Done()
}

func (sync *Sync) WaitReady() {
    sync.ready.Wait()
}
func (sync *Sync) signalDone() {
    sync.done.Done()
}
func (sync *Sync) WaitDone() {
    sync.done.Wait()
}
func (sync *Sync) Go() {
    sync.start.Done()
}

func (sync *Sync) awaitGo() {
    sync.start.Wait()
}

func NewSync( clients int) (sync Sync) {
    sync.done.Add(clients)
    sync.ready.Add(clients)
    sync.start.Add(1)
    return
}