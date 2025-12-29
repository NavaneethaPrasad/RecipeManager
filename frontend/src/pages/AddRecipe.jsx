import { useState, useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom'; // Import useSearchParams
import api from '../api/axios';
import toast from 'react-hot-toast';
import { Plus, Minus, Save } from 'lucide-react';

const AddRecipe = () => {
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();
    const editId = searchParams.get('edit'); // Check if ?edit=123 exists
    
    // Form State
    const [form, setForm] = useState({
        name: '', description: '', prep_time: 0, cook_time: 0, servings: 2, category: '', instructions: ''
    });

    const [ingredients, setIngredients] = useState([
        { name: '', amount: '', unit: '' }
    ]);

    // --- EFFECT: FETCH DATA IF EDITING ---
    useEffect(() => {
        if (editId) {
            const fetchRecipe = async () => {
                try {
                    const res = await api.get(`/recipes/${editId}`);
                    const r = res.data;
                    
                    // Pre-fill Form
                    setForm({
                        name: r.name,
                        description: r.description,
                        prep_time: r.prep_time,
                        cook_time: r.cook_time,
                        servings: r.servings,
                        category: r.category,
                        // Handle instructions array -> string
                        instructions: Array.isArray(r.instructions) 
                            ? r.instructions.map(i => i.text).join('\n') 
                            : r.instructions
                    });

                    // Pre-fill Ingredients
                    if (r.ingredients && r.ingredients.length > 0) {
                        setIngredients(r.ingredients.map(ing => ({
                            name: ing.name || (ing.ingredient ? ing.ingredient.name : ""),
                            amount: ing.quantity || ing.amount,
                            unit: ing.unit
                        })));
                    }
                } catch (err) {
                    toast.error("Failed to load recipe for editing");
                    navigate('/');
                }
            };
            fetchRecipe();
        }
    }, [editId, navigate]);
    // -------------------------------------

    const addIngredientRow = () => {
        setIngredients([...ingredients, { name: '', amount: '', unit: '' }]);
    };

    const removeIngredientRow = (index) => {
        const list = [...ingredients];
        list.splice(index, 1);
        setIngredients(list);
    };

    const handleIngChange = (index, field, value) => {
        const list = [...ingredients];
        list[index][field] = value;
        setIngredients(list);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        
        // Prepare Payload
        const payload = {
            ...form,
            prep_time: parseInt(form.prep_time),
            cook_time: parseInt(form.cook_time),
            servings: parseInt(form.servings),
            // Convert string instructions back to Array
            instructions: form.instructions.split('\n').filter(line => line.trim() !== ""),
            ingredients: ingredients.map(ing => ({
                name: ing.name,
                amount: parseFloat(ing.amount),
                unit: ing.unit
            }))
        };

        try {
            if (editId) {
                // UPDATE EXISTING
                await api.put(`/recipes/${editId}`, payload);
                toast.success("Recipe updated successfully!");
            } else {
                // CREATE NEW
                await api.post('/recipes', payload);
                toast.success("Recipe created successfully!");
            }
            navigate('/');
        } catch (err) {
            console.error(err);
            toast.error("Failed to save recipe");
        }
    };

    const inputClass = "w-full p-2 border border-gray-300 rounded focus:ring-2 focus:ring-orange-500 outline-none";

    return (
        <div className="max-w-3xl mx-auto p-6 bg-white shadow-lg rounded-xl mt-10 border border-gray-100">
            <h1 className="text-2xl font-bold mb-6 text-gray-800">
                {editId ? "Edit Recipe" : "Create New Recipe"}
            </h1>
            
            <form onSubmit={handleSubmit} className="space-y-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="md:col-span-2">
                        <label className="block text-sm font-medium mb-1">Recipe Name</label>
                        <input className={inputClass} required 
                            value={form.name} onChange={e => setForm({...form, name: e.target.value})} />
                    </div>
                    <div className="md:col-span-2">
                        <label className="block text-sm font-medium mb-1">Description</label>
                        <textarea className={inputClass} 
                            value={form.description} onChange={e => setForm({...form, description: e.target.value})} />
                    </div>
                    <div>
                        <label className="block text-sm font-medium mb-1">Prep Time (min)</label>
                        <input type="number" className={inputClass} 
                            value={form.prep_time} onChange={e => setForm({...form, prep_time: e.target.value})} />
                    </div>
                    <div>
                        <label className="block text-sm font-medium mb-1">Cook Time (min)</label>
                        <input type="number" className={inputClass} 
                            value={form.cook_time} onChange={e => setForm({...form, cook_time: e.target.value})} />
                    </div>
                    <div>
                        <label className="block text-sm font-medium mb-1">Servings</label>
                        <input type="number" className={inputClass} 
                            value={form.servings} onChange={e => setForm({...form, servings: e.target.value})} />
                    </div>
                    <div>
                        <label className="block text-sm font-medium mb-1">Category</label>
                        {/* REPLACED INPUT WITH SELECT */}
                        <select 
                            className={inputClass} 
                            value={form.category} 
                            onChange={e => setForm({...form, category: e.target.value})}
                        >
                            <option value="">Select a category</option>
                            <option value="Breakfast">Breakfast</option>
                            <option value="Lunch">Lunch</option>
                            <option value="Dinner">Dinner</option>
                            <option value="Snack">Snack</option>
                            <option value="Dessert">Dessert</option>
                            <option value="Other">Other</option>
                        </select>
                    </div>
                </div>

                {/* Ingredients */}
                <div>
                    <div className="flex justify-between items-center mb-2">
                        <h3 className="text-lg font-semibold text-gray-700">Ingredients</h3>
                        <button type="button" onClick={addIngredientRow} className="text-sm text-orange-600 flex items-center font-bold">
                            <Plus size={16} /> Add Item
                        </button>
                    </div>
                    <div className="space-y-2">
                        {ingredients.map((ing, i) => (
                            <div key={i} className="flex gap-2 items-center">
                                <input placeholder="Item" className={inputClass} required
                                    value={ing.name} onChange={e => handleIngChange(i, 'name', e.target.value)} />
                                <input placeholder="Qty" type="number" className={`${inputClass} w-24`} required
                                    value={ing.amount} onChange={e => handleIngChange(i, 'amount', e.target.value)} />
                                <input placeholder="Unit" className={`${inputClass} w-24`} required
                                    value={ing.unit} onChange={e => handleIngChange(i, 'unit', e.target.value)} />
                                
                                {ingredients.length > 1 && (
                                    <button type="button" onClick={() => removeIngredientRow(i)} className="text-red-500 p-2">
                                        <Minus size={18} />
                                    </button>
                                )}
                            </div>
                        ))}
                    </div>
                </div>

                {/* Instructions */}
                <div>
                    <label className="block text-sm font-medium mb-1">Instructions (One per line)</label>
                    <textarea className={`${inputClass} h-32`} placeholder="Step 1..." required
                        value={form.instructions} onChange={e => setForm({...form, instructions: e.target.value})} />
                </div>

                <button className="w-full bg-orange-600 text-white py-3 rounded-lg font-bold hover:bg-orange-700 flex justify-center items-center gap-2">
                    <Save size={20} /> {editId ? "Update Recipe" : "Save Recipe"}
                </button>
            </form>
        </div>
    );
};

export default AddRecipe;