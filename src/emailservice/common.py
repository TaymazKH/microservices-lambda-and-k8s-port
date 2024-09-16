import grpc
from grpc_status._common import code_to_grpc_status_code


class GrpcError(grpc.RpcError):
    def __init__(self, code: grpc.StatusCode | int, message: str = ""):
        if code is None or code == grpc.StatusCode.OK or code == grpc.StatusCode.OK.value[0]:
            raise ValueError("Non-OK status code expected for errors")
        if type(code) == int:
            self._code = code_to_grpc_status_code(code)
        else:
            self._code = code
        self._message = message

    @property
    def code(self):
        return self._code

    @property
    def int_code(self):
        return self._code.value[0]

    @property
    def message(self):
        return self._message

    def __str__(self):
        return f"(code: {self.int_code}, message: {self.message})"
