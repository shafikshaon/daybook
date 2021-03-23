__author__ = 'Shafikur Rahman'

from django.db import models


class GistManager(models.Manager):
    def get_queryset(self):
        return super(GistManager, self).get_queryset().filter(is_delete=False)
