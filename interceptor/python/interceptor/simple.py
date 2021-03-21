import grpc


class SimpleClientInterceptor(grpc.UnaryUnaryClientInterceptor):

    def intercept_unary_unary(self, continuation, client_call_details, request):
        print("client_call_details: ", client_call_details)
        print("request: ", request)
        response = continuation(client_call_details, request)
        print("response: ", response)
        return response
