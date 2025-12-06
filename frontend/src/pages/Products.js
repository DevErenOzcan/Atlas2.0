import React from 'react';
import ProductList from '../components/ProductList';

function Products() {
  return (
    <div className="page-container">
      <h1 className="page-title">Ürünler</h1>
      <div className="page-content">
        <ProductList />
      </div>
    </div>
  );
}

export default Products;
