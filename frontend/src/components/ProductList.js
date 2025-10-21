import React, { useEffect, useState } from 'react';
import api from '../services/api';

const ProductList = () => {
    const [products, setProducts] = useState([]);

    useEffect(() => {
        api.get('/products')
            .then(response => {
                setProducts(response.data);
            })
            .catch(error => {
                console.error("There was an error fetching the products!", error);
            });
    }, []);

    return (
        <div className="list-container">
            <h2>Products</h2>
            <ul>
                {products.map(product => (
                    <li key={product.product_id}>
                        <strong>ID:</strong> {product.product_id}<br/>
                        <strong>Name:</strong> {product.product_name}<br/>
                        <strong>Description:</strong> {product.description}<br/>
                        <strong>Stock:</strong> {product.stock}<br/>
                        <strong>Price:</strong> {product.price} {product.currency}<br/>
                        <strong>Seller ID:</strong> {product.seller_id}<br/>
                        <strong>Category ID:</strong> {product.category_id}<br/>
                        <strong>Dimensions:</strong> {product.dimens_details}
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default ProductList;
