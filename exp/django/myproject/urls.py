from django.contrib import admin
from django.urls import path, include
from django.http import JsonResponse

def health(request):
    return JsonResponse({'status': 'ok'})

def search(request):
    q = request.GET.get('q')
    page = request.GET.get('page')
    return JsonResponse({'query': q, 'page': page})

urlpatterns = [
    path('admin/', admin.site.urls),
    path('api/users/', include('users.urls')),
    path('api/products/', include('products.urls')),
    path('api/health/', health, name='health'),
    path('api/search/', search, name='search'),
]

