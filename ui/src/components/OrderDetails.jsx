import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useParams, useNavigate } from 'react-router-dom';

const OrderDetails = () => {
  const [order, setOrder] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const { id } = useParams();
  const navigate = useNavigate();

  useEffect(() => {
    const fetchOrderDetails = async () => {
      try {
        const token = localStorage.getItem('token');
        if (!token) {
          navigate('/login');
          return;
        }

        const response = await axios.get(`http://localhost:8080/api/orders/${id}`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        setOrder(response.data);
        setLoading(false);
      } catch (err) {
        setError('Failed to fetch order details');
        setLoading(false);
      }
    };

    fetchOrderDetails();
  }, [id, navigate]);

  if (loading) return <div className="loading">Loading...</div>;
  if (error) return <div className="error">{error}</div>;
  if (!order) return <div>Order not found</div>;

  return (
    <div className="order-details">
      <h2>Order #{order.id}</h2>
      <div className="order-info">
        <p><strong>Date:</strong> {new Date(order.createdAt).toLocaleString()}</p>
        <p><strong>Status:</strong> {order.status}</p>
        <p><strong>Total Amount:</strong> ${order.totalAmount}</p>
        
        <h3>Order Items</h3>
        <div className="order-items">
          {order.items.map((item) => (
            <div key={item.id} className="order-item">
              <h4>{item.book.title}</h4>
              <p>Author: {item.book.author}</p>
              <p>Quantity: {item.quantity}</p>
              <p>Price: ${item.price}</p>
            </div>
          ))}
        </div>

        <h3>Shipping Information</h3>
        <div className="shipping-info">
          <p><strong>Address:</strong> {order.shippingAddress}</p>
          <p><strong>Phone:</strong> {order.phone}</p>
        </div>
      </div>
    </div>
  );
};

export default OrderDetails; 