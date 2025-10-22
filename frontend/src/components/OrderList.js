import React, { useEffect, useState } from 'react';
import api from '../services/api';

const OrderList = () => {
    const [orders, setOrders] = useState([]);

    useEffect(() => {
        api.get('/orders')
            .then(response => {
                const data = response && response.data;
                setOrders(Array.isArray(data) ? data : []);
            })
            .catch(error => {
                console.error("There was an error fetching the orders!", error);
                setOrders([]);
            });
    }, []);

    return (
        <div className="list-container">
            <h2>Orders</h2>
            <ul>
                {orders.map(order => (
                    <li key={order.order_id}>
                        <strong>Order ID:</strong> {order.order_id}<br/>
                        <strong>User ID:</strong> {order.user_id}<br/>
                        <strong>Total:</strong> {order.total}<br/>
                        <strong>Final:</strong> {order.final}<br/>
                        <strong>Shipped:</strong> {order.is_shipped ? 'Yes' : 'No'}<br/>
                        <strong>Date:</strong> {new Date(order.create_date).toLocaleString()}<br/>
                        <strong>Details:</strong>
                        <ul>
                            {Array.isArray(order.order_details) && order.order_details.map(detail => (
                                <li key={detail.detail_id}>
                                    - Product ID: {detail.product_id}, Quantity: {detail.item}, Price: {detail.final}
                                </li>
                            ))}
                        </ul>
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default OrderList;
