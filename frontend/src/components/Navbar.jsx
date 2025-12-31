import { Link, useNavigate } from 'react-router-dom';
import { ChefHat } from 'lucide-react';
import { useAuth } from '../context/AuthContext';

const Navbar = () => {
    const { logout } = useAuth();
    const navigate = useNavigate();

    return (
        <nav className="bg-white shadow-sm border-b border-gray-100 mb-2">
            <div className="container mx-auto px-6 py-5 flex justify-between items-center">
                <Link to="/" className="flex items-center gap-2 group">
                    <div className="bg-orange-600 p-2 rounded-lg">
                        <ChefHat className="text-white" size={26} />
                    </div>
                    <span className="text-2xl font-black tracking-tighter text-slate-800">
                        Meal<span className="text-orange-600">Mate</span>
                    </span>
                </Link>
                
                <div className="flex items-center gap-8 text-lg font-bold text-slate-600">
                    <Link to="/recipes" className="hover:text-orange-600 transition">Recipes</Link>
                    <Link to="/meal-planner" className="hover:text-orange-600 transition">Meal Planner</Link>
                    <Link to="/shopping-list" className="hover:text-orange-600 transition">Shopping List</Link>
                    <Link to="/scale" className="hover:text-orange-600 transition">Scale</Link>
                    <button 
                        onClick={() => { logout(); navigate('/login'); }} 
                        className="text-gray-400 hover:text-red-500 ml-4 font-black transition text-sm uppercase tracking-widest"
                    >
                        Logout
                    </button>
                </div>
            </div>
        </nav>
    );
};

export default Navbar;