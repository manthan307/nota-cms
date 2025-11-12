"use client";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Plus, Trash2, Loader2 } from "lucide-react";
import { Checkbox } from "../ui/checkbox";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../ui/table";
import { SelectType } from "./selectType";
import { useContext, useState } from "react";
import { fetch } from "@/lib/instance";
import { SchemaContext } from "@/context/schema";

export type FieldType = {
  name: string;
  type: string;
  isRequired: boolean;
};

export function SchemaCreateDialog({
  name,
  fields,
  children,
}: {
  name?: string;
  fields?: FieldType[];
  children: React.ReactNode;
}) {
  const [Fields, SetFields] = useState<FieldType[]>(fields || []);
  const [schemaName, setSchemaName] = useState<string>(name || "");
  const [error, SetError] = useState<string>();
  const [loading, setLoading] = useState<boolean>(false);
  const [open, setOpen] = useState(false);
  const { addSchema, refreshSchemas } = useContext(SchemaContext);

  function AddField() {
    SetFields((prev) => [...prev, { name: "", type: "", isRequired: false }]);
  }

  function UpdateField(index: number, key: keyof FieldType, value: any) {
    SetFields((prev) =>
      prev.map((field, i) => (i === index ? { ...field, [key]: value } : field))
    );
  }

  function DeleteField(index: number) {
    SetFields((prev) => prev.filter((_, i) => i !== index));
  }

  function BuildSchema() {
    if (schemaName.trim() === "") {
      SetError("Please set schema name.");
      return null;
    }

    for (const field of Fields) {
      if (field.name.trim() === "") {
        SetError("Please set name for all fields.");
        return null;
      }
      if (field.type.trim() === "") {
        SetError("Please set type for all fields.");
        return null;
      }
    }

    SetError(undefined);
    return Fields.map((field) => ({
      name: field.name,
      type: field.type,
      isRequired: field.isRequired,
    }));
  }

  async function HandleSubmit(e: React.FormEvent) {
    e.preventDefault();
    const schema = BuildSchema();
    if (!schema) return;

    setLoading(true);
    try {
      const res = await fetch.post("/api/v1/schemas/create", {
        name: schemaName,
        definition: schema,
      });

      if (res.status === 200) {
        addSchema({ name: schemaName, definition: schema });
        refreshSchemas();
        setOpen(false); // ✅ close dialog
      } else {
        SetError(res.data?.error || "Failed to create schema.");
      }
    } catch (err) {
      console.error("Failed to create schema:", err);
      SetError("Something went wrong.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <Dialog
      open={open}
      onOpenChange={(v) => {
        setOpen(v);
        if (!v) {
          setSchemaName("");
          SetFields([]);
          SetError(undefined);
        }
      }}
    >
      <DialogTrigger asChild>{children}</DialogTrigger>

      <DialogContent className="max-w-[900px] w-full">
        <form onSubmit={HandleSubmit}>
          <DialogHeader>
            <DialogTitle>
              <input
                placeholder="Enter Schema Name"
                className="outline-0 w-full"
                value={schemaName}
                onChange={(e) => setSchemaName(e.target.value)}
              />
            </DialogTitle>
            <DialogDescription>
              Add fields you need in your schema.
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-4">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Type</TableHead>
                  <TableHead>Required</TableHead>
                  <TableHead></TableHead>
                </TableRow>
              </TableHeader>

              <TableBody>
                <TableRow>
                  <TableCell>Id</TableCell>
                  <TableCell>
                    <SelectType name="Type" defaultValue="UUID" disabled />
                  </TableCell>
                  <TableCell>
                    <Checkbox defaultChecked disabled />
                  </TableCell>
                  <TableCell>
                    <Button variant="ghost" disabled>
                      <Trash2 className="text-red-500" />
                    </Button>
                  </TableCell>
                </TableRow>

                {Fields.map((value, index) => (
                  <TableRow key={index}>
                    <TableCell>
                      <input
                        className="outline-none w-full"
                        value={value.name}
                        onChange={(e) =>
                          UpdateField(index, "name", e.target.value)
                        }
                        placeholder="Field Name"
                      />
                    </TableCell>

                    <TableCell>
                      <SelectType
                        name="Type"
                        defaultValue={value.type}
                        update={UpdateField}
                        index={index}
                      />
                    </TableCell>

                    <TableCell>
                      <Checkbox
                        checked={value.isRequired}
                        onCheckedChange={(checked) =>
                          UpdateField(index, "isRequired", Boolean(checked))
                        }
                      />
                    </TableCell>

                    <TableCell>
                      <Button
                        variant="ghost"
                        onClick={() => DeleteField(index)}
                      >
                        <Trash2 className="text-red-500" />
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>

            {error && <p className="text-red-500">{error}</p>}

            <div className="flex justify-end my-5">
              <Button variant="outline" type="button" onClick={AddField}>
                <Plus /> Add Field
              </Button>
            </div>
          </div>

          <DialogFooter>
            <DialogClose asChild>
              <Button variant="outline" type="button" disabled={loading}>
                Cancel
              </Button>
            </DialogClose>
            <Button type="submit" disabled={loading}>
              {loading ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  {name == undefined ? "Creating..." : "Updating..."}
                </>
              ) : name == undefined ? (
                "Create"
              ) : (
                "Update"
              )}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
