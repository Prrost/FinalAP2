import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useParams, useNavigate } from 'react-router-dom';

const BookDetails = () => {
  const [book, setBook] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const { id } = useParams();
  const navigate = useNavigate();

  useEffect(() => {
    const fetchBookDetails = async () => {
      try {
        const response = await axios.get(`http://localhost:8080/api/books/${id}`);
        setBook(response.data);
        setLoading(false);
      } catch (err) {
        setError('Failed to fetch book details');
        setLoading(false);
      }
    };

    fetchBookDetails();
  }, [id]);

  const handleOrder = async () => {
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        navigate('/login');
        return;
      }

      await axios.post(
        'http://localhost:8080/api/orders',
        { bookId: id },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      alert('Order placed successfully!');
    } catch (err) {
      alert('Failed to place order. Please try again.');
    }
  };

  if (loading) return <div className="loading">Loading...</div>;
  if (error) return <div className="error">{error}</div>;
  if (!book) return <div>Book not found</div>;

  return (
    <div className="book-details">
      <h2>{book.title}</h2>
      <div className="book-info">
        <p><strong>Author:</strong> {book.author}</p>
        <p><strong>Price:</strong> ${book.price}</p>
        <p><strong>Available Quantity:</strong> {book.quantity}</p>
        <p><strong>Description:</strong> {book.description}</p>
        <p><strong>ISBN:</strong> {book.isbn}</p>
      </div>
      <button 
        className="order-button"
        onClick={handleOrder}
        disabled={book.quantity === 0}
      >
        {book.quantity === 0 ? 'Out of Stock' : 'Order Now'}
      </button>
    </div>
  );
};

export default BookDetails; 