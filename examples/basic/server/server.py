import uuid
from concurrent import futures
import logging
import grpc

from example.v1 import example_pb2
from example.v1 import example_pb2_grpc

class UserServiceServicer(example_pb2_grpc.UserServiceServicer):
    def __init__(self):
        # Simple in-memory storage using dict
        self.users = {}

    def CreateUser(self, request, context):
        # Generate a unique ID for the new user
        user_id = str(uuid.uuid4())

        # Create new user object
        user = example_pb2.User(
            id=user_id,
            name=request.name,
            email=request.email,
            roles=request.roles,
            status=example_pb2.UserStatus.USER_STATUS_ACTIVE
        )

        # Store in our "database"
        logging.info("Created user %s", user.id)
        self.users[user_id] = user

        return example_pb2.CreateUserResponse(user=user)

    def UpdateUser(self, request, context):
        # Check if user exists
        if request.id not in self.users:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f'User {request.id} not found')
            return example_pb2.UpdateUserResponse()

        # Get existing user
        logging.info("Updating user %s", request.id)
        user = self.users[request.id]

        # Update fields if provided
        if request.HasField('name'):
            user.name = request.name
        if request.HasField('email'):
            user.email = request.email
        if request.roles:
            user.roles[:] = request.roles

        # Store updated user
        logging.info("Updated user %s", user.id)
        self.users[request.id] = user

        return example_pb2.UpdateUserResponse(user=user)

    def GetUser(self, request, context):
        # Check if user exists
        if request.id not in self.users:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f'User {request.id} not found')
            return example_pb2.GetUserResponse()

        logging.info("Getting user %s", request.id)
        return example_pb2.GetUserResponse(user=self.users[request.id])

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    example_pb2_grpc.add_UserServiceServicer_to_server(
        UserServiceServicer(), server
    )
    listen_addr = '[::]:50051'
    server.add_insecure_port(listen_addr)
    server.start()
    logging.info("Server started on %s", listen_addr)
    server.wait_for_termination()

if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO)
    serve()