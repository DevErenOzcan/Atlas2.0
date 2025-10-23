import React from 'react';
import { useCart } from '../context/CartContext';
import { orderAPI } from '../services/api';

function Cart() {
  const { cartItems, removeFromCart, updateQuantity, clearCart, getCartTotal } = useCart();

  const handleCheckout = async () => {
    if (cartItems.length === 0) {
      alert('Sepetiniz boş!');
      return;
    }

    try {
      const orderData = {
        items: cartItems.map(item => ({
          product_id: item.product_id,
          quantity: item.quantity,
          price: item.price
        }))
      };

      const response = await orderAPI.create(orderData);
      alert(`Sipariş başarıyla oluşturuldu! Sipariş No: ${response.data.order_id}`);
      clearCart();
    } catch (error) {
      alert('Sipariş oluşturulurken bir hata oluştu: ' + (error.response?.data || error.message));
    }
  };

  if (cartItems.length === 0) {
    return (
      <div className="page-container">
        <h1 className="page-title">Sepetim</h1>
        <div className="empty-cart">
          <p>Sepetinizde ürün bulunmuyor.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="page-container">
      <h1 className="page-title">Sepetim</h1>
      <div className="cart-container">
        <div className="cart-items">
          {cartItems.map(item => (
            <div key={item.product_id} className="cart-item">
              <div className="cart-item-info">
                <h3>{item.product_name}</h3>
                <p className="cart-item-price">{item.price} {item.currency}</p>
              </div>
              <div className="cart-item-actions">
                <div className="quantity-controls">
                  <button
                    onClick={() => updateQuantity(item.product_id, item.quantity - 1)}
                    className="quantity-btn"
                  >
                    -
                  </button>
                  <span className="quantity">{item.quantity}</span>
                  <button
                    onClick={() => updateQuantity(item.product_id, item.quantity + 1)}
                    className="quantity-btn"
                  >
                    +
                  </button>
                </div>
                <button
                  onClick={() => removeFromCart(item.product_id)}
                  className="btn-remove"
                >
                  Kaldır
                </button>
              </div>
              <div className="cart-item-total">
                {(item.price * item.quantity).toFixed(2)} {item.currency}
              </div>
            </div>
          ))}
        </div>

        <div className="cart-summary">
          <h3>Sipariş Özeti</h3>
          <div className="summary-row">
            <span>Toplam:</span>
            <span className="summary-total">{getCartTotal().toFixed(2)} TL</span>
          </div>
          <button onClick={handleCheckout} className="btn-checkout">
            Siparişi Tamamla
          </button>
          <button onClick={clearCart} className="btn-clear-cart">
            Sepeti Temizle
          </button>
        </div>
      </div>
    </div>
  );
}

export default Cart;
