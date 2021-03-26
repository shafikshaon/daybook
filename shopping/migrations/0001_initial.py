# Generated by Django 3.1.7 on 2021-03-26 05:21

from django.conf import settings
from django.db import migrations, models
import django.db.models.deletion
import gist.models.activity
import uuid


class Migration(migrations.Migration):

    initial = True

    dependencies = [
        migrations.swappable_dependency(settings.AUTH_USER_MODEL),
    ]

    operations = [
        migrations.CreateModel(
            name='Shopping',
            fields=[
                ('is_deleted', models.BooleanField(default=False, verbose_name='delete status')),
                ('id', models.BigAutoField(primary_key=True, serialize=False)),
                ('uuid', models.UUIDField(default=uuid.uuid4, editable=False, unique=True)),
                ('code', models.CharField(blank=True, max_length=20, verbose_name='code')),
                ('created_at', models.DateTimeField(auto_now_add=True, verbose_name='created at')),
                ('changed_at', models.DateTimeField(auto_now=True, null=True, verbose_name='changed at')),
                ('deleted_at', models.DateTimeField(blank=True, editable=False, null=True, verbose_name='deleted at')),
                ('bazaar_date', models.DateField()),
                ('item_name', models.CharField(max_length=255)),
                ('item_weight', models.DecimalField(blank=True, decimal_places=2, max_digits=2, null=True)),
                ('item_price', models.DecimalField(decimal_places=2, max_digits=2)),
                ('add_by', models.ForeignKey(null=True, on_delete=models.SET(gist.models.activity.get_sentinel_user), related_name='shopping_shopping_add_by', to=settings.AUTH_USER_MODEL, verbose_name='add by')),
                ('change_by', models.ForeignKey(blank=True, null=True, on_delete=models.SET(gist.models.activity.get_sentinel_user), related_name='shopping_shopping_change_by', to=settings.AUTH_USER_MODEL, verbose_name='change by')),
                ('delete_by', models.ForeignKey(blank=True, null=True, on_delete=models.SET(gist.models.activity.get_sentinel_user), related_name='shopping_shopping_delete_by', to=settings.AUTH_USER_MODEL, verbose_name='delete by')),
                ('user', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, related_name='shopping_user', to=settings.AUTH_USER_MODEL)),
            ],
            options={
                'verbose_name': 'shopping',
                'verbose_name_plural': 'shopping',
                'db_table': 'shopping',
            },
        ),
        migrations.CreateModel(
            name='ExtraCost',
            fields=[
                ('is_deleted', models.BooleanField(default=False, verbose_name='delete status')),
                ('id', models.BigAutoField(primary_key=True, serialize=False)),
                ('uuid', models.UUIDField(default=uuid.uuid4, editable=False, unique=True)),
                ('code', models.CharField(blank=True, max_length=20, verbose_name='code')),
                ('created_at', models.DateTimeField(auto_now_add=True, verbose_name='created at')),
                ('changed_at', models.DateTimeField(auto_now=True, null=True, verbose_name='changed at')),
                ('deleted_at', models.DateTimeField(blank=True, editable=False, null=True, verbose_name='deleted at')),
                ('expense_date', models.DateField()),
                ('cost_name', models.CharField(max_length=255)),
                ('cost', models.DecimalField(decimal_places=2, max_digits=2)),
                ('add_by', models.ForeignKey(null=True, on_delete=models.SET(gist.models.activity.get_sentinel_user), related_name='shopping_extracost_add_by', to=settings.AUTH_USER_MODEL, verbose_name='add by')),
                ('change_by', models.ForeignKey(blank=True, null=True, on_delete=models.SET(gist.models.activity.get_sentinel_user), related_name='shopping_extracost_change_by', to=settings.AUTH_USER_MODEL, verbose_name='change by')),
                ('delete_by', models.ForeignKey(blank=True, null=True, on_delete=models.SET(gist.models.activity.get_sentinel_user), related_name='shopping_extracost_delete_by', to=settings.AUTH_USER_MODEL, verbose_name='delete by')),
                ('user', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, related_name='extra_cost_user', to=settings.AUTH_USER_MODEL)),
            ],
            options={
                'verbose_name': 'extra_cost',
                'verbose_name_plural': 'extra_costs',
                'db_table': 'extra_costs',
            },
        ),
    ]