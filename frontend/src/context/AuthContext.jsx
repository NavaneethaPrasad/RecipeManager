import { createContext, useState, useEffect, useContext } from 'react';
import api from '../api/axios';

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    // 1. On page load, check if user is already logged in
 useEffect(() => {
        const token = localStorage.getItem('token');
        const userData = localStorage.getItem('user');

        // CHECK: Make sure userData exists AND is not the string "undefined"
        if (token && userData && userData !== "undefined") {
            try {
                setUser(JSON.parse(userData));
            } catch (e) {
                console.error("Corrupt user data found, clearing storage.");
                localStorage.removeItem('token');
                localStorage.removeItem('user');
            }
        }
        setLoading(false);
    }, []);

    // 2. Login Function
    const login = async (email, password) => {
        try {
            const res = await api.post('/auth/login', { email, password });
            
            // --- DEBUG LOGGING ---
            console.log("SERVER RESPONSE:", res.data); 
            // ---------------------

            if (res.data.user) {
                // localStorage.setItem('token', res.data.token);
                localStorage.setItem('user', JSON.stringify(res.data.user));
                setUser(res.data.user);
                return { success: true };
            } else {
                // This is where your error is coming from
                return { success: false, error: "Invalid response from server" };
            }
        } catch (err) {
            console.error(err);
            return { success: false, error: err.response?.data?.error || "Login failed" };
        }
    };

    // 3. Logout Function
    const logout = () => {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        setUser(null);
    };

    return (
        <AuthContext.Provider value={{ user, login, logout, loading }}>
            {!loading && children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => useContext(AuthContext);