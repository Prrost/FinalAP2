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
    availableQuantity: '',
    totalQuantity: '',
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

        if (!token) {
          navigate('/login');
          return;
        }

        const response = await axios.get('http://localhost:8080/api/books', {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        
        // Ensure we always set an array, even if the response is null
        setBooks(response.data || []);
        setLoading(false);
      } catch (err) {
        setError('Failed to fetch books');
        setBooks([]); // Reset books to empty array on error
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
      console.log('Adding book', newBook);
      const token = localStorage.getItem('token');
      const userRole = localStorage.getItem('userRole');



      await axios.post('http://localhost:8080/api/books', {
        ...newBook,
        availableQuantity: parseInt(newBook.availableQuantity, 10),
        totalQuantity: parseInt(newBook.totalQuantity, 10)
      }, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      console.log('222222222', newBook);
      setOpenDialog(false);
      const response = await axios.get('http://localhost:8080/api/books', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      setBooks(response.data);
    } catch (err) {
      console.log('фввывывыфвфывфывыфвфвывфвфывфывввфвывф', newBook);
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
        {true && (
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
        {Array.isArray(books) && books.length > 0 ? (
          books.map((book) => (
            <div
              key={book.id}
              className="book-card"
              onClick={() => handleBookClick(book.id)}
            >
              <h3>{book.title}</h3>
              <p>Author: {book.author}</p>
              <p>Available: {book.availableQuantity}</p>
              {true && (
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
          ))
        ) : (
          <p>No books available.</p>
        )}
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
              label="Total Quantity"
              type="number"
              fullWidth
              value={newBook.totalQuantity}
              onChange={(e) => setNewBook({ ...newBook, totalQuantity: e.target.value })}
          />

          <TextField
            margin="dense"
            label="Available Quantity"
            type="number"
            fullWidth
            value={newBook.availableQuantity}
            onChange={(e) => setNewBook({ ...newBook, availableQuantity: e.target.value })}
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