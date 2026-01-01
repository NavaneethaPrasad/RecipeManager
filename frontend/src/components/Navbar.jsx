import { Link, useNavigate } from 'react-router-dom';
import { ChefHat,Menu,X } from 'lucide-react';
import { useAuth } from '../context/AuthContext';
import { useState } from 'react';

const Navbar = () => {
    const { logout } = useAuth();
    const navigate = useNavigate();
    const [isMenuOpen, setIsMenuOpen] = useState(false);

     const handleLogout = () => {
        logout();
        navigate('/login');
    };

    const navLinks = [
        { name: 'Recipes', path: '/recipes' },
        { name: 'Meal Planner', path: '/meal-planner' },
        { name: 'Shopping List', path: '/shopping-list' },
        { name: 'Scale', path: '/scale' },
    ];

    return (
        <nav className="bg-white shadow-sm border-b border-gray-100 mb-2 sticky top-0 z-50">
            <div className="container mx-auto px-4 md:px-6 py-4 flex justify-between items-center">
                <Link to="/" className="flex items-center gap-2 group">
                    <div className="bg-orange-600 p-2 rounded-lg">
                        <ChefHat className="text-white" size={26} />
                    </div>
                    <span className="text-xl md:text-2xl font-black tracking-tighter text-slate-800">
                        Meal<span className="text-orange-600">Mate</span>
                    </span>
                </Link>
                
                <div className="hidden md:flex items-center gap-8 text-lg font-bold text-slate-600">
                    {navLinks.map((link) => (
                        <Link key={link.path} to={link.path} className="hover:text-orange-600 transition">
                            {link.name}
                        </Link>
                    ))}
                    <button 
                        o onClick={handleLogout} 
                        className="text-gray-400 hover:text-red-500 ml-4 font-black transition text-sm uppercase tracking-widest"
                    >
                        Logout
                    </button>
                </div>
                 <div className="md:hidden">
                    <button 
                        onClick={() => setIsMenuOpen(!isMenuOpen)}
                        className="p-2 text-slate-600 hover:bg-slate-100 rounded-lg"
                    >
                        {isMenuOpen ? <X size={28} /> : <Menu size={28} />}
                    </button>
                </div>
            </div>
             {isMenuOpen && (
                <div className="md:hidden bg-white border-t border-gray-100 p-4 space-y-4 shadow-lg animate-in slide-in-from-top duration-200">
                    {navLinks.map((link) => (
                        <Link 
                            key={link.path} 
                            to={link.path} 
                            onClick={() => setIsMenuOpen(false)}
                            className="block text-lg font-bold text-slate-600 hover:text-orange-600"
                        >
                            {link.name}
                        </Link>
                    ))}
                    <button 
                        onClick={handleLogout}
                        className="block w-full text-left text-red-500 font-black text-sm uppercase py-2 border-t border-gray-50"
                    >
                        Logout
                    </button>
                </div>
            )}
        </nav>
    );
};

export default Navbar;