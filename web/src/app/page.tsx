"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { Package } from "lucide-react";
import { useAppSelector } from "@/store/hooks";
import { selectIsAuthenticated } from "@/store/auth/selectors";

export default function Home() {
  const router = useRouter();
  const isAuthenticated = useAppSelector(selectIsAuthenticated);

  useEffect(() => {
    if (isAuthenticated) {
      router.push("/dashboard");
    } else {
      router.push("/login");
    }
  }, [isAuthenticated, router]);

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-background">
      <Package className="h-16 w-16 text-primary animate-pulse" />
      <p className="mt-4 text-muted-foreground">Carregando...</p>
    </div>
  );
}
