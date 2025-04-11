const std = @import("std");
const Url = std.Uri;
const mem = std.mem;
const json = std.json;

pub fn main() !void {
    // Initialize allocator
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();

    // Make request
    var client = std.http.Client{
        .allocator = allocator,
    };
    defer client.deinit();

    // Request body
    const body = try std.json.stringifyAlloc(allocator, .{
        .model = "gpt-4o",
        .messages = .{
            .{
                .role = "user",
                .content = "Tell me a joke about germans In the style of Hans Landa from Ingolorius Bastards no introduction no preludes, the audience knows it is a joke!",
            },
        },
    }, .{});
    defer allocator.free(body);

    const openai_api_key = std.process.getEnvVarOwned(allocator, "OPENAI_API_KEY") catch {
        std.debug.print("OPENAI_API_KEY environment variable not set\n", .{});
        std.os.linux.exit(1);
    };
    defer allocator.free(openai_api_key);

    const auth_header = try std.fmt.allocPrint(allocator, "Bearer {s}", .{openai_api_key});
    defer allocator.free(auth_header);

    const headers = &[_]std.http.Header{
        .{
            .name = "Authorization",
            .value = auth_header,
        },
        .{
            .name = "Content-Type",
            .value = "application/json",
        },
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
        //std.debug.print("The whole Response Body: {s}\n", .{resp_body.items});
        const parsed = try std.json.parseFromSlice(std.json.Value, allocator, resp_body.items, .{});
        defer parsed.deinit();

        const message = parsed.value.object.get("choices").?.array.items[0].object.get("message").?.object.get("content").?.string;
        std.debug.print("\x1b[31mZee joke about zee Germans:\x1b[0m\n \x1b[32m{s}\x1b[0m\n", .{message});
    } else {
        std.debug.print("Request failed with status: {any}\n", .{response});
    }
}
