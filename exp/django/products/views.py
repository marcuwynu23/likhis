from django.http import JsonResponse

def get_products(request):
    category = request.GET.get('category')
    return JsonResponse({'products': []})

def get_product(request, product_id):
    return JsonResponse({'id': product_id})

def get_product_reviews(request, product_id):
    return JsonResponse({'reviews': []})

