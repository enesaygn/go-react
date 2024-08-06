import React, { useEffect, useState } from 'react';
import axios from 'axios';

function UserComponent() {
  const [users, setUsers] = useState([]);
  const [count, setCount] = useState(0);
  useEffect(() => {
    console.log('Fetching users');
  }, [count]);

  return (
    <div>
      <h1>Users</h1>
      <button onClick={() => setCount(count + 1)}>Fetch Users</button>
    </div>
  );
}

export default UserComponent;
