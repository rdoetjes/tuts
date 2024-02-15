const std = @import("std");
const in = std.io.getStdIn();

fn is_palindrome(s: []const u8) bool {
    for (0..s.len / 2) |i| {
        if (s[i] != s[s.len - i - 1]) {
            return false;
        }
    }
    return true;
}

pub fn main() !void {
    const std_reader = std.io.getStdIn().reader();
    var br = std.io.bufferedReader(std_reader);
    var buffer: [4096]u8 = undefined;
    while (try br.reader().readUntilDelimiterOrEof(&buffer, '\n')) |l| {
        if (is_palindrome(l)) {
            try std.io.getStdOut().writer().print("{s}\n", .{l});
        }
    }
}

test "palindrome test" {
    try std.testing.expectEqual(true, is_palindrome("abba"));
    try std.testing.expectEqual(false, is_palindrome("head"));
}
