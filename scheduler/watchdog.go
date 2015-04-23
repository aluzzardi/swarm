package scheduler

import (
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/swarm/cluster"
	"github.com/docker/swarm/state"
)

const (
	watchdogTimer = 10 * time.Second
)

type Watchdog struct {
	sync.Mutex

	cluster cluster.Cluster
	store   *state.Store

	stop chan struct{}
}

func NewWatchdog(cluster cluster.Cluster, store *state.Store) *Watchdog {
	return &Watchdog{
		cluster: cluster,
		store:   store,
	}
}

func (w *Watchdog) Start() {
	go w.balanceLoop()
}

func (w *Watchdog) Stop() {
	w.stop <- struct{}{}
}

func (w *Watchdog) balanceLoop() {
	for {
		select {
		case <-w.stop:
			break
		case <-time.After(watchdogTimer):
			w.Balance()
		}
	}
}

func (w *Watchdog) Balance() {
	w.Lock()
	defer w.Unlock()

	for _, state := range w.store.All() {
		if container := w.cluster.Container(state.ID); container == nil {
			log.Infof("Container %s is missing from the cluster, rescheduling.", state.ID)

			if err := w.rebalanceContainer(state); err != nil {
				log.Errorf("Unable to reschedule container %s: %v", state.ID, err)
				continue
			}
			log.Debugf("Successfully rescheduled container %s", state.ID)
		}
	}
}

func (w *Watchdog) rebalanceContainer(state *state.RequestedState) error {
	var (
		err       error
		container *cluster.Container
	)

	// Attempt to re-create the container on the cluster.
	if container, err = w.cluster.CreateContainer(state.Config, state.Name); err != nil {
		return err
	}

	// Remove the previous state entry from the store.
	if err = w.store.Remove(state.ID); err != nil {
		return err
	}

	// FIXME: We have no idea whether the container was running or not
	// by just looking at the store. Start it anyway.
	container.Start()

	return nil
}
