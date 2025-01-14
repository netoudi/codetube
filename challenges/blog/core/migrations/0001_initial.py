# Generated by Django 5.1.2 on 2024-10-25 13:29

import django.db.models.deletion
from django.conf import settings
from django.db import migrations, models


class Migration(migrations.Migration):

    initial = True

    dependencies = [
        migrations.swappable_dependency(settings.AUTH_USER_MODEL),
    ]

    operations = [
        migrations.CreateModel(
            name='Tag',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(max_length=50, unique=True, verbose_name='Nome')),
            ],
            options={
                'verbose_name': 'Tag',
                'verbose_name_plural': 'Tags',
            },
        ),
        migrations.CreateModel(
            name='Post',
            fields=[
                ('id', models.BigAutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('title', models.CharField(max_length=100, unique=True, verbose_name='Título')),
                ('content', models.TextField(verbose_name='Conteúdo')),
                ('created_at', models.DateTimeField(auto_now_add=True, null=True, verbose_name='Criado em')),
                ('author', models.ForeignKey(editable=False, null=True, on_delete=django.db.models.deletion.PROTECT, related_name='posts', to=settings.AUTH_USER_MODEL, verbose_name='Autor')),
                ('tags', models.ManyToManyField(related_name='posts', to='core.tag', verbose_name='Tags')),
            ],
            options={
                'verbose_name': 'Post',
                'verbose_name_plural': 'Posts',
            },
        ),
    ]
