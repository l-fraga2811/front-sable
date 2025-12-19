"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { Header } from "@/components/shared/header";
import { useAppSelector } from "@/store/hooks";
import { selectIsAuthenticated } from "@/store/auth/selectors";

interface DashboardLayoutProps {
    children: React.ReactNode;
}

export default function DashboardLayout({ children }: DashboardLayoutProps) {
    const router = useRouter();
    const isAuthenticated = useAppSelector(selectIsAuthenticated);

    useEffect(() => {
        if (!isAuthenticated) {
            router.push("/login");
        }
    }, [isAuthenticated, router]);

    if (!isAuthenticated) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-background">
            <Header />
            <main className="container py-6">{children}</main>
        </div>
    );
}
