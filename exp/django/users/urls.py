from django.urls import path
from . import views

urlpatterns = [
    path('', views.get_users, name='get_users'),
    path('<int:user_id>/', views.get_user, name='get_user'),
    path('<int:user_id>/posts/', views.get_user_posts, name='get_user_posts'),
]

