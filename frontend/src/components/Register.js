import React, { useState } from 'react';
import { useAuth } from '../context/AuthContext';

function Register({ onSwitchToLogin }) {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    ad: '',
    soyad: ''
  });
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const { register, loading } = useAuth();

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setSuccess(false);

    const result = await register(formData);
    if (result.success) {
      setSuccess(true);
      setTimeout(() => {
        onSwitchToLogin();
      }, 2000);
    } else {
      setError(result.error);
    }
  };

  return (
    <div className="auth-container">
      <div className="auth-card">
        <h2 className="auth-title">Kayıt Ol</h2>
        <form onSubmit={handleSubmit} className="auth-form">
          <div className="form-group">
            <label htmlFor="ad">Ad</label>
            <input
              id="ad"
              name="ad"
              type="text"
              value={formData.ad}
              onChange={handleChange}
              placeholder="Adınız"
              required
              className="form-input"
            />
          </div>

          <div className="form-group">
            <label htmlFor="soyad">Soyad</label>
            <input
              id="soyad"
              name="soyad"
              type="text"
              value={formData.soyad}
              onChange={handleChange}
              placeholder="Soyadınız"
              required
              className="form-input"
            />
          </div>

          <div className="form-group">
            <label htmlFor="email">Email</label>
            <input
              id="email"
              name="email"
              type="email"
              value={formData.email}
              onChange={handleChange}
              placeholder="ornek@email.com"
              required
              className="form-input"
            />
          </div>

          <div className="form-group">
            <label htmlFor="password">Şifre</label>
            <input
              id="password"
              name="password"
              type="password"
              value={formData.password}
              onChange={handleChange}
              placeholder="••••••••"
              required
              minLength="6"
              className="form-input"
            />
          </div>

          {error && <div className="error-message">{error}</div>}
          {success && <div className="success-message">Kayıt başarılı! Giriş sayfasına yönlendiriliyorsunuz...</div>}

          <button type="submit" disabled={loading} className="btn-primary">
            {loading ? 'Kayıt yapılıyor...' : 'Kayıt Ol'}
          </button>
        </form>

        <div className="auth-footer">
          Zaten hesabınız var mı?{' '}
          <button onClick={onSwitchToLogin} className="link-button">
            Giriş Yap
          </button>
        </div>
      </div>
    </div>
  );
}

export default Register;
