// import BASE_URL from "../config";
import { Category } from "../models/models";

const BASE_URL = import.meta.env.VITE_API_URL;
console.log(BASE_URL);

const fetchCategories = async (): Promise<Category[]> => {
  try {
    const response = await fetch(`${BASE_URL}/api/categories`);
    if (!response.ok) {
      throw new Error("Failed to fetch categories");
    }
    return await response.json();
  } catch (error) {
    console.error("Error fetching categories:", error);
    return [];
  }
};

const topics = await fetchCategories();

const TopicsList = () => {
  return (
    <div className="bg-white p-4 rounded-lg shadow">
      <h3 className="font-semibold text-lg mb-4">Topics For You</h3>
      <div className="flex flex-wrap gap-2">
        {topics.map((topic, index) => (
          <span
            key={index}
            className="px-3 py-1 bg-gray-200 text-gray-700 text-sm rounded-full cursor-pointer"
          >
            {topic.name}
          </span>
        ))}
      </div>
    </div>
  );
};

export default TopicsList;
