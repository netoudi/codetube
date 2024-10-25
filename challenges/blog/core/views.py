from django.shortcuts import render
from core.models import Post, Tag

# Create your views here.
def posts_view(request):
    tags = Tag.objects.all()
    posts = Post.objects.all()
    return render(request, 'core/posts.html', { "tags": tags, "posts": posts })


def posts_by_tag_view(request, id):
    tags = Tag.objects.all()
    posts = Post.objects.filter(tags__id__icontains=id)
    tag = Tag.objects.get(id=id)
    return render(request, 'core/posts.html', { "tags": tags, "posts": posts, "tag": tag, "id": id })


def post_view(request, id):
    post = Post.objects.get(id=id)
    return render(request, 'core/post.html', { "post": post })
