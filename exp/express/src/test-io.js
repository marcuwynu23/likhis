const express = require('express');
const app = express();

// This route should be ignored
app.get('/io', (req, res) => {
  res.json({ message: 'Socket.IO endpoint' });
});

// This route should also be ignored (io with path)
app.get('/io/socket', (req, res) => {
  res.json({ message: 'Socket.IO socket endpoint' });
});

// This route should be included
app.get('/api/users', (req, res) => {
  res.json({ users: [] });
});

