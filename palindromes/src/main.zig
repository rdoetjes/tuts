const std = @import("std");
const stdin = std.io.getStdIn();
const stdout = std.io.getStdOut().writer();
var gpa = std.heap.GeneralPurposeAllocator(.{}){};
const allocator = gpa.allocator();

// takes the first charachter of s and compares it to the last character of s
// when they are the same, then the 2nd charachter of s is compared to the 2nd to last character of s
// when they are the same, then the 3rd charachter of s is compared to the 3rd to last character of s etc etc
// This way half of the string is checked and when that half matches the other half, then the string is a palindrome
// when a character doesn't match, then the string is not a palindrome and we return false bailing out early.
fn is_palindrome(s: []const u8) bool {
    //see how much memory we need to allocated for the new replaced string (is either the same size or smaller)
    const new_size = std.mem.replacementSize(u8, s, " ", "");

    //allocate enough memory for the string where we removed the spaces
    const spaces_removed = allocator.alloc(u8, new_size) catch |err| {
        std.debug.print("Error: {}", .{err});
        return false;
    };
    defer allocator.free(spaces_removed);

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

test "palindrome tests" {
    try std.testing.expectEqual(true, is_palindrome("abba"));
    try std.testing.expectEqual(true, is_palindrome("Abba"));
    try std.testing.expectEqual(true, is_palindrome("abbA"));
    try std.testing.expectEqual(true, is_palindrome("aBbA"));
    try std.testing.expectEqual(true, is_palindrome("A man nam A"));
    try std.testing.expectEqual(false, is_palindrome("head"));
    try std.testing.expectEqual(true, is_palindrome("Was it a car or a cat I saw"));
}

// we read a list of words delimited by newlines from stdin
// each word is checked to see if it is a palindrome
// when it is, it's printed to the stdout
pub fn main() !void {
    const std_reader = stdin.reader();
    var br = std.io.bufferedReader(std_reader);
    var buffer: [4096]u8 = undefined;

    while (br.reader().readUntilDelimiterOrEof(&buffer, '\n') catch |err| {
        std.debug.print("Error: {s}", .{@errorName(err)});
        return;
    }) |l| {
        if (is_palindrome(l)) {
            stdout.print("{s}\n", .{l}) catch |err| {
                std.debug.print("Error: {s}", .{@errorName(err)});
                std.os.exit(1);
            };
        }
    }
}
