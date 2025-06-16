import React, { useState, useEffect } from 'react';
import './App.css';

// Add a base URL for API calls
const API_BASE_URL = 'http://localhost:8080';

function App() {
  const [passwords, setPasswords] = useState([]);
  const [newPassword, setNewPassword] = useState({ 
    title: '', 
    username: '', 
    password: '', 
    notes: '' 
  });
  const [error, setError] = useState('');
  const [visiblePasswords, setVisiblePasswords] = useState({});
  const [showMasterPasswordForm, setShowMasterPasswordForm] = useState(false);
  const [masterPasswordForm, setMasterPasswordForm] = useState({
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  });
  const [isInitialized, setIsInitialized] = useState(false);
  const [initPassword, setInitPassword] = useState({
    password: '',
    confirmPassword: ''
  });

  useEffect(() => {
    checkAuthStatus();
  }, []);

  const checkAuthStatus = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/auth/status`);
      if (!response.ok) {
        throw new Error('Failed to check auth status');
      }
      const data = await response.json();
      setIsInitialized(data.initialized);
      if (data.initialized) {
        fetchPasswords();
      }
    } catch (err) {
      console.error('Error checking auth status:', err);
      setError(err.message);
    }
  };

  const handleInitSubmit = async (e) => {
    e.preventDefault();
    setError('');

    if (initPassword.password !== initPassword.confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    try {
      const response = await fetch(`${API_BASE_URL}/api/auth/initialize`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          password: initPassword.password
        }),
      });

      if (!response.ok) {
        const errorData = await response.text();
        throw new Error(errorData);
      }

      setIsInitialized(true);
      setInitPassword({ password: '', confirmPassword: '' });
      setError('Master password initialized successfully');
      fetchPasswords();
    } catch (err) {
      setError(err.message);
    }
  };

  const handleMasterPasswordSubmit = async (e) => {
    e.preventDefault();
    setError('');

    if (masterPasswordForm.newPassword !== masterPasswordForm.confirmPassword) {
      setError('New passwords do not match');
      return;
    }

    try {
      const response = await fetch(`${API_BASE_URL}/api/auth/change`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          currentPassword: masterPasswordForm.currentPassword,
          newPassword: masterPasswordForm.newPassword,
        }),
      });

      if (!response.ok) {
        const errorData = await response.text();
        throw new Error(errorData);
      }

      setShowMasterPasswordForm(false);
      setMasterPasswordForm({
        currentPassword: '',
        newPassword: '',
        confirmPassword: '',
      });
      setError('Master password changed successfully');
    } catch (err) {
      setError(err.message);
    }
  };

  const handleMasterPasswordChange = (e) => {
    const { name, value } = e.target;
    setMasterPasswordForm((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const fetchPasswords = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/passwords`);
      if (!response.ok) {
        throw new Error('Failed to fetch passwords');
      }
      const data = await response.json();
      console.log('Received data:', data);
      setPasswords(data);
    } catch (err) {
      console.error('Error fetching passwords:', err);
      setError(err.message);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    try {
      const response = await fetch(`${API_BASE_URL}/api/passwords`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(newPassword),
      });

      if (!response.ok) {
        const errorData = await response.text();
        throw new Error(errorData);
      }

      setNewPassword({
        title: '',
        username: '',
        password: '',
        notes: '',
      });
      fetchPasswords();
    } catch (err) {
      setError(err.message);
    }
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setNewPassword((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleDelete = async (id) => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/passwords/${id}`, {
        method: 'DELETE',
      });

      if (!response.ok) {
        throw new Error('Failed to delete password');
      }

      fetchPasswords();
    } catch (err) {
      setError(err.message);
    }
  };

  const togglePasswordVisibility = (id) => {
    setVisiblePasswords((prev) => ({
      ...prev,
      [id]: !prev[id],
    }));
  };

  if (!isInitialized) {
    return (
      <div className="App">
        <h1>Vault-inator</h1>
        {error && <p style={{ color: error.includes('successfully') ? 'green' : 'red' }}>{error}</p>}
        
        <form onSubmit={handleInitSubmit} style={{ marginBottom: '20px', padding: '20px', border: '1px solid #ccc', borderRadius: '4px' }}>
          <h3>Initialize Master Password</h3>
          <div style={{ marginBottom: '10px' }}>
            <input
              type="password"
              name="password"
              placeholder="Master Password"
              value={initPassword.password}
              onChange={(e) => setInitPassword(prev => ({ ...prev, password: e.target.value }))}
              required
              style={{ width: '100%', padding: '8px', marginBottom: '10px' }}
            />
          </div>
          <div style={{ marginBottom: '10px' }}>
            <input
              type="password"
              name="confirmPassword"
              placeholder="Confirm Master Password"
              value={initPassword.confirmPassword}
              onChange={(e) => setInitPassword(prev => ({ ...prev, confirmPassword: e.target.value }))}
              required
              style={{ width: '100%', padding: '8px', marginBottom: '10px' }}
            />
          </div>
          <button
            type="submit"
            style={{
              padding: '8px 16px',
              backgroundColor: '#4CAF50',
              color: 'white',
              border: 'none',
              borderRadius: '4px',
              cursor: 'pointer'
            }}
          >
            Initialize Master Password
          </button>
        </form>
      </div>
    );
  }

  return (
    <div className="App">
      <h1>Vault-inator</h1>
      {error && <p style={{ color: error.includes('successfully') ? 'green' : 'red' }}>{error}</p>}
      
      <div style={{ marginBottom: '20px' }}>
        <button
          onClick={() => setShowMasterPasswordForm(!showMasterPasswordForm)}
          style={{
            padding: '8px 16px',
            backgroundColor: '#4CAF50',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: 'pointer'
          }}
        >
          {showMasterPasswordForm ? 'Cancel' : 'Change Master Password'}
        </button>
      </div>

      {showMasterPasswordForm && (
        <form onSubmit={handleMasterPasswordSubmit} style={{ marginBottom: '20px', padding: '20px', border: '1px solid #ccc', borderRadius: '4px' }}>
          <h3>Change Master Password</h3>
          <div style={{ marginBottom: '10px' }}>
            <input
              type="password"
              name="currentPassword"
              placeholder="Current Master Password"
              value={masterPasswordForm.currentPassword}
              onChange={handleMasterPasswordChange}
              required
              style={{ width: '100%', padding: '8px', marginBottom: '10px' }}
            />
          </div>
          <div style={{ marginBottom: '10px' }}>
            <input
              type="password"
              name="newPassword"
              placeholder="New Master Password"
              value={masterPasswordForm.newPassword}
              onChange={handleMasterPasswordChange}
              required
              style={{ width: '100%', padding: '8px', marginBottom: '10px' }}
            />
          </div>
          <div style={{ marginBottom: '10px' }}>
            <input
              type="password"
              name="confirmPassword"
              placeholder="Confirm New Master Password"
              value={masterPasswordForm.confirmPassword}
              onChange={handleMasterPasswordChange}
              required
              style={{ width: '100%', padding: '8px', marginBottom: '10px' }}
            />
          </div>
          <button
            type="submit"
            style={{
              padding: '8px 16px',
              backgroundColor: '#4CAF50',
              color: 'white',
              border: 'none',
              borderRadius: '4px',
              cursor: 'pointer'
            }}
          >
            Update Master Password
          </button>
        </form>
      )}

      <form onSubmit={handleSubmit}>
        <input
          type="text"
          name="title"
          placeholder="Title"
          value={newPassword.title}
          onChange={handleInputChange}
          required
        />
        <input
          type="text"
          name="username"
          placeholder="Username"
          value={newPassword.username}
          onChange={handleInputChange}
          required
        />
        <input
          type="password"
          name="password"
          placeholder="Password"
          value={newPassword.password}
          onChange={handleInputChange}
          required
        />
        <input
          type="text"
          name="notes"
          placeholder="Notes"
          value={newPassword.notes}
          onChange={handleInputChange}
        />
        <button type="submit">Add Password</button>
      </form>

      <div className="password-list">
        {passwords.map((p) => (
          <div key={p.id} className="password-item">
            <h3>{p.title}</h3>
            <p>Username: {p.username}</p>
            <p>
              Password:{' '}
              <span style={{ fontFamily: 'monospace' }}>
                {visiblePasswords[p.id] ? p.password : '••••••••'}
              </span>
              <button
                onClick={() => togglePasswordVisibility(p.id)}
                style={{ marginLeft: '10px' }}
              >
                {visiblePasswords[p.id] ? 'Hide' : 'Show'}
              </button>
            </p>
            {p.notes && <p>Notes: {p.notes}</p>}
            <button
              onClick={() => handleDelete(p.id)}
              style={{ color: 'red' }}
            >
              Delete
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}

export default App; 