import { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { useNavigate, Link } from 'react-router-dom';
import toast from 'react-hot-toast';
import { ChefHat } from 'lucide-react';

const Login = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const { login } = useAuth();
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        const res = await login(email, password);
        if (res.success) {
            toast.success("Welcome back!");
            navigate('/'); // Go to Home Page
        } else {
            toast.error(res.error);
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-orange-50">
            <div className="bg-white p-8 rounded-xl shadow-lg w-full max-w-md border border-orange-100">
                <div className="flex justify-center mb-6">
                    <div className="bg-orange-100 p-3 rounded-full">
                        <ChefHat className="w-8 h-8 text-orange-600" />
                    </div>
                </div>
                <h2 className="text-2xl font-bold text-center text-gray-800 mb-6">Recipe Manager Login</h2>
                
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
                        <input 
                            type="email"
                            className="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-orange-500 outline-none"
                            placeholder="you@example.com" 
                            value={email} onChange={(e) => setEmail(e.target.value)} 
                            required
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Password</label>
                        <input 
                            type="password"
                            className="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-orange-500 outline-none"
                            placeholder="••••••••" 
                            value={password} onChange={(e) => setPassword(e.target.value)} 
                            required
                        />
                    </div>
                    <button className="w-full bg-orange-600 text-white py-2 rounded font-semibold hover:bg-orange-700 transition">
                        Sign In
                    </button>
                </form>
                <p className="mt-4 text-center text-sm text-gray-600">
                    Don't have an account? <Link to="/register" className="text-orange-600 font-semibold">Sign up</Link>
                </p>
            </div>
        </div>
    );
};

export default Login;