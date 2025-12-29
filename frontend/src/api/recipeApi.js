import api from "./axios";

// Create recipe
export const createRecipe = (data) => {
  return api.post("/recipes", data);
};

// Get my recipes
export const getMyRecipes = () => {
  return api.get("/recipes");
};

// Get recipe by ID
export const getRecipeById = (id) => {
  return api.get(`/recipes/${id}`);
};

// Update recipe
export const updateRecipe = (id, data) => {
  return api.put(`/recipes/${id}`, data);
};

// Delete recipe
export const deleteRecipe = (id) => {
  return api.delete(`/recipes/${id}`);
};

// Scale recipe
export const scaleRecipe = (id, servings) => {
  return api.get(`/recipes/${id}/scale`, {
    params: { servings },
  });
};
