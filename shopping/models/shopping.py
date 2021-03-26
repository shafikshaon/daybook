from django.conf import settings
from django.db import models

from gist.models.base import Base


class Shopping(Base):
    user = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.CASCADE, related_name='shopping_user')
    bazaar_date = models.DateField(blank=False, null=False)
    item_name = models.CharField(max_length=255, null=False, blank=False)
    item_weight = models.DecimalField(null=True, blank=True, decimal_places=2, max_digits=2)
    item_price = models.DecimalField(null=False, blank=False, decimal_places=2, max_digits=2)

    class Meta:
        app_label = "shopping"
        db_table = "shopping"
        verbose_name = "shopping"
        verbose_name_plural = "shopping"
