package com.example.api.controller;

import org.springframework.web.bind.annotation.*;
import org.springframework.http.ResponseEntity;
import java.util.Map;

@RestController
@RequestMapping("/users")
public class UserController {

    @GetMapping
    public ResponseEntity<Map<String, Object>> getUsers(
        @RequestParam(required = false) Integer page,
        @RequestParam(required = false) Integer limit
    ) {
        return ResponseEntity.ok(Map.of("users", new Object[]{}));
    }

    @GetMapping("/{id}")
    public ResponseEntity<Map<String, Object>> getUser(@PathVariable String id) {
        return ResponseEntity.ok(Map.of("id", id));
    }

    @PostMapping
    public ResponseEntity<Map<String, Object>> createUser(@RequestBody Map<String, String> body) {
        return ResponseEntity.ok(Map.of("id", 1));
    }

    @PutMapping("/{id}")
    public ResponseEntity<Map<String, Object>> updateUser(
        @PathVariable String id,
        @RequestBody Map<String, String> body
    ) {
        return ResponseEntity.ok(Map.of("id", id));
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Map<String, String>> deleteUser(@PathVariable String id) {
        return ResponseEntity.ok(Map.of("message", "User deleted"));
    }

    @GetMapping("/{userId}/posts")
    public ResponseEntity<Map<String, Object>> getUserPosts(@PathVariable String userId) {
        return ResponseEntity.ok(Map.of("posts", new Object[]{}));
    }
}

