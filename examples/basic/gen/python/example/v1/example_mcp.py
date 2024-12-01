import json
import logging
from typing import Any, Sequence

from mcp.server import Server
from mcp.types import Tool, TextContent, ImageContent, EmbeddedResource

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger('example.v1-server')

app = Server('example.v1-server')

@app.list_tools()
async def list_tools() -> list[Tool]:
    return [
        Tool(
            name='CreateUser',
            description='CreateUser method from UserService service',
            inputSchema={
                'type': 'object',
                'properties': {
                    'name': {
                        'type': 'string',
                        'description': 'Field name'
                    },
                    'email': {
                        'type': 'string',
                        'description': 'Field email'
                    },
                    'roles': {
                        'type': 'string',
                        'description': 'Field roles'
                    },
                },
                'required': ['name','email','roles',]
            }
        ),
        Tool(
            name='UpdateUser',
            description='UpdateUser method from UserService service',
            inputSchema={
                'type': 'object',
                'properties': {
                    'id': {
                        'type': 'string',
                        'description': 'Field id'
                    },
                    'name': {
                        'type': 'string',
                        'description': 'Field name'
                    },
                    'email': {
                        'type': 'string',
                        'description': 'Field email'
                    },
                    'roles': {
                        'type': 'string',
                        'description': 'Field roles'
                    },
                },
                'required': ['id','roles',]
            }
        ),
        Tool(
            name='GetUser',
            description='GetUser method from UserService service',
            inputSchema={
                'type': 'object',
                'properties': {
                    'id': {
                        'type': 'string',
                        'description': 'Field id'
                    },
                },
                'required': ['id',]
            }
        ),
    ]

@app.call_tool()
async def call_tool(name: str, arguments: Any) -> Sequence[TextContent | ImageContent | EmbeddedResource]:
    if name not in [
        'CreateUser',
        'UpdateUser',
        'GetUser',
    ]:
        raise ValueError(f'Unknown tool: {name}')

    # TODO: Implement method handlers
    return [
        TextContent(
            type='text',
            text=json.dumps({'result': 'Not implemented'}, indent=2)
        )
    ]

async def main():
    from mcp.server.stdio import stdio_server
    async with stdio_server() as (read_stream, write_stream):
        await app.run(
            read_stream,
            write_stream,
            app.create_initialization_options()
        )

if __name__ == '__main__':
    import asyncio
    asyncio.run(main())
