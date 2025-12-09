package com.example.api.controller;

import org.springframework.web.bind.annotation.*;
import org.springframework.http.ResponseEntity;
import java.util.Map;

@RestController
@RequestMapping("/products")
public class ProductController {

    @GetMapping
    public ResponseEntity<Map<String, Object>> getProducts(
        @RequestParam(required = false) String category
    ) {
        return ResponseEntity.ok(Map.of("products", new Object[]{}));
    }

    @GetMapping("/{id}")
    public ResponseEntity<Map<String, Object>> getProduct(@PathVariable String id) {
        return ResponseEntity.ok(Map.of("id", id));
    }

    @PostMapping
    public ResponseEntity<Map<String, Object>> createProduct(@RequestBody Map<String, Object> body) {
        return ResponseEntity.ok(Map.of("id", 1));
    }

    @PutMapping("/{id}")
    public ResponseEntity<Map<String, Object>> updateProduct(
        @PathVariable String id,
        @RequestBody Map<String, Object> body
    ) {
        return ResponseEntity.ok(Map.of("id", id));
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Map<String, String>> deleteProduct(@PathVariable String id) {
        return ResponseEntity.ok(Map.of("message", "Product deleted"));
    }

    @GetMapping("/{productId}/reviews")
    public ResponseEntity<Map<String, Object>> getProductReviews(@PathVariable String productId) {
        return ResponseEntity.ok(Map.of("reviews", new Object[]{}));
    }
}

