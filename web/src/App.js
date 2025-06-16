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

  useEffect(() => {
    fetchPasswords();
  }, []);

  const fetchPasswords = async () => {
    try {
      console.log('Fetching passwords...');
      const response = await fetch(`${API_BASE_URL}/api/passwords`);
      if (!response.ok) {
        throw new Error('Failed to fetch passwords');
      }
      const data = await response.json();
      console.log('Received data:', data);
      
      // Handle null or invalid data
      if (!data || !Array.isArray(data)) {
        console.warn('Received invalid data format, initializing as empty array');
        setPasswords([]);
        setError('No password entries found.');
        return;
      }

      setPasswords(data);
      if (data.length === 0) {
        setError('No password entries found.');
      } else {
        setError('');
      }
    } catch (err) {
      console.error('Error fetching passwords:', err);
      setError(err.message);
      setPasswords([]); // Reset to empty array on error
    }
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setNewPassword((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch(`${API_BASE_URL}/api/passwords`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(newPassword),
      });
      if (!response.ok) {
        throw new Error('Failed to add password');
      }
      setNewPassword({ title: '', username: '', password: '', notes: '' });
      fetchPasswords();
    } catch (err) {
      setError(err.message);
    }
  };

  const handleDelete = async (id) => {
    try {
      console.log('Deleting password with ID:', id);
      const response = await fetch(`${API_BASE_URL}/api/passwords/${id}`, {
        method: 'DELETE',
      });
      
      if (!response.ok) {
        const errorData = await response.text();
        throw new Error(`Failed to delete password: ${errorData}`);
      }

      // Update the passwords state by filtering out the deleted entry
      setPasswords(prevPasswords => prevPasswords.filter(p => p.ID !== id));
      
      // Also remove the visibility state for the deleted password
      setVisiblePasswords(prev => {
        const newState = { ...prev };
        delete newState[id];
        return newState;
      });

      console.log('Successfully deleted password');
    } catch (err) {
      console.error('Error deleting password:', err);
      setError(err.message);
    }
  };

  const togglePasswordVisibility = (id) => {
    setVisiblePasswords(prev => ({
      ...prev,
      [id]: !prev[id]
    }));
  };

  const handleMasterPasswordChange = (e) => {
    const { name, value } = e.target;
    setMasterPasswordForm(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleMasterPasswordSubmit = async (e) => {
    e.preventDefault();
    setError('');

    if (masterPasswordForm.newPassword !== masterPasswordForm.confirmPassword) {
      setError('New passwords do not match');
      return;
    }

    try {
      const response = await fetch(`${API_BASE_URL}/api/master-password`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          currentPassword: masterPasswordForm.currentPassword,
          newPassword: masterPasswordForm.newPassword
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
        confirmPassword: ''
      });
      setError('Master password updated successfully');
    } catch (err) {
      setError(err.message);
    }
  };

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
      <h2>Passwords</h2>
      {passwords && passwords.length > 0 ? (
        <ul style={{ listStyle: 'none', padding: 0 }}>
          {passwords.map((password) => (
            <li key={password.ID} style={{ 
              border: '1px solid #ccc', 
              margin: '10px 0', 
              padding: '15px',
              borderRadius: '5px',
              backgroundColor: '#f9f9f9'
            }}>
              <strong>{password.Title}</strong>
              <br />
              Username: {password.Username}
              <br />
              <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                Password: {visiblePasswords[password.ID] ? password.Password : '••••••••'}
                <button
                  onClick={() => togglePasswordVisibility(password.ID)}
                  style={{
                    padding: '2px 8px',
                    backgroundColor: visiblePasswords[password.ID] ? '#666' : '#4CAF50',
                    color: 'white',
                    border: 'none',
                    borderRadius: '3px',
                    cursor: 'pointer',
                    fontSize: '0.8em'
                  }}
                >
                  {visiblePasswords[password.ID] ? 'Hide' : 'Show'}
                </button>
              </div>
              <br />
              Notes: {password.Notes || 'No notes'}
              <br />
              <button 
                onClick={() => handleDelete(password.ID)}
                style={{
                  marginTop: '10px',
                  padding: '5px 10px',
                  backgroundColor: '#ff4444',
                  color: 'white',
                  border: 'none',
                  borderRadius: '3px',
                  cursor: 'pointer'
                }}
              >
                Delete
              </button>
            </li>
          ))}
        </ul>
      ) : (
        <p>No passwords stored yet.</p>
      )}
    </div>
  );
}

export default App; 