import { Link, useNavigate } from 'react-router-dom';
import { ChefHat } from 'lucide-react';
import { useAuth } from '../context/AuthContext';

const Navbar = () => {
    const { logout } = useAuth();
    const navigate = useNavigate();

    const handleLogout = () => {
        logout();
        navigate('/login');
    };

    return (
        <nav className="bg-white shadow-sm border-b border-gray-100">
            <div className="container mx-auto px-6 py-4 flex justify-between items-center">
                <Link to="/" className="flex items-center gap-2 font-bold text-xl text-gray-800 hover:text-orange-600 transition">
                    <ChefHat className="text-orange-600" /> Recipe Manager
                </Link>
                
                <div className="flex items-center gap-6 text-sm font-medium text-gray-600">
                    <Link to="/recipes" className="hover:text-orange-600 transition">Recipes</Link>
                    <Link to="/meal-planner" className="hover:text-orange-600 transition">Meal Planner</Link>
                    <Link to="/shopping-list" className="hover:text-orange-600 transition">Shopping List</Link>
                    <Link to="/scale" className="hover:text-orange-600 transition">Scale Recipes</Link>
                    <button onClick={handleLogout} className="text-gray-500 hover:text-red-500 ml-4">Logout</button>
                </div>
            </div>
        </nav>
    );
};

export default Navbar;