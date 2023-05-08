const std = @import("std");
const net = std.net;
const os = std.os;
const stdout = std.io.getStdOut().writer();
var general_purpose_allocator = std.heap.GeneralPurposeAllocator(.{}){};
const gpa = general_purpose_allocator.allocator();

pub fn main() !void {
    var biglist = std.ArrayList(u8).init(std.heap.page_allocator);
    defer biglist.deinit();

    for (1..1000001) |index| {
        if (index % 10 == 7 or index % 7 == 0) {
            try biglist.writer().print("{s}\n|", .{"SMAC"});
        } else {
            try biglist.writer().print("{d}\n|", .{index});
        }
    }

    var server = net.StreamServer.init(.{});
    server.reuse_address = true;
    defer server.deinit();
    server.listen(net.Address.parseIp("0.0.0.0", 7979) catch unreachable);
    try stdout.print("Listening on {}\n", .{server.listen_address});

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
