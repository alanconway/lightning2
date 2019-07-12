# Archived and obsolete

This work is proceeding by contributing directly to https://github.com/cloudevents and https://github.com/knative

# A simpler way to build knative eventing sources

## Encapsulate domain knowledge

Events are being captured from some non-knative "domain": be it
management events, service notifications, messages received via
foreign protocols like Kafka, MQTT, AMQP - could be anything

"Domain knowledge" is implemented in one place and separated from
knative and k8s details: as an implementation of a simple *Transport*
interface. This is an extension of
https://godoc.org/github.com/cloudevents/sdk-go/pkg/cloudevents/transport#Transport
that also includes configuration and life-cycle operations.

## Re-usable knative components

Once a "transport" is available for some domain, everything else about
building an eventing source is quite generic and repetative. Provide
ready-to-use implementations of knative-eventing components such as
adapters and controllers etc. that just need a transport to become
full-fledged eventing sources.

## To Do

Plugins: provide generic controllers/adapters that load transports as Go plug-ins?

QoS: protocols like HTTP, AMQP, MQTT can provide varied QoS. To
support at-least-once or exactly-once the sink needs to notify sthe
source. Source and sink should also advertise their QoS range so both
sides can degrade to the lowest level for performance, since the
adapter's overall QoS will be that of the weakest link anyway.

Performance: Current Source & Sink impls do more copying than
they should, no performance benchmarks have been done to date.
