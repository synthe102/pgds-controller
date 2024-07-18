package handler

import (
	"context"
	"fmt"

	"golang.org/x/time/rate"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type PendingChangesInterface interface {
	client.Object
	IsPendingChanges() bool
}

type PrioritizedEventHandler struct {
	lowPriorityWorkqueue   workqueue.TypedInterface[reconcile.Request]
	lowPriorityRateLimiter *rate.Limiter
	controllerWorkqueue    RateLimitingInterface
}

func (t *PrioritizedEventHandler) consumeLowPriority(ctx context.Context) {
	for {
		item, _ := t.lowPriorityWorkqueue.Get()
		_ = t.lowPriorityRateLimiter.Wait(ctx)
		t.controllerWorkqueue.Add(item)

		t.lowPriorityWorkqueue.Done(item)
	}
}

type CreateEvent = event.TypedCreateEvent[PendingChangesInterface]
type DeleteEvent = event.TypedDeleteEvent[PendingChangesInterface]
type UpdateEvent = event.TypedUpdateEvent[PendingChangesInterface]
type GenericEvent = event.TypedGenericEvent[PendingChangesInterface]
type RateLimitingInterface = workqueue.TypedRateLimitingInterface[reconcile.Request]

// Create implements EventHandler.
func (t PrioritizedEventHandler) Create(ctx context.Context, e CreateEvent, q RateLimitingInterface) {
}

// Delete implements EventHandler.
func (t PrioritizedEventHandler) Delete(ctx context.Context, e DeleteEvent, q RateLimitingInterface) {
}

// Update implements EventHandler.
func (t PrioritizedEventHandler) Update(ctx context.Context, e UpdateEvent, q RateLimitingInterface) {
}

// Generic implements EventHandler.
func (t PrioritizedEventHandler) Generic(ctx context.Context, e GenericEvent, _ RateLimitingInterface) {
	log := log.FromContext(ctx)
	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name: e.Object.GetName(),
		},
	}
	if e.Object.IsPendingChanges() {
		log.Info(fmt.Sprintf("Adding to low controller workqueue: %s", e.Object.GetName()))
		t.controllerWorkqueue.Add(request)
	} else {
		log.Info(fmt.Sprintf("Adding to low priority workqueue: %s", e.Object.GetName()))
		t.lowPriorityWorkqueue.Add(request)
	}
}

func (t *PrioritizedEventHandler) GetControllerQueue(controllerName string, rateLimiter workqueue.TypedRateLimiter[reconcile.Request]) workqueue.TypedRateLimitingInterface[reconcile.Request] {
	return t.controllerWorkqueue
}

func NewPrioritizedEventHandler() PrioritizedEventHandler {
	ratelimiter := workqueue.DefaultTypedControllerRateLimiter[reconcile.Request]()
	return PrioritizedEventHandler{
		lowPriorityWorkqueue:   workqueue.NewTyped[reconcile.Request](),
		controllerWorkqueue:    workqueue.NewTypedRateLimitingQueue[reconcile.Request](ratelimiter),
		lowPriorityRateLimiter: rate.NewLimiter(rate.Limit(5), 1),
	}
}

func (t PrioritizedEventHandler) Start(ctx context.Context) error {
	t.consumeLowPriority(ctx)
	return nil
}
