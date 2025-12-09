<?php

use Illuminate\Support\Facades\Route;

// Root route
Route::get('/', function () {
    return response()->json(['message' => 'Welcome to the API']);
});

// Health check
Route::get('/health', function () {
    return response()->json(['status' => 'ok']);
});

// Search route
Route::get('/search', function () {
    $q = request('q');
    $page = request('page');
    return response()->json(['query' => $q, 'page' => $page]);
});

// Auth routes
Route::post('/auth/login', function () {
    $email = request('email');
    $password = request('password');
    return response()->json(['message' => 'Login successful']);
});

// User routes
Route::get('/users', 'UserController@index');
Route::get('/users/{id}', 'UserController@show');
Route::post('/users', 'UserController@store');
Route::put('/users/{id}', 'UserController@update');
Route::delete('/users/{id}', 'UserController@destroy');
Route::get('/users/{userId}/posts', 'UserController@getPosts');

// Product routes
Route::get('/products', 'ProductController@index');
Route::get('/products/{id}', 'ProductController@show');
Route::post('/products', 'ProductController@store');
Route::put('/products/{id}', 'ProductController@update');
Route::delete('/products/{id}', 'ProductController@destroy');
Route::get('/products/{productId}/reviews', 'ProductController@getReviews');

// Settings route
Route::put('/settings/{userId}', function ($userId) {
    $theme = request('theme');
    return response()->json(['userId' => $userId, 'theme' => $theme]);
});

// Sessions route
Route::delete('/sessions/{sessionId}', function ($sessionId) {
    return response()->json(['message' => 'Session deleted', 'sessionId' => $sessionId]);
});

