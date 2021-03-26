import os
from datetime import datetime, timedelta
from pathlib import Path

# Build paths inside the project like this: BASE_DIR / 'subdir'.
from decouple import config
from django.conf import settings

BASE_DIR = Path(__file__).resolve().parent.parent

# SECURITY WARNING: keep the secret key used in production secret!
SECRET_KEY = config('SECRET_KEY')

# SECURITY WARNING: don't run with debug turned on in production!
DEBUG = config('DEBUG', cast=bool)

ALLOWED_HOSTS = []

SITE_ID = 1

# Application definition
INSTALLED_APPS = [
    'django.contrib.sites',
    'django.contrib.auth',
    'django.contrib.contenttypes',
    'django.contrib.sessions',
    'django.contrib.messages',
    'django.contrib.staticfiles',
]

INSTALLED_APPS += [
    'gist',
    'accounts',
    'meals',
    'configuration',
    'auth_token',
    'auth_record',
    'shopping',
    'rest_framework'
]

# Middleware
MIDDLEWARE = [
    'django.middleware.security.SecurityMiddleware',
    'django.contrib.sessions.middleware.SessionMiddleware',
    'django.middleware.common.CommonMiddleware',
    'django.middleware.csrf.CsrfViewMiddleware',
    'django.contrib.auth.middleware.AuthenticationMiddleware',
    'django.contrib.messages.middleware.MessageMiddleware',
    'django.middleware.clickjacking.XFrameOptionsMiddleware',
]

if config('DEBUG', cast=bool):
    MIDDLEWARE += [
        'gist.middleware.stats_middleware.StatsMiddleware',
    ]

ROOT_URLCONF = 'daybook.urls'

WSGI_APPLICATION = 'daybook.wsgi.application'

# Database
DATABASES = {
    'default': {
        'ENGINE': config('DB_ENGINE'),
        'NAME': config('DB_NAME'),
        'USER': config('DB_USER'),
        'PASSWORD': config('DB_PASSWORD'),
        'HOST': config('DB_HOST')
    }
}

# Password validation
AUTH_PASSWORD_VALIDATORS = [
    {
        'NAME': 'django.contrib.auth.password_validation.UserAttributeSimilarityValidator',
    },
    {
        'NAME': 'django.contrib.auth.password_validation.MinimumLengthValidator',
        'OPTIONS': {
            'min_length': 6
        }
    },
    {
        'NAME': 'django.contrib.auth.password_validation.CommonPasswordValidator',
    },
    {
        'NAME': 'django.contrib.auth.password_validation.NumericPasswordValidator',
    },
]

# Authentication
AUTHENTICATION_BACKENDS = (
    'accounts.authentication.EmailOrUsernameModelBackend',
    'django.contrib.auth.backends.ModelBackend'
)

AUTH_USER_MODEL = 'accounts.SystemUser'

AUTH_TOKEN_MODEL = 'auth_token.AuthToken'

REST_FRAMEWORK = {
    'DEFAULT_AUTHENTICATION_CLASSES': [
        'auth_token.authentication.BasicAuthentication',
        'auth_token.authentication.TokenAuthentication',
    ]
}

AUTH_TOKEN_CONFIG = {
    'JWT_SECRET_KEY': settings.SECRET_KEY,
    'JWT_ALGORITHM': 'HS256',
    # HS256 - HMAC using SHA-256 hash algorithm (default)
    # HS384 - HMAC using SHA-384 hash algorithm
    # HS512 - HMAC using SHA-512 hash algorithm
    # ES256 - ECDSA signature algorithm using SHA-256 hash algorithm
    # ES384 - ECDSA signature algorithm using SHA-384 hash algorithm
    # ES512 - ECDSA signature algorithm using SHA-512 hash algorithm
    # RS256 - RSASSA-PKCS1-v1_5 signature algorithm using SHA-256 hash algorithm
    # RS384 - RSASSA-PKCS1-v1_5 signature algorithm using SHA-384 hash algorithm
    # RS512 - RSASSA-PKCS1-v1_5 signature algorithm using SHA-512 hash algorithm
    # PS256 - RSASSA-PSS signature using SHA-256 and MGF1 padding with SHA-256
    # PS384 - RSASSA-PSS signature using SHA-384 and MGF1 padding with SHA-384
    # PS512 - RSASSA-PSS signature using SHA-512 and MGF1 padding with SHA-512
    # EdDSA - Ed25519255 signature using SHA-512. Provides 128-bit security
    'JWT_EXPIRATION_DELTA': timedelta(seconds=300),
    #
    'HASH_ALGORITHM': 'HS256',
    'AUTH_TOKEN_CHARACTER_LENGTH': 64,
    'TOKEN_EXPIRY': timedelta(seconds=3000),
    'USER_SERIALIZER': 'auth_token.serializers.UserSerializer',
    'TOKEN_LIMIT_PER_USER': None,
    'AUTO_REFRESH': False,
    'MIN_REFRESH_INTERVAL': 60,
    'AUTH_HEADER_PREFIX': 'Bearer',
}

AUTH_TOKEN_SETTING = AUTH_TOKEN_CONFIG

# Templates
TEMPLATES = [
    {
        'BACKEND': 'django.template.backends.django.DjangoTemplates',
        'DIRS': [BASE_DIR / 'templates'],
        'APP_DIRS': True,
        'OPTIONS': {
            'context_processors': [
                'django.template.context_processors.debug',
                'django.template.context_processors.request',
                'django.contrib.auth.context_processors.auth',
                'django.contrib.messages.context_processors.messages',
            ],
        },
    },
]

# Internationalization
LANGUAGE_CODE = 'en-us'
TIME_ZONE = config('TIME_ZONE', default='UTC')
USE_I18N = True
USE_L10N = True
USE_TZ = True

# Static files (CSS, JavaScript, Images)
STATIC_URL = '/static/'
STATIC_ROOT = os.path.join(BASE_DIR, 'staticfiles')
MEDIA_URL = '/media/'
MEDIA_ROOT = os.path.join(BASE_DIR, 'media')
STATICFILES_DIRS = [
    os.path.join(BASE_DIR, "static")
]

# Email configuration
EMAIL_BACKEND = 'django.core.mail.backends.smtp.EmailBackend'
MAILER_EMAIL_BACKEND = EMAIL_BACKEND
EMAIL_HOST = config('EMAIL_HOST')
EMAIL_HOST_PASSWORD = config('EMAIL_HOST_PASSWORD')
EMAIL_HOST_USER = config('EMAIL_HOST_USER')
EMAIL_PORT = config('EMAIL_PORT', cast=int)
EMAIL_USE_SSL = config('EMAIL_USE_SSL', cast=bool)
DEFAULT_FROM_EMAIL = EMAIL_HOST_USER

APPEND_SLASH = True
LOGIN_REDIRECT_URL = '/'
LOGOUT_REDIRECT_URL = '/auth/login/'
LOGIN_URL = '/auth/login/'

PAGINATE_BY = config('PAGINATE_BY', default=50, cast=int)
