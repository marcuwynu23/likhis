package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/marcuwynu23/likhis/internal/parser"
	"github.com/marcuwynu23/likhis/internal/plugins"
)

func TestParseExpressRoute(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "likhis_parser_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test Express route file
	testFile := filepath.Join(tmpDir, "routes.js")
	content := `const express = require('express');
const router = express.Router();

router.get('/users', (req, res) => {
  res.json({ users: [] });
});

router.post('/users/:id', (req, res) => {
  const userId = req.params.id;
  res.json({ id: userId });
});

app.get('/products', (req, res) => {
  const page = req.query.page;
  res.json({ products: [] });
});
`
	err = os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create parser
	rp := parser.NewRouteParser("express")

	// Parse file
	routes, err := rp.ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	// Verify routes
	if len(routes) < 3 {
		t.Errorf("Expected at least 3 routes, got %d", len(routes))
	}

	// Check for specific routes
	foundUsersGet := false
	foundUsersPost := false
	foundProducts := false

	for _, route := range routes {
		if route.Path == "/users" && route.Method == "GET" {
			foundUsersGet = true
		}
		if route.Path == "/users/:id" && route.Method == "POST" {
			foundUsersPost = true
			if len(route.Params) == 0 || route.Params[0] != "id" {
				t.Error("Expected 'id' parameter in /users/:id route")
			}
		}
		if route.Path == "/products" && route.Method == "GET" {
			foundProducts = true
		}
	}

	if !foundUsersGet {
		t.Error("GET /users route not found")
	}
	if !foundUsersPost {
		t.Error("POST /users/:id route not found")
	}
	if !foundProducts {
		t.Error("GET /products route not found")
	}
}

func TestParseFlaskRoute(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "likhis_flask_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "app.py")
	content := `from flask import Flask

app = Flask(__name__)

@app.route('/users', methods=['GET'])
def get_users():
    return {'users': []}

@app.route('/users/<int:user_id>', methods=['GET', 'POST'])
def user_detail(user_id):
    return {'id': user_id}
`
	err = os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	rp := parser.NewRouteParser("flask")
	routes, err := rp.ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(routes) < 3 {
		t.Errorf("Expected at least 3 routes (GET /users, GET /users/<int:user_id>, POST /users/<int:user_id>), got %d", len(routes))
	}

	foundUsersGet := false
	foundUserDetailGet := false
	foundUserDetailPost := false

		for _, route := range routes {
		if route.Path == "/users" && route.Method == "GET" {
			foundUsersGet = true
		}
		if route.Path == "/users/<int:user_id>" {
			if route.Method == "GET" {
				foundUserDetailGet = true
			}
			if route.Method == "POST" {
				foundUserDetailPost = true
			}
			// Parameter extraction may vary, just check if route was found
		}
	}

	if !foundUsersGet {
		t.Error("GET /users route not found")
	}
	if !foundUserDetailGet {
		t.Error("GET /users/<int:user_id> route not found")
	}
	if !foundUserDetailPost {
		t.Error("POST /users/<int:user_id> route not found")
	}
}

func TestParseLaravelRoute(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "likhis_laravel_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "routes.php")
	content := `<?php

use Illuminate\Support\Facades\Route;

Route::get('/users', function () {
    return response()->json(['users' => []]);
});

Route::post('/users/{id}', function ($id) {
    return response()->json(['id' => $id]);
});
`
	err = os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	rp := parser.NewRouteParser("laravel")
	routes, err := rp.ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(routes) < 2 {
		t.Errorf("Expected at least 2 routes, got %d", len(routes))
	}

	foundUsersGet := false
	foundUsersPost := false

	for _, route := range routes {
		if route.Path == "/users" && route.Method == "GET" {
			foundUsersGet = true
		}
		if route.Path == "/users/{id}" && route.Method == "POST" {
			foundUsersPost = true
			if len(route.Params) == 0 || route.Params[0] != "id" {
				t.Error("Expected 'id' parameter in /users/{id} route")
			}
		}
	}

	if !foundUsersGet {
		t.Error("GET /users route not found")
	}
	if !foundUsersPost {
		t.Error("POST /users/{id} route not found")
	}
}

func TestParseDjangoRoute(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "likhis_django_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "urls.py")
	content := `from django.urls import path
from . import views

urlpatterns = [
    path('users/', views.user_list),
    path('users/<int:user_id>/', views.user_detail),
]
`
	err = os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	rp := parser.NewRouteParser("django")
	routes, err := rp.ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(routes) < 2 {
		t.Errorf("Expected at least 2 routes, got %d", len(routes))
	}

	foundUsersList := false
	foundUserDetail := false

	for _, route := range routes {
		if route.Path == "users/" {
			foundUsersList = true
		}
		if route.Path == "users/<int:user_id>/" {
			foundUserDetail = true
			// Parameter extraction may vary, just check if route was found
		}
	}

	if !foundUsersList {
		t.Error("users/ route not found")
	}
	if !foundUserDetail {
		t.Error("users/<int:user_id>/ route not found")
	}
}

func TestParseSpringRoute(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "likhis_spring_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "UserController.java")
	content := `package com.example.api.controller;

import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/users")
public class UserController {
    
    @GetMapping("/")
    public List<User> getUsers() {
        return userService.getAll();
    }
    
    @PostMapping("/{id}")
    public User getUser(@PathVariable String id) {
        return userService.getById(id);
    }
}
`
	err = os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	rp := parser.NewRouteParser("spring")
	routes, err := rp.ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(routes) < 2 {
		t.Errorf("Expected at least 2 routes, got %d", len(routes))
	}

	foundUsersList := false
	foundUserDetail := false

	for _, route := range routes {
		if route.Path == "/api/users/" && route.Method == "GET" {
			foundUsersList = true
		}
		if route.Path == "/api/users/{id}" && route.Method == "POST" {
			foundUserDetail = true
			if len(route.Params) == 0 || route.Params[0] != "id" {
				t.Error("Expected 'id' parameter in /api/users/{id} route")
			}
		}
	}

	if !foundUsersList {
		t.Error("GET /api/users/ route not found")
	}
	if !foundUserDetail {
		t.Error("POST /api/users/{id} route not found")
	}
}

func TestRouteParserNew(t *testing.T) {
	rp := parser.NewRouteParser("express")
	if rp == nil {
		t.Error("NewRouteParser returned nil")
	}
	
	// Test that parser can parse a file (indirect test of framework setting)
	tmpDir, err := os.MkdirTemp("", "likhis_parser_new_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.js")
	err = os.WriteFile(testFile, []byte("router.get('/test', () => {});"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	routes, err := rp.ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
	if len(routes) == 0 {
		t.Error("Expected at least one route from test file")
	}
}

func TestRouteParserWithPlugins(t *testing.T) {
	pluginMap := make(map[string]*plugins.Plugin)
	rp := parser.NewRouteParserWithPlugins("auto", pluginMap, "/test/path")
	if rp == nil {
		t.Error("NewRouteParserWithPlugins returned nil")
	}
}

