import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import api from '../api/axios';
import RecipeCard from '../components/RecipeCard';
import Navbar from '../components/Navbar';
import EmptyState from '../components/EmptyState'; // <<< ADD THIS IMPORT
import { Plus } from 'lucide-react';
import toast from 'react-hot-toast';

const RecipeList = () => {
    const [recipes, setRecipes] = useState([]);
    const [loading, setLoading] = useState(true);

    const fetchRecipes = async () => {
        try {
            const res = await api.get('/recipes');
            setRecipes(res.data || []); // Ensure it's an array even if backend returns null
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
        <div className="min-h-screen bg-gray-50">
            <Navbar />

            <div className="container mx-auto p-6">
                <div className="flex justify-between items-center mb-8">
                    <h1 className="text-3xl font-bold text-gray-800">My Recipes</h1>
                    {/* Header button: only show if we have recipes to avoid cluttering EmptyState */}
                    {recipes.length > 0 && (
                        <Link to="/add-recipe" className="bg-orange-600 text-white px-4 py-2 rounded-lg flex items-center gap-2 hover:bg-orange-700 transition shadow-md">
                            <Plus size={20} /> Add Recipe
                        </Link>
                    )}
                </div>

                {recipes.length === 0 ? (
                    <EmptyState 
                        title="Your Recipe Book is Empty" 
                        message="Start your culinary journey by adding your favorite recipes. You'll be able to scale them and add them to your meal plans later!"
                    />
                ) : (
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                        {recipes.map(recipe => (
                            <RecipeCard key={recipe.ID || recipe.id} recipe={recipe} onDelete={handleDelete} />
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
};

export default RecipeList;