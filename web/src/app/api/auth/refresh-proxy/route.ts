import { NextRequest, NextResponse } from "next/server";

const BACKEND_URL = "http://localhost:8080/auth/refresh";

export async function POST(req: NextRequest) {
  try {
    const refreshTokenCookie = req.cookies.get("refresh_token");

    if (!refreshTokenCookie) {
      return new NextResponse(
        JSON.stringify({ error: "Unauthorized: No refresh token found." }),
        { status: 401 }
      );
    }

    const backendResponse = await fetch(BACKEND_URL, {
      method: "POST",
      headers: {
        Cookie: `refresh_token=${refreshTokenCookie.value}`,
      },
    });

    const data = await backendResponse.json();

    if (!backendResponse.ok) {
      return new NextResponse(JSON.stringify(data), {
        status: backendResponse.status,
      });
    }

    const response = NextResponse.json(data);

    const newSetCookieHeader = backendResponse.headers.get("Set-Cookie");
    if (newSetCookieHeader) {
      response.headers.set("Set-Cookie", newSetCookieHeader);
    }

    return response;
  } catch (error) {
    console.error("[AUTH_REFRESH_PROXY_ERROR]", error);
    return new NextResponse(
      JSON.stringify({ error: "Internal Server Error" }),
      { status: 500 }
    );
  }
}
