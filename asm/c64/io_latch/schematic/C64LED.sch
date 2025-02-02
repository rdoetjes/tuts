EESchema Schematic File Version 4
EELAYER 30 0
EELAYER END
$Descr A4 11693 8268
encoding utf-8
Sheet 1 1
Title ""
Date ""
Rev ""
Comp ""
Comment1 ""
Comment2 ""
Comment3 ""
Comment4 ""
$EndDescr
Wire Wire Line
	4250 3350 4250 3500
Wire Wire Line
	4250 4900 4250 4750
$Comp
L LED:HDSP-4830_2 BAR1
U 1 1 5E644CA2
P 5550 4100
F 0 "BAR1" H 5550 4767 50  0000 C CNN
F 1 "HDSP-4830_2" H 5550 4676 50  0000 C CNN
F 2 "Display:HDSP-4830" H 5550 3300 50  0001 C CNN
F 3 "https://docs.broadcom.com/docs/AV02-1798EN" H 3550 4300 50  0001 C CNN
	1    5550 4100
	1    0    0    -1  
$EndComp
$Comp
L LED:HDSP-4830_2 BAR2
U 1 1 5E64BB67
P 5550 5700
F 0 "BAR2" H 5550 6367 50  0000 C CNN
F 1 "HDSP-4830_2" H 5550 6276 50  0000 C CNN
F 2 "Display:HDSP-4830" H 5550 4900 50  0001 C CNN
F 3 "https://docs.broadcom.com/docs/AV02-1798EN" H 3550 5900 50  0001 C CNN
	1    5550 5700
	1    0    0    -1  
$EndComp
Wire Wire Line
	5350 5100 5350 4850
Wire Wire Line
	5350 4850 5250 4850
Wire Wire Line
	5250 4850 5250 4500
Wire Wire Line
	5250 4500 5350 4500
Wire Wire Line
	5350 5200 5350 5150
Wire Wire Line
	5350 5150 5400 5150
Wire Wire Line
	5400 5150 5400 4800
Wire Wire Line
	5400 4800 5300 4800
Wire Wire Line
	5300 4800 5300 4600
Wire Wire Line
	5300 4600 5350 4600
$Comp
L Device:R_Network08 RN1
U 1 1 5E651E4A
P 6200 4100
F 0 "RN1" V 5583 4100 50  0000 C CNN
F 1 "R_Network08" V 5674 4100 50  0000 C CNN
F 2 "Resistor_THT:R_Array_SIP9" V 6675 4100 50  0001 C CNN
F 3 "http://www.vishay.com/docs/31509/csc.pdf" H 6200 4100 50  0001 C CNN
	1    6200 4100
	0    1    1    0   
$EndComp
$Comp
L Device:R_Network08 RN2
U 1 1 5E658308
P 6250 5500
F 0 "RN2" V 5633 5500 50  0000 C CNN
F 1 "R_Network08" V 5724 5500 50  0000 C CNN
F 2 "Resistor_THT:R_Array_SIP9" V 6725 5500 50  0001 C CNN
F 3 "http://www.vishay.com/docs/31509/csc.pdf" H 6250 5500 50  0001 C CNN
	1    6250 5500
	0    1    1    0   
$EndComp
Wire Wire Line
	5750 3700 6000 3700
Wire Wire Line
	5750 3800 6000 3800
Wire Wire Line
	5750 3900 6000 3900
Wire Wire Line
	5750 4000 6000 4000
Wire Wire Line
	5750 4100 6000 4100
Wire Wire Line
	5750 4200 6000 4200
Wire Wire Line
	5750 4300 6000 4300
Wire Wire Line
	5750 4400 6000 4400
Wire Wire Line
	5750 4500 6050 4500
Wire Wire Line
	6050 4500 6050 5100
Wire Wire Line
	5750 4600 6000 4600
Wire Wire Line
	6000 4600 6000 5200
Wire Wire Line
	6000 5200 6050 5200
$Comp
L power:GND #PWR0101
U 1 1 5E6750EE
P 6700 6350
F 0 "#PWR0101" H 6700 6100 50  0001 C CNN
F 1 "GND" H 6705 6177 50  0000 C CNN
F 2 "" H 6700 6350 50  0001 C CNN
F 3 "" H 6700 6350 50  0001 C CNN
	1    6700 6350
	1    0    0    -1  
