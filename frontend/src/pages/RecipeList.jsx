import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import api from '../api/axios';
import RecipeCard from '../components/RecipeCard';
import Navbar from '../components/Navbar';
import EmptyState from '../components/EmptyState'; 
import { Plus, ChefHat } from 'lucide-react'; 
import toast from 'react-hot-toast';

const RecipeList = () => {
    const [recipes, setRecipes] = useState([]);
    const [loading, setLoading] = useState(true);

    const fetchRecipes = async () => {
        try {
            const res = await api.get('/recipes');
            setRecipes(res.data || []); 
        } catch (err) {
            console.error("Error fetching recipes");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchRecipes();
    }, []);

    const handleDelete = async (id) => {
        if (!window.confirm("Are you sure you want to delete this recipe?")) return;
        try {
            await api.delete(`/recipes/${id}`);
            setRecipes(recipes.filter(r => (r.ID || r.id) !== id));
            toast.success("Recipe deleted");
        } catch (err) {
            toast.error("Failed to delete");
        }
    };

    if (loading) return (
        <div className="min-h-screen bg-gray-50">
            <Navbar />
            <div className="text-center mt-20 text-gray-500 font-medium">Loading recipes...</div>
        </div>
    );

return (
    <div className="min-h-screen bg-gray-50 pb-10">
        <Navbar />
        <div className="container mx-auto p-6 max-w-7xl">
            <div className="flex justify-between items-center mb-8 bg-white p-6 rounded-2xl shadow-sm border border-gray-100">
                <h1 className="text-2xl font-black flex items-center gap-3 text-slate-800">
                    <ChefHat className="text-orange-600" size={28} /> My Recipes
                </h1>
                <Link to="/add-recipe" className="bg-orange-600 text-white px-6 py-2 rounded-xl font-bold hover:bg-orange-700 transition flex items-center gap-2">
                    <Plus size={20} /> Add Recipe
                </Link>
            </div>

            <div className="bg-white p-8 rounded-2xl shadow-sm border border-gray-100 min-h-[500px]">
                {recipes.length === 0 ? (
                    <EmptyState title="No Recipes" message="Add your first recipe to start." />
                ) : (
                    <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-6">
                        {recipes.map(recipe => (
                            <RecipeCard key={recipe.ID} recipe={recipe} onDelete={handleDelete} />
                        ))}
                    </div>
                )}
            </div>
        </div>
    </div>
);
};

export default RecipeList;