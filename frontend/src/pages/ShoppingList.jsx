import { useState, useEffect } from 'react';
import Navbar from '../components/Navbar';
import api from '../api/axios';
import { ShoppingCart, RefreshCw, CheckSquare, Square, Calendar } from 'lucide-react';
import toast from 'react-hot-toast';
import EmptyState from '../components/EmptyState';

const ShoppingList = () => {
    const [shoppingList, setShoppingList] = useState(null);
    const [loading, setLoading] = useState(false);
    const [hasRecipes, setHasRecipes] = useState(true);
    
    const today = new Date();
    const nextWeek = new Date();
    nextWeek.setDate(today.getDate() + 7);

    const [dates, setDates] = useState({
        start: today.toISOString().split('T')[0],
        end: nextWeek.toISOString().split('T')[0]
    });

    useEffect(() => {
        // --- FIXED: Ensure data is an array before checking length ---
        api.get('/recipes').then(res => {
            const data = res.data || [];
            setHasRecipes(data.length > 0);
        }).catch(err => {
            console.error("Error checking recipes", err);
            setHasRecipes(false); 
        });

        const savedId = localStorage.getItem('last_shopping_list_id');
        if (savedId) {
            fetchList(savedId);
        }
    }, []);

    const fetchList = async (id) => {
        setLoading(true);
        try {
            const res = await api.get(`/shopping-lists/${id}`);
            setShoppingList(res.data);
        } catch (err) {
            console.error("List not found");
            localStorage.removeItem('last_shopping_list_id');
        } finally {
            setLoading(false);
        }
    };

    const handleGenerate = async () => {
        setLoading(true);
        try {
            const res = await api.post('/shopping-lists/generate', {
                start_date: dates.start,
                end_date: dates.end
            });
            
            setShoppingList(res.data);
            
            if (res.data?.id) {
                localStorage.setItem('last_shopping_list_id', res.data.id);
            }
            
            toast.success("Shopping List Generated!");
        } catch (err) {
            console.error(err);
            toast.error("Failed to generate. Do you have meals planned?");
        } finally {
            setLoading(false);
        }
    };

    const handleToggle = async (itemId) => {
        // --- FIXED: Null-safe check for items mapping ---
        const updatedItems = (shoppingList?.items || []).map(item => 
            item.id === itemId ? { ...item, checked: !item.checked } : item
        );
        
        setShoppingList({ ...shoppingList, items: updatedItems });

        try {
            await api.patch(`/shopping-lists/items/${itemId}/toggle`);
        } catch (err) {
            toast.error("Failed to update status");
            if (shoppingList?.id) fetchList(shoppingList.id);
        }
    };

    if (!hasRecipes) {
        return (
            <div className="min-h-screen bg-gray-50">
                <Navbar />
                <div className="container mx-auto p-6 max-w-4xl">
                    <EmptyState 
                        title="No Recipes Found" 
                        message="Add recipes first to generate a shopping list from your weekly plan. Once you have recipes, you can plan meals and generate a list of ingredients automatically." 
                    />
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gray-50">
            <Navbar />
            
            <div className="container mx-auto p-6 max-w-4xl">
                <div className="flex items-center gap-3 mb-6">
                    <div className="bg-orange-100 p-3 rounded-full">
                        <ShoppingCart className="text-orange-600" size={32} />
                    </div>
                    <div>
                        <h1 className="text-2xl font-bold text-gray-800">Shopping List</h1>
                        <p className="text-gray-500 text-sm">Aggregated ingredients from your Meal Plan</p>
                    </div>
                </div>

                <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-200 mb-8">
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-end">
                        <div>
                            <label className="block text-xs font-bold text-gray-500 uppercase mb-1">Start Date</label>
                            <input 
                                type="date" 
                                className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 outline-none"
                                value={dates.start}
                                onChange={e => setDates({...dates, start: e.target.value})}
                            />
                        </div>
                        <div>
                            <label className="block text-xs font-bold text-gray-500 uppercase mb-1">End Date</label>
                            <input 
                                type="date" 
                                className="w-full p-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-orange-500 outline-none"
                                value={dates.end}
                                onChange={e => setDates({...dates, end: e.target.value})}
                            />
                        </div>
                        <button 
                            onClick={handleGenerate}
                            disabled={loading}
                            className="bg-orange-600 text-white px-6 py-2.5 rounded-lg font-bold hover:bg-orange-700 transition flex items-center justify-center gap-2 disabled:opacity-50"
                        >
                            {loading ? "Processing..." : (
                                <>
                                    <RefreshCw size={18} /> Generate List
                                </>
                            )}
                        </button>
                    </div>
                </div>

                {!shoppingList ? (
                    <div className="text-center py-20 bg-white rounded-xl border border-dashed border-gray-300">
                        <ShoppingCart className="mx-auto text-gray-300 mb-4" size={48} />
                        <p className="text-gray-500 text-lg">No list generated yet.</p>
                        <p className="text-gray-400 text-sm">Select a date range above and click Generate.</p>
                    </div>
                ) : (
                    <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
                        <div className="p-4 bg-orange-50 border-b border-orange-100 flex justify-between items-center">
                            <h3 className="font-bold text-orange-800 flex items-center gap-2">
                                <Calendar size={16}/> 
                                {shoppingList?.start_date} to {shoppingList?.end_date}
                            </h3>
                            <span className="text-xs font-semibold bg-white text-orange-600 px-2 py-1 rounded border border-orange-200">
                                {/* FIXED: Null-safe length check */}
                                {(shoppingList?.items || []).length} Items
                            </span>
                        </div>

                        {/* FIXED: Null-safe length check */}
                        {(shoppingList?.items || []).length === 0 ? (
                            <div className="p-8 text-center text-gray-500">
                                List is empty. Did you plan any meals for these dates?
                            </div>
                        ) : (
                            <div className="divide-y divide-gray-100">
                                {/* FIXED: Null-safe map check */}
                                {(shoppingList?.items || []).map((item) => (
                                    <div 
                                        key={item.id} 
                                        onClick={() => handleToggle(item.id)}
                                        className={`p-4 flex items-center gap-4 cursor-pointer transition hover:bg-gray-50 ${item.checked ? 'bg-gray-50' : 'bg-white'}`}
                                    >
                                        <button className={`transition ${item.checked ? "text-green-500" : "text-gray-300 hover:text-orange-500"}`}>
                                            {item.checked ? <CheckSquare size={24} /> : <Square size={24} />}
                                        </button>
                                        
                                        <div className="flex-1">
                                            <p className={`font-medium text-lg ${item.checked ? 'line-through text-gray-400' : 'text-gray-800'}`}>
                                                {item.name}
                                            </p>
                                        </div>

                                        <div className="text-right">
                                            <span className={`font-bold ${item.checked ? 'text-gray-400' : 'text-orange-600'}`}>
                                                {parseFloat(item.quantity).toFixed(1)}
                                            </span>
                                            <span className="text-gray-500 text-sm ml-1">{item.unit}</span>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                )}
            </div>
        </div>
    );
};

export default ShoppingList;