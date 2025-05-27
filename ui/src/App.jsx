import { BrowserRouter as Router, Routes, Route, Navigate, Link } from 'react-router-dom';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import { AppBar, Toolbar, Typography, Button, Container, Box } from '@mui/material';
import CssBaseline from '@mui/material/CssBaseline';
import Login from './components/Login';
import Register from './components/Register';
import BookList from './components/BookList';
import BookDetails from './components/BookDetails';
import OrderList from './components/OrderList';
import OrderDetails from './components/OrderDetails';
import './App.css';

const theme = createTheme({
  palette: {
    primary: {
      main: '#1976d2',
    },
    secondary: {
      main: '#dc004e',
    },
    background: {
      default: '#f5f5f5',
    },
  },
  typography: {
    fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
    h1: {
      fontSize: '2.5rem',
      fontWeight: 500,
    },
    h2: {
      fontSize: '2rem',
      fontWeight: 500,
    },
    h3: {
      fontSize: '1.75rem',
      fontWeight: 500,
    },
    h4: {
      fontSize: '1.5rem',
      fontWeight: 500,
    },
    h5: {
      fontSize: '1.25rem',
      fontWeight: 500,
    },
    h6: {
      fontSize: '1.1rem',
      fontWeight: 500,
    },
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          textTransform: 'none',
          fontWeight: 500,
          padding: '8px 16px',
        },
      },
    },
    MuiTextField: {
      styleOverrides: {
        root: {
          '& .MuiOutlinedInput-root': {
            borderRadius: 8,
          },
        },
      },
    },
    MuiAppBar: {
      styleOverrides: {
        root: {
          boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
        },
      },
    },
  },
});

function App() {
  const token = localStorage.getItem('token');
  const userRole = localStorage.getItem('userRole');
  const isAuthenticated = !!token;
  const isAdmin = userRole === 'ADMIN';

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('userRole');
    window.location.href = '/login';
  };

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Router>
        <div className="app">
          <AppBar position="static">
            <Toolbar sx={{ maxWidth: 1400, width: '100%', margin: '0 auto', px: 2 }}>
              <Typography variant="h6" component="div" sx={{ flexGrow: 1, fontSize: '1.5rem' }}>
                Book Store
              </Typography>
              <Box sx={{ display: 'flex', gap: 2 }}>
                {isAuthenticated ? (
                  <>
                    <Button color="inherit" component={Link} to="/books" sx={{ fontSize: '1rem' }}>
                      Books
                    </Button>
                    {isAdmin && (
                      <Button color="inherit" component={Link} to="/admin/books" sx={{ fontSize: '1rem' }}>
                        Manage Books
                      </Button>
                    )}
                    <Button color="inherit" component={Link} to="/orders" sx={{ fontSize: '1rem' }}>
                      My Orders
                    </Button>
                    {isAdmin && (
                      <Button color="inherit" component={Link} to="/admin/orders" sx={{ fontSize: '1rem' }}>
                        All Orders
                      </Button>
                    )}
                    <Button color="inherit" onClick={handleLogout} sx={{ fontSize: '1rem' }}>
                      Logout
                    </Button>
                  </>
                ) : (
                  <>
                    <Button color="inherit" component={Link} to="/login" sx={{ fontSize: '1rem' }}>
                      Login
                    </Button>
                    <Button color="inherit" component={Link} to="/register" sx={{ fontSize: '1rem' }}>
                      Register
                    </Button>
                  </>
                )}
              </Box>
            </Toolbar>
          </AppBar>
          <Container className="main-content" maxWidth={false}>
            <Routes>
              <Route path="/login" element={<Login />} />
              <Route path="/register" element={<Register />} />
              <Route 
                path="/books" 
                element={isAuthenticated ? <BookList /> : <Navigate to="/login" />} 
              />
              <Route 
                path="/books/:id" 
                element={isAuthenticated ? <BookDetails /> : <Navigate to="/login" />} 
              />
              <Route 
                path="/orders" 
                element={isAuthenticated ? <OrderList /> : <Navigate to="/login" />} 
              />
              <Route 
                path="/orders/:id" 
                element={isAuthenticated ? <OrderDetails /> : <Navigate to="/login" />} 
              />
              <Route 
                path="/admin/books" 
                element={isAdmin ? <BookList isAdmin={true} /> : <Navigate to="/books" />} 
              />
              <Route 
                path="/admin/orders" 
                element={isAdmin ? <OrderList isAdmin={true} /> : <Navigate to="/orders" />} 
              />
              <Route 
                path="/" 
                element={
                  isAuthenticated 
                    ? (isAdmin ? <Navigate to="/admin/books" /> : <Navigate to="/books" />)
                    : <Navigate to="/login" />
                } 
              />
            </Routes>
          </Container>
        </div>
      </Router>
    </ThemeProvider>
  );
}

export default App;
