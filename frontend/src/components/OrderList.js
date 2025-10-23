import React, { useEffect, useState } from 'react';
import { orderAPI } from '../services/api';
import { useAuth } from '../context/AuthContext';

const OrderList = () => {
    const [orders, setOrders] = useState([]);
    const [loading, setLoading] = useState(true);
    const { user } = useAuth();

    const isAdmin = user && user.role === 'YONETICI';

    const fetchOrders = async () => {
        try {
            const response = await orderAPI.getMy();
            setOrders(Array.isArray(response.data) ? response.data : []);
        } catch (error) {
            setOrders([]);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (user) {
            fetchOrders();
        }
    }, [user?.user_id, user?.role]);

    const handleStatusChange = async (orderId, isShipped) => {
        try {
            await orderAPI.updateStatus(orderId, isShipped);
            fetchOrders();
        } catch (error) {
            alert(error.response?.data || 'Durum güncellenemedi');
        }
    };

    if (loading) return <div className="list-container">Yükleniyor...</div>;

    return (
        <div className="list-container">
            <div className="list-header">
                <h2>Siparişler</h2>
            </div>

            <table className="data-table">
                <thead>
                    <tr>
                        <th>Sipariş No</th>
                        <th>Kullanıcı ID</th>
                        <th>Toplam</th>
                        <th>Final</th>
                        <th>Durum</th>
                        <th>Tarih</th>
                        {isAdmin && <th>İşlemler</th>}
                    </tr>
                </thead>
                <tbody>
                    {orders.map(order => (
                        <tr key={order.order_id}>
                            <td>{order.order_id}</td>
                            <td>{order.user_id}</td>
                            <td>{order.total?.toFixed(2)} TL</td>
                            <td>{order.final?.toFixed(2)} TL</td>
                            <td>
                                <span className={order.is_shipped ? 'status-badge shipped' : 'status-badge pending'}>
                                    {order.is_shipped ? 'Gönderildi' : 'Bekliyor'}
                                </span>
                            </td>
                            <td>{new Date(order.create_date).toLocaleDateString('tr-TR')}</td>
                            {isAdmin && (
                                <td>
                                    {!order.is_shipped ? (
                                        <button
                                            onClick={() => handleStatusChange(order.order_id, true)}
                                            className="btn-approve"
                                        >
                                            Onayla
                                        </button>
                                    ) : (
                                        <button
                                            onClick={() => handleStatusChange(order.order_id, false)}
                                            className="btn-cancel"
                                        >
                                            İptal Et
                                        </button>
                                    )}
                                </td>
                            )}
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
};

export default OrderList;
