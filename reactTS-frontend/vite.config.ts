import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    fs: {
      // Allow serving files from the project root and node_modules/leaflet/dist/images
      allow: [
        // Project root
        "C:/Users/choym/OneDrive/Desktop/CVWO Intern/Travel Web Forum/reactTS-frontend",
        // Leaflet marker icons
        "C:/Users/choym/OneDrive/Desktop/CVWO Intern/Travel Web Forum/node_modules/leaflet/dist/images",
      ],
    },
  },
});
