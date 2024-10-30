package e2e

import (
	"context"
	"github.com/glossd/pokergloss/table/domain"
	"strconv"
	"sync"
	"time"
)

const concurrentUsers = 110 // on my laptop, 100 and less doesn't produce race condition error

type doWork func(ctx context.Context, tableID string, userID string)

func worker(work doWork, tableID string, userID string, resultWG *sync.WaitGroup, wait *sync.WaitGroup) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	wait.Wait()
	work(ctx, tableID, userID)
	resultWG.Done()
}

func runWorkers(table *domain.Table, work doWork, count int) {
	var resultWg sync.WaitGroup
	resultWg.Add(count)
	var waitFlag sync.WaitGroup
	waitFlag.Add(1)
	for i := 0; i < count; i++ {
		go worker(work, table.ID.Hex(), strconv.Itoa(i), &resultWg, &waitFlag)
	}

	// Give some time for workers to start
	time.Sleep(5 * time.Millisecond)

	// Start reservations
	waitFlag.Done()

	resultWg.Wait()
}
