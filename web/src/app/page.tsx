import { BookRide } from "~/components/book-ride";
import { RegisterDriverForm } from "~/components/register-driver-form";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "~/components/ui/tabs";

export default function HomePage() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center bg-gray-100 p-4">
      <h1 className="text-4xl font-bold mb-8">Uber Clone</h1>
      <Tabs defaultValue="book" className="w-[450px]">
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="book">Book a Ride</TabsTrigger>
          <TabsTrigger value="register">Become a Driver</TabsTrigger>
        </TabsList>
        <TabsContent value="book">
          <BookRide />
        </TabsContent>
        <TabsContent value="register">
          <RegisterDriverForm />
        </TabsContent>
      </Tabs>
    </main>
  );
}
