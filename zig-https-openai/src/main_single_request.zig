const std = @import("std");
const json = std.json;

pub fn main() !void {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();
    defer _ = gpa.deinit();

    var client = std.http.Client{
        .allocator = allocator,
    };
    defer client.deinit();

    const body = try json.stringifyAlloc(allocator, .{
        .model = "gpt-4o",
        .messages = .{
            .{
                .role = "user",
                .content = "Why is the earth round and not flat?",
            },
        },
    }, .{});
    defer allocator.free(body);

    const openai_api_key = std.process.getEnvVarOwned(allocator, "OPENAI_API_KEY") catch |err| {
        std.log.err("Error getting OPENAI_API_KEY: {}", .{err});
        std.os.linux.exit(1);
    };
    defer allocator.free(openai_api_key);

    const auth_header = try std.fmt.allocPrint(allocator, "Bearer {s}", .{openai_api_key});
    defer allocator.free(auth_header);

    const headers = &[_]std.http.Header{
        .{ .name = "Authorization", .value = auth_header },
        .{ .name = "Content-Type", .value = "application/json" },
    };

    var resp_body = std.ArrayList(u8).init(allocator);
    defer resp_body.deinit();

    const response = try client.fetch(.{
        .method = .POST,
        .location = .{ .url = "https://api.openai.com/v1/chat/completions" },
        .extra_headers = headers,
        .payload = body,
        .response_storage = .{ .dynamic = &resp_body },
    });

    if (response.status == .ok) {
        const parsed = try json.parseFromSlice(json.Value, allocator, resp_body.items, .{});
        defer parsed.deinit();

        const message = parsed.value.object.get("choices").?.array.items[0].object.get("message").?.object.get("content").?.string;
        std.debug.print("Message: {s}", .{message});
    } else {
        std.debug.print("Error: {any}", .{response.status});
    }
}
