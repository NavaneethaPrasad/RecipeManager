import { Link } from 'react-router-dom';
import { ChefHat, Plus } from 'lucide-react';

const EmptyState = ({ title, message, showAddButton = true }) => {
    return (
        <div className="flex flex-col items-center justify-center py-20 px-4 bg-white rounded-2xl border-2 border-dashed border-gray-200 mt-8 shadow-sm">
            <div className="bg-orange-100 p-4 rounded-full mb-6">
                <ChefHat size={48} className="text-orange-600" />
            </div>
            <h2 className="text-2xl font-bold text-gray-800 mb-2">{title}</h2>
            <p className="text-gray-500 text-center max-w-sm mb-8">
                {message}
            </p>
            {showAddButton && (
                <Link 
                    to="/add-recipe" 
                    className="bg-orange-600 text-white px-6 py-3 rounded-xl font-bold hover:bg-orange-700 transition flex items-center gap-2 shadow-lg shadow-orange-200"
                >
                    <Plus size={20} /> Add Your First Recipe
                </Link>
            )}
        </div>
    );
};

export default EmptyState;