import type { Metadata } from "next";
import "@fontsource/jetbrains-mono";
import { Toaster } from "@/components/ui/sonner";
import { StoreProvider } from "@/store/provider";
import "./globals.css";

export const metadata: Metadata = {
  title: "Items Manager",
  description: "Gerencie seus items com facilidade",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="pt-BR">
      <body
        className="font-mono antialiased"
      >
        <StoreProvider>
          {children}
          <Toaster position="top-right" richColors />
        </StoreProvider>
      </body>
    </html>
  );
}
