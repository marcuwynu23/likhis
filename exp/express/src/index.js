const express = require('express');
const app = express();
const usersRouter = require('./routes/users');
const productsRouter = require('./routes/products');

// Root route
app.get('/', (req, res) => {
  res.json({ message: 'Welcome to the API' });
});

// Health check
app.get('/health', (req, res) => {
  res.json({ status: 'ok' });
});

// User routes
app.use('/users', usersRouter);

// Product routes
app.use('/products', productsRouter);

// Example route with query parameters
app.get('/search', (req, res) => {
  const { q, page, limit } = req.query;
  res.json({ query: q, page, limit });
});

// Example POST route
app.post('/auth/login', (req, res) => {
  const { email, password } = req.body;
  res.json({ message: 'Login successful' });
});

// Example PUT route
app.put('/settings/:userId', (req, res) => {
  const { userId } = req.params;
  res.json({ userId });
});

// Example DELETE route
app.delete('/sessions/:sessionId', (req, res) => {
  const { sessionId } = req.params;
  res.json({ message: 'Session deleted' });
});

module.exports = app;

