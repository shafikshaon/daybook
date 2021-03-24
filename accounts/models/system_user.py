from django.contrib.auth.models import AbstractUser
from django.db import models
from django.urls import reverse
from django.utils.html import escape
from django.utils.safestring import mark_safe
from django.utils.translation import gettext_lazy as _

from gist.models.activity import Activity
from gist.models.key import Key
from gist.models.timelog import TimeLog


def user_picture_directory_path(instance, filename):
    # file will be uploaded to MEDIA_ROOT/user_<id>/<filename>
    return 'user_{0}/{1}'.format(instance.id, filename)


class SystemUser(AbstractUser, Key, TimeLog, Activity):
    email = models.EmailField(_('email'), blank=False, null=False, unique=True)
    picture = models.ImageField(_('picture'), upload_to=user_picture_directory_path, blank=True, null=True)

    class Meta:
        app_label = 'accounts'
        db_table = 'accounts'
        ordering = ['-created_at']
        verbose_name = 'account'
        verbose_name_plural = 'accounts'

    def __str__(self):
        return mark_safe(
            '<a href="%s">%s</a>' % (
                reverse("profile", args=(self.username,)),
                escape(self.get_full_name() if self.get_full_name() else self.username)
            )
        )
