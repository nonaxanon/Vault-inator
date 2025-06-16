import React, { useState, useEffect } from 'react';
import './App.css';

// Add a base URL for API calls
const API_BASE_URL = 'http://localhost:8080';

function App() {
  const [passwords, setPasswords] = useState([]);
  const [error, setError] = useState('');
  const [showInitForm, setShowInitForm] = useState(false);
  const [masterPassword, setMasterPassword] = useState('');
  const [newPassword, setNewPassword] = useState({
    title: '',
    username: '',
    password: '',
    url: '',
    notes: ''
  });
  const [showPasswordForm, setShowPasswordForm] = useState(false);
  const [isInitialized, setIsInitialized] = useState(false);
  const [showChangePasswordForm, setShowChangePasswordForm] = useState(false);
  const [currentPassword, setCurrentPassword] = useState('');
  const [newMasterPassword, setNewMasterPassword] = useState('');
  const [confirmNewMasterPassword, setConfirmNewMasterPassword] = useState('');
  const [visiblePasswords, setVisiblePasswords] = useState({});
  const [searchTerm, setSearchTerm] = useState('');
  const [sortBy, setSortBy] = useState('title');
  const [sortOrder, setSortOrder] = useState('asc');

  useEffect(() => {
    checkAuthStatus();
  }, []);

  const checkAuthStatus = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/auth/status`);
      const data = await response.json();
      setIsInitialized(data.initialized);
      if (data.initialized) {
        fetchPasswords();
      } else {
        setShowInitForm(true);
      }
    } catch (error) {
      setError('Failed to check authentication status');
    }
  };

  const fetchPasswords = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/passwords`);
      const data = await response.json();
      setPasswords(data || []);
    } catch (error) {
      setError('Failed to fetch passwords');
    }
  };

  const handleInitialize = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch(`${API_BASE_URL}/api/auth/initialize`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ password: masterPassword })
      });
      if (response.ok) {
        setShowInitForm(false);
        setIsInitialized(true);
        fetchPasswords();
      } else {
        setError('Failed to initialize master password');
      }
    } catch (error) {
      setError('Failed to initialize master password');
    }
  };

  const handleAddPassword = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch(`${API_BASE_URL}/api/passwords`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newPassword)
      });
      if (response.ok) {
        setShowPasswordForm(false);
        setNewPassword({ title: '', username: '', password: '', url: '', notes: '' });
        fetchPasswords();
      } else {
        setError('Failed to add password');
      }
    } catch (error) {
      setError('Failed to add password');
    }
  };

  const handleDeletePassword = async (id) => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/passwords/${id}`, {
        method: 'DELETE'
      });
      if (response.ok) {
        fetchPasswords();
      } else {
        setError('Failed to delete password');
      }
    } catch (error) {
      setError('Failed to delete password');
    }
  };

  const handleChangeMasterPassword = async (e) => {
    e.preventDefault();
    if (newMasterPassword !== confirmNewMasterPassword) {
      setError('New passwords do not match');
      return;
    }
    try {
      const response = await fetch(`${API_BASE_URL}/api/auth/change`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          currentPassword,
          newPassword: newMasterPassword
        })
      });
      if (response.ok) {
        setShowChangePasswordForm(false);
        setCurrentPassword('');
        setNewMasterPassword('');
        setConfirmNewMasterPassword('');
      } else {
        setError('Failed to change master password');
      }
    } catch (error) {
      setError('Failed to change master password');
    }
  };

  const togglePasswordVisibility = (id) => {
    setVisiblePasswords(prev => ({
      ...prev,
      [id]: !prev[id]
    }));
  };

  const copyToClipboard = (text) => {
    navigator.clipboard.writeText(text);
  };

  const filteredPasswords = passwords
    .filter(pwd => 
      pwd.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
      pwd.username.toLowerCase().includes(searchTerm.toLowerCase()) ||
      pwd.url.toLowerCase().includes(searchTerm.toLowerCase())
    )
    .sort((a, b) => {
      const aValue = a[sortBy].toLowerCase();
      const bValue = b[sortBy].toLowerCase();
      return sortOrder === 'asc' 
        ? aValue.localeCompare(bValue)
        : bValue.localeCompare(aValue);
    });

  if (showInitForm) {
    return (
      <div className="init-container">
        <div className="init-card">
          <h2>Initialize Master Password</h2>
          <form onSubmit={handleInitialize}>
            <div className="form-group">
              <label>Master Password:</label>
              <input
                type="password"
                value={masterPassword}
                onChange={(e) => setMasterPassword(e.target.value)}
                required
              />
            </div>
            <button type="submit" className="btn-primary">Initialize</button>
          </form>
        </div>
      </div>
    );
  }

  return (
    <div className="app-container">
      <header className="app-header">
        <h1>Vault-inator</h1>
        <div className="header-actions">
          <button 
            className="btn-secondary"
            onClick={() => setShowChangePasswordForm(true)}
          >
            Change Master Password
          </button>
          <button 
            className="btn-primary"
            onClick={() => setShowPasswordForm(true)}
          >
            Add New Password
          </button>
        </div>
      </header>

      {error && <div className="error-message">{error}</div>}

      <div className="search-sort-container">
        <div className="search-box">
          <input
            type="text"
            placeholder="Search passwords..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
        <div className="sort-controls">
          <select 
            value={sortBy} 
            onChange={(e) => setSortBy(e.target.value)}
          >
            <option value="title">Title</option>
            <option value="username">Username</option>
            <option value="url">URL</option>
          </select>
          <button 
            className="btn-icon"
            onClick={() => setSortOrder(prev => prev === 'asc' ? 'desc' : 'asc')}
          >
            {sortOrder === 'asc' ? '↑' : '↓'}
          </button>
        </div>
      </div>

      <div className="passwords-grid">
        {filteredPasswords.map((pwd) => (
          <div key={pwd.id} className="password-card">
            <div className="card-header">
              <h3>{pwd.title}</h3>
              <div className="card-actions">
                <button 
                  className="btn-icon"
                  onClick={() => togglePasswordVisibility(pwd.id)}
                >
                  {visiblePasswords[pwd.id] ? '👁️' : '👁️‍🗨️'}
                </button>
                <button 
                  className="btn-icon"
                  onClick={() => handleDeletePassword(pwd.id)}
                >
                  🗑️
                </button>
              </div>
            </div>
            <div className="card-content">
              <div className="info-row">
                <span className="label">Username:</span>
                <div className="value-with-copy">
                  <span>{pwd.username}</span>
                  <button 
                    className="btn-icon"
                    onClick={() => copyToClipboard(pwd.username)}
                  >
                    📋
                  </button>
                </div>
              </div>
              <div className="info-row">
                <span className="label">Password:</span>
                <div className="value-with-copy">
                  <span>{visiblePasswords[pwd.id] ? pwd.password : '••••••••'}</span>
                  <button 
                    className="btn-icon"
                    onClick={() => copyToClipboard(pwd.password)}
                  >
                    📋
                  </button>
                </div>
              </div>
              {pwd.url && (
                <div className="info-row">
                  <span className="label">URL:</span>
                  <div className="value-with-copy">
                    <a href={pwd.url} target="_blank" rel="noopener noreferrer">
                      {pwd.url}
                    </a>
                    <button 
                      className="btn-icon"
                      onClick={() => copyToClipboard(pwd.url)}
                    >
                      📋
                    </button>
                  </div>
                </div>
              )}
              {pwd.notes && (
                <div className="info-row">
                  <span className="label">Notes:</span>
                  <span className="notes">{pwd.notes}</span>
                </div>
              )}
            </div>
          </div>
        ))}
      </div>

      {showPasswordForm && (
        <div className="modal">
          <div className="modal-content">
            <h2>Add New Password</h2>
            <form onSubmit={handleAddPassword}>
              <div className="form-group">
                <label htmlFor="title">Title:</label>
                <input
                  type="text"
                  id="title"
                  value={newPassword.title}
                  onChange={(e) => setNewPassword({ ...newPassword, title: e.target.value })}
                  required
                />
              </div>
              <div className="form-group">
                <label htmlFor="username">Username:</label>
                <input
                  type="text"
                  id="username"
                  value={newPassword.username}
                  onChange={(e) => setNewPassword({ ...newPassword, username: e.target.value })}
                  required
                />
              </div>
              <div className="form-group">
                <label htmlFor="password">Password:</label>
                <input
                  type="password"
                  id="password"
                  value={newPassword.password}
                  onChange={(e) => setNewPassword({ ...newPassword, password: e.target.value })}
                  required
                />
              </div>
              <div className="form-group">
                <label htmlFor="url">URL:</label>
                <input
                  type="url"
                  id="url"
                  value={newPassword.url}
                  onChange={(e) => setNewPassword({ ...newPassword, url: e.target.value })}
                  placeholder="https://example.com"
                />
              </div>
              <div className="form-group">
                <label htmlFor="notes">Notes:</label>
                <textarea
                  id="notes"
                  value={newPassword.notes}
                  onChange={(e) => setNewPassword({ ...newPassword, notes: e.target.value })}
                />
              </div>
              <div className="form-actions">
                <button type="button" className="btn-secondary" onClick={() => setShowPasswordForm(false)}>
                  Cancel
                </button>
                <button type="submit" className="btn-primary">Add Password</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {showChangePasswordForm && (
        <div className="modal">
          <div className="modal-content">
            <h2>Change Master Password</h2>
            <form onSubmit={handleChangeMasterPassword}>
              <div className="form-group">
                <label>Current Password:</label>
                <input
                  type="password"
                  value={currentPassword}
                  onChange={(e) => setCurrentPassword(e.target.value)}
                  required
                />
              </div>
              <div className="form-group">
                <label>New Password:</label>
                <input
                  type="password"
                  value={newMasterPassword}
                  onChange={(e) => setNewMasterPassword(e.target.value)}
                  required
                />
              </div>
              <div className="form-group">
                <label>Confirm New Password:</label>
                <input
                  type="password"
                  value={confirmNewMasterPassword}
                  onChange={(e) => setConfirmNewMasterPassword(e.target.value)}
                  required
                />
              </div>
              <div className="form-actions">
                <button type="button" className="btn-secondary" onClick={() => setShowChangePasswordForm(false)}>
                  Cancel
                </button>
                <button type="submit" className="btn-primary">Change Password</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

export default App; 