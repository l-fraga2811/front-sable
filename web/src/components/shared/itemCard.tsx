"use client";

import { useState } from "react";
import { toast } from "sonner";
import { Pencil, Trash2, Check, X, DollarSign } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";
import { useAppDispatch } from "@/store/hooks";
import { deleteItem, updateItem } from "@/store/items/actions";
import type { Item } from "@/types";

interface ItemCardProps {
    item: Item;
    onEdit: (item: Item) => void;
}

export function ItemCard({ item, onEdit }: ItemCardProps) {
    const dispatch = useAppDispatch();
    const [showDeleteDialog, setShowDeleteDialog] = useState(false);
    const [isDeleting, setIsDeleting] = useState(false);

    const handleDelete = async () => {
        setIsDeleting(true);
        const result = await dispatch(deleteItem(item.id));
        setIsDeleting(false);
        setShowDeleteDialog(false);

        if (deleteItem.fulfilled.match(result)) {
            toast.success("Item deletado com sucesso!");
        } else {
            toast.error("Erro ao deletar item");
        }
    };

    const handleToggleComplete = async () => {
        const result = await dispatch(
            updateItem({
                id: item.id,
                data: { completed: !item.completed },
            })
        );

        if (updateItem.fulfilled.match(result)) {
            toast.success(item.completed ? "Item marcado como pendente" : "Item marcado como concluído");
        }
    };

    const formatPrice = (price: number) => {
        return new Intl.NumberFormat("pt-BR", {
            style: "currency",
            currency: "BRL",
        }).format(price);
    };

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString("pt-BR", {
            day: "2-digit",
            month: "2-digit",
            year: "numeric",
        });
    };

    return (
        <>
            <Card className={`transition-all hover:shadow-md ${item.completed ? "opacity-75" : ""}`}>
                <CardHeader className="pb-2">
                    <div className="flex items-start justify-between gap-2">
                        <CardTitle className={`text-lg ${item.completed ? "line-through text-muted-foreground" : ""}`}>
                            {item.title}
                        </CardTitle>
                        <Badge variant={item.completed ? "secondary" : "default"}>
                            {item.completed ? "Concluído" : "Pendente"}
                        </Badge>
                    </div>
                </CardHeader>
                <CardContent className="pb-2">
                    {item.description && (
                        <p className="text-sm text-muted-foreground mb-3">{item.description}</p>
                    )}
                    <div className="flex items-center gap-2 text-sm">
                        <DollarSign className="h-4 w-4 text-muted-foreground" />
                        <span className="font-medium">{formatPrice(item.price)}</span>
                    </div>
                    <p className="text-xs text-muted-foreground mt-2">
                        Criado em {formatDate(item.createdAt)}
                    </p>
                </CardContent>
                <CardFooter className="pt-2 gap-2">
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={handleToggleComplete}
                        className="flex-1"
                    >
                        {item.completed ? (
                            <>
                                <X className="h-4 w-4 mr-1" />
                                Desfazer
                            </>
                        ) : (
                            <>
                                <Check className="h-4 w-4 mr-1" />
                                Concluir
                            </>
                        )}
                    </Button>
                    <Button variant="outline" size="sm" onClick={() => onEdit(item)}>
                        <Pencil className="h-4 w-4" />
                    </Button>
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setShowDeleteDialog(true)}
                        className="text-destructive hover:text-destructive"
                    >
                        <Trash2 className="h-4 w-4" />
                    </Button>
                </CardFooter>
            </Card>

            <Dialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Confirmar exclusão</DialogTitle>
                        <DialogDescription>
                            Tem certeza que deseja excluir o item &quot;{item.title}&quot;? Esta ação não pode ser desfeita.
                        </DialogDescription>
                    </DialogHeader>
                    <DialogFooter>
                        <Button
                            variant="outline"
                            onClick={() => setShowDeleteDialog(false)}
                            disabled={isDeleting}
                        >
                            Cancelar
                        </Button>
                        <Button
                            variant="destructive"
                            onClick={handleDelete}
                            disabled={isDeleting}
                        >
                            {isDeleting ? "Excluindo..." : "Excluir"}
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </>
    );
}
