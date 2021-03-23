from django.conf import settings
from django.contrib.auth import get_user_model
from django.db import models
from django.utils.translation import gettext_lazy as _


def get_sentinel_user():
    return get_user_model().objects.get_or_create(
        email='deleted@gmail.com',
        username='deleted@gmail.com',
        is_delete=True
    )[0]


class Activity(models.Model):
    is_delete = models.BooleanField(_('delete status'), default=False, null=False, blank=False)
    add_by = models.ForeignKey(
        verbose_name=_('add by'),
        to=settings.AUTH_USER_MODEL,
        on_delete=models.SET(get_sentinel_user),
        related_name='%(app_label)s_%(class)s_add_by',
        null=True,
        blank=False
    )
    change_by = models.ForeignKey(
        verbose_name=_('change by'),
        to=settings.AUTH_USER_MODEL,
        on_delete=models.SET(get_sentinel_user),
        related_name='%(app_label)s_%(class)s_change_by',
        null=True,
        blank=True
    )
    delete_by = models.ForeignKey(
        verbose_name=_('delete by'),
        to=settings.AUTH_USER_MODEL,
        on_delete=models.SET(get_sentinel_user),
        related_name='%(app_label)s_%(class)s_delete_by',
        null=True,
        blank=True
    )

    class Meta:
        abstract = True
