package awg

import (
	"sync"
	"sync/atomic"

	"github.com/rydesun/awesome-github/exch/github"
)

type Reporter struct {
	con        int64
	total      int64
	finished   int64
	waiting    int64
	invalid    []github.RepoID
	invalidLck sync.Mutex
}

func (r *Reporter) ConReqNum(num int) {
	atomic.StoreInt64(&r.con, int64(num))
}

func (r *Reporter) GetConReqNum() int {
	return int(atomic.LoadInt64(&r.con))
}

func (r *Reporter) TotalRepoNum(num int) {
	atomic.StoreInt64(&r.total, int64(num))
}

func (r *Reporter) GetTotalRepoNum() int {
	return int(atomic.LoadInt64(&r.total))
}

func (r *Reporter) Done() {
	atomic.AddInt64(&r.finished, 1)
}

func (r *Reporter) GetFinishedRepoNum() int {
	return int(atomic.LoadInt64(&r.finished))
}

func (r *Reporter) RepoWaiting() {
	atomic.AddInt64(&r.waiting, -1)
}

func (r *Reporter) GetWaitingRepo() int {
	return int(atomic.LoadInt64(&r.waiting))
}

func (r *Reporter) InvalidRepo(id github.RepoID) {
	r.invalidLck.Lock()
	r.invalid = append(r.invalid, id)
	r.invalidLck.Unlock()
}

func (r *Reporter) GetInvalidRepo() []github.RepoID {
	return r.invalid
}
