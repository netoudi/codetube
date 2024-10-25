from django.contrib import admin
from core.models import Post, Tag

class PostAdmin(admin.ModelAdmin):
    list_display = ('title', 'tags_list', 'author', 'created_at')

    def tags_list(self, obj):
        return u", ".join(o.name for o in obj.tags.all())

    tags_list.short_description = "Tags"

    def save_model(self, request, obj, form, change):
        if not obj.pk:
            obj.author = request.user
        super().save_model(request, obj, form, change)

    def __str__(self):
        return self.title

# Register your models here.
admin.site.register(Post, PostAdmin)
admin.site.register(Tag)
