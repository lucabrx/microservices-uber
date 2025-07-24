// src/services/api.ts
const API_BASE_URL = "http://localhost:8080"; // Your Go backend URL

export interface TripRequest {
  rider_id: string;
  start_lat: number;
  start_lon: number;
  end_lat: number;
  end_lon: number;
}

export interface TripResponse {
  id: string;
  rider_id: string;
  driver_id: string;
  status: string;
  price: number;
}

export const bookTrip = async (
  tripRequest: TripRequest
): Promise<TripResponse> => {
  const response = await fetch(`${API_BASE_URL}/trips`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(tripRequest),
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to book trip: ${errorText}`);
  }

  return response.json();
};
