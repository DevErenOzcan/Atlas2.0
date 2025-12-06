import React from 'react';
import UserList from '../components/UserList';

function Users() {
  return (
    <div className="page-container">
      <h1 className="page-title">Kullanıcılar</h1>
      <div className="page-content">
        <UserList />
      </div>
    </div>
  );
}

export default Users;
