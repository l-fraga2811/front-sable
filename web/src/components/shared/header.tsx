"use client";

import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { LogOut, User, Package } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { logout } from "@/store/auth/actions";
import { selectUser } from "@/store/auth/selectors";

export function Header() {
    const router = useRouter();
    const dispatch = useAppDispatch();
    const user = useAppSelector(selectUser);

    const handleLogout = async () => {
        await dispatch(logout());
        toast.success("Logout realizado com sucesso!");
        router.push("/login");
    };

    const getInitials = (name: string) => {
        return name
            .split(" ")
            .map((n) => n[0])
            .join("")
            .toUpperCase()
            .slice(0, 2);
    };

    return (
        <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60">
            <div className="container flex h-16 items-center justify-between">
                <div className="flex items-center gap-2">
                    <Package className="h-6 w-6 text-primary" />
                    <span className="text-xl font-bold">Items Manager</span>
                </div>

                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                        <Button variant="ghost" className="relative h-10 w-10 rounded-full">
                            <Avatar className="h-10 w-10">
                                <AvatarFallback className="bg-primary text-primary-foreground">
                                    {user?.username ? getInitials(user.username) : "U"}
                                </AvatarFallback>
                            </Avatar>
                        </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent className="w-56" align="end" forceMount>
                        <DropdownMenuLabel className="font-normal">
                            <div className="flex flex-col space-y-1">
                                <p className="text-sm font-medium leading-none">{user?.username}</p>
                                <p className="text-xs leading-none text-muted-foreground">
                                    {user?.email}
                                </p>
                            </div>
                        </DropdownMenuLabel>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem onClick={() => router.push("/dashboard/profile")}>
                            <User className="mr-2 h-4 w-4" />
                            <span>Perfil</span>
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem onClick={handleLogout} className="text-destructive">
                            <LogOut className="mr-2 h-4 w-4" />
                            <span>Sair</span>
                        </DropdownMenuItem>
                    </DropdownMenuContent>
                </DropdownMenu>
            </div>
        </header>
    );
}
