"use client";

import { useState, useEffect, type FormEvent } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { login } from "@/store/auth/actions";
import { selectIsAuthenticated, selectAuthLoading, selectAuthError } from "@/store/auth/selectors";
import { clearError } from "@/store/auth/reducers";
import { LoginFormCard } from "@/features/login/components/loginFormCard";
import { motion } from "framer-motion";
import { DottedGlowBackground } from "@/components/ui/dotted-glow-background";

export default function LoginPage() {
    const router = useRouter();
    const dispatch = useAppDispatch();
    const isAuthenticated = useAppSelector(selectIsAuthenticated);
    const isLoading = useAppSelector(selectAuthLoading);
    const error = useAppSelector(selectAuthError);

    const [email, setEmail] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const [phone, setPhone] = useState<string>("");

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

    const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        if (!email || !password) {
            toast.error("Preencha todos os campos");
            return;
        }

        const result = await dispatch(login({ email, password }));
        if (login.fulfilled.match(result)) {
            toast.success("Login realizado com sucesso!");
            router.push("/dashboard");
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-linear-to-br from-background to-muted p-4">
            <DottedGlowBackground  glowColor="beige" radius={3} className="mask-radial-to-60% mask-radial-at-center dark:opacity-100"/>
            <LoginFormCard
                email={email}
                password={password}
                phone={phone}
                showPassword={showPassword}
                isLoading={isLoading}
                onEmailChange={setEmail}
                onPasswordChange={setPassword}
                onPhoneChange={setPhone}
                onToggleShowPassword={() => setShowPassword((current) => !current)}
                onSubmit={handleSubmit}
            />
        </div>
    );
}
