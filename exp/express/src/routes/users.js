const express = require('express');
const router = express.Router();

// GET all users
router.get('/', (req, res) => {
  const { page, limit, sort } = req.query;
  res.json({ users: [] });
});

// GET user by ID
router.get('/:id', (req, res) => {
  const { id } = req.params;
  res.json({ id });
});

// POST create user
router.post('/', (req, res) => {
  const { name, email } = req.body;
  res.json({ id: 1 });
});

// PUT update user
router.put('/:id', (req, res) => {
  const { id } = req.params;
  res.json({ id });
});

// DELETE user
router.delete('/:id', (req, res) => {
  const { id } = req.params;
  res.json({ message: 'User deleted' });
});

// GET user's posts
router.get('/:userId/posts', (req, res) => {
  const { userId } = req.params;
  res.json({ posts: [] });
});

module.exports = router;

