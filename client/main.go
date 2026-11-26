package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
)

func main() {
	c, err := client.NewClient()

	if err != nil {
		log.Fatal(err)
	}

	req := protocol.AcquireRequest()
	resp := protocol.AcquireResponse()
	defer protocol.ReleaseRequest(req)
	defer protocol.ReleaseResponse(resp)

	req.SetMethod(consts.MethodGet)
	req.SetRequestURI("http://127.0.0.1:8080/health")

	// if err := c.Do(context.Background(), req, resp); err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("status=%d body=%s\n", resp.StatusCode(), resp.Body())

	status, body, err := c.Get(context.Background(), nil, "http://127.0.0.1:8080/ping")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("status=%d body=%s\n", status, body)
}
