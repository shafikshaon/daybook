from django.db import models
from django.utils.translation import gettext_lazy as _


class TimeLog(models.Model):
    created_at = models.DateTimeField(_('created at'), auto_now_add=True, editable=False)
    changed_at = models.DateTimeField(_('changed at'), auto_now=True, null=True, blank=True, editable=False)
    deleted_at = models.DateTimeField(_('deleted at'), null=True, blank=True, editable=False)

    class Meta:
        abstract = True
