const std = @import("std");
const openai = @import("openai_v1.zig");

pub fn main() !void {
    // Initialize allocator
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();

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

    v1.set_model("gpt-3.5-turbo");

    while (true) {

        // ask openai a question
        const response = try v1.ask("Tell me a joke about germans In the style of Hans Landa from Ingolorius Bastards no introduction no preludes, the audience knows it is a joke!");

        //print the content of the JSON response body if http status is ok
        if (response.status == .ok) {
            std.debug.print("\x1b[31mZee joke about zee Germans:\x1b[0m\n \x1b[32m{s}\x1b[0m\n", .{v1.answer});
        } else {
            std.debug.print("Request failed with status: {any}\n", .{response});
        }
        std.time.sleep(5 * 1e9);
    }
}
