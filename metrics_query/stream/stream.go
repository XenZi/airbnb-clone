package stream

import events "example/metrics_events"

type EventStream interface {
	Process(func(events.Event) error)
}
