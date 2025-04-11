const std = @import("std");
const json = std.json;

pub const OpenAI_v1 = struct {
    model: []const u8,
    auth_header: []const u8,
    url: []const u8 = "https://api.openai.com/v1/chat/completions",
    client: std.http.Client,
    headers: std.ArrayList(std.http.Header),
    allocator: std.mem.Allocator,
    answer: []const u8,

    pub fn init(allocator: std.mem.Allocator, api_key: []const u8) !*OpenAI_v1 {
        const instance = try allocator.create(OpenAI_v1);
        errdefer allocator.destroy(instance);

        const auth_header = try std.fmt.allocPrint(allocator, "Bearer {s}", .{api_key});
        const model = try std.mem.Allocator.dupe(allocator, u8, "gpt-4o");

        var headers = std.ArrayList(std.http.Header).init(allocator);
        errdefer headers.deinit();

        try headers.append(.{ .name = try allocator.dupe(u8, "Authorization"), .value = auth_header });
        try headers.append(.{ .name = try allocator.dupe(u8, "Content-Type"), .value = try allocator.dupe(u8, "application/json") });

        const client = std.http.Client{ .allocator = allocator };
        const answer = try std.mem.Allocator.dupe(allocator, u8, "");

        instance.* = .{
            .model = model,
            .allocator = allocator,
            .client = client,
            .auth_header = auth_header,
            .headers = headers,
            .answer = answer,
        };

        return instance;
    }

    pub fn set_model(self: *OpenAI_v1, model: []const u8) void {
        if (self.model.len > 0) self.allocator.free(self.model);
        self.model = std.mem.Allocator.dupe(self.allocator, u8, model) catch |err| {
            std.debug.print("ABORT! Failed to duplicate model name: {}\n", .{err});
            std.os.linux.exit(1);
        };
    }

    pub fn ask(self: *OpenAI_v1, question: []const u8) !std.http.Client.FetchResult {
        const payload = try make_body(self, self.allocator, question);
        defer self.allocator.free(payload);

        var response_body = std.ArrayList(u8).init(self.allocator);
        defer response_body.deinit();

        const response = self.client.fetch(.{
            .method = .POST,
            .location = .{ .url = self.url },
            .extra_headers = self.headers.items,
            .payload = payload,
            .response_storage = .{ .dynamic = &response_body },
        }) catch {
            return std.http.Client.FetchResult{ .status = .gone };
        };

        if (response.status == .ok) {
            parse_response(self, response_body.items) catch |err| {
                self.answer = try std.fmt.allocPrint(self.allocator, "{s}", .{"ERROR: DURING PARSING JSON"});
                std.debug.print("ERROR: Failed to parse response: {}\n", .{err});
            };
        }

        return response;
    }

    fn parse_response(self: *OpenAI_v1, response_body: []const u8) !void {
        const parsed = try std.json.parseFromSlice(std.json.Value, self.allocator, response_body, .{});
        defer parsed.deinit();

        self.allocator.free(self.answer);
        const answer = parsed.value.object.get("choices").?.array.items[0].object.get("message").?.object.get("content").?.string;
        self.answer = try std.fmt.allocPrint(self.allocator, "{s}", .{answer});
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

    pub fn deinit(self: *OpenAI_v1) void {
        if (self.model.len > 0) self.allocator.free(self.model);
        if (self.answer.len > 0) self.allocator.free(self.answer);

        for (self.headers.items) |header| {
            self.allocator.free(header.name);
            self.allocator.free(header.value);
        }
        self.allocator.free(self.auth_header);

        self.headers.deinit();
        self.client.deinit();
    }
};
