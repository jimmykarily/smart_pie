/*

Smart home with arduino

This is the arduino code that connects to the mosquitto server and subscribes
to the channels the control the state of the output pins.

The valid actions for pins are:

- high
  The specified pin will go to HIGH state no matter what the current state is
- low
  The specified pin will go to LOW state no matter what the current state is
- toggle
  When the pin is in HIGH state it will become LOW
  When the pin is in LOW state is will become HIGH

Examples:

- Digital pin 4 must go to HIGH state. The message "high" should be received
  on the <NODE>/pin/4
- Digital pin 3 must go to LOW state. The message "low" should be received
  on the <NODE>/pin/3
- Digital pin 5 must toggle state. The message "toggle" should be received
  on the <NODE>/pin/5

The mosquitto server requires the clients to authenticate with username and
password. The user name of each node is the value stored in NODE. This variable
is set in config.h file. The WiFi credentials and SSID are also stored in
the same file.

*/

#include <ESP8266WiFi.h>
#include <PubSubClient.h>
#include "config.h"

// This is a heartbeat. How often should this node advertise its enabled ports?
#define ADVERTISE_INTERVAL_MILLIS 10000

WiFiClient espClient;
PubSubClient client(espClient);
long last_heartbeat = 0;
char heartbeat_text[50];

void setup() {
  // pinMode(BUILTIN_LED, OUTPUT);     // Initialize the BUILTIN_LED pin as an output

  int number_of_enabled_pins = sizeof(ENABLED_PINS) / sizeof(int);
  for (int i=0; i<=number_of_enabled_pins; i++) {
    pinMode(ENABLED_PINS[i], OUTPUT);
  }

  Serial.begin(115200);
  setup_wifi();
  client.setServer(MQTT_SERVER, 1883);
  client.setCallback(callback);
}

void setup_wifi() {
  delay(10); // TODO: Do we need this?
  // We start by connecting to a WiFi network
  Serial.println();
  Serial.print("Connecting to ");
  Serial.println(WIFI_SSID);

  WiFi.mode(WIFI_STA);
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  Serial.println("");
  Serial.println("WiFi connected");
  Serial.println("IP address: ");
  Serial.println(WiFi.localIP());
}

int extract_pin_from_topic(char* topic) {
  int i = 0;
  int len = strlen(topic);
  int slashes=0;
  char pin_string[3];
  int pin_i=0;

  // Skip two '/' and the rest is the pin number
  for (int i=0;i<=len;i++) {
    if(slashes < 2) { // more to skip
      if(topic[i] == '/') { slashes++; }
    } else { // keep everything after the second '/'
      pin_string[pin_i] = topic[i];
      pin_i++;
    }
  }

  return atoi(pin_string);
}

String action_from_payload(byte* payload, unsigned int length) {
  char action[length];
  memcpy(action, payload, length);
  action[length] = '\0';

  return action;
}

String enabled_ports_string() {
  char text[21];
  char pin_num[3];
  int number_of_enabled_pins = sizeof(ENABLED_PINS) / sizeof(int);

  strcpy(text, "");
  for(int i=0; i<number_of_enabled_pins; i++) {
    itoa(ENABLED_PINS[i], pin_num, 10);
    strcat(text, pin_num);
    if(digitalRead(ENABLED_PINS[i])){
      strcat(text, "/up");
    } else {
      strcat(text, "/down");
    }
    // Skip the comma for the last element
    if(i != number_of_enabled_pins - 1){
      strcat(text, ",");
    }
  }

  return text;
}

void callback(char* topic, byte* payload, unsigned int length) {
  String action = action_from_payload(payload, length);
  int pin = extract_pin_from_topic(topic);

  Serial.println();
  Serial.print("Applying ");
  Serial.print(action);
  Serial.print(" to pin ");
  Serial.println(pin);

  if(strcmp(action.c_str(),"high") == 0) {
    digitalWrite(pin, HIGH);
  }
  else if(strcmp(action.c_str(), "low") == 0) {
    digitalWrite(pin, LOW);
  }
  else if(strcmp(action.c_str(), "toggle") == 0) {
    digitalWrite(pin, !digitalRead(pin));
  }
}

void reconnect() {
  // Loop until we're reconnected
  while (!client.connected()) {
    Serial.print("Attempting MQTT connection...");
    // Attempt to connect
    if (client.connect("ESP8266Client", NODE, MQTT_PASSWORD)) {
      Serial.println("connected");
      // Once connected, publish an announcement...
      // client.publish("outTopic", "hello world");

      // Subscribe to one channel for each enabled pin
      int number_of_enabled_pins = sizeof(ENABLED_PINS) / sizeof(int);
      for (int i=0; i<number_of_enabled_pins; i++) {
        char topic[256];
        snprintf(topic, sizeof(topic), "%s%s%d", NODE, "/pin/", ENABLED_PINS[i]);
        Serial.println(topic);
        client.subscribe(topic);
      }
    } else {
      Serial.print("failed, rc=");
      Serial.print(client.state());
      Serial.println(" try again in 5 seconds");
      // Wait 5 seconds before retrying
      delay(5000);
    }
  }
}
void loop() {
  if (!client.connected()) {
    reconnect();
  }
  client.loop();

  long now = millis();
  if (now - last_heartbeat > ADVERTISE_INTERVAL_MILLIS) {
    last_heartbeat = now;
    sprintf (heartbeat_text, "%s:%s", NODE, enabled_ports_string().c_str());
    Serial.print("Sending heartbeat: ");
    Serial.println(heartbeat_text);
    client.publish("heartbeats", heartbeat_text);
  }
}
