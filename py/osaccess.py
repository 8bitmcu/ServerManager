import sys

class Osaccess: 
    def __init__(self, ):
        pass

    def is_unix(self):
        return sys.platform != "win32"
