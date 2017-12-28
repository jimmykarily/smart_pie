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
	"strconv"
	"strings"
)

type Node struct {
	Name string
	Pins map[int]bool
}

type NodeManager struct {
	Nodes      map[string]Node
	mqttClient MQTT.Client
	logger     Logger
}

func NewNodeManager(brokerUri string, brokerUsername string, brokerPassword string) *NodeManager {
	clientOptions := MQTT.NewClientOptions().AddBroker(brokerUri)

	if brokerUsername != "" {
		clientOptions = clientOptions.SetClientID(brokerUsername).
			SetUsername(brokerUsername)
	}

	if brokerPassword != "" {
		clientOptions = clientOptions.SetPassword(brokerPassword)
	}

	// Create the manager to be able to use its method as a callback
	manager := &NodeManager{
		logger: Logger{"Node Manager", os.Stdout},
		Nodes:  make(map[string]Node),
	}

	clientOptions.SetDefaultPublishHandler(
		MQTT.MessageHandler(manager.heartbeatsCallback))
	manager.mqttClient = MQTT.NewClient(clientOptions)

	return manager
}

/*
msg is something like this:
	my_node: 1,2,3
where 1,2,3 are the enabled pins and my_node is the name of the node. This
method populates the Nodes array.
*/
func (m *NodeManager) heartbeatsCallback(client MQTT.Client, msg MQTT.Message) {
	if msg.Topic() == "heartbeats" {
		m.addNode(string(msg.Payload()))
	} else {
		m.logger.Log(fmt.Sprintf("Unhandled msg for topic: %s / %s", msg.Topic(), msg.Payload()))
	}
	fmt.Println(m.Nodes)
}

/*
A heartbeat message looks something like this:
my_node:1/up,2/down,4/up

and should result in a node with name "my_node" and Pins {1: true, 2: down, 4:true}
*/
func (m *NodeManager) addNode(msg string) {
	msgParts := strings.Split(msg, ":")

	var node Node
	var exists bool

	if node, exists = m.Nodes[msgParts[0]]; !exists {
		node = Node{Name: msgParts[0]}
	}

	pins := make(map[int]bool)

	for _, pinInfo := range strings.Split(msgParts[1], ",") {
		var data []string = strings.Split(pinInfo, "/")

		pinNumber, err := strconv.Atoi(strings.TrimSpace(data[0]))
		if err != nil {
			panic(err)
		}

		pins[pinNumber] = (data[1] == "up")
	}

	node.Pins = pins
	m.Nodes[msgParts[0]] = node
}

func (m *NodeManager) Subscribe(channels []string) {

	if token := m.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for _, channel := range channels {
		// TODO: Use SubscribeMultiple here
		if token := m.mqttClient.Subscribe(channel, 0, func(client MQTT.Client, msg MQTT.Message) {
			m.heartbeatsCallback(client, msg)
		}); token.Wait() && token.Error() != nil {
			m.logger.Log(token.Error().Error())
			os.Exit(1)
		}
	}
}

func (m *NodeManager) SendMessage(channel string, message string) {
	m.logger.Log(channel)
	m.logger.Log(message)
	if token := m.mqttClient.Publish(channel, 0, false, message); token.Wait() && token.Error() != nil {
		m.logger.Log(token.Error().Error())
		os.Exit(1)
	}
}
