import React, { useEffect, useRef } from "react";
import L from "leaflet";
import "leaflet/dist/leaflet.css";
import { Box } from "@mui/material";

// Fix for Leaflet marker icons
const markerIcon = new L.Icon({
  iconUrl: "/marker-icon.png",
  iconRetinaUrl: "/marker-icon-2x.png",
  shadowUrl: "/marker-shadow.png",
  iconSize: [25, 41],
  iconAnchor: [12, 41],
  popupAnchor: [1, -34],
  shadowSize: [41, 41],
});

const InteractiveMap: React.FC = () => {
  const mapRef = useRef<L.Map | null>(null); // Ref to store the map instance
  const mapContainerRef = useRef<HTMLDivElement | null>(null); // Ref to store the map container

  useEffect(() => {
    // Check if the map is already initialized
    if (mapRef.current || !mapContainerRef.current) return;

    // Initialize the map
    const map = L.map(mapContainerRef.current).setView([20, 0], 2);
    L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
      maxZoom: 18,
    }).addTo(map);

    const destinations = [
      {
        lat: -8.4095,
        lng: 115.1889,
        title: "Bali, Indonesia",
        link: "/blog/bali",
      },
      {
        lat: 48.8566,
        lng: 2.3522,
        title: "Paris, France",
        link: "/blog/paris",
      },
      {
        lat: 40.7128,
        lng: -74.006,
        title: "New York, USA",
        link: "/blog/new-york",
      },
    ];

    destinations.forEach((dest) => {
      L.marker([dest.lat, dest.lng], { icon: markerIcon })
        .addTo(map)
        .bindPopup(
          `<b>${dest.title}</b><br><a href="${dest.link}">Read More</a>`
        );
    });

    // Store the map instance in the ref
    mapRef.current = map;

    // Cleanup function to remove the map when the component unmounts
    return () => {
      if (mapRef.current) {
        mapRef.current.remove();
        mapRef.current = null;
      }
    };
  }, []);

  return (
    <Box
      id="map"
      ref={mapContainerRef}
      sx={{
        height: "500px",
        width: "100%",
        my: 4,
        borderRadius: 2,
        overflow: "hidden",
      }}
    ></Box>
  );
};

export default InteractiveMap;
