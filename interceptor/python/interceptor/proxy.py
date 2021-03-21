import grpc


class ClientCallDetails(grpc.ClientCallDetails):

    def __init__(self,
                 method,
                 timeout,
                 metadata,
                 credentials,
                 wait_for_ready,
                 compression):
        self.method = method
        self.timeout = timeout
        self.metadata = metadata
        self.credentials = credentials
        self.wait_for_ready = wait_for_ready
        self.compression = compression

    def __str__(self):
        return f"ClientCallDetails(method={self.method}, timeout={self.timeout}, metadata={self.metadata}, credentials={self.credentials}, wait_for_ready={self.wait_for_ready}, compression={self.compression})"

    __repr__ = __str__


class ProxyClientInterceptor(grpc.UnaryUnaryClientInterceptor):

    def __init__(self, path_prefix):
        self.path_prefix = path_prefix

    def intercept_unary_unary(self, continuation, client_call_details, request):
        client_call_details = ClientCallDetails(
            method=self.path_prefix + client_call_details.method,
            timeout=client_call_details.timeout,
            metadata=client_call_details.metadata,
            credentials=client_call_details.credentials,
            wait_for_ready=client_call_details.wait_for_ready,
            compression=client_call_details.compression,
        )
        print("client_call_details: ", client_call_details)
        print("request: ", request)
        response = continuation(client_call_details, request)
        print("response: ", response)
        return response
