from rest_framework import serializers

from accounts.models import SystemUser


class SystemUserSerializer(serializers.ModelSerializer):
    class Meta:
        model = SystemUser
        fields = ['id', 'username']
