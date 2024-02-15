const std = @import("std");
const in = std.io.getStdIn();
const stdout = std.io.getStdOut().writer();

fn char_to_lower(c: u8) u8 {
    if (c >= 65 and c <= 90) {
        return c | 32;
    }
    return c;
}

test "to_lower test" {
    try std.testing.expectEqual('u', char_to_lower('U'));
    try std.testing.expectEqual('a', char_to_lower('a'));
}

fn is_palindrome(s: []const u8) bool {
    for (0..s.len / 2) |i| {
        const ts = char_to_lower(s[i]);
        const te = char_to_lower(s[s.len - i - 1]);

        if (ts != te) {
            return false;
        }
    }
    return true;
}

test "palindrome test" {
    try std.testing.expectEqual(true, is_palindrome("abba"));
    try std.testing.expectEqual(false, is_palindrome("head"));
}

pub fn main() !void {
    const std_reader = std.io.getStdIn().reader();
    var br = std.io.bufferedReader(std_reader);
    var buffer: [8192]u8 = undefined;
    while (try br.reader().readUntilDelimiterOrEof(&buffer, '\n')) |l| {
        if (is_palindrome(l)) {
            try stdout.print("{s}\n", .{l});
        }
    }
}
