"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { toast } from "sonner";
import { UserPlus, Eye, EyeOff } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { register } from "@/store/auth/actions";
import { selectIsAuthenticated, selectAuthLoading, selectAuthError } from "@/store/auth/selectors";
import { clearError } from "@/store/auth/reducers";

export default function RegisterPage() {
    const router = useRouter();
    const dispatch = useAppDispatch();
    const isAuthenticated = useAppSelector(selectIsAuthenticated);
    const isLoading = useAppSelector(selectAuthLoading);
    const error = useAppSelector(selectAuthError);

    const [username, setUsername] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [showPassword, setShowPassword] = useState(false);

    useEffect(() => {
        if (isAuthenticated) {
            router.push("/dashboard");
        }
    }, [isAuthenticated, router]);

    useEffect(() => {
        if (error) {
            toast.error(error);
            dispatch(clearError());
        }
    }, [error, dispatch]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!username || !email || !password || !confirmPassword) {
            toast.error("Preencha todos os campos");
            return;
        }

        if (password !== confirmPassword) {
            toast.error("As senhas não coincidem");
            return;
        }

        if (password.length < 6) {
            toast.error("A senha deve ter pelo menos 6 caracteres");
            return;
        }

        const result = await dispatch(register({ username, email, password }));
        if (register.fulfilled.match(result)) {
            toast.success("Conta criada com sucesso! Faça login para continuar.");
            router.push("/login");
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-linear-to-br from-background to-muted p-4">
            <Card className="w-full max-w-md">
                <CardHeader className="space-y-1">
                    <CardTitle className="text-2xl font-bold text-center">Criar Conta</CardTitle>
                    <CardDescription className="text-center">
                        Preencha os dados abaixo para criar sua conta
                    </CardDescription>
                </CardHeader>
                <form onSubmit={handleSubmit}>
                    <CardContent className="space-y-4">
                        <div className="space-y-2">
                            <Label htmlFor="username">Usuário</Label>
                            <Input
                                id="username"
                                type="text"
                                placeholder="Digite seu usuário"
                                value={username}
                                onChange={(e) => setUsername(e.target.value)}
                                disabled={isLoading}
                            />
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="email">E-mail</Label>
                            <Input
                                id="email"
                                type="email"
                                placeholder="Digite seu e-mail"
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                disabled={isLoading}
                            />
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="password">Senha</Label>
                            <div className="relative">
                                <Input
                                    id="password"
                                    type={showPassword ? "text" : "password"}
                                    placeholder="Digite sua senha"
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                    disabled={isLoading}
                                />
                                <Button
                                    type="button"
                                    variant="ghost"
                                    size="icon"
                                    className="absolute right-0 top-0 h-full px-3 hover:bg-transparent"
                                    onClick={() => setShowPassword(!showPassword)}
                                >
                                    {showPassword ? (
                                        <EyeOff className="h-4 w-4 text-muted-foreground" />
                                    ) : (
                                        <Eye className="h-4 w-4 text-muted-foreground" />
                                    )}
                                </Button>
                            </div>
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="confirmPassword">Confirmar Senha</Label>
                            <Input
                                id="confirmPassword"
                                type={showPassword ? "text" : "password"}
                                placeholder="Confirme sua senha"
                                value={confirmPassword}
                                onChange={(e) => setConfirmPassword(e.target.value)}
                                disabled={isLoading}
                            />
                        </div>
                    </CardContent>
                    <CardFooter className="flex flex-col space-y-4">
                        <Button type="submit" className="w-full" disabled={isLoading}>
                            {isLoading ? (
                                <span className="flex items-center gap-2">
                                    <span className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                                    Criando conta...
                                </span>
                            ) : (
                                <span className="flex items-center gap-2">
                                    <UserPlus className="h-4 w-4" />
                                    Criar Conta
                                </span>
                            )}
                        </Button>
                        <p className="text-sm text-muted-foreground text-center">
                            Já tem uma conta?{" "}
                            <Link href="/login" className="text-primary hover:underline font-medium">
                                Entrar
                            </Link>
                        </p>
                    </CardFooter>
                </form>
            </Card>
        </div>
    );
}
