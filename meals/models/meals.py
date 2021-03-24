import django
from django.db import models
from django.utils import timezone

from accounts.models import SystemUser
from gist.models.base import Base


class Meal(Base):
    member = models.ForeignKey(SystemUser, on_delete=models.CASCADE, null=False, related_name='meals')
    meal_date = models.DateField(default=timezone.now, blank=False, null=False)
    breakfast = models.SmallIntegerField(default=1)
    lunch = models.SmallIntegerField(default=1)
    dinner = models.SmallIntegerField(default=1)

    class Meta:
        app_label = "meals"
        db_table = "meals"
        verbose_name = "meal"
        verbose_name_plural = "meals"
