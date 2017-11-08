/*
  NodeManager is the "class" that handles the auto discovery of new Nodes.
	It subscribes to the "heartbeats" channel on Mosquitto server. Whenever
	a new Node appears, it is written in the nodes list. For every Node in the
	list, we know the enabled pins. The web app serves these nodes with their
	pins.
*/
package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"os"
)

type Node struct {
}

type NodeManager struct {
	nodes      []Node
	mqttClient MQTT.Client
	logger     Logger
}

func NewNodeManager(brokerUri string, brokerUsername string, brokerPassword string) *NodeManager {
	clientOptions := MQTT.NewClientOptions().
		AddBroker(brokerUri).
		SetClientID(brokerUsername).
		SetUsername(brokerUsername).
		SetPassword(brokerPassword)

	// Create the manager to be able to user it's method as a callback
	manager := &NodeManager{
		logger: Logger{"Node Manager", os.Stdout},
	}

	clientOptions.SetDefaultPublishHandler(
		MQTT.MessageHandler(manager.heartbeatsCallback))
	manager.mqttClient = MQTT.NewClient(clientOptions)

	return manager
}

func (m *NodeManager) heartbeatsCallback(client MQTT.Client, msg MQTT.Message) {
	m.logger.Log(fmt.Sprintf("TOPIC: %s", msg.Topic()))
	m.logger.Log(fmt.Sprintf("MSG: %s", msg.Payload()))
}

func (m *NodeManager) Subscribe() {
	if token := m.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := m.mqttClient.Subscribe("heartbeats", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}
