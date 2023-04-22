const std = @import("std");
const stdout = std.io.getStdOut().writer();

pub fn main() !void {
    var biglist = std.ArrayList(u8).init(std.heap.page_allocator);
    defer biglist.deinit();

    var index: u32 = 0;
    while (index < 1000001) : (index += 1) {
        if (index % 10 == 7 or index % 7 == 0) {
            try biglist.writer().print("{s}\n", .{"SMAC"});
        } else {
            try biglist.writer().print("{d}\n", .{index});
        }
    }

    try print_list(&biglist);
    biglist.clearAndFree();
}

fn print_list(list: *std.ArrayList(u8)) !void {
    var iter = std.mem.split(u8, list.items, "\n");
    while (iter.next()) |item| {
        try stdout.print("{s}\n", .{item});
    }
}
