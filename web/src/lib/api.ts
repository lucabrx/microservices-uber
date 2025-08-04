const API_BASE_URL = "http://localhost:8080";

export interface User {
  id: string;
  name: string;
  email: string;
}

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

let getAccessToken: () => string | null = () => null;
let setAccessToken: (token: string | null) => void = () => {};
let logout: () => void = () => {};

export const initializeApi = (
  _getAccessToken: () => string | null,
  _setAccessToken: (token: string | null) => void,
  _logout: () => void
) => {
  getAccessToken = _getAccessToken;
  setAccessToken = _setAccessToken;
  logout = _logout;
};

let refreshPromise: Promise<string | null> | null = null;
const refreshToken = async (): Promise<string | null> => {
  if (refreshPromise) {
    return refreshPromise;
  }
  refreshPromise = new Promise(async (resolve, reject) => {
    try {
      const response = await fetch(`/api/auth/refresh-proxy`, {
        method: "POST",
      });

      if (!response.ok) {
        logout();
        throw new Error("Session expired. Please log in again.");
      }

      const data = await response.json();
      const newAccessToken = data.access_token;
      setAccessToken(newAccessToken);
      resolve(newAccessToken);
    } catch (error) {
      reject(error);
    } finally {
      refreshPromise = null;
    }
  });
  return refreshPromise;
};

const apiClient = async (
  endpoint: string,
  options: RequestInit = {}
): Promise<Response> => {
  const token = getAccessToken();
  const headers = new Headers(options.headers || {});
  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }
  headers.set("Content-Type", "application/json");

  let response = await fetch(`${API_BASE_URL}/${endpoint}`, {
    ...options,
    headers,
  });

  if (response.status === 401) {
    try {
      const newToken = await refreshToken();
      if (newToken) {
        headers.set("Authorization", `Bearer ${newToken}`);
        response = await fetch(`${API_BASE_URL}/${endpoint}`, {
          ...options,
          headers,
        });
      }
    } catch (error) {
      throw error;
    }
  }

  return response;
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

export const getMe = async (): Promise<User> => {
  const response = await apiClient("me");
  if (!response.ok) throw new Error("Failed to fetch user data");
  return response.json();
};

export const registerDriver = async (
  name: string,
  lat: number,
  lon: number
): Promise<Driver> => {
  const response = await apiClient("drivers", {
    method: "POST",
    body: JSON.stringify({ name, lat, lon }),
  });
  if (!response.ok) throw new Error("Failed to register driver");
  return response.json();
};

export const bookTrip = async (
  tripRequest: TripRequest
): Promise<TripResponse> => {
  const response = await apiClient("trips", {
    method: "POST",
    body: JSON.stringify(tripRequest),
  });
  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to book trip: ${errorText}`);
  }
  return response.json();
};
