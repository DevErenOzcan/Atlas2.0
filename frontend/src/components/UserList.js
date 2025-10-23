import React, { useEffect, useState } from 'react';
import { userAPI } from '../services/api';
import { useAuth } from '../context/AuthContext';

const UserList = () => {
    const [users, setUsers] = useState([]);
    const [loading, setLoading] = useState(true);
    const [showModal, setShowModal] = useState(false);
    const [editingUser, setEditingUser] = useState(null);
    const [formData, setFormData] = useState({
        email: '',
        password: '',
        ad: '',
        soyad: '',
        role: 'KULLANICI',
        address_id: null
    });
    const { user } = useAuth();

    const isAdmin = user && user.role === 'YONETICI';

    const fetchUsers = async () => {
        try {
            const response = await userAPI.getAll();
            setUsers(response.data);
        } catch (error) {
            setUsers([]);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchUsers();
    }, []);

    const getRoleName = (role) => {
        switch (role) {
            case 'KULLANICI': return 'Kullanıcı';
            case 'SATICI': return 'Satıcı';
            case 'YONETICI': return 'Yönetici';
            default: return 'Bilinmiyor';
        }
    };

    const getRoleClass = (role) => {
        switch (role) {
            case 'KULLANICI': return 'role-badge user';
            case 'SATICI': return 'role-badge seller';
            case 'YONETICI': return 'role-badge admin';
            default: return 'role-badge';
        }
    };

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        setFormData({
            ...formData,
            [name]: name === 'address_id'
                ? (value === '' ? null : Number(value))
                : value
        });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            if (editingUser) {
                const updateData = { ...formData };
                if (!updateData.password) {
                    delete updateData.password;
                }
                await userAPI.update(editingUser.user_id, updateData);
            }
            setShowModal(false);
            setEditingUser(null);
            resetForm();
            fetchUsers();
        } catch (error) {
            alert(error.response?.data || 'İşlem başarısız');
        }
    };

    const handleEdit = (user) => {
        setEditingUser(user);
        setFormData({
            email: user.email,
            password: '',
            ad: user.ad,
            soyad: user.soyad,
            role: user.role,
            address_id: user.address_id || null
        });
        setShowModal(true);
    };

    const handleDelete = async (id) => {
        if (window.confirm('Bu kullanıcıyı silmek istediğinizden emin misiniz?')) {
            try {
                await userAPI.delete(id);
                fetchUsers();
            } catch (error) {
                alert(error.response?.data || 'Silme başarısız');
            }
        }
    };

    const resetForm = () => {
        setFormData({
            email: '',
            password: '',
            ad: '',
            soyad: '',
            role: 'KULLANICI',
            address_id: null
        });
    };

    if (loading) return <div className="list-container">Yükleniyor...</div>;

    return (
        <div className="list-container">
            <div className="list-header">
                <h2>Kullanıcılar</h2>
            </div>

            <table className="data-table">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Ad Soyad</th>
                        <th>Email</th>
                        <th>Rol</th>
                        {isAdmin && <th>İşlemler</th>}
                    </tr>
                </thead>
                <tbody>
                    {users.map(user => (
                        <tr key={user.user_id}>
                            <td>{user.user_id}</td>
                            <td>{user.ad} {user.soyad}</td>
                            <td>{user.email}</td>
                            <td>
                                <span className={getRoleClass(user.role)}>
                                    {getRoleName(user.role)}
                                </span>
                            </td>
                            {isAdmin && (
                                <td>
                                    <button onClick={() => handleEdit(user)} className="btn-edit">
                                        Düzenle
                                    </button>
                                    <button onClick={() => handleDelete(user.user_id)} className="btn-delete">
                                        Sil
                                    </button>
                                </td>
                            )}
                        </tr>
                    ))}
                </tbody>
            </table>

            {showModal && (
                <div className="modal-overlay" onClick={() => setShowModal(false)}>
                    <div className="modal-content" onClick={(e) => e.stopPropagation()}>
                        <h2>Kullanıcı Düzenle</h2>
                        <form onSubmit={handleSubmit} className="product-form">
                            <div className="form-row">
                                <div className="form-group">
                                    <label>Ad</label>
                                    <input
                                        type="text"
                                        name="ad"
                                        value={formData.ad}
                                        onChange={handleInputChange}
                                        required
                                        className="form-input"
                                    />
                                </div>
                                <div className="form-group">
                                    <label>Soyad</label>
                                    <input
                                        type="text"
                                        name="soyad"
                                        value={formData.soyad}
                                        onChange={handleInputChange}
                                        required
                                        className="form-input"
                                    />
                                </div>
                            </div>
                            <div className="form-group">
                                <label>Email</label>
                                <input
                                    type="email"
                                    name="email"
                                    value={formData.email}
                                    onChange={handleInputChange}
                                    required
                                    className="form-input"
                                />
                            </div>
                            <div className="form-group">
                                <label>Şifre (Boş bırakılırsa değişmez)</label>
                                <input
                                    type="password"
                                    name="password"
                                    value={formData.password}
                                    onChange={handleInputChange}
                                    className="form-input"
                                    placeholder="Değiştirmek için yeni şifre girin"
                                />
                            </div>
                            <div className="form-group">
                                <label>Rol</label>
                                <select
                                    name="role"
                                    value={formData.role}
                                    onChange={handleInputChange}
                                    className="form-input"
                                >
                                    <option value="KULLANICI">Kullanıcı</option>
                                    <option value="SATICI">Satıcı</option>
                                    <option value="YONETICI">Yönetici</option>
                                </select>
                            </div>
                            <div className="modal-actions">
                                <button type="submit" className="btn-primary">
                                    Güncelle
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

export default UserList;
