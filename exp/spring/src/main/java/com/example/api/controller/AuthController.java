package com.example.api.controller;

import org.springframework.web.bind.annotation.*;
import org.springframework.http.ResponseEntity;
import java.util.Map;

@RestController
public class AuthController {

    @GetMapping("/")
    public ResponseEntity<Map<String, String>> index() {
        return ResponseEntity.ok(Map.of("message", "Welcome to the API"));
    }

    @GetMapping("/health")
    public ResponseEntity<Map<String, String>> health() {
        return ResponseEntity.ok(Map.of("status", "ok"));
    }

    @GetMapping("/search")
    public ResponseEntity<Map<String, Object>> search(
        @RequestParam(required = false) String q,
        @RequestParam(required = false) Integer page
    ) {
        return ResponseEntity.ok(Map.of("query", q != null ? q : "", "page", page != null ? page : 1));
    }

    @PostMapping("/auth/login")
    public ResponseEntity<Map<String, String>> login(@RequestBody Map<String, String> body) {
        return ResponseEntity.ok(Map.of("message", "Login successful"));
    }

    @PutMapping("/settings/{userId}")
    public ResponseEntity<Map<String, Object>> updateSettings(
        @PathVariable String userId,
        @RequestBody Map<String, String> body
    ) {
        return ResponseEntity.ok(Map.of("userId", userId, "theme", body.getOrDefault("theme", "")));
    }

    @DeleteMapping("/sessions/{sessionId}")
    public ResponseEntity<Map<String, Object>> deleteSession(@PathVariable String sessionId) {
        return ResponseEntity.ok(Map.of("message", "Session deleted", "sessionId", sessionId));
    }
}

