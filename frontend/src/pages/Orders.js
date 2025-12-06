import React from 'react';
import OrderList from '../components/OrderList';

function Orders() {
  return (
    <div className="page-container">
      <h1 className="page-title">Siparişler</h1>
      <div className="page-content">
        <OrderList />
      </div>
    </div>
  );
}

export default Orders;
