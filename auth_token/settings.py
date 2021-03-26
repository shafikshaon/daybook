from datetime import timedelta

from django.conf import settings
from rest_framework.settings import APISettings, api_settings

USER_SETTINGS = getattr(settings, 'AUTH_TOKEN_SETTING', None)

DEFAULTS = {
    'HASH_ALGORITHM': 'HS256',
    'JWT_SECRET_KEY': settings.SECRET_KEY,
    'AUTH_TOKEN_CHARACTER_LENGTH': 64,
    'TOKEN_EXPIRY': timedelta(seconds=3000),
    'USER_SERIALIZER': 'auth_token.serializers.UserSerializer',
    'TOKEN_LIMIT_PER_USER': None,
    'AUTO_REFRESH': False,
    'MIN_REFRESH_INTERVAL': 60,
    'AUTH_HEADER_PREFIX': 'Token',
    'EXPIRY_DATETIME_FORMAT': api_settings.DATETIME_FORMAT,
}

IMPORT_STRINGS = {
    'USER_SERIALIZER',
}

auth_token_settings = APISettings(user_settings=USER_SETTINGS, defaults=DEFAULTS, import_strings=IMPORT_STRINGS)
