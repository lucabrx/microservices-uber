"use client";

import { Car } from "lucide-react";
import { Driver } from "~/lib/api";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "./ui/card";

interface AvailableDriversSidebarProps {
  drivers: Driver[];
}

export function AvailableDriversSidebar({
  drivers,
}: AvailableDriversSidebarProps) {
  return (
    <Card className="hidden md:block">
      <CardHeader>
        <CardTitle>Live Driver Feed</CardTitle>
        <CardDescription>
          Available drivers will appear here in real-time.
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-2">
        {drivers?.length > 0 ? (
          drivers?.map((driver) => (
            <div
              key={driver.id}
              className="flex items-center justify-between p-2 border rounded-md"
            >
              <div className="flex items-center gap-2">
                <Car className="h-5 w-5 text-gray-600" />
                <span>{driver.name} is available</span>
              </div>
            </div>
          ))
        ) : (
          <p className="text-sm text-gray-500">
            Searching for available drivers...
          </p>
        )}
      </CardContent>
    </Card>
  );
}
