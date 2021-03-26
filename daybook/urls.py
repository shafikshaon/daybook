from django.urls import path
from rest_framework.routers import DefaultRouter

from accounts.views.system_user import AccountViewSet
from auth_token.views import GetAuthToken

router = DefaultRouter()
router.register(r'users', AccountViewSet, basename='user')
urlpatterns = router.urls
urlpatterns += [
    path('login/', GetAuthToken.as_view()),
]
