const std = @import("std");

fn to_roman(alloc: std.mem.Allocator, year: u32) ![]const u8 {
    const RomanEntry = struct {
        numeral: []const u8,
        value: u16,
    };

    const numerals = [_]RomanEntry{
        .{ .numeral = "M", .value = 1000 },
        .{ .numeral = "CM", .value = 900 },
        .{ .numeral = "D", .value = 500 },
        .{ .numeral = "CD", .value = 400 },
        .{ .numeral = "C", .value = 100 },
        .{ .numeral = "XC", .value = 90 },
        .{ .numeral = "L", .value = 50 },
        .{ .numeral = "XL", .value = 40 },
        .{ .numeral = "X", .value = 10 },
        .{ .numeral = "IX", .value = 9 },
        .{ .numeral = "V", .value = 5 },
        .{ .numeral = "IV", .value = 4 },
        .{ .numeral = "I", .value = 1 },
    };

    var result = try alloc.alloc(u8, 0); // start empty
    var remainder = year;
    var i: usize = 0;

    while (remainder != 0) {
        if (remainder < numerals[i].value) {
            i += 1;
        } else {
            result = try std.fmt.allocPrint(alloc, "{s}{s}", .{ result, numerals[i].numeral }); // not most effient but nice and short (no need to arraylist and flatten)
            remainder -= numerals[i].value;
        }
    }
    return result;
}

fn write_error_and_exit(msg: []const u8, err: anyerror) void {
    std.debug.print("{s}\n Error: {any}\n\n", .{ msg, err });
    std.posix.exit(1);
}

pub fn main() !void {
    const MAX_LENGTH = 5;
    const stdin = std.io.getStdIn().reader();
    const stdout = std.io.getStdOut().writer();
    const alloc = std.heap.page_allocator;

    const line = stdin.readUntilDelimiterOrEofAlloc(alloc, '\n', MAX_LENGTH + 1) catch |err| {
        write_error_and_exit("Input too long", err);
        unreachable;
    } orelse unreachable;
    defer alloc.free(line);

    const year = std.fmt.parseInt(u32, line, 10) catch |err| {
        write_error_and_exit("Error parsing input.", err);
        unreachable;
    };

    const roman = to_roman(alloc, year) catch |err| {
        write_error_and_exit("Error converting Roman numeral.", err);
        unreachable;
    };
    defer alloc.free(roman);

    try stdout.print("{s}\n", .{roman});
}
