const std = @import("std");
const net = std.net;
const os = std.os;
var general_purpose_allocator = std.heap.GeneralPurposeAllocator(.{}){};
const gpa = general_purpose_allocator.allocator();

pub fn main() !void {
    var biglist = std.ArrayList(u8).init(std.heap.page_allocator);
    defer biglist.deinit();

    var i: u32 = 1;
    while (i <= 1_000_000) : (i += 1) {
        if (i % 10 == 7 or i % 7 == 0) {
            try biglist.writer().print("{s}\n|", .{"SMAC"});
        } else {
            try biglist.writer().print("{d}\n|", .{i});
        }
    }

    var server = net.StreamServer.init(.{});
    server.reuse_address = true;
    defer server.deinit();
    try server.listen(net.Address.parseIp("0.0.0.0", 7979) catch unreachable);
    std.debug.print("Listening on {}\n", .{server.listen_address});

    while (true) {
        var conn = try server.accept();
        _ = try std.Thread.spawn(.{}, print_list, .{ &biglist, conn });
    }
}

fn print_list(biglist: *std.ArrayList(u8), conn: net.StreamServer.Connection) !void {
    var iter = std.mem.split(u8, biglist.items, "|");
    while (iter.next()) |item| {
        _ = try conn.stream.writeAll(item);
    }
    conn.stream.close();
}
