from django.db import models
from django.utils.translation import gettext_lazy as _


class TimeLog(models.Model):
    add_at = models.DateTimeField(_('code'), auto_now_add=True, editable=False)
    change_at = models.DateTimeField(_('change at'), null=True, blank=True, editable=False)
    delete_at = models.DateTimeField(_('delete at'), null=True, blank=True, editable=False)

    class Meta:
        abstract = True
