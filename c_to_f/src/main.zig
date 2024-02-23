const std = @import("std");
const stdin = std.io.getStdIn();
const stdout = std.io.getStdOut().writer();
const BUFFER_SIZE = 50;
const allocator = std.heap.page_allocator;

fn convert_to_celsius(fahrenheit: f64) f64 {
    return (fahrenheit - 32) * 5.0 / 9.0;
}

test "convert fahrenheit to celsius" {
    try std.testing.expectEqual(0.0, convert_to_celsius(32));
    try std.testing.expectApproxEqAbs(-0.888888, convert_to_celsius(30.4), 0.000001);
    try std.testing.expectEqual(-40.0, convert_to_celsius(-40));
    try std.testing.expectApproxEqAbs(100, convert_to_celsius(212), 0.000001);
    try std.testing.expectApproxEqAbs(-273.15, convert_to_celsius(-459.67), 0.000001);
}

// l takes a record that contains a temperature and a unit f.i. 30.0 F
// the temperature is converted to celsius if the unit is F.
// otherwise the original temperature is printed.
fn convert_record(l: []const u8, converted_line: []u8) !void {
    var it = std.mem.splitAny(u8, l, " ");

    // get the temperature part of the <temp> <unit> pair
    if (it.next()) |temperature| {
        // check if the unit part is not F for fahernheit
        if (!std.mem.eql(u8, it.peek() orelse return, "F")) {
            // if it's temperature is not in F then print the C record
            _ = try std.fmt.bufPrint(converted_line, "{s}\n", .{l});
            return;
        }

        // the unit is F, therefore convert the fahrenheit string to float
        const fahrenheit = std.fmt.parseFloat(f64, temperature) catch |err| {
            std.debug.print("Error: {s} value: {s} cannot convert\n", .{ @errorName(err), temperature });
            return;
        };
        //create a string with the converted temperature in C with 1 decimal place
        _ = try std.fmt.bufPrint(converted_line, "{d:.1} C\n", .{convert_to_celsius(fahrenheit)});
    }
}

test "convert record to celsius" {
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

test "convert faulty record " {
    //should keep original value
    var converted_line = std.mem.zeroes([BUFFER_SIZE]u8);
    try convert_record("32", &converted_line);
    try std.testing.expectStringStartsWith(&converted_line, "");

    converted_line = std.mem.zeroes([BUFFER_SIZE]u8);
    try convert_record("hello world", &converted_line);
    try std.testing.expectStringStartsWith(&converted_line, "hello world\n");

    //should throw error
    converted_line = std.mem.zeroes([BUFFER_SIZE]u8);
    convert_record("far toooooooooooooooooooooo long for the record soooo what now?", &converted_line) catch {
        try std.testing.expect(true);
        return;
    };
    try std.testing.expect(false);
}

// we read a list of records delimited by newlines from stdin
// records contain tenperature and unit.
pub fn main() !u8 {
    const std_reader = stdin.reader();
    var br = std.io.bufferedReader(std_reader);

    while (br.reader().readUntilDelimiterOrEofAlloc(allocator, '\n', 10 * BUFFER_SIZE) catch |err| {
        std.debug.print("Error: {s}", .{@errorName(err)});
        return 1;
    }) |l| {
        defer allocator.free(l);

        const converted_line = allocator.alloc(u8, l.len + 10) catch |err| {
            std.debug.print("Error: {s}", .{@errorName(err)});
            return 1;
        };
        defer allocator.free(converted_line);
        @memset(converted_line, 0);

        convert_record(l, converted_line) catch |err| {
            std.debug.print("Error: {s} converting record, skipping record\n", .{@errorName(err)});
            return 1;
        };

        stdout.print("{s}", .{converted_line}) catch |err| {
            std.debug.print("Error: {s} skipping record\n", .{@errorName(err)});
            return 1;
        };
    }
    return 0;
}
