from django.conf import settings
from django.db import models
from django.utils import timezone
from django.utils.translation import gettext_lazy as _

from auth_token.settings import auth_token_settings

app_setting = auth_token_settings


def generate_jwt_key(payload):
    import jwt
    secret_key = app_setting.JWT_SECRET_KEY
    algorithm = app_setting.HASH_ALGORITHM
    expiry = timezone.now() + app_setting.TOKEN_EXPIRY
    payload['exp'] = int(expiry.timestamp())
    return jwt.encode(payload, secret_key, algorithm=algorithm)


class AuthToken(models.Model):
    """
    The default authorization token model.
    """
    key = models.CharField(_("Key"), max_length=255, db_index=True)
    user = models.ForeignKey(
        settings.AUTH_USER_MODEL, related_name='auth_tokens',
        on_delete=models.CASCADE, verbose_name=_("User")
    )
    created = models.DateTimeField(_("Created"), auto_now_add=True)
    expiry = models.DateTimeField(_("Expiry"), null=True, blank=True)

    class Meta:
        abstract = 'auth_token' not in settings.INSTALLED_APPS
        verbose_name = _("Token")
        verbose_name_plural = _("Tokens")
        unique_together = ['key', 'expiry']

    def save(self, *args, **kwargs):
        if not self.key:
            self.key = generate_jwt_key({})
            self.expiry = timezone.now() + app_setting.TOKEN_EXPIRY
        return super().save(*args, **kwargs)

    def __str__(self):
        return self.key
