"""
This is a simple API wrapper around the 9p command line utility that allows
for programmatic control of acme windows from python.
It is assumed that you have pla9port installed and that you have the plan9
bin directory on your path.

See the following for more details (middle mouse in acme to open):
    acme(1)
    acme(4)
    acmeevent(1)

TODO:
  - Write Event class for parsing event strings
  - Allow event handler to register events to listen to
  - Ensure that we correctly clean up the EventHandler class when we
    stop listening so that we don't kill event handling for the window.
"""
import os
from subprocess import run, Popen, PIPE


def _ensure_winid(winid):
    if winid is None:
        # Try to pull a default id from the environment under the assumption
        # that we have been called from within acme itself
        winid = os.environ.get("winid")
        if winid is None:
            raise AcmeError("Unable to determine active window")

    return winid


class AcmeError(Exception):
    pass


class AcmeWindow:
    def __init__(self, winid):
        self.id = _ensure_winid(winid)

    def _read(self, fname):
        """
        Read the contents of an acme control file and return the content
        as a utf-8 string.

        NOTE: Depending on the file, there may be restrictions on the text that
              will be accepted. See acme(4) for details of the available files,
              their effect on the state of acme when written to and any
              restrictions on writing to them.
        """
        proc = run(f"9p read acme/{self.id}/{fname}", shell=True, stdout=PIPE)
        return proc.stdout.decode()

    def _write(self, fname, text):
        """
        Write a string to an acme control file.

        NOTE: Depending on the file, there may be restrictions on the text that
              will be accepted. See acme(4) for details of the available files,
              their effect on the state of acme when written to and any
              restrictions on writing to them.
        """
        run(f"echo {text} | 9p write acme/{self.id}/{fname}", shell=True)

    def mark_clean(self):
        """
        Mark a window as clean.

        NOTE: This removes the indicator in the UI but does not write any
              outstanding changes to the buffer to disk.
        """
        self._write("ctl", "clean", self.id)

    def mark_dirty(self):
        """
        Mark a window as dirty.

        NOTE: This sets the indicator in the UI but does not otherwise affect the
              content of the buffer or the file on disk.
        """
        self._write("ctl", "dirty", self.id)

    def clear_tags(self):
        """
        Remove all custom tags on a window.

        Custom tags are defined as those following the pipe character (|). It is
        not possible for the user to delete the default tags.
        """
        self._write("ctl", "cleartag", self.id)

    def reload_window(self):
        self._write("ctl", "get", self.id)

    def save(self):
        self._write("ctl", "put", self.id)

    def fname_and_tags(self):
        """Fetch a window's current file/directory name and any custom tags"""
        tag = self._read("tag", self.id)
        fname, _, tags = tag.partition("|")
        return fname, tags.strip().split()

    def get_window_name(self):
        """Find the current window's name, assuming no spaces"""
        return self._read("tag", self.id).split()[0]


class Event:
    """
    A minimal wrapper around acmeevent strings

    See acmeevent(1) for details.
    """

    def __init__(self, estr):
        self._event = estr

    def __repr__(self):
        return f"EVENT: {self._event}"


class EventHandler:
    """
    A simple control class for managing acmeevent format events for an acme
    window.
    """

    def __init__(self, winid=None):
        self._id = _ensure_winid(winid)
        self._gen_events = None

    def _init_event_gen(self):
        """
        Create a generator that listens to the window events file and yields
        them as Event objects.
        """
        proc = Popen(
            f"9p read acme/{self._id}/event | acmeevent",
            shell=True,
            stdout=PIPE,
            stderr=PIPE,
            text=True,
        )

        def _gen():
            while True:
                continue_listening = yield
                if not continue_listening:
                    # We're done so kill the process
                    proc.kill()
                    return

                line = proc.stdout.readline().strip()
                yield Event(line)

        self._gen_events = _gen()
        self._stop = lambda: self._gen_events.send(False)

    def get_event(self):
        if self._gen_events is None:
            self._init_event_gen()
            next(self._gen_events)

        evt = self._gen_events.send(True)
        next(self._gen_events)
        return evt

    def stop(self):
        if self._gen_events is None:
            return
        try:
            self._stop()
        except StopIteration:
            pass

        self._gen_events = None
