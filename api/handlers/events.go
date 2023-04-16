package handlers

import (
	"fmt"
	"net/http"

	"github.com/rugwirobaker/hermes"
	"github.com/rugwirobaker/hermes/api/render"
)

// SubscribeHandler handles user subscriptions to delivery notifications
func SubscribeHandler(events hermes.Pubsub) http.HandlerFunc {
	// const op = "handlers.SubscribeHandler"

	return func(w http.ResponseWriter, r *http.Request) {

		// ctx, span := observ.StartSpan(r.Context(), op)
		// defer span.End()

		// id := chi.URLParam(r, "id")

		// h := w.Header()
		// h.Set("Content-Type", "text/event-stream")
		// h.Set("Cache-Control", "no-cache")
		// h.Set("Connection", "keep-alive")
		// h.Set("X-Accel-Buffering", "no")

		// f, ok := w.(http.Flusher)
		// if !ok {
		// 	log.Println("could not start stream")
		// 	return
		// }

		// ctx, cancel := context.WithCancel(ctx)
		// defer cancel()

		// event, err := events.Subscribe(ctx, id)
		// if err != nil {
		// 	log.Println(err)
		// 	render.Flush(w, f, NewError(err.Error()))
		// 	return
		// }

		// for {
		// 	select {
		// 	case <-ctx.Done():
		// 		log.Println("event: stream canceled")
		// 		render.Flush(w, f, NewError("context canceled"))
		// 		return

		// 	case <-time.After(time.Second * 10):
		// 		log.Println("event: stream timeout")
		// 		render.Flush(w, f, NewError("connection timeout"))
		// 		return

		// 	case res := <-event:
		// 		render.Flush(w, f, res)
		// 		events.Done(ctx, res.ID)
		// 		return
		// 	}
		// }
		render.HttpError(w, fmt.Errorf("not implemented"))
	}
}

// type event struct {
// 	MsgRef     string `json:"msgRef"`
// 	Recipient  string `json:"recipient"`
// 	GatewayRef string `json:"gatewayRef"`
// 	Status     int    `json:"status"`
// }

// func convertEvent(event *event) hermes.Event {
// 	return hermes.Event{
// 		ID:        event.MsgRef,
// 		Recipient: event.Recipient,
// 		Status:    hermes.St(event.Status),
// 	}
// }
