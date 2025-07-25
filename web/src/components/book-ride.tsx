"use client";

import { useState } from "react";
import { Car, User } from "lucide-react";
import { Button } from "./ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "./ui/card";
import { bookTrip, Driver, TripResponse } from "~/lib/api";

type View = "drivers" | "confirmation";

interface BookRideProps {
  drivers: Driver[];
}

export function BookRide({ drivers }: BookRideProps) {
  const [view, setView] = useState<View>("drivers");
  const [trip, setTrip] = useState<TripResponse | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [selectedDriver, setSelectedDriver] = useState<Driver | null>(null);

  const handleBookTrip = async () => {
    if (!selectedDriver) return;
    setIsLoading(true);
    try {
      const result = await bookTrip({
        rider_id: "rider-xyz",
        driver_id: selectedDriver.id,
        start_lat: 34.06,
        start_lon: -118.26,
        end_lat: 34.07,
        end_lon: -118.27,
      });
      setTrip(result);
      setView("confirmation");
      console.log({
        title: "Trip Booked!",
        description: `Your driver is on the way.`,
      });
    } catch (error) {
      console.log({
        title: "Error Booking Trip",
        description: error instanceof Error ? error.message : "",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  const reset = () => {
    setView("drivers");
    setTrip(null);
    setSelectedDriver(null);
  };

  if (view === "confirmation" && trip) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="text-green-600">Trip Confirmed!</CardTitle>
          <CardDescription>Your driver is on the way.</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <p>
            <strong>Trip ID:</strong> {trip.id}
          </p>
          <p>
            <strong>Assigned Driver:</strong> {trip.driver_id}
          </p>
          <p>
            <strong>Status:</strong> {trip.status}
          </p>
          <p>
            <strong>Price:</strong> ${trip.price.toFixed(2)}
          </p>
        </CardContent>
        <CardFooter>
          <Button onClick={reset} className="w-full">
            Book Another Ride
          </Button>
        </CardFooter>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Available Drivers Near You</CardTitle>
        <CardDescription>
          {drivers?.length > 0
            ? "Please select a driver to book your trip."
            : "Searching for drivers..."}
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-2">
        {drivers?.map((driver) => (
          <div
            key={driver.id}
            className={`flex items-center justify-between p-2 border rounded-md cursor-pointer ${
              selectedDriver?.id === driver.id
                ? "bg-blue-100 border-blue-500"
                : ""
            }`}
            onClick={() => setSelectedDriver(driver)}
          >
            <div className="flex items-center gap-2">
              <Car className="h-5 w-5 text-gray-600" />
              <span>{driver.name}</span>
            </div>
          </div>
        ))}
      </CardContent>
      <CardFooter className="flex-col gap-2">
        <Button
          className="w-full"
          onClick={handleBookTrip}
          disabled={isLoading || !selectedDriver || drivers?.length === 0}
        >
          {isLoading ? "Booking..." : "Book Ride Now"}
        </Button>
        <Button variant="outline" className="w-full" onClick={reset}>
          Cancel
        </Button>
      </CardFooter>
    </Card>
  );
}