$EndComp
Wire Wire Line
	6400 3700 6700 3700
Wire Wire Line
	6700 3700 6700 5100
Wire Wire Line
	6450 5100 6700 5100
Connection ~ 6700 5100
Wire Wire Line
	6700 5100 6700 6350
$Comp
L Device:C C2
U 1 1 5E67DFA2
P 3350 2650
F 0 "C2" H 3465 2696 50  0000 L CNN
F 1 ".1uF" H 3465 2605 50  0000 L CNN
F 2 "Capacitor_THT:C_Disc_D5.0mm_W2.5mm_P2.50mm" H 3388 2500 50  0001 C CNN
F 3 "~" H 3350 2650 50  0001 C CNN
	1    3350 2650
	1    0    0    -1  
$EndComp
$Comp
L Device:C C3
U 1 1 5E67EAF0
P 3650 2650
F 0 "C3" H 3765 2696 50  0000 L CNN
F 1 ".1uF" H 3765 2605 50  0000 L CNN
F 2 "Capacitor_THT:C_Disc_D5.0mm_W2.5mm_P2.50mm" H 3688 2500 50  0001 C CNN
F 3 "~" H 3650 2650 50  0001 C CNN
	1    3650 2650
	1    0    0    -1  
$EndComp
Connection ~ 3350 2500
Wire Wire Line
	3350 2500 3000 2500
$Comp
L Device:C C1
U 1 1 5E67CF03
P 3000 2650
F 0 "C1" H 3115 2696 50  0000 L CNN
F 1 ".1uF" H 3115 2605 50  0000 L CNN
F 2 "Capacitor_THT:C_Disc_D6.0mm_W2.5mm_P5.00mm" H 3038 2500 50  0001 C CNN
F 3 "~" H 3000 2650 50  0001 C CNN
	1    3000 2650
	1    0    0    -1  
$EndComp
Connection ~ 3000 2500
Entry Wire Line
	2800 3800 2700 3900
Entry Wire Line
	2800 3900 2700 4000
Entry Wire Line
	2800 4000 2700 4100
Entry Wire Line
	2800 4100 2700 4200
Entry Wire Line
	2800 4200 2700 4300
Entry Wire Line
	2800 4300 2700 4400
Entry Wire Line
	2800 4400 2700 4500
Entry Wire Line
	2800 4500 2700 4600
Entry Wire Line
	2800 4600 2700 4700
$Comp
L power:GND #PWR0104
U 1 1 5E6B2FD2
P 2350 5200
F 0 "#PWR0104" H 2350 4950 50  0001 C CNN
F 1 "GND" H 2355 5027 50  0000 C CNN
F 2 "" H 2350 5200 50  0001 C CNN
F 3 "" H 2350 5200 50  0001 C CNN
	1    2350 5200
	1    0    0    -1  
$EndComp
Wire Wire Line
	2200 5100 2350 5100
Wire Wire Line
	2350 5100 2350 5200
Wire Wire Line
	3000 2800 3350 2800
Connection ~ 3350 2800
Wire Wire Line
	5350 5900 5350 6000
Wire Wire Line
	5350 6100 5350 6200
Wire Wire Line
	5350 6000 5350 6100
Connection ~ 5350 6000
Connection ~ 5350 6100
Wire Wire Line
	5750 5900 5750 6000
Wire Wire Line
	5750 6000 5750 6100
Connection ~ 5750 6000
Wire Wire Line
	5750 6100 5750 6200
Connection ~ 5750 6100
Entry Wire Line
	2800 3850 2900 3750
Text Label 2900 3750 0    50   ~ 0
PA2
$Comp
L Connector:Conn_01x10_Male J1
U 1 1 5E728117
P 2000 4300
F 0 "J1" H 2108 4881 50  0000 C CNN
F 1 "Conn_01x10_Male" H 2108 4790 50  0000 C CNN
F 2 "Connector_PinHeader_2.54mm:PinHeader_1x10_P2.54mm_Vertical" H 2000 4300 50  0001 C CNN
F 3 "~" H 2000 4300 50  0001 C CNN
	1    2000 4300
	1    0    0    -1  
$EndComp
Wire Wire Line
	2200 4800 2200 5100
Wire Wire Line
	2900 2500 3000 2500
