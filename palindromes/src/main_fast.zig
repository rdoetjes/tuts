const std = @import("std");
const stdin = std.io.getStdIn();
const stdout = std.io.getStdOut().writer();
const BUFFER_SIZE = 1024;

// takes the first charachter of s and compares it to the last character of s
// when they are the same, then the 2nd charachter of s is compared to the 2nd to last character of s
// when they are the same, then the 3rd charachter of s is compared to the 3rd to last character of s etc etc
// This way half of the string is checked and when that half matches the other half, then the string is a palindrome
// when a character doesn't match, then the string is not a palindrome and we return false bailing out early.
fn is_palindrome(s: []const u8, spaces_removed: []u8) bool {
    //see how much memory we need to allocated for the new replaced string (is either the same size or smaller)
    const new_size = std.mem.replacementSize(u8, s, " ", "");

    //remove the spaces from the string so that Was it a car or a cat I saw, is also seen as a palindrome
    _ = std.mem.replace(u8, s, " ", "", spaces_removed);

    //ieterate through the string and compare the first and last character and then the 2nd and last character and so on
    for (0..new_size / 2) |i| {
        const ts = std.ascii.toLower(spaces_removed[i]);
        const te = std.ascii.toLower(spaces_removed[new_size - i - 1]);

        // if chracters don't match, then the string is not a palindrome and we return false early
        if (ts != te) {
            return false;
        }
    }
    return true;
}

test "these are palindromess" {
    var spaces_removed: [BUFFER_SIZE]u8 = undefined;

    try std.testing.expectEqual(true, is_palindrome("abba", &spaces_removed));
    try std.testing.expectEqual(true, is_palindrome("Abba", &spaces_removed));
    try std.testing.expectEqual(true, is_palindrome("abbA", &spaces_removed));
    try std.testing.expectEqual(true, is_palindrome("aBbA", &spaces_removed));
    try std.testing.expectEqual(true, is_palindrome("A man nam A", &spaces_removed));
    try std.testing.expectEqual(true, is_palindrome("Was it a car or a cat I saw", &spaces_removed));
}

test "these are NOT palindromess" {
    var spaces_removed: [BUFFER_SIZE]u8 = undefined;

    try std.testing.expectEqual(false, is_palindrome("head", &spaces_removed));
    try std.testing.expectEqual(false, is_palindrome("Was it a car or a cat I saw?", &spaces_removed));
}

// we read a list of words delimited by newlines from stdin
// each word is checked to see if it is a palindrome
// when it is, it's printed to the stdout
pub fn main() !void {
    const std_reader = stdin.reader();
    var br = std.io.bufferedReader(std_reader);
    var buffer: [BUFFER_SIZE]u8 = undefined;
    var spaces_removed: [BUFFER_SIZE]u8 = undefined;

    while (br.reader().readUntilDelimiterOrEof(&buffer, '\n') catch |err| {
        std.debug.print("Error: {s}", .{@errorName(err)});
        return;
    }) |l| {
        if (is_palindrome(l, &spaces_removed)) {
            stdout.print("{s}\n", .{l}) catch |err| {
                std.debug.print("Error: {s}", .{@errorName(err)});
                std.os.exit(1);
            };
        }
    }
}
