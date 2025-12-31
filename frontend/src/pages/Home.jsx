import Navbar from '../components/Navbar';
import { ChefHat, Calendar, ShoppingCart, Scale } from 'lucide-react';
import { Link } from 'react-router-dom';

const Home = () => {
    return (
        <div className="min-h-screen bg-slate-100">
            <Navbar />
            <div className="flex flex-col items-center justify-center pt-16 pb-12 px-4 text-center">
                <div className="bg-orange-600 p-4 rounded-3xl mb-6 shadow-xl shadow-orange-200">
                    <ChefHat size={64} className="text-white mx-auto" />
                </div>
                <h1 className="text-5xl font-black text-slate-800 mb-4 tracking-tighter">
                    Meal<span className="text-orange-600">Mate</span>
                </h1>
                <p className="text-xl text-slate-500 font-bold max-w-2xl">
                    Your intelligent kitchen assistant for recipes and planning.
                </p>
            </div>

            <div className="container mx-auto px-6 pb-20 max-w-6xl">
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
                    <DashboardCard to="/recipes" icon={<ChefHat/>} title="Recipes" desc="Manage your collection" />
                    <DashboardCard to="/meal-planner" icon={<Calendar/>} title="Planner" desc="Organize your week" />
                    <DashboardCard to="/shopping-list" icon={<ShoppingCart/>} title="Shopping" desc="Auto-generate lists" />
                    <DashboardCard to="/scale" icon={<Scale/>} title="Scale" desc="Portion control" />
                </div>
            </div>
        </div>
    );
};

const DashboardCard = ({ to, icon, title, desc }) => (
    <Link to={to} className="bg-white border border-gray-100 rounded-3xl p-8 hover:shadow-xl transition-all duration-300 group flex flex-col items-center text-center">
        <div className="mb-4 text-orange-600 group-hover:scale-110 transition-transform">{icon}</div>
        <h3 className="text-xl font-black text-slate-800 mb-2">{title}</h3>
        <p className="text-slate-400 font-medium text-sm mb-6">{desc}</p>
        <span className="text-orange-600 font-black text-xs uppercase tracking-widest">Open</span>
    </Link>
);

export default Home;