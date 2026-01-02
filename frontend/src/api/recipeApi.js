import api from "./axios";

export const createRecipe = (data) => {
  return api.post("/recipes", data);
};

export const getMyRecipes = () => {
  return api.get("/recipes");
};

export const getRecipeById = (id) => {
  return api.get(`/recipes/${id}`);
};

export const updateRecipe = (id, data) => {
  return api.put(`/recipes/${id}`, data);
};

export const deleteRecipe = (id) => {
  return api.delete(`/recipes/${id}`);
};

export const scaleRecipe = (id, servings) => {
  return api.get(`/recipes/${id}/scale`, {
    params: { servings },
  });
};
