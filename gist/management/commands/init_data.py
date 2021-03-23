from decouple import config
from django.contrib.sites.models import Site
from django.core.management import BaseCommand

from accounts.models.system_user import SystemUser


class Command(BaseCommand):

    def handle(self, *args, **options):
        self.generate_site_info()
        self.stdout.write(self.style.SUCCESS('Site info generated successfully.'))

        # Generate role, permissions for super organization
        self.generate_super_user()
        self.stdout.write(self.style.SUCCESS('Super user generated successfully.'))

    def generate_super_user(self):
        """
        Generate super user
        :return: user object
        """
        super_user = SystemUser.objects.filter(username='shafikshaon')
        if not super_user:
            super_user = SystemUser.objects.create(
                first_name='Shafikur',
                last_name='Rahman',
                username='shafikshaon@gmail.com',
                email='shafikshaon@gmail.com',
                is_superuser=True,
                is_active=True,
                is_organization_admin=True,
                is_staff=True,
                code='U-00001'
            )
            super_user.set_password('p@ss1234')
            super_user.save()
        return super_user

    def generate_site_info(self):
        site = Site.objects.first()
        site.domain = config('DOMAIN')
        site.name = config('SITE_NAME')
        site.save()
