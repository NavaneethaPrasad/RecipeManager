import { useState, useEffect } from 'react';
import api from '../api/axios';
import { X } from 'lucide-react';
import toast from 'react-hot-toast';

const AddMealModal = ({ date, mealType, onClose, onSave }) => {
    const [recipes, setRecipes] = useState([]);
    const [selectedRecipeId, setSelectedRecipeId] = useState('');
    const [servings, setServings] = useState(2);

    useEffect(() => {
        api.get('/recipes').then(res => setRecipes(res.data)).catch(console.error);
    }, []);

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (!selectedRecipeId) return toast.error("Select a recipe");

        try {
            await api.post('/meal-plans', {
                date,
                meal_type: mealType,
                recipe_id: parseInt(selectedRecipeId),
                target_servings: parseInt(servings)
            });
            toast.success("Added!");
            onSave();
            onClose();
        } catch (err) {
            toast.error("Failed to add meal");
        }
    };

    return (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
            <div className="bg-white p-6 rounded-lg shadow-xl w-96 relative">
                <button onClick={onClose} className="absolute top-4 right-4"><X size={20}/></button>
                <h2 className="text-xl font-bold mb-4 capitalize">Add {mealType}</h2>
                
                <form onSubmit={handleSubmit} className="space-y-4">
                    <select 
                        className="w-full p-2 border rounded"
                        onChange={e => setSelectedRecipeId(e.target.value)}
                        required
                    >
                        <option value="">Select Recipe</option>
                         {(recipes || []).map(r => (
                                <option key={r.id || r.ID} value={r.id || r.ID}>{r.name}</option>
                            ))}
                    </select>

                    <input 
                        type="number" 
                        className="w-full p-2 border rounded"
                        value={servings}
                        onChange={e => setServings(e.target.value)}
                        min="1"
                        placeholder="Servings"
                    />

                    <button className="w-full bg-orange-600 text-white py-2 rounded">Save</button>
                </form>
            </div>
        </div>
    );
};

export default AddMealModal;