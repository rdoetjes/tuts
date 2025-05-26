const std = @import("std");

fn to_roman(alloc: std.mem.Allocator, year: u32) ![]const u8 {
    const dec_roman = struct {
        roman: []const u8,
        dec: u16,
    };

    const roman = [_]dec_roman{
        dec_roman{ .roman = "M", .dec = 1000 },
        dec_roman{ .roman = "CM", .dec = 900 },
        dec_roman{ .roman = "D", .dec = 500 },
        dec_roman{ .roman = "CD", .dec = 400 },
        dec_roman{ .roman = "C", .dec = 100 },
        dec_roman{ .roman = "XC", .dec = 90 },
        dec_roman{ .roman = "L", .dec = 50 },
        dec_roman{ .roman = "XL", .dec = 40 },
        dec_roman{ .roman = "X", .dec = 10 },
        dec_roman{ .roman = "IX", .dec = 9 },
        dec_roman{ .roman = "V", .dec = 5 },
        dec_roman{ .roman = "IV", .dec = 4 },
        dec_roman{ .roman = "I", .dec = 1 },
    };

    var i: usize = 0;
    var co = year;
    var result: []u8 = "";
    while (co != 0) {
        if (co < roman[i].dec) {
            i += 1;
        } else {
            // not the most efficient but the least amount of code
            result = try std.fmt.allocPrint(alloc, "{s}{s}", .{ result, roman[i].roman });
            co -= roman[i].dec;
        }
    }
    return result;
}

pub fn main() !void {
    const stdin = std.io.getStdIn().reader();
    const stdout = std.io.getStdOut().writer();
    const alloc = std.heap.page_allocator;
    const line_opt = stdin.readUntilDelimiterOrEofAlloc(alloc, '\n', 20) catch {
        std.debug.print("string too long (no more than 20 chars)", .{});
        std.posix.exit(1);
    };
    defer if (line_opt) |line| alloc.free(line);

    const line = line_opt orelse {
        std.posix.exit(1);
    };

    const co = std.fmt.parseInt(u32, line, 10) catch |err| {
        std.debug.print("Error trying to convert to int: {any}", .{err});
        std.posix.exit(1);
    };

    const result = to_roman(alloc, co) catch |err| {
        std.debug.print("Error trying to convert to roman: {any}", .{err});
        std.posix.exit(1);
    };
    defer alloc.free(result);

    try stdout.print("{s}\n", .{result});
}
