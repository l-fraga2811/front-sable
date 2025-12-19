"use client";

import { useEffect, useState } from "react";
import { Plus, Package, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ItemCard } from "@/components/shared/itemCard";
import { ItemForm } from "@/components/shared/itemForm";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { fetchItems } from "@/store/items/actions";
import { selectItems, selectItemsLoading } from "@/store/items/selectors";
import type { Item } from "@/types";

export default function DashboardPage() {
    const dispatch = useAppDispatch();
    const items = useAppSelector(selectItems);
    const isLoading = useAppSelector(selectItemsLoading);

    const [showForm, setShowForm] = useState(false);
    const [editingItem, setEditingItem] = useState<Item | null>(null);

    useEffect(() => {
        dispatch(fetchItems());
    }, [dispatch]);

    const handleEdit = (item: Item) => {
        setEditingItem(item);
        setShowForm(true);
    };

    const handleCloseForm = () => {
        setShowForm(false);
        setEditingItem(null);
    };

    const handleNewItem = () => {
        setEditingItem(null);
        setShowForm(true);
    };

    if (isLoading && items.length === 0) {
        return (
            <div className="flex flex-col items-center justify-center py-20">
                <Loader2 className="h-8 w-8 animate-spin text-primary" />
                <p className="mt-4 text-muted-foreground">Carregando items...</p>
            </div>
        );
    }

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold">Meus Items</h1>
                    <p className="text-muted-foreground">
                        Gerencie seus items de forma simples e eficiente
                    </p>
                </div>
                <Button onClick={handleNewItem}>
                    <Plus className="h-4 w-4 mr-2" />
                    Novo Item
                </Button>
            </div>

            {items.length === 0 ? (
                <div className="flex flex-col items-center justify-center py-20 border-2 border-dashed rounded-lg">
                    <Package className="h-12 w-12 text-muted-foreground" />
                    <h3 className="mt-4 text-lg font-medium">Nenhum item encontrado</h3>
                    <p className="mt-2 text-sm text-muted-foreground">
                        Comece criando seu primeiro item
                    </p>
                    <Button className="mt-4" onClick={handleNewItem}>
                        <Plus className="h-4 w-4 mr-2" />
                        Criar primeiro item
                    </Button>
                </div>
            ) : (
                <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                    {items.map((item) => (
                        <ItemCard key={item.id} item={item} onEdit={handleEdit} />
                    ))}
                </div>
            )}

            <ItemForm item={editingItem} open={showForm} onClose={handleCloseForm} />
        </div>
    );
}
