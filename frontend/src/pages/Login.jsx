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
            navigate('/'); 
        } else {
            toast.error(res.error);
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-slate-100 p-4 md:p-6">
            <div className="bg-white p-6 sm:p-10 rounded-[2.5rem] shadow-xl w-full max-w-md border border-gray-100 animate-in fade-in zoom-in duration-300">
                <div className="flex flex-col items-center mb-8">
                    <div className="bg-orange-600 p-4 rounded-2xl shadow-lg shadow-orange-200 mb-4">
                        <ChefHat className="w-8 h-8 text-white" />
                    </div>
                    <h2 className="text-3xl font-black text-center text-slate-800 tracking-tighter">
                        Meal<span className="text-orange-600">Mate</span>
                    </h2>
                    <p className="text-slate-400 font-bold text-sm mt-1 uppercase tracking-widest">Sign in to continue</p>
                </div>
                
                <form onSubmit={handleSubmit} className="space-y-5">
                    <div>
                        <label className="block text-xs font-black text-slate-500 uppercase mb-2 ml-1">Email Address</label>
                        <input 
                            type="email"
                            className="w-full p-4 bg-slate-50 border border-gray-100 rounded-2xl font-bold text-lg text-slate-700 outline-none focus:ring-2 focus:ring-orange-500 transition-all"
                            placeholder="name@example.com" 
                            value={email} onChange={(e) => setEmail(e.target.value)} 
                            required
                        />
                    </div>
                    <div>
                        <label className="block text-xs font-black text-slate-500 uppercase mb-2 ml-1">Password</label>
                        <input 
                            type="password"
                            className="w-full p-4 bg-slate-50 border border-gray-100 rounded-2xl font-bold text-lg text-slate-700 outline-none focus:ring-2 focus:ring-orange-500 transition-all"
                            placeholder="••••••••" 
                            value={password} onChange={(e) => setPassword(e.target.value)} 
                            required
                        />
                    </div>
                    
                    <button 
                        type="submit"
                        className="w-full bg-orange-600 text-white py-4 rounded-2xl font-black text-lg hover:bg-orange-700 transition-all shadow-lg shadow-orange-100 active:scale-95"
                    >
                        Login
                    </button>
                </form>

                <div className="mt-8 pt-6 border-t border-gray-50 text-center">
                    <p className="text-slate-500 font-bold">
                        New here? <Link to="/register" className="text-orange-600 hover:text-orange-700 underline decoration-2 underline-offset-4">Create an account</Link>
                    </p>
                </div>
            </div>
        </div>
    );
};

export default Login;