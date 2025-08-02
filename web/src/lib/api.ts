const API_BASE_URL = "http://localhost:8080";

export interface Driver {
  id: string;
  name: string;
  lat: number;
  lon: number;
}

export interface TripRequest {
  rider_id: string;
  driver_id: string;
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

const getAuthHeaders = () => {
  const token = localStorage.getItem("authToken");
  const headers: HeadersInit = {
    "Content-Type": "application/json",
  };
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }
  return headers;
};

export const registerDriver = async (
  name: string,
  lat: number,
  lon: number
): Promise<Driver> => {
  const response = await fetch(`${API_BASE_URL}/drivers`, {
    method: "POST",
    headers: { "Content-Type": "application/json", ...getAuthHeaders() },
    body: JSON.stringify({ name, lat, lon }),
  });
  if (!response.ok) throw new Error("Failed to register driver");
  return response.json();
};

export const findAvailableDrivers = async (
  lat: number,
  lon: number
): Promise<Driver[]> => {
  const response = await fetch(
    `${API_BASE_URL}/drivers/available?lat=${lat}&lon=${lon}`
  );
  if (!response.ok) throw new Error("Failed to find drivers");
  return response.json();
};

export const bookTrip = async (
  tripRequest: TripRequest
): Promise<TripResponse> => {
  const response = await fetch(`${API_BASE_URL}/trips`, {
    method: "POST",
    headers: { "Content-Type": "application/json", ...getAuthHeaders() },
    body: JSON.stringify(tripRequest),
  });
  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to book trip: ${errorText}`);
  }
  return response.json();
};
