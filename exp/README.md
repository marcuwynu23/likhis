# Framework Examples

This folder contains example route structures for different frameworks. These are minimal examples showing only the route definitions and structure - no complete implementations or installation files.

## Structure

- **express/** - Node.js Express.js example
  - `src/index.js` - Main app with routes
  - `src/routes/users.js` - User routes
  - `src/routes/products.js` - Product routes

- **flask/** - Python Flask example
  - `app.py` - Main Flask app with routes
  - `routes/users.py` - User routes using Blueprint

- **django/** - Python Django example
  - `myproject/urls.py` - Main URL configuration
  - `users/urls.py` - User URL patterns
  - `users/views.py` - User views
  - `products/urls.py` - Product URL patterns
  - `products/views.py` - Product views

- **laravel/** - PHP Laravel example
  - `routes/api.php` - API route definitions
  - `app/Http/Controllers/UserController.php` - User controller
  - `app/Http/Controllers/ProductController.php` - Product controller

- **spring/** - Java Spring Boot example
  - `src/main/java/com/example/api/ApiApplication.java` - Main application
  - `src/main/java/com/example/api/controller/UserController.java` - User controller
  - `src/main/java/com/example/api/controller/ProductController.java` - Product controller
  - `src/main/java/com/example/api/controller/AuthController.java` - Auth controller

## Testing with likhis

You can test the API mapper tool on these examples:

```bash
# Express
likhis -p exp/express -o postman -F express

# Flask
likhis -p exp/flask -o postman -F flask

# Django
likhis -p exp/django -o postman -F django

# Laravel
likhis -p exp/laravel -o postman -F laravel

# Spring
likhis -p exp/spring -o postman -F spring
```

## Note

These are minimal examples showing only route structures. They are not complete applications and cannot be run without proper framework setup and dependencies.

