const express = require('express');
const router = express.Router();

// GET all products
router.get('/', (req, res) => {
  const { category, minPrice, maxPrice } = req.query;
  res.json({ products: [] });
});

// GET product by ID
router.get('/:id', (req, res) => {
  const { id } = req.params;
  res.json({ id });
});

// POST create product
router.post('/', (req, res) => {
  const { name, price } = req.body;
  res.json({ id: 1 });
});

// PUT update product
router.put('/:id', (req, res) => {
  const { id } = req.params;
  res.json({ id });
});

// DELETE product
router.delete('/:id', (req, res) => {
  const { id } = req.params;
  res.json({ message: 'Product deleted' });
});

module.exports = router;

