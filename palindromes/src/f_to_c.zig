const std = @import("std");
const stdin = std.io.getStdIn();
const stdout = std.io.getStdOut().writer();
const BUFFER_SIZE = 1024;

fn convert_to_celsius(fahrenheit: f64) f64 {
    return (fahrenheit - 32) * 5.0 / 9.0;
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
        var it = std.mem.splitAny(u8, l, " C");

        // get the temperature part of the <temp> <unit> pair
        if (it.next()) |temperature| {
            // check if the unit part is F for fahernheit
            if (std.mem.eql(u8, it.peek() orelse "", "F")) {
                // convert the fharenheit temperature to celsius and print the new value in C
                const fahrenheit = std.fmt.parseFloat(f64, temperature) catch |err| {
                    std.debug.print("Error: {s} value: {s} cannot convert\n", .{ @errorName(err), temperature });
                    continue;
                };
                try stdout.print("{d:.1} C\n", .{convert_to_celsius(fahrenheit)});
            } else {
                //if the unit is not F then print the original value
                std.debug.print("{s}\n", .{l});
            }
        }
    }
}
