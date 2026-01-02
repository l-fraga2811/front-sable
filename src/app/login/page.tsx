"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { login } from "@/store/auth/actions";
import {
  selectIsAuthenticated,
  selectAuthLoading,
  selectAuthError,
} from "@/store/auth/selectors";
import { clearError } from "@/store/auth/reducers";
import { LoginFormCard } from "@/features/login/components/loginFormCard";
import { DottedGlowBackground } from "@/components/ui/dotted-glow-background";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import type { LoginFormData } from "@/types/authTypes";

const loginSchema = z.object({
  email: z.string().email("E-mail inválido"),
  password: z.string().min(6, "Senha deve ter pelo menos 6 caracteres"),
  phone: z.string(),
});

export default function LoginPage() {
  const router = useRouter();
  const dispatch = useAppDispatch();
  const isAuthenticated = useAppSelector(selectIsAuthenticated);
  const isLoading = useAppSelector(selectAuthLoading);
  const error = useAppSelector(selectAuthError);

  const loginForm = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: "",
      password: "",
      phone: "",
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

  const handleSubmit = async (data: LoginFormData) => {
    const result = await dispatch(
      login({ email: data.email, password: data.password })
    );
    if (login.fulfilled.match(result)) {
      toast.success("Login realizado com sucesso!");
      router.push("/dashboard");
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-linear-to-br from-background to-muted p-4">
      <DottedGlowBackground
        glowColor="beige"
        radius={3}
        className="mask-radial-to-60% mask-radial-at-center dark:opacity-100"
      />
      <LoginFormCard
        form={loginForm}
        showPassword={showPassword}
        isLoading={isLoading}
        onToggleShowPassword={() => setShowPassword((current) => !current)}
        onSubmit={handleSubmit}
      />
    </div>
  );
}
