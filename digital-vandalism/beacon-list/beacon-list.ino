#include <ESP8266WiFi.h>
#include <EEPROM.h>
#include <ESP8266WiFi.h>
#include <ESP8266WebServer.h>

#define MAX_SSID 12

/* Put your SSID & Password */
const char* ssid = "hacker";  // Enter SSID here
const char* password = "hacker1234";  //Enter Password here

/* Put IP Address details */
IPAddress local_ip(192,168,1,1);
IPAddress gateway(192,168,1,1);
IPAddress subnet(255,255,255,0);

ESP8266WebServer server(80);

byte channel;
String spoof_ssid="I❤️PUSSY!";
bool spamming=false;

// uint8_t packet[128] = { 0x80, 0x00, 0x00, 0x00, 
//                         0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 
//                         0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
//                         0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 
//                         0xc0, 0x6c, 
//                         0x83, 0x51, 0xf7, 0x8f, 0x0f, 0x00, 0x00, 0x00, 
//                         0x64, 0x00, 
//                         0x01, 0x04, 
                
//                         0x00, 0x0a, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72,
//                         0x01, 0x08, 0x82, 0x84,
//                         0x8b, 0x96, 0x24, 0x30, 0x48, 0x6c, 0x03, 0x01, 
//                         0x04
// };                       

uint8_t packet[] = {
    // Frame Control
    0x80, 0x00,
    // Duration
    0x00, 0x00,
    // Destination (Broadcast)
    0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
    // Source (AP MAC)
    0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc,
    // BSSID (AP MAC)
    0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc,
    // Sequence Control
    0x10, 0x00,
    // Timestamp
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    // Beacon Interval
    0x64, 0x00,
    // Capability Info
    0x01, 0x04,
    // SSID Tag
    0x00, 0x14, ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', 
    ' ', ' ', ' ', ' ', ' ', ' ', ' ',
    // Supported Rates Tag
    0x01, 0x08, 0x82, 0x84, 0x8b, 0x96, 0x24, 0x30, 0x48, 0x6c,
    // DS Parameter Set Tag
    0x03, 0x01, 0x06
};

extern "C" {
  #include "user_interface.h"
}

void setup() {
  Serial.begin(115200);
  delay(500);

  WiFi.softAP(ssid, password);
  WiFi.softAPConfig(local_ip, gateway, subnet);
  delay(100);
  
  server.on("/", handle_OnConnect);
  server.on("/spam", handle_Spam);
  server.on("/stop", handle_Stop);

  server.begin();
}

void set_ssid(const char *ssid, uint8_t *packet, int channel){
  //set mac address to random values
  packet[10] = packet[16] = random(256);
  packet[11] = packet[17] = random(256);
  packet[12] = packet[18] = random(256);
  packet[13] = packet[19] = random(256);
  packet[14] = packet[20] = random(256);
  packet[15] = packet[21] = random(256);

  // set WPA2
  packet[34] = 0x31;

  //set lentgh of the SSID
  int ssid_len = strlen(ssid);
  if (ssid_len > 19) ssid_len = 19;
  packet[37] = 19;

  //set channel
  packet[82] = channel;

  //clean SSID
  for (int i=0; i<19; i++){
    packet[38+i]=' ';
  }

  for (int i=0; i<strlen(ssid); i++){
    packet[38+i]=ssid[i];
  }
}

void handle_Spam(){
    spoof_ssid = server.arg("ssid");
    spoof_ssid = spoof_ssid.substring(0,spoof_ssid.length()-1); 
    String html = "<!DOCTYPE html>"
              "<html>"
              "<head>"
              "<meta charset=\"UTF-8\">"
              "<title>SSID Sender</title>"
              "<style>"
              "body { font-size: 24px; }"
              "</style>"
              "</head>"
              "<body>"
              "SENDING BEACON WITH SSID"+spoof_ssid+
              "</body>"
              "</html>";
    server.send(200, "text/html", html);
    //wifi_set_opmode(STATION_MODE);
    //wifi_promiscuous_enable(1); 
    spamming=true;
}

void handle_OnConnect() { 
  String ssids = "";
  for (int i=0; i < MAX_SSID; i++){
    ssids+="<input type='text' maxlength='19' id='" + String(i) + "' placeholder='Enter SSID'><br/>";
  }

  String f;
  for (int i=0; i < MAX_SSID; i++){
    f +="ssid += document.getElementById("+String(i)+").value; ssid+=\"|\";";
  }

  String html = "<!DOCTYPE html>"
              "<html>"
              "<head>"
              "<title>SSID Sender</title>"
              "<style>"
              "body { font-size: 24px; }"
              "input { font-size: 24px; padding: 10px; width: 80%; margin: 20px; }"
              "button { font-size: 24px; padding: 10px 20px; margin: 20px; }"
              "</style>"
              "</head>"
              "<body>" + ssids +
              "<button onclick='sendSSID()'>Send SSID</button>"
              "<script>"
              "function sendSSID() { var ssid=\"\";" + f +
              "const baseUrl = window.location.origin;"
              "window.location.href = baseUrl + '/spam?ssid=' + encodeURIComponent(ssid);"
              "}"
              "</script>"
              "</body>"
              "</html>";
    server.send(200, "text/html", html);
}

void handle_Stop() {
    server.send(200, "text/html", "STOPPED SPAMMING");
    spamming=false;
}

void loop(){
  server.handleClient();

   if(spamming){
    char * pch;
    String temp = spoof_ssid.c_str();
    pch = strtok ((char *)temp.c_str() ,"|");
    while (pch != NULL)
    {
      channel = random(1,12); 
      wifi_set_channel(channel);

      set_ssid(pch, (uint8_t *)&packet, channel);
          
      wifi_send_pkt_freedom(packet, 57, 0);
      wifi_send_pkt_freedom(packet, 57, 0);
      wifi_send_pkt_freedom(packet, 57, 0);
      delay(1);

      pch = strtok (NULL, "|");     
    }
  }
}

