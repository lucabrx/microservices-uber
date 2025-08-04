"use client";

import { useEffect } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { useAuth } from "~/components/auth-context";
import { getMe, initializeApi } from "~/lib/api";

export default function AuthCallback() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { login, logout } = useAuth();

  useEffect(() => {
    const token = searchParams.get("token");
    if (token) {
      initializeApi(
        () => token,
        () => {},
        logout
      );

      getMe()
        .then((user) => {
          login(token, user);
          router.push("/");
        })
        .catch((err) => {
          console.error("Failed to fetch user after login:", err);
          router.push("/");
        });
    } else {
      router.push("/");
    }
  }, [searchParams, router, login, logout]);

  return (
    <div className="flex min-h-screen flex-col items-center justify-center">
      <p>Authenticating, please wait...</p>
    </div>
  );
}
