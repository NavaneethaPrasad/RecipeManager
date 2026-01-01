import { useState, useEffect } from 'react';
import api from '../api/axios';
import { X, ChefHat, Users } from 'lucide-react';
import toast from 'react-hot-toast';

const AddMealModal = ({ date, mealType, onClose, onSave }) => {
    const [recipes, setRecipes] = useState([]);
    const [selectedId, setSelectedId] = useState('');
    const [servings, setServings] = useState(2);

    useEffect(() => {
        api.get('/recipes')
            .then(res => setRecipes(res.data || []))
            .catch(console.error);
    }, []);

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            await api.post('/meal-plans', { 
                date, 
                meal_type: mealType, 
                recipe_id: parseInt(selectedId), 
                target_servings: parseInt(servings) 
            });
            toast.success("Planned successfully!"); 
            onSave(); 
            onClose();
        } catch { 
            toast.error("Please select a recipe"); 
        }
    };

    return (
        <div className="fixed inset-0 bg-slate-900/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
            <div className="bg-white rounded-[2rem] shadow-2xl w-full max-w-md overflow-hidden animate-in fade-in zoom-in duration-200 flex flex-col max-h-[90vh]">

                <div className="bg-orange-600 p-5 md:p-6 text-white flex justify-between items-center shrink-0">
                    <div className="flex items-center gap-3">
                        <div className="bg-white/20 p-2 rounded-lg">
                            <ChefHat size={20} className="text-white" />
                        </div>
                        <div>
                            <h2 className="text-lg md:text-2xl font-black uppercase tracking-tighter leading-none">
                                Add {mealType}
                            </h2>
                            <p className="text-orange-100 font-bold text-xs md:text-sm mt-1">
                                {new Date(date).toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' })}
                            </p>
                        </div>
                    </div>
                    <button 
                        onClick={onClose} 
                        className="p-2 hover:bg-orange-700 rounded-xl transition-colors active:scale-90"
                    >
                        <X size={24} />
                    </button>
                </div>
                
                <form onSubmit={handleSubmit} className="p-6 md:p-8 space-y-5 md:space-y-6 overflow-y-auto">
                    <div>
                        <label className="text-[10px] md:text-xs font-black text-slate-400 uppercase mb-2 block ml-1 tracking-widest">
                            Pick a Recipe
                        </label>
                        <select 
                            className="w-full p-3 md:p-4 bg-slate-50 border border-slate-100 rounded-2xl font-bold text-base md:text-lg text-slate-700 outline-none focus:ring-2 focus:ring-orange-500 transition-all appearance-none cursor-pointer"
                            onChange={e => setSelectedId(e.target.value)} 
                            required
                        >
                            <option value="">-- Choose Recipe --</option>
                            {(recipes || []).map(r => (
                                <option key={r.id || r.ID} value={r.id || r.ID}>{r.name}</option>
                            ))}
                        </select>
                    </div>

                    <div>
                        <label className="text-[10px] md:text-xs font-black text-slate-400 uppercase mb-2 block ml-1 tracking-widest">
                            Desired Servings
                        </label>
                        <div className="relative">
                            <input 
                                type="number" 
                                className="w-full p-3 md:p-4 bg-slate-50 border border-slate-100 rounded-2xl font-bold text-base md:text-lg text-slate-700 outline-none focus:ring-2 focus:ring-orange-500" 
                                value={servings} 
                                onChange={e => setServings(e.target.value)} 
                                min="1" 
                            />
                            <div className="absolute right-4 top-1/2 -translate-y-1/2 text-slate-300">
                                <Users size={20} />
                            </div>
                        </div>
                    </div>

                    <div className="pt-2">
                        <button 
                            type="submit"
                            className="w-full bg-slate-800 text-white py-4 rounded-2xl font-black text-base md:text-lg hover:bg-slate-900 transition-all shadow-lg shadow-slate-200 active:scale-[0.98]"
                        >
                            Confirm Plan
                        </button>
                        <button 
                            type="button"
                            onClick={onClose}
                            className="w-full mt-3 text-slate-400 font-bold text-sm uppercase tracking-widest py-2 hover:text-slate-600 md:hidden"
                        >
                            Cancel
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default AddMealModal;