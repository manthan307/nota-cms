"use client";

import { Separator } from "@/components/ui/separator";
import { Spinner } from "@/components/ui/spinner";
import { SchemaContext } from "@/context/schema";
import { useContext, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Edit, Trash2 } from "lucide-react";
import { fetch } from "@/lib/instance";
import { toast } from "sonner";
import { SchemaCreateDialog } from "@/components/schema/create";

export default function PageSchemaDashboard() {
  const { schemas, loading, setSchemas } = useContext(SchemaContext);
  const [selectedSchema, setSelectedSchema] = useState<any | null>(null);
  const [deletingId, setDeletingId] = useState<string | null>(null);
  const [editingId, setEditingId] = useState<string | null>(null);

  if (loading) {
    return (
      <div className="h-screen w-screen flex justify-center items-center">
        <Spinner className="size-8" />
      </div>
    );
  }

  if (!schemas || schemas.length === 0) {
    return (
      <div className="flex h-full w-full justify-center items-center">
        <p>No Schema Found.</p>
      </div>
    );
  }

  async function handleDelete(id: string) {
    if (!confirm("Are you sure you want to delete this schema?")) return;
    setDeletingId(id);

    try {
      const res = await fetch.delete(`/api/v1/schemas/delete/${id}`);
      if (res.status === 200) {
        setSchemas((prev) => prev.filter((schema) => schema.ID !== id));
        if (selectedSchema?.ID === id) setSelectedSchema(null);
        toast?.success?.("Schema deleted successfully!");
      } else {
        toast?.error?.("Failed to delete schema.");
      }
    } catch (err) {
      console.error("Error deleting schema:", err);
      toast?.error?.("Something went wrong.");
    } finally {
      setDeletingId(null);
    }
  }

  function handleEdit(schema: any) {
    setEditingId(schema.ID);
    // TODO: open a dialog (like your SchemaCreateDialog) in edit mode
    console.log("Edit schema:", schema);
  }

  return (
    <section className="flex h-full w-full">
      {/* Sidebar */}
      <div className="h-full w-1/3 border-r p-4 overflow-y-auto">
        <h2 className="text-lg font-semibold mb-4">Schemas</h2>
        <div className="flex flex-col gap-2">
          {schemas.map((schema) => (
            <div
              key={schema.ID}
              className={`flex items-center justify-between rounded-md border px-3 py-2 transition-colors ${
                selectedSchema?.ID === schema.ID
                  ? "bg-accent"
                  : "hover:bg-muted/50"
              }`}
            >
              <button
                onClick={() => setSelectedSchema(schema)}
                className="flex-1 text-left truncate"
              >
                {schema.Name}
              </button>
            </div>
          ))}
        </div>
      </div>

      <Separator orientation="vertical" className="mx-2" />

      {/* Details Panel */}
      <div className="flex-1 p-4 overflow-y-auto">
        {!selectedSchema ? (
          <div className="h-full w-full flex items-center justify-center">
            <p className="text-muted-foreground">
              Select a schema to view details.
            </p>
          </div>
        ) : (
          <Card className="shadow-sm">
            <CardHeader>
              <CardTitle className="text-xl flex items-center justify-between">
                {selectedSchema.Name}
                <Button
                  size="default"
                  variant="destructive"
                  onClick={() => handleDelete(selectedSchema.ID)}
                  disabled={deletingId === selectedSchema.ID}
                >
                  {deletingId === selectedSchema.ID ? (
                    <Spinner className="h-4 w-4" />
                  ) : (
                    <>
                      <Trash2 className="h-4 w-4" /> Delete
                    </>
                  )}
                </Button>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Field Name</TableHead>
                    <TableHead>Type</TableHead>
                    <TableHead>Required</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {selectedSchema.Definition?.map((field: any, idx: number) => (
                    <TableRow key={idx}>
                      <TableCell>{field.name}</TableCell>
                      <TableCell>{field.type}</TableCell>
                      <TableCell>{field.isRequired ? "Yes" : "No"}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        )}
      </div>
    </section>
  );
}
