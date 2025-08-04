"use client";

import { useState } from "react";
import { Button } from "./ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "./ui/card";
import { Input } from "./ui/input";
import { Label } from "./ui/label";
import { registerDriver } from "~/lib/api";

export function RegisterDriverForm() {
  const [name, setName] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleRegister = async () => {
    if (!name) {
      console.log({ title: "Please enter a name.", variant: "destructive" });
      return;
    }
    setIsLoading(true);
    try {
      const result = await registerDriver(name, 34.05, -118.25);
      console.log({
        title: "Driver Registered!",
        description: `Welcome, ${result.name}! Your ID is ${result.id}`,
      });
      setName("");
    } catch (error) {
      console.log({ title: "Registration Failed", variant: "destructive" });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Become a Driver</CardTitle>
        <CardDescription>Register to start giving rides.</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="flex flex-col space-y-1.5">
          <Label htmlFor="name">Full Name</Label>
          <Input
            id="name"
            placeholder="John Doe"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </div>
      </CardContent>
      <CardFooter>
        <Button
          className="w-full"
          onClick={handleRegister}
          disabled={isLoading}
        >
          {isLoading ? "Registering..." : "Register Now"}
        </Button>
      </CardFooter>
    </Card>
  );
}
