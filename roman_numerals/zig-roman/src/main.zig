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
            // Append the numeral to result, not most effient but nice and short (no need to arraylist and flatten)
            result = try std.fmt.allocPrint(alloc, "{s}{s}", .{ result, numerals[i].numeral });
            remainder -= numerals[i].value;
        }
    }
    return result;
}

pub fn main() !void {
    const MAX_LENGTH = 6;
    const stdin = std.io.getStdIn().reader();
    const stdout = std.io.getStdOut().writer();
    const alloc = std.heap.page_allocator;

    // Read input line
    const line = stdin.readUntilDelimiterOrEofAlloc(alloc, '\n', MAX_LENGTH) catch {
        std.debug.print("Input too long (max {d} chars).\n", .{MAX_LENGTH - 1});
        std.posix.exit(1);
    } orelse unreachable;

    defer alloc.free(line);

    // Convert input to number
    const year = std.fmt.parseInt(u32, line, 10) catch |err| {
        std.debug.print("Error parsing input: {any}\n", .{err});
        std.posix.exit(1);
    };

    // Convert to Roman numeral
    const roman = to_roman(alloc, year) catch |err| {
        std.debug.print("Error converting to Roman: {any}\n", .{err});
        std.posix.exit(1);
    };
    defer alloc.free(roman);

    try stdout.print("{s}\n", .{roman});
}