Connection ~ 3650 2500
Wire Wire Line
	3650 2500 4850 2500
Wire Wire Line
	3750 2900 4850 2900
Wire Wire Line
	4850 2900 4850 2500
Entry Wire Line
	4100 3800 4200 3700
Entry Wire Line
	4100 3900 4200 3800
Entry Wire Line
	4100 4000 4200 3900
Entry Wire Line
	4100 4100 4200 4000
Entry Wire Line
	4100 4200 4200 4100
Entry Wire Line
	4100 4300 4200 4200
Entry Wire Line
	4100 4400 4200 4300
Entry Wire Line
	4100 4500 4200 4400
Entry Wire Line
	4100 5200 4200 5100
Entry Wire Line
	4100 5300 4200 5200
Entry Wire Line
	4100 5400 4200 5300
Entry Wire Line
	4100 5500 4200 5400
Entry Wire Line
	4100 5600 4200 5500
Entry Wire Line
	4100 5700 4200 5600
Entry Wire Line
	4100 5800 4200 5700
Entry Wire Line
	4100 5900 4200 5800
Wire Wire Line
	4200 3700 4250 3700
Wire Wire Line
	4200 3800 4250 3800
Wire Wire Line
	4200 3900 4250 3900
Wire Wire Line
	4200 4000 4250 4000
Wire Wire Line
	4200 4100 4250 4100
Wire Wire Line
	4200 4200 4250 4200
Wire Wire Line
	4200 4300 4250 4300
Wire Wire Line
	4200 4400 4250 4400
Wire Wire Line
	4200 5100 4250 5100
Wire Wire Line
	4200 5200 4250 5200
Wire Wire Line
	4200 5300 4250 5300
Wire Wire Line
	4200 5400 4250 5400
Wire Wire Line
	4200 5500 4250 5500
Wire Wire Line
	4200 5600 4250 5600
Wire Wire Line
	4200 5800 4250 5800
Wire Wire Line
	4200 5700 4250 5700
Text Label 4200 3700 0    50   ~ 0
PB0
Text Label 4200 3800 0    50   ~ 0
PB1
Text Label 4200 3900 0    50   ~ 0
PB2
Text Label 4200 4000 0    50   ~ 0
PB3
Text Label 4200 4100 0    50   ~ 0
PB4
Text Label 4200 4200 0    50   ~ 0
PB5
Text Label 4200 4300 0    50   ~ 0
PB6
Text Label 4200 4400 0    50   ~ 0
PB7
Text Label 4200 5100 0    50   ~ 0
PB0
Text Label 4200 5200 0    50   ~ 0
PB1
Text Label 4200 5300 0    50   ~ 0
PB2
Text Label 4200 5400 0    50   ~ 0
PB3
Text Label 4200 5500 0    50   ~ 0
PB4
Text Label 4200 5600 0    50   ~ 0
PB5
Text Label 4200 5700 0    50   ~ 0
PB6
Text Label 4200 5800 0    50   ~ 0
PB7
Wire Wire Line
	2200 3900 2700 3900
Wire Wire Line
	2200 4000 2700 4000
Wire Wire Line
	2200 4100 2700 4100
Wire Wire Line
	2200 4200 2700 4200
Wire Wire Line
	2200 4300 2700 4300
Wire Wire Line
	2200 4400 2700 4400
Wire Wire Line
	2200 4500 2700 4500
Wire Wire Line
	2200 4600 2700 4600
Wire Wire Line
	2200 4700 2700 4700
Text Label 2350 3900 0    50   ~ 0
PB0
Text Label 2350 4000 0    50   ~ 0
PB1
Text Label 2350 4100 0    50   ~ 0
PB2
Text Label 2350 4200 0    50   ~ 0
PB3
Text Label 2350 4300 0    50   ~ 0
PB4
Text Label 2350 4400 0    50   ~ 0
PB5
Text Label 2350 4500 0    50   ~ 0
PB6
Text Label 2350 4600 0    50   ~ 0
PB7
Text Label 2350 4700 0    50   ~ 0
PA2
$Comp
L Transistor_BJT:BC547 Q1
U 1 1 5E7D5AC0
P 3650 3750
F 0 "Q1" H 3841 3796 50  0000 L CNN
F 1 "BC547" H 3841 3705 50  0000 L CNN
F 2 "Package_TO_SOT_THT:TO-92_Inline" H 3850 3675 50  0001 L CIN
F 3 "http://www.fairchildsemi.com/ds/BC/BC547.pdf" H 3650 3750 50  0001 L CNN
	1    3650 3750
	1    0    0    -1  
