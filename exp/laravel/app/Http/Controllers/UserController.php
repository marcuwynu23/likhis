<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;

class UserController extends Controller
{
    public function index(Request $request)
    {
        $page = $request->query('page');
        $limit = $request->query('limit');
        return response()->json(['users' => []]);
    }

    public function show($id)
    {
        return response()->json(['id' => $id]);
    }

    public function store(Request $request)
    {
        $name = $request->input('name');
        $email = $request->input('email');
        return response()->json(['id' => 1]);
    }

    public function update(Request $request, $id)
    {
        $name = $request->input('name');
        return response()->json(['id' => $id]);
    }

    public function destroy($id)
    {
        return response()->json(['message' => 'User deleted']);
    }

    public function getPosts($userId)
    {
        return response()->json(['posts' => []]);
    }
}

