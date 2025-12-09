<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;

class ProductController extends Controller
{
    public function index(Request $request)
    {
        $category = $request->query('category');
        return response()->json(['products' => []]);
    }

    public function show($id)
    {
        return response()->json(['id' => $id]);
    }

    public function store(Request $request)
    {
        $name = $request->input('name');
        $price = $request->input('price');
        return response()->json(['id' => 1]);
    }

    public function update(Request $request, $id)
    {
        $name = $request->input('name');
        return response()->json(['id' => $id]);
    }

    public function destroy($id)
    {
        return response()->json(['message' => 'Product deleted']);
    }

    public function getReviews($productId)
    {
        return response()->json(['reviews' => []]);
    }
}

