import * as React from "react";

import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { InputTypes } from "@/lib/types";
import { FieldType } from "./create";

export type SelectTypeProps = {
  name: string;
  defaultValue?: string;
  disabled?: boolean;
  update?: (index: number, key: keyof FieldType, value: any) => void;
  index?: number;
};
export function SelectType({
  name,
  defaultValue,
  disabled,
  update,
  index,
}: SelectTypeProps) {
  const [value, setValue] = React.useState(defaultValue || "");

  const updateValues = (v: string) => {
    if (update != undefined && index != undefined) update(index, "type", v);
    setValue(v);
  };

  return (
    <Select value={value} onValueChange={updateValues} disabled={disabled}>
      <SelectTrigger className="w-[100px]">
        <SelectValue placeholder="Type" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel>Type</SelectLabel>
          {InputTypes.map((value, index) => (
            <SelectItem value={value} key={index}>
              {value}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>

      {/* Hidden input for form submission */}
      <input type="hidden" name={name} value={value} />
    </Select>
  );
}
