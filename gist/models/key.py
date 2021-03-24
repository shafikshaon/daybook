import uuid

from django.db import models
from django.db.models import Manager
from django.utils.translation import gettext_lazy as _

from gist.manager.gist_manager import GistManager


class Key(models.Model):
    id = models.BigAutoField(primary_key=True)
    uuid = models.UUIDField(unique=True, default=uuid.uuid4, editable=False)
    code = models.CharField(_('code'), max_length=20, null=False, blank=True)

    objects = GistManager()
    all_objects = Manager()

    class Meta:
        abstract = True

    def save(self, *args, **kwargs):
        if not self.code:
            self.code = None
        super(Key, self).save(*args, **kwargs)
