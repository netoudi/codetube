"""
URL configuration for blog project.

The `urlpatterns` list routes URLs to views. For more information please see:
    https://docs.djangoproject.com/en/5.1/topics/http/urls/
Examples:
Function views
    1. Add an import:  from my_app import views
    2. Add a URL to urlpatterns:  path('', views.home, name='home')
Class-based views
    1. Add an import:  from other_app.views import Home
    2. Add a URL to urlpatterns:  path('', Home.as_view(), name='home')
Including another URLconf
    1. Import the include() function: from django.urls import include, path
    2. Add a URL to urlpatterns:  path('blog/', include('blog.urls'))
"""
from django.contrib import admin
from django.urls import path
from core.views import posts_view, posts_by_tag_view, post_view

urlpatterns = [
    path('admin/', admin.site.urls),

    path("blog/", posts_view, name="core_posts_view"),
    path("blog/tag/<int:id>/", posts_by_tag_view, name="core_posts_by_tag_view"),
    path("blog/post/<int:id>/", post_view, name="core_post_view"),
]
