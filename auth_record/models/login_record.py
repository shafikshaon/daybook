from django.db import models

from gist import httpagentparser


class LogRecord(models.Model):
    username = models.CharField(max_length=255, null=True, blank=True)
    ip_address = models.CharField(max_length=40, null=True, blank=True, verbose_name="IP")
    os = models.CharField(max_length=40, null=True, blank=True, verbose_name="Operating System")
    forwarded_by = models.CharField(max_length=1000, null=True, blank=True)
    user_agent = models.CharField(max_length=1000, null=True, blank=True)
    login_status = models.BooleanField(null=False, blank=False)
    timestamp = models.DateTimeField(auto_now_add=True)

    class Meta:
        app_label = 'auth_record'
        ordering = ['-timestamp']

    def __str__(self):
        return '%s|%s|%s|%s|%s|%s' % (
            self.username, self.ip_address, self.forwarded_by, self.user_agent, self.os, self.timestamp)


class LoginLogger(object):

    @classmethod
    def log_failed_login(cls, username, request):
        fields = cls.extract_log_info(username, request)
        fields['login_status'] = False
        LogRecord.objects.create(**fields)

    @classmethod
    def log_login(cls, username, request):
        fields = cls.extract_log_info(username, request)
        fields['login_status'] = True
        LogRecord.objects.create(**fields)

    @classmethod
    def extract_log_info(cls, username, request):
        USER_AGENT_MAX_LENGTH = LogRecord._meta.get_field('user_agent').max_length
        if request:
            ip_address, proxies = cls.extract_ip_address(request)
            user_agent = request.META.get('HTTP_USER_AGENT')
        else:
            ip_address = None
            proxies = None
            user_agent = None

        if user_agent is not None and len(user_agent) > USER_AGENT_MAX_LENGTH:
            user_agent = user_agent[:USER_AGENT_MAX_LENGTH]
        try:
            os = httpagentparser.detect(request['HTTP_USER_AGENT'])
        except:
            os = ''

        return {
            'username': username,
            'ip_address': ip_address,
            'os': os,
            'user_agent': user_agent,
            'forwarded_by': ",".join(proxies or [])
        }

    @classmethod
    def extract_ip_address(cls, request):
        client_ip = request.META.get('REMOTE_ADDR')
        proxies = None
        forwarded_for = request.META.get('HTTP_X_FORWARDED_FOR')
        if forwarded_for is not None:
            closest_proxy = client_ip
            forwarded_for_ips = [ip.strip() for ip in forwarded_for.split(',')]
            client_ip = forwarded_for_ips.pop(0)
            forwarded_for_ips.reverse()
            proxies = [closest_proxy] + forwarded_for_ips

        return (client_ip, proxies)