$EndComp
$Comp
L Device:R R2
U 1 1 5E7D82FC
P 3750 3250
F 0 "R2" H 3820 3296 50  0000 L CNN
F 1 "1K" H 3820 3205 50  0000 L CNN
F 2 "Resistor_THT:R_Axial_DIN0204_L3.6mm_D1.6mm_P5.08mm_Horizontal" V 3680 3250 50  0001 C CNN
F 3 "~" H 3750 3250 50  0001 C CNN
	1    3750 3250
	1    0    0    -1  
$EndComp
Wire Wire Line
	3750 2900 3750 3100
Wire Wire Line
	3750 3400 3750 3500
Wire Wire Line
	3750 3500 4250 3500
Connection ~ 3750 3500
Wire Wire Line
	3750 3500 3750 3550
$Comp
L Device:R R1
U 1 1 5E7F2E86
P 3250 3750
F 0 "R1" V 3043 3750 50  0000 C CNN
F 1 "100K" V 3134 3750 50  0000 C CNN
F 2 "Resistor_THT:R_Axial_DIN0204_L3.6mm_D1.6mm_P5.08mm_Horizontal" V 3180 3750 50  0001 C CNN
F 3 "~" H 3250 3750 50  0001 C CNN
	1    3250 3750
	0    1    1    0   
$EndComp
Wire Wire Line
	3450 3750 3400 3750
Wire Wire Line
	2900 3750 3000 3750
Wire Bus Line
	2800 6000 4100 6000
Connection ~ 3000 2800
Connection ~ 2600 2800
Wire Wire Line
	2300 2800 2600 2800
Wire Wire Line
	2600 2800 3000 2800
Wire Wire Line
	2300 2700 2300 2800
$Comp
L Connector:Barrel_Jack J2
U 1 1 5E86DF3A
P 1300 2600
F 0 "J2" H 1071 2558 50  0000 R CNN
F 1 "Barrel_Jack" H 1071 2649 50  0000 R CNN
F 2 "Connector_BarrelJack:BarrelJack_CUI_PJ-063AH_Horizontal_CircularHoles" H 1350 2560 50  0001 C CNN
F 3 "~" H 1350 2560 50  0001 C CNN
	1    1300 2600
	1    0    0    1   
$EndComp
Wire Wire Line
	3350 2500 3650 2500
Wire Wire Line
	3350 2800 3650 2800
Wire Wire Line
	1600 2700 2300 2700
$Comp
L power:GND #PWR0105
U 1 1 5E8D7755
P 3750 5100
F 0 "#PWR0105" H 3750 4850 50  0001 C CNN
F 1 "GND" H 3755 4927 50  0000 C CNN
F 2 "" H 3750 5100 50  0001 C CNN
F 3 "" H 3750 5100 50  0001 C CNN
	1    3750 5100
	1    0    0    -1  
$EndComp
$Comp
L 74xx_IEEE:74LS541 U1
U 1 1 5E910F75
P 4800 3950
F 0 "U1" H 4800 4816 50  0000 C CNN
F 1 "74LS541" H 4800 4725 50  0000 C CNN
F 2 "Package_DIP:DIP-20_W7.62mm_Socket" H 4800 3950 50  0001 C CNN
F 3 "" H 4800 3950 50  0001 C CNN
	1    4800 3950
	1    0    0    -1  
$EndComp
Connection ~ 4250 3500
$Comp
L 74xx_IEEE:74LS541 U2
U 1 1 5E91DBB1
P 4800 5350
F 0 "U2" H 4800 6216 50  0000 C CNN
F 1 "74LS541" H 4800 6125 50  0000 C CNN
F 2 "Package_DIP:DIP-20_W7.62mm_Socket" H 4800 5350 50  0001 C CNN
F 3 "" H 4800 5350 50  0001 C CNN
	1    4800 5350
	1    0    0    -1  
$EndComp
Connection ~ 4250 4750
Wire Wire Line
	5750 5300 6050 5300
Wire Wire Line
	5750 5400 6050 5400
