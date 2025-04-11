**Zig HTTP Client Example**

We are using the `std.http` module to make an HTTP request to the OpenAI API to generate a joke bout Germans.

We are using Zig 0.14.0 here, a lot as changed since the last version and probably will change again in the future. As it seems the interfaces are not stable yet.

In the main_single_request.zig is a very simple example of how to use the `std.http` module to make an HTTP request to the OpenAI API to generate a joke bout Germans. It doesn't abstract the OpenAI API complexity. This is a good starting point for learning how to use the `std.http` module.

In the main.zig is a more complex example of how to use the `std.http` module to make multiple HTTP requests to the OpenAI API to generate multiple jokes about Germans. It abstracts the OpenAI API complexity. This is a good starting point for learning how to use the `std.http` module.
