const std = @import("std");
const stdin = std.io.getStdIn();
const stdout = std.io.getStdOut().writer();

fn is_palindrome(s: []const u8) bool {
    for (0..s.len / 2) |i| {
        const ts = std.ascii.toLower(s[i]);
        const te = std.ascii.toLower(s[s.len - i - 1]);

        if (ts != te) {
            return false;
        }
    }
    return true;
}

test "palindrome test" {
    try std.testing.expectEqual(true, is_palindrome("abba"));
    try std.testing.expectEqual(true, is_palindrome("Abba"));
    try std.testing.expectEqual(true, is_palindrome("abbA"));
    try std.testing.expectEqual(true, is_palindrome("aBbA"));
    try std.testing.expectEqual(false, is_palindrome("head"));
}

pub fn main() !void {
    const std_reader = stdin.reader();
    var br = std.io.bufferedReader(std_reader);
    var buffer: [8192]u8 = undefined;

    while (try br.reader().readUntilDelimiterOrEof(&buffer, '\n')) |l| {
        if (is_palindrome(l)) {
            try stdout.print("{s}\n", .{l});
        }
    }
}
