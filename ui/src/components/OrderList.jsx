import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { Button, Select, MenuItem, FormControl, InputLabel } from '@mui/material';

const OrderList = ({ isAdmin = false }) => {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchOrders = async () => {
      try {
        const token = localStorage.getItem('token');
        const userRole = localStorage.getItem('userRole');
        
        if (isAdmin && userRole !== 'ADMIN') {
          navigate('/orders');
          return;
        }

        if (!token) {
          navigate('/login');
          return;
        }

        const endpoint = isAdmin 
          ? 'http://localhost:8080/admin/orders'
          : 'http://localhost:8080/orders';

        const response = await axios.get(endpoint, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        setOrders(response.data);
        setLoading(false);
      } catch (err) {
        setError('Failed to fetch orders');
        setLoading(false);
      }
    };

    fetchOrders();
  }, [navigate, isAdmin]);

  const handleOrderClick = (orderId) => {
    if (isAdmin) {
      navigate(`/admin/orders/${orderId}`);
    } else {
      navigate(`/orders/${orderId}`);
    }
  };

  const handleStatusChange = async (orderId, newStatus) => {
    try {
      const token = localStorage.getItem('token');
      const userRole = localStorage.getItem('userRole');
      
      if (userRole !== 'ADMIN') {
        alert('Only administrators can update order status');
        return;
      }

      await axios.patch(
        `http://localhost:8080/api/admin/orders/${orderId}/status`,
        { status: newStatus },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      const response = await axios.get('http://localhost:8080/api/admin/orders', {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      setOrders(response.data);
    } catch (err) {
      alert('Failed to update order status');
    }
  };

  if (loading) return <div className="loading">Loading...</div>;
  if (error) return <div className="error">{error}</div>;

  return (
    <div className={`order-list ${isAdmin ? 'admin-view' : ''}`}>
      <h2>{isAdmin ? 'All Orders' : 'Your Orders'}</h2>
      {orders.length === 0 ? (
        <p>No orders found.</p>
      ) : (
        <div className="orders-grid">
          {orders.map((order) => (
            <div
              key={order.id}
              className="order-card"
              onClick={() => handleOrderClick(order.id)}
            >
              <h3>Order #{order.id}</h3>
              <p>Date: {new Date(order.createdAt).toLocaleDateString()}</p>
              <p>Status: {order.status}</p>
              <p>Total: ${order.totalAmount}</p>
              {isAdmin && (
                <div className="order-actions">
                  <FormControl fullWidth>
                    <InputLabel>Status</InputLabel>
                    <Select
                      value={order.status}
                      label="Status"
                      onChange={(e) => {
                        e.stopPropagation();
                        handleStatusChange(order.id, e.target.value);
                      }}
                    >
                      <MenuItem value="PENDING">Pending</MenuItem>
                      <MenuItem value="PROCESSING">Processing</MenuItem>
                      <MenuItem value="SHIPPED">Shipped</MenuItem>
                      <MenuItem value="DELIVERED">Delivered</MenuItem>
                      <MenuItem value="CANCELLED">Cancelled</MenuItem>
                    </Select>
                  </FormControl>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default OrderList; 