const std = @import("std");
const mem = std.mem;
const json = std.json;

pub fn main() !void {
    // Initialize allocator
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();

    // Make client used to make requests using fetch
    var client = std.http.Client{
        .allocator = allocator,
    };
    defer client.deinit();

    // Request body for the OpenAI API to tell a nice joke about Germans
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

    // fetch my OpenAI API key from environment variable OPENAI_API_KEY (std.os.getenv no longer works it seems)
    const openai_api_key = std.process.getEnvVarOwned(allocator, "OPENAI_API_KEY") catch {
        std.debug.print("OPENAI_API_KEY environment variable not set\n", .{});
        std.os.linux.exit(1);
    };
    defer allocator.free(openai_api_key);

    // create the authorization header with the API key as a Bearer token
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

    // Initialize the response body buffer (ths is of course a u8 arraylist as per standard in zig)
    var resp_body = std.ArrayList(u8).init(allocator);
    defer resp_body.deinit();

    // perform the fetch request using the header, body and response storage
    const response = try client.fetch(.{
        .method = .POST,
        .location = .{ .url = "https://api.openai.com/v1/chat/completions" },
        .extra_headers = headers,
        .payload = body,
        .response_storage = .{ .dynamic = &resp_body },
    });

    //print the content of the JSON response body if http status is ok
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
