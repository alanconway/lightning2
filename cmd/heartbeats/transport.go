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
	"encoding/json"
	"errors"
	"strconv"
	"time"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/transport"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
)

type Duration struct{ time.Duration }

func (d Duration) MarshalJSON() ([]byte, error) { return json.Marshal(d.String()) }

func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	var v interface{}
	if err = json.Unmarshal(b, &v); err != nil {
		return
	}
	switch v := v.(type) {
	case float64:
		d.Duration = time.Duration(v)
	case string:
		d.Duration, err = time.ParseDuration(v)
	default:
		err = errors.New("invalid duration")
	}
	return
}

type Transport struct {
	Source types.URLRef
	Delay  Duration
	id     int
	r      transport.Receiver
}

func (t *Transport) ConfigJSON(b []byte) error { return json.Unmarshal(b, t) }

func (t *Transport) Send(context.Context, cloudevents.Event) (*cloudevents.Event, error) {
	return nil, errors.New("heartbeat transport: send not supported")
}

func (t *Transport) SetReceiver(r transport.Receiver) { t.r = r }

func (t *Transport) StartReceiver(ctx context.Context) error {
	for ts := range time.Tick(t.Delay.Duration) {
		t.id++
		e := cloudevents.Event{
			Context: cloudevents.EventContextV02{
				Type:   "dev.knative.eventing.samples.heartbeat",
				Source: t.Source,
				ID:     strconv.Itoa(t.id),
				Time:   &types.Timestamp{ts},
			}.AsV02(),
		}
		t.r.Receive(ctx, e, nil)
	}
	return nil
}
