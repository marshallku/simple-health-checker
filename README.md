
# Statusy

![statusy](https://github.com/user-attachments/assets/5d2ed59e-5a5e-4584-a258-9dbdd972ddb6)

A lightweight, configurable health checker for web services written in Go. This tool performs HTTP requests to specified URLs and checks for expected status codes, response times, and content inclusion. It provides real-time status monitoring through a web interface and can send notifications to Discord when checks fail.

## Features

- Web interface for real-time status monitoring
- Real-time updates via WebSocket
- History tracking of the last 10 events
- Configurable health checks via YAML file
- Checks for HTTP status codes
- Checks for response time
- Checks for specific text inclusion in responses
- Custom HTTP methods and headers for requests
- Discord notifications for failed checks

## Installation

1. Ensure you have Go installed on your system (version 1.23 or later recommended).
2. Clone this repository:

   ```bash
   git clone https://github.com/marshallku/statusy.git
   cd statusy
   ```

3. Install dependencies:

   ```bash
   go mod tidy
   ```

## Configuration

Create a `config.yaml` file in the project root directory. Here's an example configuration:

```yaml
webhook_url: https://discord.com/api/webhooks/your_webhook_url_here
timeout: 5000  # Global timeout in milliseconds
checkInterval: 60  # Check interval in seconds

pages:
  - url: https://example.com
    status: 200
    text_to_include: Welcome
    speed: 2000
  - url: https://api.example.com
    status: 200
  - url: https://example.com/about
    text_to_include: About Us
```

### Configuration Options

- `webhook_url`: Discord webhook URL for notifications
- `timeout`: Global timeout for all requests in milliseconds
- `checkInterval`: Interval between health checks in seconds
- `pages`: List of pages to check
  - `url`: URL to check (required)
  - `status`: Expected HTTP status code (default: 200)
  - `text_to_include`: String to look for in the response body (optional)
  - `speed`: Maximum acceptable response time in milliseconds (optional)
  - `request`: Custom request options (optional)
    - `method`: HTTP method (GET, POST, etc.)
    - `headers`: Custom HTTP headers
    - `body`: Request body for POST/PUT requests

## Usage

### Running Locally

Run the health checker with:

```bash
go run .
```

Or, to specify a different config file:

```bash
go run . --config path/to/your/config.yaml
```

If you want to run statusy just once without the web interface:

```bash
go run . --mode cli
```

### Using Docker

Build and run with Docker:

```bash
docker build -t statusy .
docker run -p 8080:8080 statusy
```

### Accessing the Web Interface

Once running, access the web interface at:

- Status Dashboard: <http://localhost:8080/>
- History Page: <http://localhost:8080/history>

The web interface features:

- Real-time status updates via WebSocket
- Current status of all monitored pages
- History of the last 10 events
- Automatic updates without page refresh

## Building

To build an executable:

```bash
go build -o statusy
```

Then run the executable:

```bash
./statusy
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
