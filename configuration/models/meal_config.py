from django.db import models
from django.utils.translation import gettext_lazy as _

__author__ = 'Shafikur Rahman'

from gist.models.base import Base


class MealConfig(Base):
    breakfast = models.DecimalField(
        max_digits=5,
        decimal_places=2,
        blank=False,
        null=False,
        default=0.5,
        help_text=_('This will consider as meal unit i.e. for breakfast if value set to 0.5 then breakfast consider as '
                    'half(0.5)')
    )
    lunch = models.DecimalField(
        max_digits=5,
        decimal_places=2,
        blank=False,
        null=False,
        default=1,
        help_text=_('This will consider as meal unit i.e. for lunch if value set to 1 then lunch consider as one(1)')
    )
    dinner = models.DecimalField(
        max_digits=5,
        decimal_places=2,
        blank=False,
        null=False,
        default=1,
        help_text=_('This will consider as meal unit i.e. for dinner if value set to 1 then dinner consider as '
                    'one(1)')
    )

    class Meta:
        app_label = 'configuration'
        db_table = 'meal_configuration'
        ordering = ['-created_at']
        verbose_name = "meal_configuration"
        verbose_name_plural = 'meal_configuration'

    def __str__(self):
        return self.code
