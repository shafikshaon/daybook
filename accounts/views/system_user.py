from rest_framework import viewsets

from accounts.models import SystemUser
from accounts.serializers.system_user import SystemUserSerializer


class AccountViewSet(viewsets.ModelViewSet):
    """
    A simple ViewSet for viewing and editing accounts.
    """
    queryset = SystemUser.objects.all()
    serializer_class = SystemUserSerializer
    permission_classes = []
