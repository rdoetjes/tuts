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
fn convert_record(l: []const u8) ![]const u8 {
    var it = std.mem.splitAny(u8, l, " ");

    // get the temperature part of the <temp> <unit> pair
    if (it.next()) |temperature| {
        // check if the unit part is not F for fahernheit
        if (!std.mem.eql(u8, it.peek() orelse "", "F")) {
            // if it's temperature is not in F then print the C record
            return std.fmt.allocPrint(allocator, "{s}\n", .{l});
        }

        // the unit is F, therefore convert the fahrenheit string to float
        const fahrenheit = std.fmt.parseFloat(f64, temperature) catch |err| {
            std.debug.print("Error: {s} value: {s} cannot convert\n", .{ @errorName(err), temperature });
            return "";
        };
        //create a string with the converted temperature in C with 1 decimal place
        return std.fmt.allocPrint(allocator, "{d:.1} C\n", .{convert_to_celsius(fahrenheit)});
    }
    return "";
}

test "convert record to celsius" {
    try std.testing.expectStringStartsWith(try convert_record("32 C"), "32 C\n");
    try std.testing.expectStringStartsWith(try convert_record("30 F"), "-1.1 C\n");
    try std.testing.expectStringStartsWith(try convert_record("72.2 F"), "22.3 C\n");
}

test "convert faulty record " {
    //should keep original value
    try std.testing.expectStringStartsWith(try convert_record("32"), "");
    try std.testing.expectStringStartsWith(try convert_record("hello world"), "hello world\n");
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

        const converted_line = convert_record(l) catch |err| {
            std.debug.print("Error: {s} converting record, {s}\n", .{ @errorName(err), l });
            return 1;
        };
        defer allocator.free(converted_line);

        stdout.print("{s}", .{converted_line}) catch |err| {
            std.debug.print("Error: {s}\n", .{@errorName(err)});
            return 1;
        };
    }
    return 0;
}
