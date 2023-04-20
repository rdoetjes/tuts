const std = @import("std");
const rnd_gen = std.rand.DefaultPrng;
const heap_alloc = std.heap.page_allocator;
const stdout = std.io.getStdOut().writer();

pub fn main() !void {
    const x = 4;
    var t: [x]std.Thread = undefined;
    var list: [x]std.ArrayList(u8) = undefined;

    for (0..x) |i| {
        list[i] = std.ArrayList(u8).init(heap_alloc);
        defer list[i].deinit();
        t[i] = try std.Thread.spawn(.{}, count_it, .{ @truncate(u32, i), &list[i] });
    }

    //we join all te threads to wait till they're done counting
    for (0..x) |i| {
        t[i].join();
        try print_list(&list[i]);
        //let's clean and free the indiviual list after copying
        //because we also have immediate control, you gotta love Zig
        list[i].clearAndFree();
    }
}

fn print_list(list: *std.ArrayList(u8)) !void {
    var iter = std.mem.split(u8, list.items, "|");
    while (iter.next()) |item| {
        try stdout.print("{s}", .{item});
    }
}

fn count_it(id: u32, list: *std.ArrayList(u8)) !void {
    var addr: usize = 0;
    if (id == 3) addr += 1;

    //var rnd = rnd_gen.init(0);
    for (id * 250000..(id * 250000) + 250000 + addr) |index| {
        if (index == 0) continue;

        if (index % 7 == 0 or index % 10 == 7) {
            try list.writer().print("{s}\n|", .{"SMAC"});
        } else {
            try list.writer().print("{d}\n|", .{index});
        }
    }
}
