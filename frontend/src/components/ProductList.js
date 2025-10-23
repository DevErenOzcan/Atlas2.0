import React, { useEffect, useState } from 'react';
import { productAPI } from '../services/api';
import { useAuth } from '../context/AuthContext';
import { useCart } from '../context/CartContext';

const ProductList = () => {
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [showModal, setShowModal] = useState(false);
    const [editingProduct, setEditingProduct] = useState(null);
    const { user } = useAuth();
    const { addToCart } = useCart();

    const [formData, setFormData] = useState({
        product_name: '',
        description: '',
        stock: 0,
        price: 0,
        currency: 'TRY',
        seller_id: user?.user_id || 1,
        category_id: 1,
        dimens_details: ''
    });

    const isAdmin = user && user.role === 'YONETICI';
    const isSeller = user && user.role === 'SATICI';
    const isCustomer = user && user.role === 'KULLANICI';
    const canManageProducts = isAdmin || isSeller;

    const fetchProducts = async () => {
        try {
            const response = canManageProducts
                ? await productAPI.getMy()
                : await productAPI.getAll();
            setProducts(Array.isArray(response.data) ? response.data : []);
        } catch (error) {
            setProducts([]);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (user) {
            fetchProducts();
        }
    }, [user?.user_id, user?.role]);

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        setFormData({
            ...formData,
            [name]: name === 'stock' || name === 'price' || name === 'seller_id' || name === 'category_id'
                ? Number(value)
                : value
        });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const productData = {
                ...formData,
                seller_id: user.user_id
            };

            if (editingProduct) {
                await productAPI.update(editingProduct.product_id, productData);
            } else {
                await productAPI.create(productData);
            }
            setShowModal(false);
            setEditingProduct(null);
            resetForm();
            fetchProducts();
        } catch (error) {
            alert(error.response?.data || 'İşlem başarısız');
        }
    };

    const handleEdit = (product) => {
        setEditingProduct(product);
        setFormData({
            product_name: product.product_name,
            description: product.description,
            stock: product.stock,
            price: product.price,
            currency: product.currency,
            seller_id: product.seller_id,
            category_id: product.category_id,
            dimens_details: product.dimens_details
        });
        setShowModal(true);
    };

    const handleDelete = async (id) => {
        if (window.confirm('Bu ürünü silmek istediğinizden emin misiniz?')) {
            try {
                await productAPI.delete(id);
                fetchProducts();
            } catch (error) {
                alert(error.response?.data || 'Silme başarısız');
            }
        }
    };

    const resetForm = () => {
        setFormData({
            product_name: '',
            description: '',
            stock: 0,
            price: 0,
            currency: 'TRY',
            seller_id: user?.user_id || 1,
            category_id: 1,
            dimens_details: ''
        });
    };

    const openAddModal = () => {
        setEditingProduct(null);
        resetForm();
        setShowModal(true);
    };

    if (loading) return <div className="list-container">Yükleniyor...</div>;

    return (
        <div className="list-container">
            <div className="list-header">
                <h2>Ürünler</h2>
                {canManageProducts && (
                    <button onClick={openAddModal} className="btn-add">
                        + Yeni Ürün
                    </button>
                )}
            </div>

            <div className="product-grid">
                {products.map(product => (
                    <div key={product.product_id} className="product-card">
                        <h3>{product.product_name}</h3>
                        <p className="product-description">{product.description}</p>
                        <div className="product-info">
                            <span className="product-price">{product.price} {product.currency}</span>
                            <span className="product-stock">Stok: {product.stock}</span>
                        </div>
                        {canManageProducts && (
                            <div className="product-actions">
                                <button onClick={() => handleEdit(product)} className="btn-edit">
                                    Düzenle
                                </button>
                                <button onClick={() => handleDelete(product.product_id)} className="btn-delete">
                                    Sil
                                </button>
                            </div>
                        )}
                        {isCustomer && product.stock > 0 && (
                            <div className="product-actions">
                                <button
                                    onClick={() => {
                                        addToCart(product, 1);
                                        alert(`${product.product_name} sepete eklendi!`);
                                    }}
                                    className="btn-add-cart"
                                >
                                    Sepete Ekle
                                </button>
                            </div>
                        )}
                    </div>
                ))}
            </div>

            {showModal && (
                <div className="modal-overlay" onClick={() => setShowModal(false)}>
                    <div className="modal-content" onClick={(e) => e.stopPropagation()}>
                        <h2>{editingProduct ? 'Ürün Düzenle' : 'Yeni Ürün Ekle'}</h2>
                        <form onSubmit={handleSubmit} className="product-form">
                            <div className="form-group">
                                <label>Ürün Adı</label>
                                <input
                                    type="text"
                                    name="product_name"
                                    value={formData.product_name}
                                    onChange={handleInputChange}
                                    required
                                    className="form-input"
                                />
                            </div>
                            <div className="form-group">
                                <label>Açıklama</label>
                                <textarea
                                    name="description"
                                    value={formData.description}
                                    onChange={handleInputChange}
                                    required
                                    className="form-input"
                                    rows="3"
                                />
                            </div>
                            <div className="form-row">
                                <div className="form-group">
                                    <label>Stok</label>
                                    <input
                                        type="number"
                                        name="stock"
                                        value={formData.stock}
                                        onChange={handleInputChange}
                                        required
                                        className="form-input"
                                    />
                                </div>
                                <div className="form-group">
                                    <label>Fiyat</label>
                                    <input
                                        type="number"
                                        step="0.01"
                                        name="price"
                                        value={formData.price}
                                        onChange={handleInputChange}
                                        required
                                        className="form-input"
                                    />
                                </div>
                            </div>
                            <div className="form-group">
                                <label>Para Birimi</label>
                                <select
                                    name="currency"
                                    value={formData.currency}
                                    onChange={handleInputChange}
                                    className="form-input"
                                >
                                    <option value="TRY">TRY</option>
                                    <option value="USD">USD</option>
                                    <option value="EUR">EUR</option>
                                </select>
                            </div>
                            <div className="form-group">
                                <label>Boyutlar</label>
                                <input
                                    type="text"
                                    name="dimens_details"
                                    value={formData.dimens_details}
                                    onChange={handleInputChange}
                                    className="form-input"
                                    placeholder="Örn: 10x10x10 cm"
                                />
                            </div>
                            <div className="modal-actions">
                                <button type="submit" className="btn-primary">
                                    {editingProduct ? 'Güncelle' : 'Ekle'}
                                </button>
                                <button type="button" onClick={() => setShowModal(false)} className="btn-secondary">
                                    İptal
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default ProductList;
