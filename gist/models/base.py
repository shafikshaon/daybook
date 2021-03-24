from gist.models.activity import Activity
from gist.models.key import Key
from gist.models.timelog import TimeLog


class Base(Key, TimeLog, Activity):
    class Meta:
        abstract = True
