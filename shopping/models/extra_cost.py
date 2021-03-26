from django.conf import settings
from django.db import models

from gist.models.base import Base


class ExtraCost(Base):
    user = models.ForeignKey(settings.AUTH_USER_MODEL, on_delete=models.CASCADE, related_name='extra_cost_user')
    expense_date = models.DateField(blank=False, null=False)
    cost_name = models.CharField(max_length=255, null=False, blank=False)
    cost = models.DecimalField(null=False, blank=False, decimal_places=2, max_digits=2)

    class Meta:
        app_label = "shopping"
        db_table = "extra_costs"
        verbose_name = "extra_cost"
        verbose_name_plural = "extra_costs"