Wire Wire Line
	5750 5500 6050 5500
Wire Wire Line
	5750 5600 6050 5600
Wire Wire Line
	5750 5700 6050 5700
Wire Wire Line
	5750 5800 6050 5800
$Comp
L power:VCC #PWR0103
U 1 1 5E95C329
P 3650 2350
F 0 "#PWR0103" H 3650 2200 50  0001 C CNN
F 1 "VCC" H 3667 2523 50  0000 C CNN
F 2 "" H 3650 2350 50  0001 C CNN
F 3 "" H 3650 2350 50  0001 C CNN
	1    3650 2350
	1    0    0    -1  
$EndComp
Wire Wire Line
	3650 2350 3650 2500
$Comp
L power:GND #PWR0102
U 1 1 5E96B587
P 2600 2900
F 0 "#PWR0102" H 2600 2650 50  0001 C CNN
F 1 "GND" H 2605 2727 50  0000 C CNN
F 2 "" H 2600 2900 50  0001 C CNN
F 3 "" H 2600 2900 50  0001 C CNN
	1    2600 2900
	1    0    0    -1  
$EndComp
Wire Wire Line
	2600 2800 2600 2850
Wire Wire Line
	2600 2850 2450 2850
Wire Wire Line
	2450 2850 2450 2950
Connection ~ 2600 2850
Wire Wire Line
	2600 2850 2600 2900
$Comp
L power:PWR_FLAG #FLG0102
U 1 1 5E96F830
P 2450 2950
F 0 "#FLG0102" H 2450 3025 50  0001 C CNN
F 1 "PWR_FLAG" H 2450 3123 50  0000 C CNN
F 2 "" H 2450 2950 50  0001 C CNN
F 3 "~" H 2450 2950 50  0001 C CNN
	1    2450 2950
	-1   0    0    1   
$EndComp
$Comp
L Device:D D1
U 1 1 5E64C9E5
P 1850 2500
F 0 "D1" H 1850 2284 50  0000 C CNN
F 1 "GP0240" H 1850 2375 50  0000 C CNN
F 2 "Diode_THT:D_DO-34_SOD68_P7.62mm_Horizontal" H 1850 2500 50  0001 C CNN
F 3 "~" H 1850 2500 50  0001 C CNN
	1    1850 2500
	-1   0    0    1   
$EndComp
Wire Wire Line
	2000 2500 2100 2500
$Comp
L Regulator_Linear:L7805 U3
U 1 1 5E662957
P 2600 2500
F 0 "U3" H 2600 2742 50  0000 C CNN
F 1 "L7805" H 2600 2651 50  0000 C CNN
F 2 "Package_TO_SOT_THT:TO-220-3_Vertical" H 2625 2350 50  0001 L CIN
F 3 "http://www.st.com/content/ccc/resource/technical/document/datasheet/41/4f/b3/b0/12/d4/47/88/CD00000444.pdf/files/CD00000444.pdf/jcr:content/translations/en.CD00000444.pdf" H 2600 2450 50  0001 C CNN
	1    2600 2500
	1    0    0    -1  
$EndComp
$Comp
L power:PWR_FLAG #FLG0101
U 1 1 5E664A91
P 1650 2200
F 0 "#FLG0101" H 1650 2275 50  0001 C CNN
F 1 "PWR_FLAG" H 1650 2373 50  0000 C CNN
F 2 "" H 1650 2200 50  0001 C CNN
F 3 "~" H 1650 2200 50  0001 C CNN
	1    1650 2200
	1    0    0    -1  
$EndComp
Wire Wire Line
	1650 2200 1650 2300
Wire Wire Line
	1650 2300 2100 2300
Wire Wire Line
	2100 2300 2100 2500
Connection ~ 2100 2500
Wire Wire Line
	2100 2500 2300 2500
Wire Wire Line
	1600 2500 1700 2500
Wire Wire Line
	3000 3750 3000 4750
Wire Wire Line
	3000 4750 4250 4750
Connection ~ 3000 3750
Wire Wire Line
	3000 3750 3100 3750
Wire Wire Line
	3750 3900 3750 3950
Wire Bus Line
	2800 3800 2800 6000
Wire Bus Line
	4100 3700 4100 6000
Connection ~ 3750 3950
Wire Wire Line
	3750 3950 3750 5100
$EndSCHEMATC
