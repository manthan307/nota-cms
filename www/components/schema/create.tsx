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
import { Plus, Trash2 } from "lucide-react";
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
import { useState } from "react";

type FieldType = {
  name: string;
  type: string;
  isRequired: boolean;
};

export function SchemaCreateDialog() {
  const [Fields, SetFields] = useState<FieldType[]>([]);

  function AddField() {
    const DefaultField: FieldType = {
      name: "",
      type: "",
      isRequired: false,
    };

    SetFields((prev) => [...prev, DefaultField]);
  }

  function UpdateField(index: number, key: keyof FieldType, value: any) {
    SetFields((prev) =>
      prev.map((field, i) => (i === index ? { ...field, [key]: value } : field))
    );
  }

  function DeleteField(index: number) {
    SetFields((prev) => prev.filter((_, i) => i !== index));
  }

  return (
    <Dialog>
      <form>
        <DialogTrigger asChild>
          <Button variant="outline">
            <Plus /> Create
          </Button>
        </DialogTrigger>

        <DialogContent className="max-w-[900px] w-full">
          <DialogHeader>
            <DialogTitle>
              <input placeholder="Enter Schema Name" className="outline-0" />
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
                {/* Static ID Row */}
                <TableRow>
                  <TableCell>Id</TableCell>
                  <TableCell>
                    <SelectType name="Type" defaultValue="UUID" disable />
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

                {/* Dynamic Fields */}
                {Fields.map((value, index) => (
                  <TableRow key={index}>
                    <TableCell>
                      <input
                        className="outline-none"
                        value={value.name}
                        onChange={(e) =>
                          UpdateField(index, "name", e.target.value)
                        }
                        placeholder="Field Name"
                      />
                    </TableCell>

                    <TableCell>
                      <SelectType name="Type" defaultValue={value.type} />
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

            <div className="flex justify-end">
              <Button variant="outline" type="button" onClick={AddField}>
                <Plus /> Add Field
              </Button>
            </div>
          </div>

          <DialogFooter>
            <DialogClose asChild>
              <Button variant="outline">Cancel</Button>
            </DialogClose>
            <Button type="submit">Create</Button>
          </DialogFooter>
        </DialogContent>
      </form>
    </Dialog>
  );
}
