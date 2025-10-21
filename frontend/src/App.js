import React from 'react';
import './App.css';
import UserList from './components/UserList';
import ProductList from './components/ProductList';
import OrderList from './components/OrderList';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>Gateway Test</h1>
      </header>
      <div className="container">
        <UserList />
        <ProductList />
        <OrderList />
      </div>
    </div>
  );
}

export default App;

