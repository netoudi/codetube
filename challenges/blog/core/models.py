from django.db import models

# Create your models here.
class Post(models.Model):
    title = models.CharField(max_length=100, unique=True, verbose_name='Título')
    content = models.TextField(verbose_name='Conteúdo')
    created_at = models.DateTimeField(verbose_name='Criado em', null=True, editable=False, auto_now_add=True)
    tags = models.ManyToManyField('Tag', verbose_name='Tags', related_name='posts')
    author = models.ForeignKey('auth.User', on_delete=models.PROTECT, verbose_name='Autor', related_name='posts', editable=False, null=True)

    class Meta:
        verbose_name = 'Post'
        verbose_name_plural = 'Posts'

    def __str__(self):
        return self.title


class Tag(models.Model):
    name = models.CharField(max_length=50, unique=True, verbose_name='Nome')

    class Meta:
        verbose_name = 'Tag'
        verbose_name_plural = 'Tags'

    def __str__(self):
        return self.name
