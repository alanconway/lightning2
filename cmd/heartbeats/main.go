/*
Copyright 2019 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/knative/eventing-sources/pkg/kncloudevents"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
)

var (
	hostname, _ = os.Hostname()
	source      = flag.String("source", hostname, "source attribute for events")
	sink        = flag.String("sink", "", "send events to this URL")
	delay       = flag.Duration("delay", 5*time.Second, "delay between events")
)

func check(err error, format string, args ...interface{}) {
	if err != nil {
		log.Fatalf("%v: %v", fmt.Sprintf(format, args), err)
	}
}

type forwarder struct{ client.Client }

func (f forwarder) Receive(ctx context.Context, event cloudevents.Event, resp *cloudevents.EventResponse) (err error) {
	log.Printf("received event %+v", event)
	if _, err = f.Send(ctx, event); err != nil {
		log.Printf("failed to send %+v: %s", event, err.Error())
	}
	return
}

func main() {
	flag.Parse()
	t := Transport{Source: *types.ParseURLRef(*source), Delay: Duration{*delay}}
	c, err := kncloudevents.NewDefaultClient(*sink)
	check(err, "cannot connect to %#v", *sink)
	log.Printf("heartbeats from %#v to %#v every %v", *source, *sink, *delay)
	f := forwarder{c}
	t.SetReceiver(f)
	t.StartReceiver(context.Background())
}
