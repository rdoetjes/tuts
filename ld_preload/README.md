* Building shared lib for Zig

```
zig build-lib -dynamic -isystem . strncmp.zig   
```

** c style strings in Zig
```
const c_string: [*c]const u8 = "some c string";
```

* loading the ld_preload
```
export LD_PRELOAD=<absolute_path>/libstrncmp.so;./demo my_argument
```
