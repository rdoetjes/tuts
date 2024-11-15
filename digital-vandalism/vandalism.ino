#include <ESP8266WiFi.h>
#include <EEPROM.h>

extern "C" {
  #include "user_interface.h"
}

byte channel;
String ssid="I❤️PUSSY!";

uint8_t packet[128] = { 0x80, 0x00, 0x00, 0x00, 
                        0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 
                        0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
                        0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 
                        0xc0, 0x6c, 
                        0x83, 0x51, 0xf7, 0x8f, 0x0f, 0x00, 0x00, 0x00, 
                        0x64, 0x00, 
                        0x01, 0x04, 
                
                        0x00, 0x0a, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72, 0x72,
                        0x01, 0x08, 0x82, 0x84,
                        0x8b, 0x96, 0x24, 0x30, 0x48, 0x6c, 0x03, 0x01, 
                        0x04
};                       

void setup() {
  Serial.begin(115200);
  delay(500);
  Serial.println("YOUR SSID?!");
  Serial.setTimeout(20*1000);
  String t=Serial.readString();
  if (t.length()!=0) {
    ssid=t;
  } 

  wifi_set_opmode(STATION_MODE);
  wifi_promiscuous_enable(1); 
}

void set_ssid(const char *ssid, uint8_t *packet){
  static int nr_spaces=0;

  //set mac address to random values
  packet[10] = packet[16] = random(256);
  packet[11] = packet[17] = random(256);
  packet[12] = packet[18] = random(256);
  packet[13] = packet[19] = random(256);
  packet[14] = packet[20] = random(256);
  packet[15] = packet[21] = random(256);

  //set lentgh of the SSID
  packet[37]=strlen(ssid)+2;

  //clean SSID
  for (int i=0; i<32; i++){
    packet[38+i]=0;
  }

  for (int i=0; i<strlen(ssid); i++){
    packet[38+i]=ssid[i];
  }

  packet[38+strlen(ssid)+1]=0xC2;
  packet[38+strlen(ssid)+2]=0xA0;

  if (strlen(ssid)+nr_spaces>=30) nr_spaces = 0;
  nr_spaces++;

  for(int i=0; i<nr_spaces; i++){
      packet[38+strlen(ssid)+i]=random(31);
  }
}

void loop() { 
    channel = random(1,12); 
    wifi_set_channel(channel);

    set_ssid(ssid.c_str(), (uint8_t *)&packet);
    
    packet[82] = channel;
    
    wifi_send_pkt_freedom(packet, 57, 0);
    wifi_send_pkt_freedom(packet, 57, 0);
    wifi_send_pkt_freedom(packet, 57, 0);
    delay(1);
}
