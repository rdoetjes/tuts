const std = @import("std");
const stdout = std.io.getStdOut().writer();
const openai = @import("openai_v1.zig");

pub fn main() !void {
    // Initialize allocator
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();
    var retry: u8 = 0;

    // fetch my OpenAI API key from environment variable OPENAI_API_KEY (std.os.getenv no longer works it seems)
    const openai_api_key = std.process.getEnvVarOwned(allocator, "OPENAI_API_KEY") catch {
        std.debug.print("OPENAI_API_KEY environment variable not set\n", .{});
        std.os.linux.exit(1);
    };
    defer allocator.free(openai_api_key);

    // Initialize the OpenAI client
    const v1 = openai.OpenAI_v1.init(allocator, openai_api_key) catch |err| {
        std.debug.print("Failed to initialize OpenAI client: {any}\n", .{err});
        std.os.linux.exit(1);
    };
    defer v1.deinit();

    v1.set_model("gpt-3.5-turbo"); // I am cheap, so I use the cheaper model

    while (true) {
        try stdout.print("\n\x1b[31mAsk me anything: \x1b[0m", .{});
        const question = try std.io.getStdIn().reader().readUntilDelimiterAlloc(allocator, '\n', 4096);
        defer allocator.free(question);

        // ask openai a questione
        var response = try v1.ask(question);

        if (response.status == .gone) {
            response = try v1.ask(question);
        }

        //print the content of the JSON response body if http status is ok
        if (response.status == .ok and std.mem.containsAtLeast(u8, v1.answer, 0, "ERROR:")) {
            try stdout.print("\n\x1b[31mMy answer:\x1b[0m \x1b[32m{s}\x1b[0m\n", .{v1.answer});
            retry = 0;
        } else {
            std.debug.print("Request failed with status: {any}\n", .{response});
            std.os.linux.exit(1);
        }
        std.time.sleep(2 * 1e9); // avoid spamming and consuming your prepaid api fee ;)
    }
}
