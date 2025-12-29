"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { register } from "@/store/auth/actions";
import { selectIsAuthenticated, selectAuthLoading, selectAuthError } from "@/store/auth/selectors";
import { clearError } from "@/store/auth/reducers";
import { RegisterFormCard } from "@/features/register/components/registerFormCard";
import { DottedGlowBackground } from "@/components/ui/dotted-glow-background";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import type { RegisterFormData } from "@/types/auth";

const registerSchema = z.object({
    username: z.string().min(3, "Usuário deve ter pelo menos 3 caracteres"),
    email: z.string().email("E-mail inválido"),
    phone: z.string(),
    password: z.string().min(6, "Senha deve ter pelo menos 6 caracteres"),
    confirmPassword: z.string().min(6, "Confirme sua senha"),
}).refine((data) => data.password === data.confirmPassword, {
    message: "As senhas não coincidem",
    path: ["confirmPassword"],
});

export default function RegisterPage() {
    const router = useRouter();
    const dispatch = useAppDispatch();
    const isAuthenticated = useAppSelector(selectIsAuthenticated);
    const isLoading = useAppSelector(selectAuthLoading);
    const error = useAppSelector(selectAuthError);

    const registerForm = useForm<RegisterFormData>({
        resolver: zodResolver(registerSchema),
        defaultValues: {
            username: "",
            email: "",
            phone: "",
            password: "",
            confirmPassword: "",
        },
    });

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

    const handleSubmit = async (data: RegisterFormData) => {
        const result = await dispatch(register({ username: data.username, email: data.email, password: data.password }));
        if (register.fulfilled.match(result)) {
            toast.success("Conta criada com sucesso! Faça login para continuar.");
            router.push("/login");
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-linear-to-br from-background to-muted p-4">
            <DottedGlowBackground glowColor="beige" radius={3} className="mask-radial-to-60% mask-radial-at-center dark:opacity-100" />
            <RegisterFormCard
                form={registerForm}
                showPassword={showPassword}
                isLoading={isLoading}
                onToggleShowPassword={() => setShowPassword((current) => !current)}
                onSubmit={handleSubmit}
            />
        </div>
    );
}
