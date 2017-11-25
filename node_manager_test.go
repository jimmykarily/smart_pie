package main_test

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	. "github.com/jimmykarily/smart_pie"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("NodeManager", func() {
	var (
		manager *NodeManager
		client  mqtt.Client
	)

	BeforeEach(func() {
		manager = NewNodeManager(BROKER_URI, "admin", "30062002")

		opts := mqtt.NewClientOptions().AddBroker(BROKER_URI)
		client = mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			Fail(token.Error().Error())
		}
	})

	Describe("addNode", func() {
		It("adds a node based on the message", func() {
			manager.Subscribe([]string{"heartbeats"})

			if token := client.Publish("heartbeats", 0, false, "mynode:4,3,1"); token.Wait() && token.Error() != nil {
				Fail(token.Error().Error())
			}

			// Do something better here than waiting a random amount of time
			// We wait until the message is processed by the callback.
			time.Sleep(50 * time.Millisecond)

			Expect(manager.Nodes["mynode"].Name).To(Equal("mynode"))
			Expect(manager.Nodes["mynode"].Pins).To(Equal([]int{4, 3, 1}))
		})
	})
})
