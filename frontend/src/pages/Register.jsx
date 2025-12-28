import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import api from '../api/axios'; // Use your configured axios instance
import toast from 'react-hot-toast';
import { ChefHat, UserPlus } from 'lucide-react';

const Register = () => {
    const navigate = useNavigate();
    
    // Form State
    const [formData, setFormData] = useState({
        name: '',
        email: '',
        password: ''
    });

    const handleChange = (e) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            // Call Backend
            await api.post('/auth/register', formData);
            
            // Success
            toast.success("Account created! Please login.");
            navigate('/login'); // Redirect to Login page
        } catch (err) {
            // Error Handling
            const errorMsg = err.response?.data?.error || "Registration failed";
            toast.error(errorMsg);
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
                <h2 className="text-2xl font-bold text-center text-gray-800 mb-2">Create Account</h2>
                <p className="text-center text-gray-500 mb-6">Join Recipe Manager today</p>
                
                <form onSubmit={handleSubmit} className="space-y-4">
                    {/* Name Input */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Full Name</label>
                        <input 
                            name="name"
                            type="text"
                            className="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-orange-500 outline-none"
                            placeholder="John Doe" 
                            value={formData.name} 
                            onChange={handleChange} 
                            required
                        />
                    </div>

                    {/* Email Input */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
                        <input 
                            name="email"
                            type="email"
                            className="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-orange-500 outline-none"
                            placeholder="you@example.com" 
                            value={formData.email} 
                            onChange={handleChange} 
                            required
                        />
                    </div>

                    {/* Password Input */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Password</label>
                        <input 
                            name="password"
                            type="password"
                            className="w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-orange-500 outline-none"
                            placeholder="Min 6 chars" 
                            value={formData.password} 
                            onChange={handleChange} 
                            required
                            minLength={6}
                        />
                    </div>

                    <button className="w-full bg-orange-600 text-white py-2 rounded font-semibold hover:bg-orange-700 transition flex items-center justify-center gap-2">
                        <UserPlus size={18} />
                        Sign Up
                    </button>
                </form>
                
                <p className="mt-4 text-center text-sm text-gray-600">
                    Already have an account? <Link to="/login" className="text-orange-600 font-semibold">Login</Link>
                </p>
            </div>
        </div>
    );
};

export default Register;