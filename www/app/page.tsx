"use client";

import { LoginForm } from "@/components/login-form";
import { Verify } from "@/lib/verify";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function Home() {
  const router = useRouter();
  useEffect(() => {
    (async () => {
      const data = await Verify();
      if (data.auth) {
        router.replace("/dashboard");
      }
    })();
  }, []);

  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm">
        <LoginForm />
      </div>
    </div>
  );
}
