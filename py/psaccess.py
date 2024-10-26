import sys
import time
from subprocess import PIPE, Popen
from threading  import Thread
from queue import Queue
ON_POSIX = 'posix' in sys.builtin_module_names

def enqueue_output(out, queue):
    for line in iter(out.readline, b''):
        queue.put(line)
    out.close()

class Psaccess: 

    def __init__(self, dba, fsa):
        self.dba = dba
        self.fsa = fsa

        self.process = None
        self.queue = None
        self.thread = None
        self.lines = ""

    def start_server(self, osa):
        if self.is_running():
            return

        self.lines = ""

        acserverexe = self.fsa.get_serverexe(osa)
        cwd = self.fsa.get_serverpath()
        self.process = Popen(acserverexe, stdout=PIPE, bufsize=1, close_fds=ON_POSIX, cwd=cwd)
        self.queue = Queue()
        self.thread = Thread(target=enqueue_output, args=(self.process.stdout, self.queue))
        self.thread.daemon = True
        self.thread.start()

    def is_running(self):
        if self.process is None: 
            return False
        return self.process.poll() is None

    def stop_server(self):
        self.process.kill()
        self.process = None

    def process_content(self):

        now = time.time()            # get the time
        while 1:
            try:
                line = self.queue.get_nowait()
                self.lines = self.lines + "<br />" + line.decode("utf-8")
            except:
                pass
            elapsed = time.time() - now

            if elapsed > 0.05:
                return self.lines

