from django.db import models
from django.utils import timezone


class LoginAttemptRecord(models.Model):
    username = models.CharField(max_length=255, null=True, blank=True, unique=True)
    count = models.PositiveIntegerField(null=True, blank=True, default=0)
    timestamp = models.DateTimeField(auto_now_add=True)

    class Meta:
        app_label = 'auth_record'


class LoginAttemptLogger(object):
    @classmethod
    def reset(cls, username):
        defaults = {
            'count': 0,
            'timestamp': timezone.now()
        }
        LoginAttemptRecord.objects.update_or_create(username=username, defaults=defaults)

    @classmethod
    def increment(cls, username):
        obj, created = LoginAttemptRecord.objects.get_or_create(username=username)
        obj.count += 1
        obj.timestamp = timezone.now()
        obj.save()
