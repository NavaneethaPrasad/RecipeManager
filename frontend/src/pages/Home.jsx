import Navbar from '../components/Navbar';
import { ChefHat, Calendar, ShoppingCart, Scale } from 'lucide-react';
import { Link } from 'react-router-dom';

const Home = () => {
    return (
        <div className="min-h-screen bg-white">
            <Navbar />

            {/* Hero Section */}
            <div className="flex flex-col items-center justify-center pt-16 pb-12 px-4 text-center">
                <div className="mb-4">
                    <ChefHat size={64} className="text-orange-700 mx-auto" />
                </div>
                <h1 className="text-4xl font-extrabold text-gray-900 mb-4 tracking-tight">
                    Recipe Manager
                </h1>
                <p className="text-lg text-gray-500 max-w-2xl">
                    Organize your recipes, plan your meals, and generate shopping lists with ease
                </p>
            </div>

            {/* Cards Grid */}
            <div className="container mx-auto px-6 pb-20">
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
                    
                    {/* 1. Recipes Card */}
                    <div className="border border-gray-200 rounded-xl p-8 hover:shadow-xl transition duration-300 flex flex-col bg-white">
                        <div className="mb-4">
                            <ChefHat size={40} className="text-orange-600" />
                        </div>
                        <h3 className="text-xl font-bold text-gray-900 mb-2">Recipes</h3>
                        <p className="text-gray-500 text-sm mb-8 flex-1">
                            Create, edit, and manage your recipe collection
                        </p>
                        <Link to="/recipes" className="w-full block text-center border border-gray-300 text-gray-700 font-semibold py-2 rounded-lg hover:bg-gray-50 transition">
                            View Recipes
                        </Link>
                    </div>

                    {/* 2. Meal Planner Card */}
                    <div className="border border-gray-200 rounded-xl p-8 hover:shadow-xl transition duration-300 flex flex-col bg-white">
                        <div className="mb-4">
                            <Calendar size={40} className="text-orange-600" />
                        </div>
                        <h3 className="text-xl font-bold text-gray-900 mb-2">Meal Planner</h3>
                        <p className="text-gray-500 text-sm mb-8 flex-1">
                            Plan your meals for the week ahead
                        </p>
                        <Link to="/meal-planner" className="w-full block text-center border border-gray-300 text-gray-700 font-semibold py-2 rounded-lg hover:bg-gray-50 transition">
                            Plan Meals
                        </Link>
                    </div>

                    {/* 3. Shopping List Card */}
                    <div className="border border-gray-200 rounded-xl p-8 hover:shadow-xl transition duration-300 flex flex-col bg-white">
                        <div className="mb-4">
                            <ShoppingCart size={40} className="text-orange-600" />
                        </div>
                        <h3 className="text-xl font-bold text-gray-900 mb-2">Shopping List</h3>
                        <p className="text-gray-500 text-sm mb-8 flex-1">
                            Generate lists from your meal plans
                        </p>
                        <Link to="/shopping-list" className="w-full block text-center border border-gray-300 text-gray-700 font-semibold py-2 rounded-lg hover:bg-gray-50 transition">
                            View List
                        </Link>
                    </div>

                    {/* 4. Scale Recipes Card */}
                    <div className="border border-gray-200 rounded-xl p-8 hover:shadow-xl transition duration-300 flex flex-col bg-white">
                        <div className="mb-4">
                            <Scale size={40} className="text-orange-600" />
                        </div>
                        <h3 className="text-xl font-bold text-gray-900 mb-2">Scale Recipes</h3>
                        <p className="text-gray-500 text-sm mb-8 flex-1">
                            Adjust ingredient quantities easily
                        </p>
                        <Link to="/scale" className="w-full block text-center border border-gray-300 text-gray-700 font-semibold py-2 rounded-lg hover:bg-gray-50 transition">
                            Scale Recipes
                        </Link>
                    </div>

                </div>
            </div>
        </div>
    );
};

export default Home;