import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { Button, Dialog, DialogTitle, DialogContent, DialogActions, TextField } from '@mui/material';

const BookList = ({ isAdmin = false }) => {
  const [books, setBooks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [openDialog, setOpenDialog] = useState(false);
  const [newBook, setNewBook] = useState({
    title: '',
    author: '',
    price: '',
    quantity: '',
    description: '',
    isbn: ''
  });
  const navigate = useNavigate();

  useEffect(() => {
    const fetchBooks = async () => {
      try {
        const token = localStorage.getItem('token');
        const userRole = localStorage.getItem('userRole');
        
        if (isAdmin && userRole !== 'ADMIN') {
          navigate('/books');
          return;
        }

        const response = await axios.get('http://localhost:8080/api/books', {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        setBooks(response.data);
        setLoading(false);
      } catch (err) {
        setError('Failed to fetch books');
        setLoading(false);
      }
    };

    fetchBooks();
  }, [navigate, isAdmin]);

  const handleBookClick = (bookId) => {
    if (isAdmin) {
      navigate(`/admin/books/${bookId}`);
    } else {
      navigate(`/books/${bookId}`);
    }
  };

  const handleAddBook = async () => {
    try {
      const token = localStorage.getItem('token');
      const userRole = localStorage.getItem('userRole');
      
      if (userRole !== 'ADMIN') {
        alert('Only administrators can add books');
        return;
      }

      await axios.post('http://localhost:8080/api/books', newBook, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      setOpenDialog(false);
      const response = await axios.get('http://localhost:8080/api/books', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      setBooks(response.data);
    } catch (err) {
      alert('Failed to add book');
    }
  };

  const handleDeleteBook = async (bookId) => {
    try {
      const token = localStorage.getItem('token');
      const userRole = localStorage.getItem('userRole');
      
      if (userRole !== 'ADMIN') {
        alert('Only administrators can delete books');
        return;
      }

      if (window.confirm('Are you sure you want to delete this book?')) {
        await axios.delete(`http://localhost:8080/api/books/${bookId}`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        const response = await axios.get('http://localhost:8080/api/books', {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        setBooks(response.data);
      }
    } catch (err) {
      alert('Failed to delete book');
    }
  };

  if (loading) return <div className="loading">Loading...</div>;
  if (error) return <div className="error">{error}</div>;

  return (
    <div className={`book-list ${isAdmin ? 'admin-view' : ''}`}>
      <div className="book-list-header">
        <h2>{isAdmin ? 'Manage Books' : 'Available Books'}</h2>
        {isAdmin && (
          <Button 
            variant="contained" 
            color="primary" 
            onClick={() => setOpenDialog(true)}
          >
            Add New Book
          </Button>
        )}
      </div>
      <div className="books-grid">
        {books.map((book) => (
          <div
            key={book.id}
            className="book-card"
            onClick={() => handleBookClick(book.id)}
          >
            <h3>{book.title}</h3>
            <p>Author: {book.author}</p>
            <p>Price: ${book.price}</p>
            <p>Available: {book.quantity}</p>
            {isAdmin && (
              <div className="book-actions">
                <Button
                  variant="contained"
                  color="error"
                  onClick={(e) => {
                    e.stopPropagation();
                    handleDeleteBook(book.id);
                  }}
                >
                  Delete
                </Button>
              </div>
            )}
          </div>
        ))}
      </div>

      <Dialog open={openDialog} onClose={() => setOpenDialog(false)}>
        <DialogTitle>Add New Book</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Title"
            fullWidth
            value={newBook.title}
            onChange={(e) => setNewBook({ ...newBook, title: e.target.value })}
          />
          <TextField
            margin="dense"
            label="Author"
            fullWidth
            value={newBook.author}
            onChange={(e) => setNewBook({ ...newBook, author: e.target.value })}
          />
          <TextField
            margin="dense"
            label="Price"
            type="number"
            fullWidth
            value={newBook.price}
            onChange={(e) => setNewBook({ ...newBook, price: e.target.value })}
          />
          <TextField
            margin="dense"
            label="Quantity"
            type="number"
            fullWidth
            value={newBook.quantity}
            onChange={(e) => setNewBook({ ...newBook, quantity: e.target.value })}
          />
          <TextField
            margin="dense"
            label="Description"
            fullWidth
            multiline
            rows={4}
            value={newBook.description}
            onChange={(e) => setNewBook({ ...newBook, description: e.target.value })}
          />
          <TextField
            margin="dense"
            label="ISBN"
            fullWidth
            value={newBook.isbn}
            onChange={(e) => setNewBook({ ...newBook, isbn: e.target.value })}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenDialog(false)}>Cancel</Button>
          <Button onClick={handleAddBook} color="primary">
            Add Book
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};

export default BookList; 