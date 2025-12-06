import React from 'react';
import { useAuth } from '../context/AuthContext';
import { useCart } from '../context/CartContext';

function Sidebar({ currentPage, onNavigate }) {
  const { user } = useAuth();
  const { getCartCount } = useCart();

  // Role-based menu items
  const getMenuItems = () => {
    const isAdmin = user?.role === 'YONETICI';
    const isSeller = user?.role === 'SATICI';
    const isCustomer = user?.role === 'KULLANICI';

    if (isAdmin) {
      // Admin: Tüm sayfaları görebilir
      return [
        { id: 'dashboard', label: 'Dashboard', icon: '🏠' },
        { id: 'users', label: 'Kullanıcılar', icon: '👥' },
        { id: 'products', label: 'Ürünler', icon: '📦' },
        { id: 'orders', label: 'Siparişler', icon: '🛒' },
      ];
    } else if (isSeller) {
      // Satıcı: Sadece ürünler
      return [
        { id: 'products', label: 'Ürünlerim', icon: '📦' },
      ];
    } else if (isCustomer) {
      // Kullanıcı: Sadece ürünler, sepet ve siparişler
      return [
        { id: 'products', label: 'Ürünler', icon: '📦' },
        { id: 'cart', label: 'Sepetim', icon: '🛒', badge: getCartCount() },
        { id: 'orders', label: 'Siparişlerim', icon: '📋' },
      ];
    }
    return [];
  };

  const menuItems = getMenuItems();

  return (
    <div className="sidebar">
      <div className="sidebar-header">
        <h2>Atlas</h2>
      </div>
      <nav className="sidebar-nav">
        {menuItems.map((item) => (
          <button
            key={item.id}
            className={`nav-item ${currentPage === item.id ? 'active' : ''}`}
            onClick={() => onNavigate(item.id)}
          >
            <span className="nav-icon">{item.icon}</span>
            <span className="nav-label">{item.label}</span>
            {item.badge > 0 && <span className="cart-badge">{item.badge}</span>}
          </button>
        ))}
      </nav>
    </div>
  );
}

export default Sidebar;
