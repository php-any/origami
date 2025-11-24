# API Gateway Example

This example demonstrates how to build a powerful API Gateway using Origami.
The gateway routes requests to different applications based on the `Host` header (domain name), similar to Nginx, but with the ability to dynamically load and execute Origami applications directly.

## Features

- **Dynamic Routing**: Dispatches requests based on domain names.
- **Application Isolation**: Each application runs in its own VM context (via `Net\Http\app`).
- **Unified Entry Point**: Handles cross-cutting concerns (logging, error handling) centrally.
- **Full MVC Support**: Routes to fully structured MVC applications using `@Application` and `@Controller`.

## Structure

- `main.zy`: The gateway server and routing logic.
- `apps/demo/`: A complete MVC application structure.
  - `src/main.zy`: App entry point with `@Application`.
  - `src/controllers/`: Contains `HelloController`.
- `../spring/`: References the existing Spring-like example app.

## Usage

1. Start the gateway:
   ```bash
   origami examples/gateway/main.zy
   ```

2. Test with different domains (using curl to simulate Host header):

   **Route to Demo App (Hello Controller):**
   ```bash
   curl -H "Host: demo.local" http://localhost:8080/hello
   ```
   *Response:*
   ```json
   {"service": "Demo Service", "message": "Hello from the Controller!", ...}
   ```

   **Route to Demo App (Echo):**
   ```bash
   curl -H "Host: demo.local" http://localhost:8080/echo/origami
   ```

   **Route to Spring App:**
   ```bash
   curl -H "Host: spring.local" http://localhost:8080/
   ```

   **Default Route (Demo App):**
   ```bash
   curl http://localhost:8080/hello
   ```

## Configuration

The routing table is defined in `main.zy`. It maps domain names to the application's entry file and function name.

```zy
$routes = {
    "demo.local": {
        "path": "apps/demo/src/main.zy",
        "func": "DemoApp\\main"
    },
    // ...
};
```
