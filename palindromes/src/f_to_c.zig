const std = @import("std");
const stdin = std.io.getStdIn();
const stdout = std.io.getStdOut().writer();
const BUFFER_SIZE = 50;

fn convert_to_celsius(fahrenheit: f64) f64 {
    return (fahrenheit - 32) * 5.0 / 9.0;
}

test "to celsius" {
    try std.testing.expectEqual(0.0, convert_to_celsius(32));
    try std.testing.expectApproxEqAbs(-0.888888, convert_to_celsius(30.4), 0.000001);
    try std.testing.expectEqual(-40.0, convert_to_celsius(-40));
    try std.testing.expectApproxEqAbs(100, convert_to_celsius(212), 0.000001);
    try std.testing.expectApproxEqAbs(-273.15, convert_to_celsius(-459.67), 0.000001);
}

fn convert_record(l: []const u8, converted_line: []u8) !void {
    var it = std.mem.splitAny(u8, l, " C");

    // get the temperature part of the <temp> <unit> pair
    if (it.next()) |temperature| {
        // check if the unit part is F for fahernheit
        if (std.mem.eql(u8, it.peek() orelse "", "F")) {
            // convert the fharenheit temperature to celsius and print the new value in C
            const fahrenheit = std.fmt.parseFloat(f64, temperature) catch |err| {
                std.debug.print("Error: {s} value: {s} cannot convert\n", .{ @errorName(err), temperature });
                return;
            };
            //if the unit is not F then print the original value
            _ = try std.fmt.bufPrint(converted_line, "{d:.1} C\n", .{convert_to_celsius(fahrenheit)});
        } else {
            _ = try std.fmt.bufPrint(converted_line, "{s}\n", .{l});
        }
    }
}

test "convert record C" {
    var converted_line = std.mem.zeroes([BUFFER_SIZE]u8);
    try convert_record("32 C", &converted_line);
    try std.testing.expectStringStartsWith(&converted_line, "32 C\n");

    converted_line = std.mem.zeroes([BUFFER_SIZE]u8);
    try convert_record("30 F", &converted_line);
    try std.testing.expectStringStartsWith(&converted_line, "-1.1 C\n");

    converted_line = std.mem.zeroes([BUFFER_SIZE]u8);
    try convert_record("72.2 F", &converted_line);
    try std.testing.expectStringStartsWith(&converted_line, "22.3 C\n");
}

// we read a list of words delimited by newlines from stdin
// each word is checked to see if it is a palindrome
// when it is, it's printed to the stdout
pub fn main() !void {
    const std_reader = stdin.reader();
    var br = std.io.bufferedReader(std_reader);
    var buffer: [BUFFER_SIZE]u8 = undefined;

    while (br.reader().readUntilDelimiterOrEof(&buffer, '\n') catch |err| {
        std.debug.print("Error: {s}", .{@errorName(err)});
        return;
    }) |l| {
        if (l.len > BUFFER_SIZE) {
            continue;
        }

        var converted_line = std.mem.zeroes([BUFFER_SIZE]u8);

        convert_record(l, &converted_line) catch |err| {
            std.debug.print("Error: {s}", .{@errorName(err)});
            continue;
        };

        stdout.print("{s}", .{converted_line}) catch |err| {
            std.debug.print("Error: {s}\n", .{@errorName(err)});
            continue;
        };
    }
}