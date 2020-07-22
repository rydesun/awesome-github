package awg

import (
	"sync"

	"github.com/rydesun/awesome-github/exch/github"
)

type Reporter struct {
	con         int
	total       int
	finished    int
	waiting     int
	invalid     []github.RepoID
	finishedLck sync.RWMutex
	waitingLck  sync.RWMutex
	invalidLck  sync.Mutex
}

func (r *Reporter) ConReqNum(num int) {
	r.con = num
}

func (r *Reporter) GetConReqNum() int {
	return r.con
}

func (r *Reporter) TotalRepoNum(num int) {
	r.total = num
}

func (r *Reporter) GetTotalRepoNum() int {
	return r.total
}

func (r *Reporter) Done() {
	r.finishedLck.Lock()
	r.finished += 1
	r.finishedLck.Unlock()
}

func (r *Reporter) GetFinishedRepoNum() int {
	r.finishedLck.RLock()
	num := r.finished
	r.finishedLck.RUnlock()
	return num
}

func (r *Reporter) RepoWaiting() {
	r.waitingLck.Lock()
	r.waiting -= 1
	r.waitingLck.Unlock()
}

func (r *Reporter) GetWaitingRepo() int {
	r.waitingLck.RLock()
	num := r.waiting
	r.waitingLck.RUnlock()
	return num
}

func (r *Reporter) InvalidRepo(id github.RepoID) {
	r.invalidLck.Lock()
	r.invalid = append(r.invalid, id)
	r.invalidLck.Unlock()
}

func (r *Reporter) GetInvalidRepo() []github.RepoID {
	return r.invalid
}
