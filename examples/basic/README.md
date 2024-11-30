# protoc-gen-mcp/examples/basic

Basic example of using protoc-gen-mcp to generate a python grpc mcp server.

## Usage

1. Sync & Activate venv
2. gRPC server
    1. `python server/server.py`
    2. `python test_client.py`
3. MCP server [wip]
    1. `python mcp_query "Create Admin user John Doe with email john@example.com"`
