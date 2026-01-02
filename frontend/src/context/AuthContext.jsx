import { createContext, useState, useEffect, useContext } from 'react';
import api from '../api/axios';

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true); 

    useEffect(() => {
    const checkSession = async () => {
            const authHint = localStorage.getItem('auth_hint');
            
            if (!authHint) {
                setLoading(false);
                return;
            }
        try {
            const res = await api.get('/profile');
            if (res.data && res.data.id) {
                setUser(res.data);
            } else {
                localStorage.removeItem('auth_hint');
            }
        } catch (err) {
            localStorage.removeItem('auth_hint');
                setUser(null);
        } finally {
            setLoading(false);
        }
    };
    checkSession();
}, []);

    const login = async (email, password) => {
        try {
            const res = await api.post('/auth/login', { email, password });
            localStorage.setItem('auth_hint', 'true');
            setUser(res.data.user);
            return { success: true };
        } catch (err) {
            return { success: false, error: err.response?.data?.error || "Login failed" };
        }
    };

const logout = async () => {
    try {
        await api.post('/auth/logout');
    } catch (err) {
        console.error("Logout request failed", err);
    } finally {
            localStorage.removeItem('auth_hint');
            setUser(null);
            window.location.href = '/login'; 
        }
};
    return (
        <AuthContext.Provider value={{ user, login, logout, loading }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => useContext(AuthContext);