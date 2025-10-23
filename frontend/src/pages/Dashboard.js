import React from 'react';
import { useAuth } from '../context/AuthContext';

function Dashboard() {
  const { user } = useAuth();

  const roleText = {
    0: 'Kullanıcı',
    1: 'Satıcı',
    2: 'Yönetici'
  };

  return (
    <div className="page-container">
      <h1 className="page-title">Dashboard</h1>
      <div className="dashboard-grid">
        <div className="dashboard-card">
          <div className="card-icon">👤</div>
          <h3>Hoş Geldiniz</h3>
          <p className="user-email">{user.email}</p>
          <span className="user-role">{roleText[user.role]}</span>
        </div>

        <div className="dashboard-card">
          <div className="card-icon">📊</div>
          <h3>İstatistikler</h3>
          <p>Sisteminizi yönetmek için yan menüyü kullanın</p>
        </div>

        <div className="dashboard-card">
          <div className="card-icon">🚀</div>
          <h3>Hızlı Başlangıç</h3>
          <ul className="quick-links">
            <li>• Kullanıcıları Görüntüle</li>
            <li>• Ürünleri Yönet</li>
            <li>• Siparişleri İncele</li>
          </ul>
        </div>
      </div>
    </div>
  );
}

export default Dashboard;
