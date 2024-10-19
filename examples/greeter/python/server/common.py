from grpc import RpcError, StatusCode
from grpc_status._common import code_to_grpc_status_code


class GrpcError(RpcError):
    def __init__(self, code: StatusCode | int, message: str = ""):
        if code is None or code == StatusCode.OK or code == StatusCode.OK.value[0]:
            raise ValueError("Non-OK status code expected for errors")
        if isinstance(code, int):
            self._code = code_to_grpc_status_code(code)
        elif isinstance(code, StatusCode):
            self._code = code
        else:
            raise TypeError(f"Status code must be grpc.StatusCode or int, not {type(code)}")
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
