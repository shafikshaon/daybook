# Generated by Django 3.1.7 on 2021-03-26 13:19

from django.db import migrations, models


class Migration(migrations.Migration):

    initial = True

    dependencies = [
    ]

    operations = [
        migrations.CreateModel(
            name='LoginAttemptRecord',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('username', models.CharField(blank=True, max_length=255, null=True, unique=True)),
                ('count', models.PositiveIntegerField(blank=True, default=0, null=True)),
                ('timestamp', models.DateTimeField(auto_now_add=True)),
            ],
        ),
        migrations.CreateModel(
            name='LogRecord',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('username', models.CharField(blank=True, max_length=255, null=True)),
                ('ip_address', models.CharField(blank=True, max_length=40, null=True, verbose_name='IP')),
                ('os', models.CharField(blank=True, max_length=40, null=True, verbose_name='Operating System')),
                ('forwarded_by', models.CharField(blank=True, max_length=1000, null=True)),
                ('user_agent', models.CharField(blank=True, max_length=1000, null=True)),
                ('login_status', models.BooleanField()),
                ('timestamp', models.DateTimeField(auto_now_add=True)),
            ],
            options={
                'ordering': ['-timestamp'],
            },
        ),
    ]
