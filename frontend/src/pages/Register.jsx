import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import api from '../api/axios'; 
import toast from 'react-hot-toast';
import { ChefHat, UserPlus } from 'lucide-react';

const Register = () => {
    const navigate = useNavigate();
    
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
            await api.post('/auth/register', formData);
            toast.success("Account created! Please login.");
            navigate('/login'); 
        } catch (err) {
            const errorMsg = err.response?.data?.error || "Registration failed";
            toast.error(errorMsg);
        }
    };

    const inputClass = "w-full p-4 bg-slate-50 border border-gray-100 rounded-2xl font-bold text-lg text-slate-700 outline-none focus:ring-2 focus:ring-orange-500 transition-all";
    const labelClass = "block text-xs font-black text-slate-500 uppercase mb-2 ml-1";

    return (
        <div className="min-h-screen flex items-center justify-center bg-slate-100 p-4 md:p-6">
            <div className="bg-white p-6 sm:p-10 rounded-[2.5rem] shadow-xl w-full max-w-md border border-gray-100 animate-in fade-in zoom-in duration-300">

                <div className="flex flex-col items-center mb-8">
                    <div className="bg-orange-600 p-4 rounded-2xl shadow-lg shadow-orange-200 mb-4">
                        <ChefHat className="w-8 h-8 text-white" />
                    </div>
                    <h2 className="text-3xl font-black text-center text-slate-800 tracking-tighter">
                        Join <span className="text-orange-600">MealMate</span>
                    </h2>
                    <p className="text-slate-400 font-bold text-sm mt-1 uppercase tracking-widest text-center">
                        Start your culinary journey
                    </p>
                </div>
                
                <form onSubmit={handleSubmit} className="space-y-5">
                    <div>
                        <label className={labelClass}>Full Name</label>
                        <input 
                            name="name"
                            type="text"
                            className={inputClass}
                            placeholder="Navaneetha Prasad" 
                            value={formData.name} 
                            onChange={handleChange} 
                            required
                        />
                    </div>

                    <div>
                        <label className={labelClass}>Email Address</label>
                        <input 
                            name="email"
                            type="email"
                            className={inputClass}
                            placeholder="name@example.com" 
                            value={formData.email} 
                            onChange={handleChange} 
                            required
                        />
                    </div>

                    <div>
                        <label className={labelClass}>Password</label>
                        <input 
                            name="password"
                            type="password"
                            className={inputClass}
                            placeholder="Min. 6 characters" 
                            value={formData.password} 
                            onChange={handleChange} 
                            required
                            minLength={6}
                        />
                    </div>

                    <button 
                        type="submit"
                        className="w-full bg-orange-600 text-white py-4 rounded-2xl font-black text-lg hover:bg-orange-700 transition-all shadow-lg shadow-orange-100 flex items-center justify-center gap-2 active:scale-95"
                    >
                        <UserPlus size={22} strokeWidth={3} />
                        Create Account
                    </button>
                </form>
                
                <div className="mt-8 pt-6 border-t border-gray-50 text-center">
                    <p className="text-slate-500 font-bold">
                        Already a member? <Link to="/login" className="text-orange-600 hover:text-orange-700 underline decoration-2 underline-offset-4">Login here</Link>
                    </p>
                </div>
            </div>
        </div>
    );
};

export default Register;