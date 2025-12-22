"use client";

import { useEffect } from "react";
import { User, Mail, Calendar } from "lucide-react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Separator } from "@/components/ui/separator";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { getProfile } from "@/store/auth/actions";
import { selectUser, selectAuthLoading } from "@/store/auth/selectors";

export default function ProfilePage() {
    const dispatch = useAppDispatch();
    const user = useAppSelector(selectUser);
    const isLoading = useAppSelector(selectAuthLoading);

    useEffect(() => {
        dispatch(getProfile());
    }, [dispatch]);

    const getInitials = (name: string) => {
        return name
            .split(" ")
            .map((n) => n[0])
            .join("")
            .toUpperCase()
            .slice(0, 2);
    };

    if (isLoading) {
        return (
            <div className="flex items-center justify-center py-20">
                <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
            </div>
        );
    }

    return (
        <div className="max-w-2xl mx-auto space-y-6">
            <div>
                <h1 className="text-3xl font-bold">Meu Perfil</h1>
                <p className="text-muted-foreground">
                    Visualize suas informações de conta
                </p>
            </div>

            <Card>
                <CardHeader className="flex flex-row items-center gap-4">
                    <Avatar className="h-20 w-20">
                        <AvatarFallback className="bg-primary text-primary-foreground text-2xl">
                            {user?.username ? getInitials(user.username) : "U"}
                        </AvatarFallback>
                    </Avatar>
                    <div>
                        <CardTitle className="text-2xl">{user?.username}</CardTitle>
                        <CardDescription>Usuário do sistema</CardDescription>
                    </div>
                </CardHeader>
                <Separator />
                <CardContent className="pt-6 space-y-4">
                    <div className="flex items-center gap-3">
                        <div className="flex items-center justify-center h-10 w-10 rounded-full bg-muted">
                            <User className="h-5 w-5 text-muted-foreground" />
                        </div>
                        <div>
                            <p className="text-sm text-muted-foreground">Nome de usuário</p>
                            <p className="font-medium">{user?.username}</p>
                        </div>
                    </div>
                    <div className="flex items-center gap-3">
                        <div className="flex items-center justify-center h-10 w-10 rounded-full bg-muted">
                            <Mail className="h-5 w-5 text-muted-foreground" />
                        </div>
                        <div>
                            <p className="text-sm text-muted-foreground">E-mail</p>
                            <p className="font-medium">{user?.email}</p>
                        </div>
                    </div>
                    <div className="flex items-center gap-3">
                        <div className="flex items-center justify-center h-10 w-10 rounded-full bg-muted">
                            <Calendar className="h-5 w-5 text-muted-foreground" />
                        </div>
                        <div>
                            <p className="text-sm text-muted-foreground">ID do usuário</p>
                            <p className="font-medium text-sm font-mono">{user?.id}</p>
                        </div>
                    </div>
                </CardContent>
            </Card>
        </div>
    );
}
