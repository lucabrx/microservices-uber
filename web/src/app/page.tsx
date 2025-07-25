"use client";

import { useEffect, useState } from "react";
import { BookRide } from "~/components/book-ride";
import { RegisterDriverForm } from "~/components/register-driver-form";
import { AvailableDriversSidebar } from "~/components/available-drivers-sidebar";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "~/components/ui/tabs";
import { Driver } from "~/lib/api";

export default function HomePage() {
  const [drivers, setDrivers] = useState<Driver[]>([]);

  useEffect(() => {
    const ws = new WebSocket(
      "ws://localhost:8080/ws/drivers/available?lat=34.06&lon=-118.26"
    );

    ws.onopen = () => {
      console.log("WebSocket connection established");
    };

    ws.onmessage = (event) => {
      const driverData = JSON.parse(event.data);
      setDrivers(driverData);
    };

    ws.onclose = () => {
      console.log("WebSocket connection closed");
    };

    ws.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    return () => {
      ws.close();
    };
  }, []);

  return (
    <main className="flex min-h-screen flex-col items-center justify-center bg-gray-100 p-4">
      <h1 className="text-4xl font-bold mb-8">Uber Clone</h1>
      <div className="grid md:grid-cols-2 gap-8 w-full max-w-6xl">
        <div className="md:col-span-1">
          <AvailableDriversSidebar drivers={drivers} />
        </div>

        <div className="md:col-span-1">
          <Tabs defaultValue="book" className="w-full">
            <TabsList className="grid w-full grid-cols-2">
              <TabsTrigger value="book">Book a Ride</TabsTrigger>
              <TabsTrigger value="register">Become a Driver</TabsTrigger>
            </TabsList>
            <TabsContent value="book">
              <BookRide drivers={drivers} />
            </TabsContent>
            <TabsContent value="register">
              <RegisterDriverForm />
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </main>
  );
}
