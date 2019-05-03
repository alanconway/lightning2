/*
Copyright 2018 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// A generic event receiver that logs each event.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/knative/eventing-sources/pkg/kncloudevents"
)

func check(err error, format string, args ...interface{}) {
	if err != nil {
		log.Fatalf("%s: %s", fmt.Sprintf(format, args...), err)
	}
}

func receive(ctx context.Context, e cloudevents.Event, er *cloudevents.EventResponse) error {
	log.Printf("CloudEvent:\n%s", e)
	return nil
}

func main() {
	c, err := kncloudevents.NewDefaultClient()
	check(err, "can't create client")
	log.Printf("%s listening on default port 8080", os.Args[0])
	ctx := context.TODO()
	check(c.StartReceiver(ctx, receive), "failed to start receiver")
	<-ctx.Done()

	// FIXME aconway 2019-05-01:
	runtime.FuncForPC(0)
}
