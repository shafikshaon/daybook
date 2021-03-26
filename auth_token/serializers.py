from django.contrib.auth import authenticate, get_user_model
from django.utils.translation import gettext_lazy as _

from rest_framework import serializers

from auth_record.models.login_attempt_record import LoginAttemptLogger
from auth_record.models.login_record import LoginLogger

User = get_user_model()

username_field = User.USERNAME_FIELD if hasattr(User, 'USERNAME_FIELD') else 'username'


class UserSerializer(serializers.ModelSerializer):
    class Meta:
        model = User
        fields = ('id', username_field,)


class AuthTokenSerializer(serializers.Serializer):
    username = serializers.CharField(
        label=_("Username"),
        write_only=True
    )
    password = serializers.CharField(
        label=_("Password"),
        style={'input_type': 'password'},
        trim_whitespace=False,
        write_only=True
    )
    token = serializers.CharField(
        label=_("Token"),
        read_only=True
    )

    def validate(self, attrs):
        username = attrs.get('username')
        password = attrs.get('password')

        if username and password:
            user = authenticate(request=self.context.get('request'),
                                username=username, password=password)

            # The authenticate call simply returns None for is_active=False
            # users. (Assuming the default ModelBackend authentication
            # backend.)
            if not user:
                msg = _('Unable to log in with provided credentials.')
                LoginLogger.log_failed_login(username=username, request=self.context['request'])
                LoginAttemptLogger.increment(username)
                raise serializers.ValidationError(msg, code='authorization')
        else:
            msg = _('Must include "username" and "password".')
            LoginLogger.log_failed_login(username=username, request=self.context['request'])
            LoginAttemptLogger.increment(username)
            raise serializers.ValidationError(msg, code='authorization')

        attrs['user'] = user
        return attrs
