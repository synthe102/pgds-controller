package main

import (
	"context"

	"github.com/synthe102/pgds-controller/internal/datastore"
	"github.com/synthe102/pgds-controller/internal/handler"

	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func main() {
	log.SetLogger(zap.New())
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		panic(err)
	}
	prioritizedHandler := handler.NewPrioritizedEventHandler()
	mgr.Add(prioritizedHandler)
	ctrl, err := controller.New("item-controller", mgr, controller.Options{
		Reconciler: reconciler{},
		NewQueue:   prioritizedHandler.GetControllerQueue,
	})
	if err != nil {
		panic(err)
	}

	events := make(chan event.TypedGenericEvent[handler.PendingChangesInterface])
	ds := datastore.Datastore{
		EventChan: events,
	}
	mgr.Add(ds)
	err = ctrl.Watch(source.Channel[handler.PendingChangesInterface](events, prioritizedHandler))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	err = mgr.Start(ctx)
	if err != nil {
		panic(err)
	}
}

type reconciler struct{}

func (r reconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	log := log.FromContext(ctx)
	log.Info("reconciling", "req", req)
	return reconcile.Result{}, nil
}
