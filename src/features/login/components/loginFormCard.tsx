"use client";

import Link from "next/link";
import type { FormEvent } from "react";
import { Eye, EyeOff, LogIn } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { IMaskInput } from "react-imask";
import { motion } from "framer-motion";

type LoginFormCardProps = {
    email: string;
    password: string;
    phone: string;
    showPassword: boolean;
    isLoading: boolean;
    onEmailChange: (value: string) => void;
    onPasswordChange: (value: string) => void;
    onPhoneChange: (value: string) => void;
    onToggleShowPassword: () => void;
    onSubmit: (e: FormEvent<HTMLFormElement>) => void;
};

export function LoginFormCard({
    email,
    password,
    phone,
    showPassword,
    isLoading,
    onEmailChange,
    onPasswordChange,
    onPhoneChange,
    onToggleShowPassword,
    onSubmit,
}: LoginFormCardProps) {
    return (
        <motion.div className="flex items-center justify-center w-full h-full" initial={{ opacity: 0 }} animate={{ opacity: 100 }} transition={{ duration: 1.2 }}>
            <Card className="w-full max-w-md">
                <CardHeader className="space-y-1">
                    <CardTitle className="text-2xl font-bold text-center">Entrar</CardTitle>
                    <CardDescription className="text-center">Digite suas credenciais para acessar sua conta</CardDescription>
                </CardHeader>
                <form onSubmit={onSubmit}>
                    <CardContent className="space-y-4">
                        <div className="space-y-2">
                            <Label htmlFor="email">E-mail</Label>
                            <Input
                                id="email"
                                type="email"
                                placeholder="Digite seu e-mail"
                                value={email}
                                onChange={(e) => onEmailChange(e.target.value)}
                                disabled={isLoading}
                            />
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="phone">Telefone</Label>
                            <IMaskInput
                                id="phone"
                                mask="(00) 00000-0000"
                                placeholder="(XX) XXXXX-XXXX"
                                value={phone}
                                onAccept={(value: string) => onPhoneChange(value)}
                                disabled={isLoading}
                                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
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
                                    onChange={(e) => onPasswordChange(e.target.value)}
                                    disabled={isLoading}
                                />
                                <Button
                                    type="button"
                                    variant="ghost"
                                    size="icon"
                                    className="absolute right-0 top-0 h-full px-3 hover:bg-transparent"
                                    onClick={onToggleShowPassword}
                                >
                                    {showPassword ? (
                                        <EyeOff className="h-4 w-4 text-muted-foreground" />
                                    ) : (
                                        <Eye className="h-4 w-4 text-muted-foreground" />
                                    )}
                                </Button>
                            </div>
                        </div>
                    </CardContent>
                    <CardFooter className="flex flex-col space-y-4">
                        <Button type="submit" className="w-full mt-4" disabled={isLoading}>
                            {isLoading ? (
                                <span className="flex items-center gap-2">
                                    <span className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                                    Entrando...
                                </span>
                            ) : (
                                <span className="flex items-center gap-2">
                                    <LogIn className="h-4 w-4" />
                                    Entrar
                                </span>
                            )}
                        </Button>
                        <p className="text-sm text-muted-foreground text-center">
                            Não tem uma conta?{" "}
                            <Link href="/register" className="text-primary hover:underline font-medium">
                                Registre-se
                            </Link>
                        </p>
                    </CardFooter>
                </form>
            </Card>
        </motion.div>
    );
}
