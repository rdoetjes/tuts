const std = @import("std");
const json = std.json;

pub const OpenAI_v1 = struct {
    model: []const u8,
    auth_header: []const u8,
    url: []const u8 = "https://api.openai.com/v1/chat/completions",
    client: std.http.Client,
    headers: std.ArrayList(std.http.Header),
    allocator: std.mem.Allocator,

    pub fn set_model(self: *OpenAI_v1, model: []const u8) void {
        if (self.model.len > 0) self.allocator.free(self.model);
        self.model = std.mem.Allocator.dupe(self.allocator, u8, model) catch |err| {
            std.debug.print("ABORT! Failed to duplicate model name: {}\n", .{err});
            std.os.linux.exit(1);
        };
    }

    pub fn init(allocator: std.mem.Allocator, api_key: []const u8) !*OpenAI_v1 {
        const instance = try allocator.create(OpenAI_v1);

        const auth_header = try std.fmt.allocPrint(allocator, "Bearer {s}", .{api_key});
        const model = try std.fmt.allocPrint(allocator, "gpt-4o", .{});
        var headers = std.ArrayList(std.http.Header).init(allocator);
        try headers.append(.{ .name = "Authorization", .value = auth_header });
        try headers.append(.{ .name = "Content-Type", .value = "application/json" });
        const client = std.http.Client{ .allocator = allocator };

        instance.* = .{
            .model = model,
            .allocator = allocator,
            .client = client,
            .auth_header = auth_header,
            .headers = headers,
        };
        return instance;
    }

    fn make_body(self: *OpenAI_v1, allocator: std.mem.Allocator, question: []const u8) ![]u8 {
        const Message = struct {
            role: []const u8,
            content: []const u8,
        };

        const Body = struct {
            model: []const u8,
            messages: [1]Message,
        };

        const body = Body{
            .model = self.model,
            .messages = [_]Message{
                .{ .role = "user", .content = question },
            },
        };

        return try json.stringifyAlloc(allocator, body, .{});
    }

    pub fn ask(
        self: *OpenAI_v1,
        question: []const u8,
        response_body: *std.ArrayList(u8),
    ) !std.http.Client.FetchResult {
        const payload = try make_body(self, self.client.allocator, question);
        defer self.client.allocator.free(payload);

        const response = try self.client.fetch(.{
            .method = .POST,
            .location = .{ .url = self.url },
            .extra_headers = self.headers.items,
            .payload = payload,
            .response_storage = .{ .dynamic = response_body },
        });

        return response;
    }

    pub fn deinit(self: *OpenAI_v1) void {
        self.allocator.free(self.model);
        self.allocator.free(self.auth_header);
        self.headers.deinit();
        self.client.deinit();
    }
};
