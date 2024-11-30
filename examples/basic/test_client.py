import logging
import grpc
from example.v1 import example_pb2
from example.v1 import example_pb2_grpc

def run():
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = example_pb2_grpc.UserServiceStub(channel)

        # Test CreateUser
        print("\nTesting CreateUser...")
        create_response = stub.CreateUser(example_pb2.CreateUserRequest(
            name="John Doe",
            email="john@example.com",
            roles=["USER", "ADMIN"]
        ))
        print(f"Created user: {create_response.user}")

        # Store the user ID for subsequent operations
        user_id = create_response.user.id

        # Test GetUser
        print("\nTesting GetUser...")
        get_response = stub.GetUser(example_pb2.GetUserRequest(
            id=user_id
        ))
        print(f"Retrieved user: {get_response.user}")

        # Test UpdateUser
        print("\nTesting UpdateUser...")
        update_response = stub.UpdateUser(example_pb2.UpdateUserRequest(
            id=user_id,
            name="John Updated Doe",
            roles=["USER"]
        ))
        print(f"Updated user: {update_response.user}")

        # Verify update with GetUser
        print("\nVerifying update...")
        get_response = stub.GetUser(example_pb2.GetUserRequest(
            id=user_id
        ))
        print(f"Retrieved updated user: {get_response.user}")

        # Test error handling - try to get non-existent user
        print("\nTesting error handling...")
        try:
            stub.GetUser(example_pb2.GetUserRequest(
                id="non-existent-id"
            ))
        except grpc.RpcError as e:
            print(f"Expected error received: {e.code()}: {e.details()}")

if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO)
    run()
