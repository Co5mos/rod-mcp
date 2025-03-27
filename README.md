# Rod MCP Server

<div align="center">

<img src="assets/logo2.png" alt="logo" width="400" height="400">


<strong>Wow! It's awesome, now you can use the MCP server of Rod!🚀</strong>

<br>

<strong>Rod-MCP provides browser automation capabilities for your applications by using [Rod](https://github.com/go-rod/rod). The server provides many useful mcp tools enable LLMs to interact with the web pages, like click, take screenshot, save page as pdf etc.</strong>

</div>


<h3>Engilsh | <a href='./README_CN.md'> 中文 </a></h3>


## Features

- 🚀 Browser automation powered by Rod
- 🎯 Rich web interaction capabilities
  - Element clicking
  - Screenshot capture
  - PDF generation
  - And more...
- 🎨 Headless/GUI mode support
- ⚡ High performance and stability
- 🔧 Easy to configure and extend
- 🤖 Designed for LLMs interaction

## Installation

### Prerequisites

- Go 1.23 or higher
- Chrome/Chromium browser

### Steps

1. Clone the repository:
```bash
git clone https://github.com/go-rod/rod-mcp.git
cd rod-mcp
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the project:
```bash
go build
```

## Usage

### Basic Usage

1. Configure MCP:
```json
{
    "mcpServers": {
        "rod-mcp": {
            "command": "rod-mcp",
            "args": [
                "-c", "rod-mcp.yaml"
            ]
        }
    }
}
```

### Configuration

The configuration file supports the following options:
- serverName: Server name, default is "Rod Server"
- browserBinPath: Browser executable file path, use system default browser if empty
- headless: Whether to run the browser in headless mode, default is false
- browserTempDir: Browser temporary file directory, default is "./rod/browser"
- noSandbox: Whether to disable sandbox mode, default is false
- proxy: Proxy server settings, supports socks5 proxy

## Project Structure

```
rod-mcp/
├── assets/          # Static resources
├── banner/          # Banner resources
├── cmd.go           # Command line processing
├── main.go          # Program entry
├── resources/       # Resource files
├── server.go        # Server implementation
├── tools/           # Tool implementation
├── types/           # Type definitions
└── utils/           # Utility functions
```

## Contribution Guidelines

Welcome to submit Pull Request or create Issue!

## Tool List

Rod-MCP provides the following tools:

| Tool Name | Description | Parameters |
|-----------|-------------|------------|
| rod_navigate | Navigate to a URL | url: URL to navigate to |
| rod_go_back | Go back to the previous page | None |
| rod_go_forward | Go forward to the next page | None |
| rod_reload | Reload the current page | None |
| rod_press_key | Press a key on the keyboard | key: Name of the key to press (like "a" or "ArrowLeft") |
| rod_click | Click an element on the page | selector: CSS selector |
| rod_fill | Fill out an input field | selector: CSS selector, value: Text to fill |
| rod_selector | Select an option in a dropdown | selector: CSS selector, value: Value to select |
| rod_hover | Hover over an element | selector: CSS selector |
| rod_drag | Perform drag and drop operation | source_selector: Source element CSS selector, target_selector: Target element CSS selector |
| rod_screenshot | Take a screenshot | name: Screenshot name, selector: (optional) Element CSS selector, width: (optional) Width, height: (optional) Height |
| rod_pdf | Save page as PDF | file_path: File save path, file_name: File name |
| rod_wait | Wait for a specified time | time: Wait time in seconds (max 10s) |
| rod_snapshot | Capture accessibility snapshot | None |
| rod_close | Close the browser | None |
| rod_evaluate | Execute JavaScript code | script: JavaScript code to execute |

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file
