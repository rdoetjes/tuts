: TEST-IF 0 = IF 111 ELSE 222 THEN . ;
0 TEST-IF
5 TEST-IF

: IS-NEG 0 < IF 333 . THEN ;
0 5 - IS-NEG
5 IS-NEG

: TEST-PRINT ." Hello world" ;
TEST-PRINT

." Interpretation test"
