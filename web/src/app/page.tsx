// src/app/page.tsx
"use client";

import { useState } from "react";
import { Button } from "~/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
import { bookTrip, TripRequest } from "~/lib/api";

export default function HomePage() {
  const [riderId, setRiderId] = useState("rider-123");
  const [isLoading, setIsLoading] = useState(false);

  const handleBookTrip = async () => {
    setIsLoading(true);
    const tripRequest: TripRequest = {
      rider_id: riderId,
      // For this example, we'll use hardcoded coordinates.
      // A real app would get these from a map interface.
      start_lat: 34.06,
      start_lon: -118.26,
      end_lat: 34.07,
      end_lon: -118.27,
    };

    try {
      const result = await bookTrip(tripRequest);
      console.log({
        title: "Trip Booked Successfully!",
        description: `Driver ${
          result.driver_id
        } is on their way. Price: $${result.price.toFixed(2)}`,
      });
    } catch (error) {
      console.log({
        title: "Error Booking Trip",
        description:
          error instanceof Error ? error.message : "An unknown error occurred.",
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <main className="flex min-h-screen items-center justify-center bg-gray-100">
      <Card className="w-[400px]">
        <CardHeader>
          <CardTitle>Book a Ride</CardTitle>
          <CardDescription>Enter your details to find a ride.</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid w-full items-center gap-4">
            <div className="flex flex-col space-y-1.5">
              <Label htmlFor="riderId">Rider ID</Label>
              <Input
                id="riderId"
                placeholder="Enter your Rider ID"
                value={riderId}
                onChange={(e) => setRiderId(e.target.value)}
              />
            </div>
          </div>
        </CardContent>
        <CardFooter>
          <Button
            className="w-full"
            onClick={handleBookTrip}
            disabled={isLoading}
          >
            {isLoading ? "Finding your ride..." : "Book Now"}
          </Button>
        </CardFooter>
      </Card>
    </main>
  );
}
