import grpc
from pyprotos import common_pb2, helloworld_pb2_grpc


if __name__ == "__main__":
    with grpc.insecure_channel("localhost:50051") as channel:
        stub = helloworld_pb2_grpc.GreeterStub(channel)
        request = common_pb2.HelloRequest(name="gxg")
        response = stub.SayHello(request)
        print(response.message)
