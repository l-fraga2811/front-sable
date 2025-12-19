"use client";

import { useState } from "react";
import { toast } from "sonner";
import { Save, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { createItem, updateItem } from "@/store/items/actions";
import { selectItemsLoading } from "@/store/items/selectors";
import type { Item } from "@/types";

interface ItemFormProps {
    item?: Item | null;
    open: boolean;
    onClose: () => void;
}

export function ItemForm({ item, open, onClose }: ItemFormProps) {
    const dispatch = useAppDispatch();
    const isLoading = useAppSelector(selectItemsLoading);

    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [price, setPrice] = useState("");

    const handleOpenChange = (isOpen: boolean) => {
        if (isOpen && item) {
            setTitle(item.title);
            setDescription(item.description || "");
            setPrice(item.price?.toString() || "");
        } else if (isOpen) {
            setTitle("");
            setDescription("");
            setPrice("");
        }
        if (!isOpen) {
            onClose();
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!title.trim()) {
            toast.error("O título é obrigatório");
            return;
        }

        const data = {
            title: title.trim(),
            description: description.trim(),
            price: parseFloat(price) || 0,
        };

        if (item) {
            const result = await dispatch(updateItem({ id: item.id, data }));
            if (updateItem.fulfilled.match(result)) {
                toast.success("Item atualizado com sucesso!");
                onClose();
            } else {
                toast.error("Erro ao atualizar item");
            }
        } else {
            const result = await dispatch(createItem(data));
            if (createItem.fulfilled.match(result)) {
                toast.success("Item criado com sucesso!");
                onClose();
            } else {
                toast.error("Erro ao criar item");
            }
        }
    };

    return (
        <Dialog open={open} onOpenChange={handleOpenChange}>
            <DialogContent className="sm:max-w-md">
                <DialogHeader>
                    <DialogTitle>{item ? "Editar Item" : "Novo Item"}</DialogTitle>
                </DialogHeader>
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div className="space-y-2">
                        <Label htmlFor="title">Título *</Label>
                        <Input
                            id="title"
                            value={title}
                            onChange={(e) => setTitle(e.target.value)}
                            placeholder="Digite o título do item"
                            disabled={isLoading}
                        />
                    </div>
                    <div className="space-y-2">
                        <Label htmlFor="description">Descrição</Label>
                        <Input
                            id="description"
                            value={description}
                            onChange={(e) => setDescription(e.target.value)}
                            placeholder="Digite uma descrição (opcional)"
                            disabled={isLoading}
                        />
                    </div>
                    <div className="space-y-2">
                        <Label htmlFor="price">Preço (R$)</Label>
                        <Input
                            id="price"
                            type="number"
                            step="0.01"
                            min="0"
                            value={price}
                            onChange={(e) => setPrice(e.target.value)}
                            placeholder="0,00"
                            disabled={isLoading}
                        />
                    </div>
                    <div className="flex justify-end gap-2 pt-4">
                        <Button type="button" variant="outline" onClick={onClose} disabled={isLoading}>
                            <X className="h-4 w-4 mr-1" />
                            Cancelar
                        </Button>
                        <Button type="submit" disabled={isLoading}>
                            {isLoading ? (
                                <span className="flex items-center gap-2">
                                    <span className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                                    Salvando...
                                </span>
                            ) : (
                                <>
                                    <Save className="h-4 w-4 mr-1" />
                                    Salvar
                                </>
                            )}
                        </Button>
                    </div>
                </form>
            </DialogContent>
        </Dialog>
    );
}
