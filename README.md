# Smartpie

**Warning**: The following README described the desired state which is not yet reached.
This project is still a work in progress. This warning will be removed as soon as
the project reaches a "working" state.

This is DIY Smart home platform that lets you replace switches to your house
with WiFi enabled Arduinos (ESPduino). It consists of the following entities:

- The nodes
  These are the Arduinos. Each board has a number of of pins that can be
  controlled remotely. The platform will simply let you change their state
  (HIGH, LOW). What that does is up to you. In the simplest case, you would
  control a relay to turn on a light/door/computer/etc.

- The admin node
  This is a web application that provides the interface to control the nodes.
  You can give meaningful names to the pins of each node and you can change
  their state from the nice web UI.

- The mosquitto server (https://mosquitto.org/)
  This is one implementation of an MQTT broker. The MQTT protocol is designed
  to be lightweight and suitable for IoT applications. This is a pub/sub server
  to which messages are sent from the admin node to the nodes. These message
  control the state of the pins.

Below you can find information on how to setup each of the components.

## Nodes

We assume that each node is an ESPduino. That is an Arduino with an ESP8266
module built in. You can find all the code needed to setup a node under the
`node` directory.

Before you setup a new node you need to copy `node/config.sample.h` to
`node/config.h` and edit the values in the file as needed. After that
you will need to compile and upload the sketch to your board. You can use the
[default Arduino IDE](https://www.arduino.cc/en/Main/Software) for this.
You will have to add the ESP8266 support (https://github.com/esp8266/Arduino#installing-with-boards-manager).

## Admin node

This is a web application written in Go. To compile it you will need a working Go development environment.

If you are planning to run the admin application on a RasberryPi you can use something like

```
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -a -ldflags '-extldflags "-static"' .
docker build -t smart_pie .
```
This will create a statically linked executable for Rasberry Pi (v3) and will then create a docker image
which you can run on your Rasberry using something like this:

```
docker run -dit --restart unless-stopped -e "MQTT_USERNAME=your_admin_username" -e "MQTT_PASSWORD=your_admin_password" -e "MQTT_BROKER_URI=tcp://the_ip_of_your_broker:1883"  -p 8080:8080 smart_pie
```

## Mosquitto

We want to authenticate the nodes (arduinos and admin) to prevent unauthorized
access you our pins (and house). A nice and quick read about MQTT security can
be found here: https://www.hivemq.com/blog/mqtt-security-fundamentals/

In the `mosquitto` directory you can find 2 files. 
`mosquitto.conf` is the configuration for the server and `mosquitto_acl` configures
the authorization of each client. No node should be allowed to read another node's
channels and only the admin node should be allowed to write to channels controlling
the pins. You will need to create a `passwords` file which holds the password
of the nodes. You can either do this manually of using the mosquitto_passwd tool.
More about this here: https://mosquitto.org/man/mosquitto-conf-5.html

We recommend using docker to start the mosquitto server. It can be achieved with
a single commmand (issued from within the `mosquitto` directory):

```
 docker run --rm -it -p 1883:1883 -p 90:9001 -v $PWD/mosquitto.conf:/mosquitto/config/mosquitto.conf -v $PWD/mosquitto_acl:/mosquitto_acl -v $PWD/passwords:/passwords eclipse-mosquitto
```

TODO: Setup SSL communication.
