const std = @import("std");
const stdin = std.io.getStdIn();
const stdout = std.io.getStdOut().writer();

// takes the first charachter of s and compares it to the last character of s
// when they are the same, then the 2nd charachter of s is compared to the 2nd to last character of s
// when they are the same, then the 3rd charachter of s is compared to the 3rd to last character of s etc etc
// This way half of the string is checked and when that half matches the other half, then the string is a palindrome
// when a character doesn't match, then the string is not a palindrome and we return false bailing out early.
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

// we read a list of words delimited by newlines from stdin
// each word is checked to see if it is a palindrome
// when it is, it's printed to the stdout
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
