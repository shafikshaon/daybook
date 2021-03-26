from django.utils import timezone
from rest_framework import parsers, renderers
from rest_framework.response import Response
from rest_framework.views import APIView

from auth_token.serializers import AuthTokenSerializer
from auth_token.settings import auth_token_settings

app_setting = auth_token_settings



class GetAuthToken(APIView):
    throttle_classes = ()
    permission_classes = ()
    parser_classes = (parsers.FormParser, parsers.MultiPartParser, parsers.JSONParser,)
    renderer_classes = (renderers.JSONRenderer,)
    serializer_class = AuthTokenSerializer

    def get_serializer_context(self):
        return {
            'request': self.request,
            'format': self.format_kwarg,
            'view': self
        }

    def get_serializer(self, *args, **kwargs):
        kwargs['context'] = self.get_serializer_context()
        return self.serializer_class(*args, **kwargs)

    def get_post_response_data(self, request, instance):
        UserSerializer = app_setting.USER_SERIALIZER

        data = {
            'expiry': instance.expiry,
            'token': instance.key
        }
        if UserSerializer is not None:
            serialize_data = UserSerializer(instance.user).data
            data["user"] = serialize_data
        return data

    def post(self, request, *args, **kwargs):
        from auth_token.utils import get_auth_token_model
        serializer = self.get_serializer(data=request.data)
        serializer.is_valid(raise_exception=True)
        user = serializer.validated_data['user']
        # payload = get_decoded_data_from_encoded_key()
        token_expiry = timezone.now() + app_setting.TOKEN_EXPIRY
        instance, created = get_auth_token_model().objects.get_or_create(user=user, expiry=token_expiry)
        data = self.get_post_response_data(request, instance)
        return Response(data)
