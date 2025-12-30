import { useState, useEffect } from 'react';
import Navbar from '../components/Navbar';
import AddMealModal from '../components/AddMealModal'; 
import api from '../api/axios';
import { ChevronLeft, ChevronRight, Plus, Trash2, Calendar } from 'lucide-react';
import toast from 'react-hot-toast';
import EmptyState from '../components/EmptyState';

const MealPlanner = () => {
    const [currentDate, setCurrentDate] = useState(new Date());
    const [weekDates, setWeekDates] = useState([]);
    const [plans, setPlans] = useState([]);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [modalData, setModalData] = useState({ date: '', type: '' });
    const [hasRecipes, setHasRecipes] = useState(true);

    // Calculate Week
    useEffect(() => {
        const startOfWeek = new Date(currentDate);
        const day = startOfWeek.getDay(); 
        const diff = startOfWeek.getDate() - day + (day === 0 ? -6 : 1); 
        startOfWeek.setDate(diff);

        const days = [];
        for (let i = 0; i < 7; i++) {
            const d = new Date(startOfWeek);
            d.setDate(startOfWeek.getDate() + i);
            days.push(d);
        }

        // Check for recipes once when week changes (Moved out of the loop for performance)
        api.get('/recipes').then(res => {
            const data = res.data || []; // Ensure data is at least an empty array
            if (data.length === 0) {
                setHasRecipes(false);
            } else {
                setHasRecipes(true);
            }
        });

        setWeekDates(days);
        if(days.length > 0) fetchPlans(days[0], days[6]);
    }, [currentDate]);

    // Fetch Plans
    const fetchPlans = async (start, end) => {
        try {
            const s = start.toISOString().split('T')[0];
            const e = end.toISOString().split('T')[0];
            const res = await api.get(`/meal-plans?start_date=${s}&end_date=${e}`);
            setPlans(res.data || []);
        } catch (err) {
            console.error("Failed to fetch plans");
        }
    };

    const handleDelete = async (id) => {
        if(!confirm("Remove meal?")) return;
        try {
            await api.delete(`/meal-plans/${id}`);
            toast.success("Meal removed");
            fetchPlans(weekDates[0], weekDates[6]);
        } catch (err) {
            toast.error("Delete failed");
        }
    };

    const getPlan = (date, type) => {
        const dateStr = date.toISOString().split('T')[0];
        return plans.find(p => p.date.startsWith(dateStr) && p.meal_type === type);
    };

    const MEAL_TYPES = ['breakfast', 'lunch', 'dinner', 'snack'];

    // --- ADDED THIS SECTION TO SHOW EMPTY STATE ---
    if (!hasRecipes) {
        return (
            <div className="min-h-screen bg-gray-50">
                <Navbar />
                <div className="container mx-auto p-6 max-w-4xl">
                    <EmptyState 
                        title="Add Recipes First" 
                        message="You can't plan meals until you have recipes in your collection. Add your first recipe to get started!" 
                    />
                </div>
            </div>
        );
    }
    // ----------------------------------------------

    return (
        <div className="min-h-screen bg-gray-50">
            <Navbar />
            <div className="container mx-auto p-6">
                <div className="flex justify-between items-center mb-6 bg-white p-4 rounded-xl shadow-sm">
                    <h1 className="text-2xl font-bold flex items-center gap-2">
                        <Calendar className="text-orange-600" /> Meal Planner
                    </h1>
                    <div className="flex gap-4 items-center">
                        <button onClick={() => {
                            const d = new Date(currentDate); 
                            d.setDate(d.getDate() - 7); 
                            setCurrentDate(d);
                        }} className="p-2 hover:bg-gray-100 rounded-full"><ChevronLeft /></button>
                        
                        <span className="font-semibold">
                            {weekDates[0]?.toLocaleDateString()} - {weekDates[6]?.toLocaleDateString()}
                        </span>

                        <button onClick={() => {
                            const d = new Date(currentDate); 
                            d.setDate(d.getDate() + 7); 
                            setCurrentDate(d);
                        }} className="p-2 hover:bg-gray-100 rounded-full"><ChevronRight /></button>
                    </div>
                </div>

                <div className="grid grid-cols-7 gap-4 min-w-[1000px]">
                    {weekDates.map((date, i) => (
                        <div key={i} className="flex flex-col gap-3">
                            <div className={`p-2 text-center border-b-4 rounded-t ${
                                date.toDateString() === new Date().toDateString() ? 'border-orange-500 bg-orange-50' : 'border-transparent bg-white'
                            }`}>
                                <p className="font-bold">{date.toLocaleDateString('en-US', { weekday: 'short' })}</p>
                                <p>{date.getDate()}</p>
                            </div>

                            <div className="flex flex-col gap-2 bg-white p-2 rounded shadow h-full min-h-[400px]">
                                {MEAL_TYPES.map(type => {
                                    const plan = getPlan(date, type);
                                    return (
                                        <div key={type} className="group min-h-[80px]">
                                            <span className="text-xs text-gray-400 uppercase">{type}</span>
                                            {plan ? (
                                                <div className="bg-orange-100 p-2 rounded text-sm relative">
                                                    <p className="font-semibold truncate">{plan.recipe?.name}</p>
                                                    <p className="text-xs">{plan.target_servings} ppl</p>
                                                    <button 
                                                        onClick={() => handleDelete(plan.id)}
                                                        className="absolute top-1 right-1 text-red-500 opacity-0 group-hover:opacity-100"
                                                    >
                                                        <Trash2 size={14} />
                                                    </button>
                                                </div>
                                            ) : (
                                                <button 
                                                    onClick={() => {
                                                        setModalData({ date: date.toISOString().split('T')[0], type });
                                                        setIsModalOpen(true);
                                                    }}
                                                    className="w-full h-full border border-dashed border-gray-200 rounded flex items-center justify-center text-gray-300 hover:text-orange-500 hover:border-orange-300"
                                                >
                                                    <Plus size={16} />
                                                </button>
                                            )}
                                        </div>
                                    )
                                })}
                            </div>
                        </div>
                    ))}
                </div>
            </div>

            {isModalOpen && (
                <AddMealModal 
                    date={modalData.date}
                    mealType={modalData.type}
                    onClose={() => setIsModalOpen(false)}
                    onSave={() => fetchPlans(weekDates[0], weekDates[6])}
                />
            )}
        </div>
    );
};

export default MealPlanner;