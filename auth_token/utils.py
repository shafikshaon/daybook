def get_auth_token_model():
    from django.conf import settings
    if hasattr(settings, 'AUTH_TOKEN_MODEL'):
        from django.apps import apps
        return apps.get_model(settings.AUTH_TOKEN_MODEL)
    from auth_token.models import AuthToken
    return AuthToken


def get_decoded_data_from_encoded_key(encoded=''):
    from django.utils.translation import gettext_lazy as _
    from rest_framework import exceptions
    from auth_token.settings import auth_token_settings
    import jwt

    app_setting = auth_token_settings
    secret_key = app_setting.JWT_SECRET_KEY
    algorithm = app_setting.HASH_ALGORITHM
    try:
        return jwt.decode(encoded, key=secret_key, algorithms=[algorithm])
    except jwt.ExpiredSignatureError:
        raise exceptions.AuthenticationFailed(_('Token expired.'))
