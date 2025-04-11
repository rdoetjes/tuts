ORG $801

jmp main

// 256-entry sine lookup table (values from 0 to 255, centered at 128)
sine_table:
.byte 128, 131, 134, 137, 140, 143, 146, 149, 152, 156, 159, 162, 165, 168, 171, 174
.byte 177, 180, 183, 186, 189, 191, 194, 197, 200, 202, 205, 207, 210, 212, 214, 217
.byte 219, 221, 223, 225, 227, 229, 231, 232, 234, 236, 237, 239, 240, 241, 243, 244
.byte 245, 246, 247, 248, 249, 250, 250, 251, 251, 252, 252, 252, 252, 252, 252, 252
.byte 252, 252, 251, 251, 250, 250, 249, 248, 247, 246, 245, 244, 243, 241, 240, 239
.byte 237, 236, 234, 232, 231, 229, 227, 225, 223, 221, 219, 217, 214, 212, 210, 207
.byte 205, 202, 200, 197, 194, 191, 189, 186, 183, 180, 177, 174, 171, 168, 165, 162
.byte 159, 156, 152, 149, 146, 143, 140, 137, 134, 131, 128, 125, 122, 119, 116, 113
.byte 110, 107, 104, 100, 97, 94, 91, 88, 85, 82, 79, 76, 73, 70, 67, 65, 62, 59
.byte 56, 54, 51, 49, 46, 44, 42, 39, 37, 35, 33, 31, 29, 27, 25, 24
.byte 22, 20, 19, 17, 16, 15, 13, 12, 11, 10, 9, 8, 7, 6, 6, 5
.byte 5, 4, 4, 4, 4, 4, 4, 4, 4, 5, 5, 6, 6, 7, 8, 9
.byte 10, 11, 12, 13, 15, 16, 17, 19, 20, 22, 24, 25, 27, 29, 31, 33
.byte 35, 37, 39, 42, 44, 46, 49, 51, 54, 56, 59, 62, 65, 67, 70, 73
.byte 76, 79, 82, 85, 88, 91, 94, 97, 100, 104, 107, 110, 113, 116, 119, 122
.byte 125

// cosine is just sine shifted by 90 degrees (64 index positions)
cosine_table:
.byte sine_table + 64  //g Cosine is the same as sine but offset by 64

//Usage Example**
//To get sine or cosine:
//- Load the angle into the index register (`X`).
//- Read from the table.

main:
    LDX angle         //g Load angle into X register
    LDA sine_table,X  //g Get sine value
    STA result        //g Store result

    LDX angle         //g Load angle again
    LDA cosine_table,X //g Get cosine value
    STA result_cos    //g Store result

