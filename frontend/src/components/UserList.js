import React, { useEffect, useState } from 'react';
import api from '../services/api';

const UserList = () => {
    const [users, setUsers] = useState([]);

    useEffect(() => {
        api.get('/users')
            .then(response => {
                const data = response && response.data;
                setUsers(Array.isArray(data) ? data : []);
            })
            .catch(error => {
                console.error("There was an error fetching the userr!", error);
                setUsers([]);
            });
    }, []);

    return (
        <div className="list-container">
            <h2>Users</h2>
            <ul>
                {users.map(user => (
                    <li key={user.user_id}>
                        <strong>ID:</strong> {user.user_id}<br/>
                        <strong>Name:</strong> {user.ad} {user.soyad}<br/>
                        <strong>Email:</strong> {user.email}<br/>
                        <strong>Role:</strong> {user.role}<br/>
                        <strong>Address ID:</strong> {user.address_id}
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default UserList;
