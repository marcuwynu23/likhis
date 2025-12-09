from django.http import JsonResponse

def get_users(request):
    page = request.GET.get('page')
    limit = request.GET.get('limit')
    return JsonResponse({'users': []})

def get_user(request, user_id):
    return JsonResponse({'id': user_id})

def get_user_posts(request, user_id):
    return JsonResponse({'posts': []})

