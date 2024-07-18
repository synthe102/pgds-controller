package datastore

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/synthe102/pgds-controller/internal/handler"
	"github.com/synthe102/pgds-controller/internal/model"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Datastore struct {
	EventChan chan event.TypedGenericEvent[handler.PendingChangesInterface]
}

func (d Datastore) Start(ctx context.Context) error {
	log := log.FromContext(ctx)

	// Fetch all items before starting to watch
	resp, err := http.Get("http://localhost:3000/items")
	if err != nil {
		log.Error(err, "error listing all items")
	}
	var itemList map[string]model.Item
	err = json.NewDecoder(resp.Body).Decode(&itemList)
	resp.Body.Close()
	if err != nil {
		panic(err)
	}

	for _, item := range itemList {
		item.ObjectMeta.Name = item.ID
		d.EventChan <- event.TypedGenericEvent[handler.PendingChangesInterface]{
			Object: &item,
		}
	}
	log.Info("watchItems started")

Loop:
	for {
		select {
		case <-time.After(5 * time.Second):
			resp, err := http.Get("http://localhost:3000/items/watch")
			if err != nil {
				log.Error(err, "error watching for items")
			}
			if resp.StatusCode == http.StatusRequestTimeout {
				log.Info("timeout in watchItems")
				continue
			}
			var item model.Item
			err = json.NewDecoder(resp.Body).Decode(&item)
			resp.Body.Close()
			if err != nil {
				panic(err)
			}
			item.ObjectMeta.Name = item.ID
			d.EventChan <- event.TypedGenericEvent[handler.PendingChangesInterface]{
				Object: &item,
			}
		case <-ctx.Done():
			log.Info("watchItems stopped")
			break Loop
		}
	}
	return nil
}
