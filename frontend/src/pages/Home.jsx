import React from 'react'; 
import Navbar from '../components/Navbar';
import { ChefHat, Calendar, ShoppingCart, Scale } from 'lucide-react';
import { Link } from 'react-router-dom';

const Home = () => {
    return (
        <div className="min-h-screen bg-slate-100">
            <Navbar />
            
            <div className="flex flex-col items-center justify-center pt-10 pb-8 md:pt-16 md:pb-12 px-4 text-center">
                <div className="bg-orange-600 p-3 md:p-4 rounded-[1.5rem] md:rounded-3xl mb-4 md:mb-6 shadow-xl shadow-orange-200">
                    <ChefHat className="text-white mx-auto w-10 h-10 md:w-16 md:h-16" />
                </div>

                <h1 className="text-3xl sm:text-4xl md:text-5xl font-black text-slate-800 mb-3 md:mb-4 tracking-tighter">
                    Meal<span className="text-orange-600">Mate</span>
                </h1>
                
                <p className="text-base md:text-xl text-slate-500 font-bold max-w-2xl px-2">
                    Your intelligent kitchen assistant for recipes and planning.
                </p>
            </div>

            <div className="container mx-auto px-4 md:px-6 pb-16 md:pb-20 max-w-6xl">
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 md:gap-8">
                    <DashboardCard 
                        to="/recipes" 
                        icon={<ChefHat size={32} />} 
                        title="Recipes" 
                        desc="Manage your collection" 
                    />
                    <DashboardCard 
                        to="/meal-planner" 
                        icon={<Calendar size={32} />} 
                        title="Planner" 
                        desc="Organize your week" 
                    />
                    <DashboardCard 
                        to="/shopping-list" 
                        icon={<ShoppingCart size={32} />} 
                        title="Shopping" 
                        desc="Auto-generate lists" 
                    />
                    <DashboardCard 
                        to="/scale" 
                        icon={<Scale size={32} />} 
                        title="Scale" 
                        desc="Portion control" 
                    />
                </div>
            </div>
        </div>
    );
};
t
const DashboardCard = ({ to, icon, title, desc }) => (
    <Link 
        to={to} 
        className="bg-white border border-gray-100 rounded-[2rem] p-6 md:p-8 hover:shadow-xl transition-all duration-300 group flex flex-col items-center text-center active:scale-95"
    >
        <div className="mb-3 md:mb-4 text-orange-600 group-hover:scale-110 transition-transform duration-300">
            {icon}
        </div>
        
        <h3 className="text-lg md:text-xl font-black text-slate-800 mb-1 md:mb-2 tracking-tight">
            {title}
        </h3>
        
        <p className="text-slate-400 font-medium text-xs md:text-sm mb-4 md:mb-6">
            {desc}
        </p>
        
        <span className="text-orange-600 font-black text-[10px] md:text-xs uppercase tracking-widest bg-orange-50 px-3 py-1 rounded-full group-hover:bg-orange-600 group-hover:text-white transition-colors duration-300">
            Open
        </span>
    </Link>
);

export default Home;