0 data 33, 17, 65, -1
5 s=54272 
10 for l=s to s+24 : poke l,0 : next 
20 poke s+5,0 : poke s+6,64 
40 poke s+1,25 : poke s,177
50 poke s+2,25 : poke s+3,255
80 read a
90 if a = -1 goto 1000
100 print a
110 poke s+4, a
120 poke s+24,15
130 for t=1to2500:next t
140 poke s+24,0
150 for t=1to1000:next t
160 poke s+23,241:poke s+21,255:poke s+22,255:poke s+24,31
170 for t=0to255
180 poke s+21,t
190 poke s+22,t
200 next t
210 poke s+24,0:poke s+23,0
220 for t=1to1000:next t
860 goto 80
1000 poke s+24,0
