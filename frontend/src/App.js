import React, { useState, useEffect } from 'react';
import './App.css';
import { AuthProvider, useAuth } from './context/AuthContext';
import { CartProvider } from './context/CartContext';
import Login from './components/Login';
import Register from './components/Register';
import Sidebar from './components/Sidebar';
import Dashboard from './pages/Dashboard';
import Users from './pages/Users';
import Products from './pages/Products';
import Orders from './pages/Orders';
import Cart from './pages/Cart';

function AppContent() {
  const { user, logout } = useAuth();
  const [showRegister, setShowRegister] = useState(false);

  // Default sayfayı role göre belirle
  const getDefaultPage = () => {
    if (user?.role === 'YONETICI') return 'dashboard';
    if (user?.role === 'SATICI') return 'products';
    if (user?.role === 'KULLANICI') return 'products';
    return 'products';
  };

  const [currentPage, setCurrentPage] = useState(getDefaultPage());

  // User değiştiğinde default sayfaya yönlendir
  useEffect(() => {
    if (user) {
      setCurrentPage(getDefaultPage());
    }
  }, [user]);

  if (!user) {
    return showRegister ? (
      <Register onSwitchToLogin={() => setShowRegister(false)} />
    ) : (
      <Login onSwitchToRegister={() => setShowRegister(true)} />
    );
  }

  const isAdmin = user?.role === 'YONETICI';
  const isSeller = user?.role === 'SATICI';
  const isCustomer = user?.role === 'KULLANICI';

  const renderPage = () => {
    // Role-based page access control
    switch (currentPage) {
      case 'dashboard':
        return isAdmin ? <Dashboard /> : <Products />;
      case 'users':
        return isAdmin ? <Users /> : <Products />;
      case 'products':
        return <Products />;
      case 'cart':
        return isCustomer ? <Cart /> : <Products />;
      case 'orders':
        return (isAdmin || isCustomer) ? <Orders /> : <Products />;
      default:
        return <Products />;
    }
  };

  const getRoleLabel = () => {
    if (isAdmin) return 'Yönetici';
    if (isSeller) return 'Satıcı';
    if (isCustomer) return 'Müşteri';
    return '';
  };

  return (
    <div className="App">
      <div className="app-layout">
        <Sidebar currentPage={currentPage} onNavigate={setCurrentPage} />
        <div className="main-content">
          <header className="App-header">
            <h1>Atlas Dashboard</h1>
            <div className="user-info">
              <span className="user-role">{getRoleLabel()}</span>
              <span>Hoş geldin, {user.email}</span>
              <button onClick={logout} className="btn-logout">
                Çıkış Yap
              </button>
            </div>
          </header>
          <div className="content-area">
            {renderPage()}
          </div>
        </div>
      </div>
    </div>
  );
}

function App() {
  return (
    <AuthProvider>
      <CartProvider>
        <AppContent />
      </CartProvider>
    </AuthProvider>
  );
}

export default App;

