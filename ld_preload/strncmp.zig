const std = @import("std");
export fn strncmp(s1: [*c]u8, s2: [*c]u8, len: usize) u32 {
    std.debug.print("{s} {s} {d}\n", .{ s1, s2, len });
    return 0;
}
