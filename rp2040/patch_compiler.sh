sed -i 's/add r3, r0; adds r3, #3; lsrs r3, #2; lsls r3, #2; str r3, \[r2\]/add r3, r0; adds r3, #3; lsrs r3, #2; lsls r3, #2; ldr r6, =0x46c0b500; str r6, [r3]; adds r3, #4; str r3, [r2]/g' forth.s
sed -i 's/ldr r2, =0x4770; strh r2, \[r1\]; adds r1, #2/ldr r2, =0x46c0bd00; str r2, [r1]; adds r1, #4/g' forth.s
